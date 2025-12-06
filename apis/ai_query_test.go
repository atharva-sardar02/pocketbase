package apis_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/stretchr/testify/require"
)

// getSuperuserToken gets a superuser auth token for testing
func getSuperuserToken(t *testing.T, app *tests.TestApp) string {
	authRecord, err := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")
	if err != nil {
		// If not found, create one
		collection, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		require.NoError(t, err)
		superuser := core.NewRecord(collection)
		superuser.SetEmail("test@example.com")
		superuser.SetPassword("1234567890")
		err = app.Save(superuser)
		require.NoError(t, err)
		authRecord = superuser
	}
	require.NotNil(t, authRecord)
	// Generate a new auth token (JWT)
	token, err := authRecord.NewAuthToken()
	require.NoError(t, err)
	return token
}

func setupTestAppWithAI(t *testing.T, mockOpenAIServer *httptest.Server) *tests.TestApp {
	app, err := tests.NewTestApp()
	require.NoError(t, err)

	// Enable AI settings
	settings := app.Settings()
	settings.AI.Enabled = true
	settings.AI.Provider = "openai"
	if mockOpenAIServer != nil {
		settings.AI.BaseURL = mockOpenAIServer.URL
	} else {
		settings.AI.BaseURL = "https://api.openai.com/v1"
	}
	settings.AI.APIKey = "test-key"
	settings.AI.Model = "gpt-4o-mini"
	settings.AI.Temperature = 0.1
	app.Save(settings)

	// Create a test superuser if one doesn't exist
	_, err = app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "test@example.com")
	if err != nil {
		collection, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		require.NoError(t, err)
		require.NotNil(t, collection)
		
		superuser := core.NewRecord(collection)
		superuser.SetEmail("test@example.com")
		superuser.SetPassword("1234567890")
		err = app.Save(superuser)
		require.NoError(t, err)
	}

	return app
}

