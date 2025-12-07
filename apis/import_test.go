package apis_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func TestImportPreview_CSV(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			URL:             "/api/import/preview",
			Body:            strings.NewReader(`{"data":"name,email\nJohn,john@test.com","format":"csv"}`),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodPost,
			URL:    "/api/import/preview",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			Body:            strings.NewReader(`{"data":"name,email\nJohn,john@test.com","format":"csv"}`),
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser with valid CSV",
			Method: http.MethodPost,
			URL:    "/api/import/preview",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{"data":"name,email,age\nJohn,john@test.com,30\nJane,jane@test.com,25\nBob,bob@test.com,35","format":"csv"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"headers":["name","email","age"]`,
				`"sampleRows":`,
				`"totalRows":3`,
				`"format":"csv"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser with tab-delimited CSV",
			Method: http.MethodPost,
			URL:    "/api/import/preview",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{"data":"name\temail\nJohn\tjohn@test.com","format":"csv","delimiter":"\\t"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"headers":["name","email"]`,
				`"totalRows":1`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "missing data",
			Method: http.MethodPost,
			URL:    "/api/import/preview",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body:            strings.NewReader(`{"format":"csv"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message"`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestImportPreview_JSON(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:   "authorized as superuser with valid JSON array",
			Method: http.MethodPost,
			URL:    "/api/import/preview",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{"data":"[{\"name\":\"John\",\"email\":\"john@test.com\"},{\"name\":\"Jane\",\"email\":\"jane@test.com\"}]","format":"json"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"headers":`,
				`"sampleRows":`,
				`"totalRows":2`,
				`"format":"json"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "auto-detect JSON format",
			Method: http.MethodPost,
			URL:    "/api/import/preview",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{"data":"[{\"name\":\"Test\"}]"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"format":"json"`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "invalid JSON",
			Method: http.MethodPost,
			URL:    "/api/import/preview",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{"data":"{invalid json}","format":"json"}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"errors":`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestImportValidate(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			URL:             "/api/import/validate",
			Body:            strings.NewReader(`{"collection":"demo1","mapping":{"name":"title"}}`),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as superuser with valid collection",
			Method: http.MethodPost,
			URL:    "/api/import/validate",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{"collection":"demo1","mapping":{"name":"title"}}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"valid":`,
				`"fieldTypes":`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
		{
			Name:   "invalid collection",
			Method: http.MethodPost,
			URL:    "/api/import/validate",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body:            strings.NewReader(`{"collection":"nonexistent","mapping":{}}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message"`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "missing collection",
			Method: http.MethodPost,
			URL:    "/api/import/validate",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body:            strings.NewReader(`{"mapping":{}}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message"`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "invalid field mapping",
			Method: http.MethodPost,
			URL:    "/api/import/validate",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{"collection":"demo1","mapping":{"col1":"nonexistent_field"}}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"valid":false`,
				`"errors":`,
			},
			ExpectedEvents: map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

func TestImportExecute(t *testing.T) {
	t.Parallel()

	scenarios := []tests.ApiScenario{
		{
			Name:            "unauthorized",
			Method:          http.MethodPost,
			URL:             "/api/import/execute",
			Body:            strings.NewReader(`{"collection":"demo1","data":"name\nTest","format":"csv","mapping":{"name":"title"}}`),
			ExpectedStatus:  401,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "authorized as regular user",
			Method: http.MethodPost,
			URL:    "/api/import/execute",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6IjRxMXhsY2xtZmxva3UzMyIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoiX3BiX3VzZXJzX2F1dGhfIiwiZXhwIjoyNTI0NjA0NDYxLCJyZWZyZXNoYWJsZSI6dHJ1ZX0.ZT3F0Z3iM-xbGgSG3LEKiEzHrPHr8t8IuHLZGGNuxLo",
			},
			Body:            strings.NewReader(`{"collection":"demo1","data":"name\nTest","format":"csv","mapping":{"name":"title"}}`),
			ExpectedStatus:  403,
			ExpectedContent: []string{`"data":{}`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "missing collection",
			Method: http.MethodPost,
			URL:    "/api/import/execute",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body:            strings.NewReader(`{"data":"name\nTest","format":"csv","mapping":{"name":"title"}}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message"`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "missing data",
			Method: http.MethodPost,
			URL:    "/api/import/execute",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body:            strings.NewReader(`{"collection":"demo1","format":"csv","mapping":{"name":"title"}}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message"`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "missing mapping",
			Method: http.MethodPost,
			URL:    "/api/import/execute",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body:            strings.NewReader(`{"collection":"demo1","data":"name\nTest","format":"csv"}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message"`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
		{
			Name:   "successful CSV import",
			Method: http.MethodPost,
			URL:    "/api/import/execute",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{"collection":"demo1","data":"title\nImported Record 1\nImported Record 2","format":"csv","mapping":{"title":"title"}}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"totalRows":2`,
				`"successCount":`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordCreate":         2,
				"OnRecordCreateExecute":  2,
				"OnModelCreate":          2,
				"OnModelCreateExecute":   2,
				"OnRecordValidate":       2,
			},
		},
		{
			Name:   "successful JSON import",
			Method: http.MethodPost,
			URL:    "/api/import/execute",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body: strings.NewReader(`{"collection":"demo1","data":"[{\"title\":\"JSON Import 1\"},{\"title\":\"JSON Import 2\"}]","format":"json","mapping":{"title":"title"}}`),
			ExpectedStatus: 200,
			ExpectedContent: []string{
				`"totalRows":2`,
				`"successCount":`,
			},
			ExpectedEvents: map[string]int{
				"OnRecordCreate":         2,
				"OnRecordCreateExecute":  2,
				"OnModelCreate":          2,
				"OnModelCreateExecute":   2,
				"OnRecordValidate":       2,
			},
		},
		{
			Name:   "invalid collection",
			Method: http.MethodPost,
			URL:    "/api/import/execute",
			Headers: map[string]string{
				"Authorization": "eyJhbGciOiJIUzI1NiJ9.eyJpZCI6InN5d2JoZWNuaDQ2cmhtMCIsInR5cGUiOiJhdXRoIiwiY29sbGVjdGlvbklkIjoicGJjXzMxNDI2MzU4MjMiLCJleHAiOjI1MjQ2MDQ0NjEsInJlZnJlc2hhYmxlIjp0cnVlfQ.UXgO3j-0BumcugrFjbd7j0M4MQvbrLggLlcu_YNGjoY",
			},
			Body:            strings.NewReader(`{"collection":"nonexistent","data":"name\nTest","format":"csv","mapping":{"name":"title"}}`),
			ExpectedStatus:  400,
			ExpectedContent: []string{`"message"`},
			ExpectedEvents:  map[string]int{"*": 0},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
