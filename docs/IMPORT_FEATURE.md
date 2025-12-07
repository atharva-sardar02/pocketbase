# Data Import Feature

## Overview

The Data Import Wizard is a feature that enables bulk data import from CSV and JSON files into PocketBase collections. It provides a step-by-step wizard interface for selecting files, previewing data, mapping columns to collection fields, and executing imports with detailed progress tracking.

## Features

### Core Features (V3)
- **Multi-Format Support** - Import from CSV and JSON files
- **4-Step Wizard** - Guided import process with visual progress indicator
- **Data Preview** - Preview parsed data before import
- **Field Mapping** - Map file columns to collection fields with auto-detection
- **Progress Tracking** - Real-time import progress with success/failure counts
- **Error Reporting** - Detailed error messages for failed rows
- **CSV Options** - Support for comma, tab, and custom delimiters

## Access

### Admin UI

1. Log in to the PocketBase Admin dashboard as a superuser
2. Click **Import** in the main sidebar (upload icon)
3. Or navigate directly to `http://127.0.0.1:8090/_/#/import`

### Requirements

- **Authentication**: Superuser access required
- **Target Collection**: Must exist before import
- **File Formats**: CSV or JSON

## Import Wizard Steps

### Step 1: Select & Upload

1. **Select Collection** - Choose the target collection from the dropdown
2. **Upload File** - Drag & drop or click to select a CSV/JSON file
3. **Preview** - Automatic preview of parsed data

### Step 2: Preview Data

- Review the parsed headers and sample rows
- Verify data looks correct before proceeding
- Go back to re-upload if needed

### Step 3: Map Fields

- Map each source column to a target collection field
- Auto-detection suggests mappings based on column names
- Select "-" to skip unmapped columns
- Field types are displayed for reference

### Step 4: Import

- Click "Start Import" to begin
- Progress bar shows completion percentage
- Success/failure counts update in real-time
- Error details shown for failed rows
- Option to "View Collection" or "Start Over" when complete

## File Format Requirements

### CSV Format

```csv
name,email,age,active
John Doe,john@example.com,30,true
Jane Smith,jane@example.com,25,false
```

**Requirements:**
- First row must contain headers
- Consistent column count across rows
- UTF-8 encoding recommended

**Supported Delimiters:**
- Comma (`,`) - default
- Tab (`\t`)
- Semicolon (`;`)
- Pipe (`|`)

### JSON Format

```json
[
  {
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "active": true
  },
  {
    "name": "Jane Smith",
    "email": "jane@example.com",
    "age": 25,
    "active": false
  }
]
```

**Requirements:**
- Must be an array of objects
- All objects should have consistent keys
- Keys become column headers

## API Reference

All import endpoints require superuser authentication.

### Preview Data

Parse file content and return headers + sample rows.

```
POST /api/import/preview
```

**Request Body:**
```json
{
  "data": "name,email,age\nJohn,john@test.com,30",
  "format": "csv",
  "delimiter": ","
}
```

**Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `data` | string | Yes | Raw file content |
| `format` | string | No | `csv` or `json` (auto-detected if omitted) |
| `delimiter` | string | No | CSV delimiter (default: `,`) |

**Response:**
```json
{
  "headers": ["name", "email", "age"],
  "sampleRows": [
    ["John", "john@test.com", "30"]
  ],
  "totalRows": 1,
  "format": "csv",
  "errors": []
}
```

---

### Validate Mapping

Validate field mapping against collection schema.

```
POST /api/import/validate
```

**Request Body:**
```json
{
  "collection": "users",
  "mapping": {
    "name": "username",
    "email": "email",
    "age": "age"
  }
}
```

**Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `collection` | string | Yes | Target collection name or ID |
| `mapping` | object | No | Source column → target field mapping |

**Response:**
```json
{
  "valid": true,
  "errors": [],
  "fieldTypes": {
    "id": "text",
    "username": "text",
    "email": "email",
    "age": "number",
    "created": "autodate",
    "updated": "autodate"
  },
  "requiredFields": ["username"]
}
```

---

### Execute Import

Perform the bulk import operation.

```
POST /api/import/execute
```

**Request Body:**
```json
{
  "collection": "users",
  "data": "name,email,age\nJohn,john@test.com,30\nJane,jane@test.com,25",
  "format": "csv",
  "delimiter": ",",
  "mapping": {
    "name": "username",
    "email": "email",
    "age": "age"
  }
}
```

**Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `collection` | string | Yes | Target collection name or ID |
| `data` | string | Yes | Raw file content |
| `format` | string | Yes | `csv` or `json` |
| `delimiter` | string | No | CSV delimiter |
| `mapping` | object | Yes | Source column → target field mapping |
| `skipHeader` | boolean | No | Skip first row for CSV (default: true) |

**Response:**
```json
{
  "totalRows": 2,
  "successCount": 2,
  "failureCount": 0,
  "errors": []
}
```

**Error Response (partial success):**
```json
{
  "totalRows": 3,
  "successCount": 2,
  "failureCount": 1,
  "errors": [
    {
      "row": 2,
      "message": "validation failed: email must be valid",
      "data": {
        "username": "Invalid User",
        "email": "not-an-email",
        "age": "25"
      }
    }
  ]
}
```

---

## Example Usage

### JavaScript (using fetch)

