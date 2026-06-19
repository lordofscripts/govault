/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * QUERY ENGINE :: UPDATE
 * Parses the SQL-like string and executes the search of the form
 * `UPDATE FROM folders SET field = value [, field = value] WHERE condition`
 *-----------------------------------------------------------------*/
package fql

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// UpdateFolders finds matching folders and applies the requested field changes.
// · UPDATE folders SET Encrypted = true WHERE Name = "Legal"
// · UPDATE folders SET Encrypted = true, Tags = "secret" WHERE Name = "Computers"
// Returns the number of updated records (Folders)
func UpdateFolders(
	folders []fs.Folder,
	query string,
	operators map[string]OperatorFunc,
) uint {
	// 1. Regex to capture the SET part and the WHERE part
	// Syntax: UPDATE folders SET {assignments} WHERE {condition}
	re := regexp.MustCompile(`(?i)UPDATE\s+.*?\s+SET\s+(.*?)\s+WHERE\s+(.*)`)
	matches := re.FindStringSubmatch(query)
	if len(matches) < 3 {
		return 0
	}

	setPart := strings.TrimSpace(matches[1])
	wherePart := strings.TrimSpace(matches[2])

	// 2. Parse the SET assignments
	// Handles: Name = 'Test', Encrypted = true
	assignments := make(map[FolderField]any)

	// Split by comma, but be mindful that values might contain commas if quoted
	// For simplicity in this implementation, we split by comma and trim
	rawAssignments := strings.Split(setPart, ",")
	for _, asgn := range rawAssignments {
		parts := strings.SplitN(strings.TrimSpace(asgn), "=", 2)
		if len(parts) != 2 {
			continue
		}

		fieldName, _ := extractQuotedString(parts[0])
		rawValue, _ := extractQuotedString(parts[1])

		if enum, ok := fieldMap[fieldName]; ok {
			assignments[enum] = rawValue
		}
	}

	// 3. Recursive Walk and Update
	var walk func([]fs.Folder)
	var updateCount uint = 0
	walk = func(current []fs.Folder) {
		// @note we iterate using for i := range current and pass &current[i]
		// to the helper, the changes are persisted directly into the
		// original slice provided by the caller.
		for i := range current {
			if evaluateWhere(current[i], wherePart, operators) {
				// Apply all assignments to this instance
				for field, val := range assignments {
					applyUpdate(&current[i], field, val.(string))
				}
				updateCount++
			}

			if len(current[i].Children) > 0 {
				walk(current[i].Children)
			}
		}
	}

	walk(folders)
	return updateCount
}

// applyUpdate handles the type conversion and assignment for specific fields
func applyUpdate(f *fs.Folder, field FolderField, value string) {
	switch field {
	case FolderName:
		f.Name = value
	case FolderPath:
		f.Path = value
	case FolderTags:
		f.Tags = value
	case FolderTemplate:
		f.Template = value
	case FolderEncrypted:
		if b, err := strconv.ParseBool(value); err == nil {
			f.Encrypted = b
		}
	default:
		// Note: Id is typically not updated via query as it is a generated identifier
		fmt.Fprintf(os.Stderr, "field '%s' not supported on UPDATE", field)
	}
}
