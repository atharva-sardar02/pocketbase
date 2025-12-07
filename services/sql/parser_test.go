package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSQL_Select_Simple(t *testing.T) {
	stmt, err := ParseSQL("SELECT * FROM users")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementSelect, stmt.Type)
	assert.Contains(t, stmt.Tables, "users")
}

func TestParseSQL_Select_WithWhere(t *testing.T) {
	stmt, err := ParseSQL("SELECT * FROM orders WHERE status = 'active'")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementSelect, stmt.Type)
	assert.Contains(t, stmt.Tables, "orders")
	assert.Contains(t, stmt.Where, "status = 'active'")
}

func TestParseSQL_Select_WithJoin(t *testing.T) {
	stmt, err := ParseSQL("SELECT o.*, c.name FROM orders o JOIN customers c ON o.customer = c.id")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementSelect, stmt.Type)
	assert.Contains(t, stmt.Tables, "orders")
	assert.Contains(t, stmt.Tables, "customers")
	assert.Len(t, stmt.Joins, 1)
	assert.Equal(t, "customers", stmt.Joins[0].Table)
}

func TestParseSQL_Select_WithGroupBy(t *testing.T) {
	stmt, err := ParseSQL("SELECT status, COUNT(*) FROM orders GROUP BY status")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementSelect, stmt.Type)
	assert.Contains(t, stmt.GroupBy, "status")
}

func TestParseSQL_Select_WithOrderBy(t *testing.T) {
	stmt, err := ParseSQL("SELECT * FROM users ORDER BY created DESC")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementSelect, stmt.Type)
	assert.Len(t, stmt.OrderBy, 1)
	assert.True(t, stmt.OrderBy[0].Desc)
}

func TestParseSQL_Insert(t *testing.T) {
	stmt, err := ParseSQL("INSERT INTO users (name, email) VALUES ('John', 'john@test.com')")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementInsert, stmt.Type)
	assert.Contains(t, stmt.Tables, "users")
	assert.Equal(t, "John", stmt.Values["name"])
	assert.Equal(t, "john@test.com", stmt.Values["email"])
}

func TestParseSQL_Update(t *testing.T) {
	stmt, err := ParseSQL("UPDATE users SET name = 'Jane', status = 'active' WHERE id = '123'")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementUpdate, stmt.Type)
	assert.Contains(t, stmt.Tables, "users")
	assert.Equal(t, "Jane", stmt.SetClauses["name"])
	assert.Equal(t, "active", stmt.SetClauses["status"])
	assert.Contains(t, stmt.Where, "id = '123'")
}

func TestParseSQL_Delete(t *testing.T) {
	stmt, err := ParseSQL("DELETE FROM users WHERE status = 'inactive'")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementDelete, stmt.Type)
	assert.Contains(t, stmt.Tables, "users")
	assert.Contains(t, stmt.Where, "status = 'inactive'")
}

func TestParseSQL_CreateTable_Simple(t *testing.T) {
	sql := `CREATE TABLE products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL,
		created DATETIME
	)`
	
	stmt, err := ParseSQL(sql)
	
	assert.NoError(t, err)
	assert.Equal(t, StatementCreateTable, stmt.Type)
	assert.Contains(t, stmt.Tables, "products")
	
	// Check columns
	assert.GreaterOrEqual(t, len(stmt.Columns), 3)
	
	// Find name column
	var nameCol *ColumnDef
	for i := range stmt.Columns {
		if stmt.Columns[i].Name == "name" {
			nameCol = &stmt.Columns[i]
			break
		}
	}
	assert.NotNil(t, nameCol)
	assert.Equal(t, "TEXT", nameCol.Type)
	assert.True(t, nameCol.Required)
}

func TestParseSQL_CreateTable_WithForeignKey(t *testing.T) {
	sql := `CREATE TABLE orders (
		id TEXT PRIMARY KEY,
		customer TEXT REFERENCES customers,
		total REAL
	)`
	
	stmt, err := ParseSQL(sql)
	
	assert.NoError(t, err)
	assert.Equal(t, StatementCreateTable, stmt.Type)
	
	// Find customer column
	var customerCol *ColumnDef
	for i := range stmt.Columns {
		if stmt.Columns[i].Name == "customer" {
			customerCol = &stmt.Columns[i]
			break
		}
	}
	assert.NotNil(t, customerCol)
	assert.Equal(t, "customers", customerCol.Reference)
}

