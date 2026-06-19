/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Folder Field enumeration used to access Folder object using
 * the Folder Query Language.
 *-----------------------------------------------------------------*/
package fql

import (
	"fmt"
	"os"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	ALL_COLUMNS string = "*" // all fields of the Folder object

	FolderAll FolderField = iota
	FolderId
	FolderName
	FolderPath
	FolderTags
	FolderTemplate
	FolderEncrypted
)

/* ----------------------------------------------------------------
 *                       L O C A L S
 *-----------------------------------------------------------------*/

// translates a Camel-cased FolderField text to its enumeration
var fieldMap = map[string]FolderField{
	"All":             FolderAll,
	"FolderAll":       FolderAll,
	"Id":              FolderId,
	"FolderId":        FolderId,
	"Name":            FolderName,
	"FolderName":      FolderName,
	"Path":            FolderPath,
	"FolderPath":      FolderPath,
	"Tags":            FolderTags,
	"FolderTags":      FolderTags,
	"Template":        FolderTemplate,
	"FolderTemplate":  FolderTemplate,
	"Encrypted":       FolderEncrypted,
	"FolderEncrypted": FolderEncrypted,
}

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ fmt.Stringer = FolderField(0)

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

// Identifies a struct field in a `fs.Folder` object.
type FolderField byte

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// implements fmt.Stringer.
// Converts FolderField enumeration value to its string form
// without the 'Folder' prefix so that FQL queries are more compact.
func (ff FolderField) String() string {
	name := ""
	switch ff {
	case FolderId:
		name = "Id"
	case FolderName:
		name = "Name"
	case FolderPath:
		name = "Path"
	case FolderEncrypted:
		name = "Encrypted"
	case FolderTags:
		name = "Tags"
	case FolderTemplate:
		name = "Template"
	case FolderAll:
		name = "All"
	}
	return name
}

// Attempts to parse the str argument into a FolderField enumeration
// value. The string can be as produced by the String() method or
// with the Folder prefix, i.e. FolderName or Name. It is case-sensitive.
// The behavior of cases where the string cannot be identified is
// controlled by panicOnError. If set to false the value upon which
// Parse was invoked is not altered.
func (ff *FolderField) Parse(str string, panicOnError bool) {
	if enum, ok := fieldMap[str]; ok {
		*ff = enum
	} else {
		fmt.Fprintf(os.Stderr, "'%s' is not a FolderField item\n", str)
		if panicOnError {
			panic("bad thing happened")
		}
	}
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
