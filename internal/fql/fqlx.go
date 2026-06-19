/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Extended Folder Query Language where instead of a fluent interface
 * pattern, our query comes from a string similar to SQL. This
 * simplified version supports SELECT and UPDATE statements.
 *-----------------------------------------------------------------*/
package fql

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lordofscripts/goapp/app"
	"github.com/lordofscripts/govault/internal/fs"
	"gopkg.in/yaml.v3"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

type FQLX struct {
	useSource string                  // JSON file containing Logical Folder hierarchy
	db        []fs.Folder             // Logical Folder Hiearchy after loading
	operators map[string]OperatorFunc // Comparison operators for WHERE
	results   []FQLXResult
}

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (Ctor) creates an instance of the Extended Folder Query Language
// that processes queries in string form to produce results.
func NewFQLX() *FQLX {
	return &FQLX{
		useSource: "",
		db:        nil,
		operators: nil,
		results:   nil,
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// Tells FQLX where the Logical Folder Tree should be loaded from.
// It returns nil if not found. The actual loading is deferred
// until Connect() is called. It clears anything that may have been
// previoulsy loaded.
// Built-in operators are loaded. Only after Use() you can call AddOperator().
func (fqlx *FQLX) Use(filename string) *FQLX {
	if app.FileExists(filename) {
		fqlx.useSource = filename
		fqlx.clear()
		fqlx.loadBuiltInOperators()
		return fqlx
	}

	fmt.Fprintf(os.Stderr, "Could not find DB file '%s'\n", filename)
	return nil
}

// reads the JSON file specified previously in Use() and load it
// in memory for a query.
func (fqlx *FQLX) Connect() *FQLX {
	if data, err := fs.LoadFolderTable(fqlx.useSource); err == nil {
		fqlx.db = data
		return fqlx
	} else {
		fmt.Fprintf(os.Stderr, "Could not load logical folders '%s': %v\n", fqlx.useSource, err)
		return nil
	}
}

// Adds an operator that can be used in WHERE clauses to complement the
// built-in =, != and CONTAINS operators.
func (fqlx *FQLX) AddOperator(operator string, evaluator OperatorFunc) *FQLX {
	fqlx.operators[operator] = evaluator
	return fqlx
}

// Execute the SELECT statement on the Folder hierarchy. The results are
// cached internally until Fetch() or Update() are called.
func (fqlx *FQLX) Select(statement string) uint {
	fqlx.results = QueryFolders(fqlx.db, statement, fqlx.operators)
	return uint(len(fqlx.results))
}

// Execute an UPDATE statement on the Folder hierarchy.
// Returns the number of updated records. Fetch() returns nil
// after an Update().
func (fqlx *FQLX) Update(statement string) uint {
	var updateCount uint = 0
	fqlx.results = nil
	if fqlx.db != nil {
		updateCount = UpdateFolders(fqlx.db, statement, fqlx.operators)
	}

	return updateCount
}

// Runs the query by recognizing whether it is a SELECT or UPDATE,
// thus calling Select() or Update() accordingly.
func (fqlx *FQLX) Query(statement string) uint {
	var count uint = 0
	parts := strings.Fields(statement)
	if len(parts) > 0 {
		parts[0] = strings.ToUpper(parts[0])
		switch parts[0] {
		case "SELECT":
			count = fqlx.Select(statement)
		case "UPDATE":
			count = fqlx.Update(statement)
		default:
			fmt.Fprintf(os.Stderr, "unknown FQL statement '%s'\n", parts[0])
		}
	}
	return count
}

// Returns all the result rows after Select()
func (fqlx *FQLX) Fetch() []FQLXResult {
	return fqlx.results
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Returns the number of root nodes in the Logical Folder structure.
// Requires Connect() .
func (fqlx *FQLX) RootNodeCount() int {
	count := 0
	if fqlx.db != nil {
		count = len(fqlx.db)
	}
	return count
}

// The number of rows in a result after Select()
func (fqlx *FQLX) RowCount() int {
	if fqlx.results != nil {
		return len(fqlx.results)
	}
	return 0
}

// Clears any results and releases memory
func (fqlx *FQLX) Clear() {
	fqlx.results = nil
}

// Commit persists the folder slice to a file.
// Uses .json or .yaml/.yml extension to determine format.
func (fqlx *FQLX) Save(folders []fs.Folder, filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	var data []byte
	var err error

	switch ext {
	case ".json":
		// Marshal with 4-space indentation
		data, err = json.MarshalIndent(folders, "", "    ")

	case ".yaml", ".yml":
		// Create an encoder to set indentation
		data, err = marshalYAML(folders)

	default:
		return os.ErrInvalid // Or a custom "unsupported extension" error
	}

	if err != nil {
		return err
	}

	// Write to file with standard permissions (rw-r--r--)
	return os.WriteFile(filename, data, 0644)
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// clears anything used in a previous query except the filename
// specified in Use()
func (fqlx *FQLX) clear() {
	fqlx.db = nil
	fqlx.results = nil
}

// Defines built-in operators for use in the WHERE clause expressions.
func (fqlx *FQLX) loadBuiltInOperators() {
	// Define built-in operators
	fqlx.operators = map[string]OperatorFunc{
		// Equality operator
		"=": func(a, b any) bool { return a == b },
		// Inequality operator
		"!=": func(a, b any) bool { return a != b },
		// Contains operator for string types
		"CONTAINS": func(a, b any) bool {
			s, _ := a.(string)
			return strings.Contains(s, b.(string))
		},
	}
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Helper to handle YAML indentation specifically
func marshalYAML(folders []fs.Folder) ([]byte, error) {
	// yaml.Marshal is standard, but if you want specific indentation:
	var b strings.Builder
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2) // Standard YAML indentation
	err := enc.Encode(folders)
	if err != nil {
		return nil, err
	}
	return []byte(b.String()), nil
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

/*
func GreaterThanT[T int | uint | byte | float32 | float64](a, b T) bool {
	return a > b
}

func DemoFQLX() {
	 // A custom Greater Than operator
	greaterThan := func(a, b any) bool {
		if !IsNumericType(a) || !IsNumericType(b) {
			panic("Cannot use non-numeric type with '>' operator.")
		} else if reflect.TypeOf(a) != reflect.TypeOf(b) {
			panic("both sides of '>' should be of the same type.")
		} else {
			switch vA := a.(type) {
			case int:
				vB := b.(int)
				return GreaterThanT(vA, vB)
			case uint:
				vB := b.(uint)
				return GreaterThanT(vA, vB)
			case byte:
				vB := b.(byte)
				return GreaterThanT(vA, vB)
			case float32:
				vB := b.(float32)
				return GreaterThanT(vA, vB)
			case float64:
				vB := b.(float64)
				return GreaterThanT(vA, vB)
			default:
				panic("this shouldn't occur")
			}
		}
	}

	db := NewFQLX().
		Use("~/.config/coralys/govault/govault-docs.json").
		Connect().
		AddOperator(">", greaterThan)

	uCnt := db.Update("UPDATE folders SET Name = 'Archived' WHERE Id = '1.2'")
	sCnt := db.Select("SELECT Id,Name,Path FROM Documents WHERE Path LIKE 'Documents/%%'")
	rows := db.Fetch()
	for n, row := range rows {
		fmt.Println("Row #", n+1)
		fmt.Println(row.String())
	}
	db.Clear()
}
*/
