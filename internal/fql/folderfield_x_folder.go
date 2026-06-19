/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Linkage between `FolderField` enumeration and `fs.Folder` member.
 *-----------------------------------------------------------------*/
package fql

import (
	"fmt"

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

/* ----------------------------------------------------------------
 *               P R I V A T E    F U N C T I O N S
 *-----------------------------------------------------------------*/

// Helper to get field value from Folder instance
func getFieldValue(f fs.Folder, field FolderField) any {
	switch field {
	case FolderId:
		return f.Id
	case FolderName:
		return f.Name
	case FolderPath:
		return f.Path
	case FolderTags:
		return f.Tags
	case FolderTemplate:
		return f.Template
	case FolderEncrypted:
		return f.Encrypted
	case FolderAll: // @note this was not in the original version
		return f
	default:
		return nil
	}
}

// given the field handle, it retrieves its value
func getFieldFrom(fld FolderField, obj fs.Folder) any { // @audit deprecate, use getFieldValue()
	var targetValue any
	switch fld {
	case FolderId:
		targetValue = obj.Id
	case FolderName:
		targetValue = obj.Name
	case FolderTags:
		targetValue = obj.Tags
	case FolderTemplate:
		targetValue = obj.Template
	case FolderPath:
		targetValue = obj.Path
	case FolderEncrypted:
		targetValue = fmt.Sprintf("%t", obj.Encrypted)
	case FolderAll:
		targetValue = obj
	}
	return targetValue
}

// convert a whole Folder object to a map indexed by exported field names.
func getObjectValue(f fs.Folder) map[string]any {
	columns := make(map[string]any)
	columns["Id"] = f.Id
	columns["Name"] = f.Name
	columns["Path"] = f.Path
	columns["Tags"] = f.Tags
	columns["Template"] = f.Template
	columns["Encrypted"] = f.Encrypted
	return columns
}
