/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Represents a Result of a Folder Query Language query execution.
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

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

type FqlResult struct {
	Fields map[FolderField]string
	Whole  *fs.Folder
}

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func newFQLResult() *FqlResult {
	return &FqlResult{
		Fields: make(map[FolderField]string),
		Whole:  nil,
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Get the entire result object from the query, not just its fields.
// Only set if using Select(FolderAll). It returns a pointer to the
// object so that it can be modified if using (future) ALTER/UPDATE
// queries.
func (fqlR *FqlResult) Object() *fs.Folder {
	return fqlR.Whole
}

func (fqlR *FqlResult) AsString(fld FolderField) string {
	result := ""
	if v, ok := fqlR.Fields[fld]; ok {
		result = v
	}
	return result
}

func (fqlR *FqlResult) AsBoolean(fld FolderField) *bool {
	var result *bool = nil
	if fld == FolderEncrypted {
		if fqlR.Whole != nil {
			result = &fqlR.Whole.Encrypted
		} else if v, ok := fqlR.Fields[FolderEncrypted]; ok {
			switch v {
			case "True":
				fallthrough
			case "true":
				True := true
				return &True
			case "False":
				fallthrough
			case "false":
				False := false
				return &False
			}
		}
	}

	return result
}

// implements fmt.Stringer for FqlResult
func (fqlR *FqlResult) String() string {
	if fqlR.Whole != nil {
		return fqlR.Whole.StringPlus()
	} else { // by fields
		getMapField := func(fld FolderField) string {
			if v, ok := fqlR.Fields[fld]; ok {
				return v
			}
			return ""
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "\t%12s: %s\n", FolderId.String(), getMapField(FolderId))
		fmt.Fprintf(&sb, "\t%12s: %s\n", FolderName.String(), getMapField(FolderName))
		fmt.Fprintf(&sb, "\t%12s: %s\n", FolderEncrypted.String(), getMapField(FolderEncrypted))
		fmt.Fprintf(&sb, "\t%12s: %s\n", FolderTags.String(), getMapField(FolderTags))
		fmt.Fprintf(&sb, "\t%12s: %s\n", FolderTemplate.String(), getMapField(FolderTemplate))
		fmt.Fprintf(&sb, "\t%12s: %s\n", FolderPath.String(), getMapField(FolderPath))

		return sb.String()
	}
}

func (fqlR *FqlResult) GetField(fld FolderField) string {
	if v, ok := fqlR.Fields[fld]; ok {
		return v
	}
	return ""
}

// returns true if storing fields only, false if storing
// the entire object
func (fqlR *FqlResult) IsPartial() bool {
	return len(fqlR.Fields) == 0
}

// the opposite of isPartial()
func (fqlR *FqlResult) IsWhole() bool {
	return !fqlR.IsPartial()
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// set a field result
func (fqlR *FqlResult) set(field FolderField, value string) *FqlResult {
	fqlR.Fields[field] = value
	return fqlR
}

// store an entire object in the result
func (fqlR *FqlResult) store(value *fs.Folder) *FqlResult {
	fqlR.Whole = value
	return fqlR
}

// clear the result
func (fqlR *FqlResult) clear() *FqlResult {
	fqlR.Fields = make(map[FolderField]string)
	fqlR.Whole = nil
	return fqlR
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