func TestParseSQL_AlterTable(t *testing.T) {
	stmt, err := ParseSQL("ALTER TABLE users ADD COLUMN age INTEGER")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementAlterTable, stmt.Type)
	assert.Contains(t, stmt.Tables, "users")
}

func TestParseSQL_DropTable(t *testing.T) {
	stmt, err := ParseSQL("DROP TABLE products")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementDropTable, stmt.Type)
	assert.Contains(t, stmt.Tables, "products")
}

func TestParseSQL_DropTable_IfExists(t *testing.T) {
	stmt, err := ParseSQL("DROP TABLE IF EXISTS products")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementDropTable, stmt.Type)
	assert.Contains(t, stmt.Tables, "products")
}

func TestParseSQL_Empty(t *testing.T) {
	_, err := ParseSQL("")
	assert.Error(t, err)
}

func TestParseSQL_Unsupported(t *testing.T) {
	_, err := ParseSQL("TRUNCATE TABLE users")
	assert.Error(t, err)
}

func TestParseSQL_WithSemicolon(t *testing.T) {
	stmt, err := ParseSQL("SELECT * FROM users;")
	
	assert.NoError(t, err)
	assert.Equal(t, StatementSelect, stmt.Type)
	assert.Contains(t, stmt.Tables, "users")
}

func TestSQLStatement_IsReadOnly(t *testing.T) {
	selectStmt := &SQLStatement{Type: StatementSelect}
	assert.True(t, selectStmt.IsReadOnly())

	insertStmt := &SQLStatement{Type: StatementInsert}
	assert.False(t, insertStmt.IsReadOnly())
}

func TestSQLStatement_IsDestructive(t *testing.T) {
	deleteStmt := &SQLStatement{Type: StatementDelete}
	assert.True(t, deleteStmt.IsDestructive())

	dropStmt := &SQLStatement{Type: StatementDropTable}
	assert.True(t, dropStmt.IsDestructive())

	selectStmt := &SQLStatement{Type: StatementSelect}
	assert.False(t, selectStmt.IsDestructive())
}

func TestSQLStatement_RequiresConfirmation(t *testing.T) {
	deleteStmt := &SQLStatement{Type: StatementDelete}
	assert.True(t, deleteStmt.RequiresConfirmation())

	updateStmt := &SQLStatement{Type: StatementUpdate}
	assert.True(t, updateStmt.RequiresConfirmation())

	selectStmt := &SQLStatement{Type: StatementSelect}
	assert.False(t, selectStmt.RequiresConfirmation())
}

func TestParseValue(t *testing.T) {
	// Test NULL
	assert.Nil(t, parseValue("NULL"))
	assert.Nil(t, parseValue("null"))

	// Test quoted strings
	assert.Equal(t, "hello", parseValue("'hello'"))
	assert.Equal(t, "world", parseValue("\"world\""))

	// Test booleans
	assert.Equal(t, true, parseValue("TRUE"))
	assert.Equal(t, true, parseValue("1"))
	assert.Equal(t, false, parseValue("FALSE"))
	assert.Equal(t, false, parseValue("0"))

	// Test unquoted value
	assert.Equal(t, "123", parseValue("123"))
}

func TestSplitValues(t *testing.T) {
	// Simple values
	values := splitValues("'a', 'b', 'c'")
	assert.Len(t, values, 3)
	assert.Equal(t, "'a'", values[0])
	assert.Equal(t, "'b'", values[1])
	assert.Equal(t, "'c'", values[2])

	// Values with commas in quotes
	values = splitValues("'a,b', 'c'")
	assert.Len(t, values, 2)
	assert.Equal(t, "'a,b'", values[0])
}

func TestSplitColumnDefs(t *testing.T) {
	// Simple columns
	defs := splitColumnDefs("id TEXT, name VARCHAR(255), age INT")
	assert.Len(t, defs, 3)

	// Columns with constraints
	defs = splitColumnDefs("id TEXT PRIMARY KEY, name TEXT NOT NULL")
	assert.Len(t, defs, 2)
}