```javascript
const token = 'YOUR_AUTH_TOKEN';
const baseUrl = 'http://127.0.0.1:8090';

// Step 1: Preview the file
const csvData = `name,email,age
John Doe,john@example.com,30
Jane Smith,jane@example.com,25`;

const previewResponse = await fetch(`${baseUrl}/api/import/preview`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    data: csvData,
    format: 'csv'
  })
});

const preview = await previewResponse.json();
console.log('Headers:', preview.headers);
console.log('Total rows:', preview.totalRows);

// Step 2: Validate mapping
const validateResponse = await fetch(`${baseUrl}/api/import/validate`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    collection: 'users',
    mapping: {
      'name': 'username',
      'email': 'email',
      'age': 'age'
    }
  })
});

const validation = await validateResponse.json();
if (!validation.valid) {
  console.error('Validation errors:', validation.errors);
}

// Step 3: Execute import
const importResponse = await fetch(`${baseUrl}/api/import/execute`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    collection: 'users',
    data: csvData,
    format: 'csv',
    mapping: {
      'name': 'username',
      'email': 'email',
      'age': 'age'
    }
  })
});

const result = await importResponse.json();
console.log(`Imported ${result.successCount}/${result.totalRows} records`);
if (result.errors.length > 0) {
  console.log('Errors:', result.errors);
}
```

### cURL

```bash
# Preview CSV data
curl -X POST "http://127.0.0.1:8090/api/import/preview" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": "name,email,age\nJohn,john@test.com,30",
    "format": "csv"
  }'

# Validate mapping
curl -X POST "http://127.0.0.1:8090/api/import/validate" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "collection": "users",
    "mapping": {"name": "username", "email": "email"}
  }'

# Execute import
curl -X POST "http://127.0.0.1:8090/api/import/execute" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "collection": "users",
    "data": "name,email,age\nJohn,john@test.com,30",
    "format": "csv",
    "mapping": {"name": "username", "email": "email", "age": "age"}
  }'
```

### JSON Import Example

```bash
curl -X POST "http://127.0.0.1:8090/api/import/execute" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "collection": "products",
    "data": "[{\"name\":\"Widget\",\"price\":29.99},{\"name\":\"Gadget\",\"price\":49.99}]",
    "format": "json",
    "mapping": {"name": "name", "price": "price"}
  }'
```

## Architecture

### Backend Components

```
apis/import.go           # 3 API endpoints
apis/import_test.go      # 17 test cases
apis/base.go             # Route registration
```

### Frontend Components

```
ui/src/pages/ImportWizard.svelte                 # Main wizard page
ui/src/stores/import.js                          # State management
ui/src/components/import/FileUpload.svelte       # Drag-drop file zone
ui/src/components/import/DataPreview.svelte      # Preview table
ui/src/components/import/FieldMapper.svelte      # Column mapping UI
ui/src/components/import/ImportProgress.svelte   # Progress + errors
ui/src/scss/_import.scss                         # Styles
```

## Best Practices

### Data Preparation

1. **Clean your data** - Remove duplicates and fix formatting issues before import
2. **Match field types** - Ensure values match the target field types (numbers, dates, etc.)
3. **Use consistent formats** - Use ISO 8601 for dates (`YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`)
4. **Validate emails** - Ensure email fields contain valid email addresses
5. **Check required fields** - Ensure all required fields have values

### Large Imports

1. **Split large files** - Import in batches of 1,000-10,000 rows for reliability
2. **Monitor progress** - Watch the progress indicator for errors
3. **Test first** - Do a test import with a small sample before large imports
4. **Check server resources** - Ensure adequate memory for large imports

### Field Mapping Tips

1. **Auto-detection** - Use matching column names for automatic field detection
2. **Skip unnecessary columns** - Map unused columns to "-" 
3. **Check field types** - Verify the mapping makes sense for field types
4. **Handle relations** - Relation fields expect record IDs

## Troubleshooting

### "Collection not found" Error

**Problem**: The selected collection doesn't exist.

**Solutions**:
1. Verify collection name is correct
2. Refresh the collections list
3. Create the collection first if it doesn't exist

### "Failed to parse CSV" Error

**Problem**: CSV parsing failed.

**Possible Causes**:
- Incorrect delimiter
- Malformed CSV (inconsistent columns)
- Encoding issues

**Solutions**:
1. Try different delimiters (comma, tab, semicolon)
2. Verify CSV structure in a spreadsheet application
3. Save as UTF-8 encoded file

### "Validation failed" Errors

**Problem**: Individual rows fail validation.

**Possible Causes**:
- Invalid email format
- Required fields missing
- Wrong data types

**Solutions**:
1. Review error messages for specific fields
2. Fix data in source file
3. Ensure mappings match field types

### Partial Import (Some Rows Failed)

**Problem**: Some rows imported successfully, others failed.

**Solutions**:
1. Review the error list for failed rows
2. Fix issues in source data
3. Re-import only the failed rows
4. Or delete imported records and retry entire import

### Slow Import Performance

**Problem**: Large imports are slow.

**Solutions**:
1. Split into smaller batches
2. Remove unnecessary columns from mapping
3. Ensure server has adequate resources

## Security

- All endpoints require superuser authentication
- Validation prevents importing to system collections
- Field validation applies PocketBase's standard rules
- No file upload to server - data is sent as JSON payload

## Limitations

### Current Limitations

1. **No file upload** - Data sent as JSON payload (limited by request size)
2. **No relation lookup** - Relation fields require explicit record IDs
3. **No upsert mode** - Each import creates new records only
4. **No async processing** - Import is synchronous (may timeout for very large imports)
5. **No rollback** - Partial imports cannot be automatically rolled back

### Future Enhancements (Planned)

- File upload for larger imports
- Upsert mode (update existing records)
- Async import with job status tracking
- Relation field lookup by display value
- Import templates/presets

## Related Documentation

- [AI Query Feature](./AI_QUERY_FEATURE.md)
- [SQL Terminal Feature](./SQL_TERMINAL_FEATURE.md)
- [Dashboard Feature](./DASHBOARD_FEATURE.md)
- [PocketBase Collections](https://pocketbase.io/docs/collections)