func TestAIQueryAPI_Success(t *testing.T) {
	// Create mock OpenAI server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": `status = "active"`,
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	app := setupTestAppWithAI(t, mockServer)
	defer app.Cleanup()

	// Create test collection
	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "status"})
	collection.ListRule = new(string) // Public collection
	*collection.ListRule = ""
	app.Save(collection)

	// Create test record
	record := core.NewRecord(collection)
	record.Set("status", "active")
	app.Save(record)

	scenarios := []tests.ApiScenario{
		{
			Name:   "successful query without execution",
			Method: http.MethodPost,
			URL:    "/api/ai/query",
			Body: strings.NewReader(`{
				"collection": "posts",
				"query": "active posts",
				"execute": false
			}`),
			Headers: map[string]string{
				"Authorization": getSuperuserToken(t, app),
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"filter":"status = \"active\""`,
			},
		},
		{
			Name:   "successful query with execution",
			Method: http.MethodPost,
			URL:    "/api/ai/query",
			Body: strings.NewReader(`{
				"collection": "posts",
				"query": "active posts",
				"execute": true
			}`),
			Headers: map[string]string{
				"Authorization": getSuperuserToken(t, app),
			},
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"filter":"status = \"active\""`,
				`"results":[{`,
				`"status":"active"`,
				`"totalItems":1`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.TestAppFactory = func(tb testing.TB) *tests.TestApp {
			return app
		}
		scenario.DisableTestAppCleanup = true
		scenario.Test(t)
	}
}

func TestAIQueryAPI_Unauthorized(t *testing.T) {
	app := setupTestAppWithAI(t, nil)
	defer app.Cleanup()

	scenarios := []tests.ApiScenario{
		{
			Name:           "no auth token",
			Method:         http.MethodPost,
			URL:            "/api/ai/query",
			Body:          strings.NewReader(`{"collection":"posts","query":"test"}`),
			ExpectedStatus: 401,
			ExpectedContent: []string{`"data":{}`},
		},
	}

	for _, scenario := range scenarios {
		scenario.TestAppFactory = func(tb testing.TB) *tests.TestApp {
			return app
		}
		scenario.DisableTestAppCleanup = true
		scenario.Test(t)
	}
}

func TestAIQueryAPI_AIDisabled(t *testing.T) {
	app, err := tests.NewTestApp()
	require.NoError(t, err)
	defer app.Cleanup()

	// AI settings disabled (default)
	settings := app.Settings()
	settings.AI.Enabled = false
	app.Save(settings)

	scenarios := []tests.ApiScenario{
		{
			Name:   "AI feature disabled",
			Method: http.MethodPost,
			URL:    "/api/ai/query",
			Body: strings.NewReader(`{
				"collection": "posts",
				"query": "test"
			}`),
			Headers: map[string]string{
				"Authorization": getSuperuserToken(t, app),
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message":"AI Query feature is not enabled."`},
		},
	}

	for _, scenario := range scenarios {
		scenario.TestAppFactory = func(tb testing.TB) *tests.TestApp {
			return app
		}
		scenario.DisableTestAppCleanup = true
		scenario.Test(t)
	}
}

func TestAIQueryAPI_InvalidCollection(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": `status = "active"`,
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	app := setupTestAppWithAI(t, mockServer)
	defer app.Cleanup()

	scenarios := []tests.ApiScenario{
		{
			Name:   "collection not found",
			Method: http.MethodPost,
			URL:    "/api/ai/query",
			Body: strings.NewReader(`{
				"collection": "nonexistent",
				"query": "test"
			}`),
			Headers: map[string]string{
				"Authorization": getSuperuserToken(t, app),
			},
			ExpectedStatus:  404,
			ExpectedContent: []string{`"message":"Collection not found."`},
		},
	}

	for _, scenario := range scenarios {
		scenario.TestAppFactory = func(tb testing.TB) *tests.TestApp {
			return app
		}
		scenario.DisableTestAppCleanup = true
		scenario.Test(t)
	}
}

func TestAIQueryAPI_EmptyQuery(t *testing.T) {
	app := setupTestAppWithAI(t, nil)
	defer app.Cleanup()

	scenarios := []tests.ApiScenario{
		{
			Name:   "empty query",
			Method: http.MethodPost,
			URL:    "/api/ai/query",
			Body: strings.NewReader(`{
				"collection": "posts",
				"query": ""
			}`),
			Headers: map[string]string{
				"Authorization": getSuperuserToken(t, app),
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message":"Query is required."`},
		},
		{
			Name:   "missing query",
			Method: http.MethodPost,
			URL:    "/api/ai/query",
			Body: strings.NewReader(`{
				"collection": "posts"
			}`),
			Headers: map[string]string{
				"Authorization": getSuperuserToken(t, app),
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message":"Query is required."`},
		},
	}

	for _, scenario := range scenarios {
		scenario.TestAppFactory = func(tb testing.TB) *tests.TestApp {
			return app
		}
		scenario.DisableTestAppCleanup = true
		scenario.Test(t)
	}
}

func TestAIQueryAPI_ValidationError(t *testing.T) {
	// Mock server that returns invalid filter
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": `invalid_field = "value"`, // Field doesn't exist
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	app := setupTestAppWithAI(t, mockServer)
	defer app.Cleanup()

	// Create test collection
	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "status"})
	collection.ListRule = new(string)
	*collection.ListRule = ""
	app.Save(collection)

	scenarios := []tests.ApiScenario{
		{
			Name:   "invalid filter from LLM",
			Method: http.MethodPost,
			URL:    "/api/ai/query",
			Body: strings.NewReader(`{
				"collection": "posts",
				"query": "test query"
			}`),
			Headers: map[string]string{
				"Authorization": getSuperuserToken(t, app),
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message":"Generated filter is invalid."`},
		},
	}

	for _, scenario := range scenarios {
		scenario.TestAppFactory = func(tb testing.TB) *tests.TestApp {
			return app
		}
		scenario.DisableTestAppCleanup = true
		scenario.Test(t)
	}
}

func TestAIQueryAPI_LLMError(t *testing.T) {
	// Mock server that returns error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		response := map[string]string{
			"error": "Invalid API key",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	app := setupTestAppWithAI(t, mockServer)
	defer app.Cleanup()

	// Create test collection
	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "status"})
	collection.ListRule = new(string)
	*collection.ListRule = ""
	app.Save(collection)

	scenarios := []tests.ApiScenario{
		{
			Name:   "LLM API error",
			Method: http.MethodPost,
			URL:    "/api/ai/query",
			Body: strings.NewReader(`{
				"collection": "posts",
				"query": "test query"
			}`),
			Headers: map[string]string{
				"Authorization": getSuperuserToken(t, app),
			},
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message":"Failed to generate filter from query."`},
		},
	}

	for _, scenario := range scenarios {
		scenario.TestAppFactory = func(tb testing.TB) *tests.TestApp {
			return app
		}
		scenario.DisableTestAppCleanup = true
		scenario.Test(t)
	}
}

func TestAIQueryAPI_RespectsAPIRules(t *testing.T) {
	// Mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": `status = "active"`,
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	app := setupTestAppWithAI(t, mockServer)
	defer app.Cleanup()

	// Create collection with listRule (only superusers can access)
	collection := core.NewCollection(core.CollectionTypeBase, "private")
	collection.Fields.Add(&core.TextField{Name: "status"})
	// ListRule is nil, meaning only superusers can access
	app.Save(collection)

	// Create regular auth record (not superuser)
	authCollection := core.NewCollection(core.CollectionTypeAuth, "test_users")
	authCollection.Fields.Add(&core.EmailField{Name: "email"})
	authCollection.Fields.Add(&core.PasswordField{Name: "password"})
	authCollection.Fields.Add(&core.TextField{Name: "username"})
	// Configure password auth identity fields
	authCollection.PasswordAuth.IdentityFields = []string{"email"}
	err := app.Save(authCollection)
	require.NoError(t, err)

	record := core.NewRecord(authCollection)
	record.Set("email", "user@example.com")
	record.Set("username", "user")
	record.SetPassword("1234567890")
	err = app.Save(record)
	require.NoError(t, err)

	// Reload the record to ensure token key is set
	authRecord, err := app.FindAuthRecordByEmail("test_users", "user@example.com")
	require.NoError(t, err)
	require.NotNil(t, authRecord)

	// Get auth token for regular user
	authToken, err := authRecord.NewAuthToken()
	require.NoError(t, err)

	scenarios := []tests.ApiScenario{
		{
			Name:   "regular user cannot access collection without listRule",
			Method: http.MethodPost,
			URL:    "/api/ai/query",
			Body: strings.NewReader(`{
				"collection": "private",
				"query": "test"
			}`),
			Headers: map[string]string{
				"Authorization": authToken,
			},
			ExpectedStatus:  403,
			ExpectedContent: []string{`"message":"Only superusers can perform this action."`},
		},
	}

	for _, scenario := range scenarios {
		scenario.TestAppFactory = func(tb testing.TB) *tests.TestApp {
			return app
		}
		scenario.DisableTestAppCleanup = true
		scenario.Test(t)
	}
}

