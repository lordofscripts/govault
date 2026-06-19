/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package fql

import (
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

type Comparer func(string, string) bool
type FolderComparer func(fs.Folder) bool
type LikeAnyComparer func(folders []fs.Folder, pattern string) []*fs.Folder

type FolderFieldPicker func(FolderField, fs.Folder) any
type FieldGetter func(FolderField, fs.Folder) any

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/
/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Exactly equal including letter-case
func IsEqual(leftStr, rightStr string) bool {
	return leftStr == rightStr
}

// Exactly equal ignoring letter-case
func Is(leftStr, rightStr string) bool {
	return strings.EqualFold(leftStr, rightStr)
}

// Whether leftStr contains rightStr, case-sensitive
func Contains(leftStr, rightStr string) bool {
	return strings.Contains(leftStr, rightStr)
}

// Whether leftStr contains rightStr, case-sensitive
func ContainsInsensitive(leftStr, rightStr string) bool {
	return strings.Contains(strings.ToLower(leftStr), strings.ToLower(rightStr))
}

// Determine whether instance of type T.
// Example: Is[Folder](someVar)
func IsType[T any](instance any) bool {
	_, ok := instance.(T)
	return ok
}
