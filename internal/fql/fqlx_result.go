/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A Result object for a single row suitable for the FQLX object.
 *-----------------------------------------------------------------*/
package fql

import (
	"fmt"
	"strings"

	"github.com/lordofscripts/govault/internal/fs"
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

// An FQLXResult holds only the columns/fields that were requested
// in the SELECT query. If all the columns were requested (FolderAll)
// or "*" then each column value is stored here rather than the object
// itself.
type FQLXResult map[string]any

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (Ctor) creates an empty map of columns.
func NewFQLXResult() FQLXResult {
	return make(map[string]any, 0)
}

// (Ctor) creates a map with all the columns
func NewFQLXResultFrom(f fs.Folder) FQLXResult {
	r := getObjectValue(f)
	return r
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

func (m FQLXResult) Get(name string) any {
	if v, exists := m[name]; exists {
		return v
	} else {
		return nil
	}
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// implements fmt.Stringer returning a key-value rendering of all
// map items of this result instance.
func (m FQLXResult) String() string {
	var sb strings.Builder
	for k, v := range m {
		fmt.Fprintf(&sb, "\t%12s: %s\n", k, v)
	}

	return sb.String()
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// whether x is any of int|uint|byte|float32|float64 types
func IsNumericType(x any) bool {
	var result bool
	switch x.(type) {
	case int:
		result = true
	case float64:
		result = true
	case float32:
		result = true
	case uint:
		result = true
	case byte:
		result = true
	default:
		result = false
	}
	return result
}
