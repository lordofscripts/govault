/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Inheritance propagation of Tags and Template in Folder struct as
 * Folder Query Language extension.
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

// InheritFromParent returns the string value assigned to (or already present in) the target.
// * If the target had a value and the parent was empty: It returns the target's original value.
// * If the target was empty and inherited a value: It returns the inherited value.
// * If it was a Tag merge: It returns the new, unique comma-separated list.
// * If the target wasn't found: It returns "".
func InheritFromParent(field FolderField, folders []fs.Folder, identifier string) string {
	if field != FolderTags && field != FolderTemplate {
		return ""
	}

	val, found := applyInheritance(folders, identifier, field, "")
	if found {
		return val
	}
	return ""
}

func applyInheritance(folders []fs.Folder, target string, field FolderField, lastParentValue string) (string, bool) {
	for i := range folders {
		f := &folders[i]

		if f.Id == target || f.Path == target {
			// Logic for Template
			if field == FolderTemplate {
				if f.Template == "" && lastParentValue != "" {
					f.Template = lastParentValue
				}
				return f.Template, true
			}

			// Logic for Tags
			if field == FolderTags {
				if f.Tags == "" {
					if lastParentValue != "" {
						f.Tags = lastParentValue
					}
				} else if lastParentValue != "" {
					f.Tags = mergeUniqueTags(f.Tags, lastParentValue)
				}
				return f.Tags, true
			}
		}

		// Determine the value to pass down to children
		nextInheritValue := lastParentValue
		if field == FolderTemplate && f.Template != "" {
			nextInheritValue = f.Template
		} else if field == FolderTags && f.Tags != "" {
			nextInheritValue = f.Tags
		}

		// Recurse
		if val, found := applyInheritance(f.Children, target, field, nextInheritValue); found {
			return val, true
		}
	}
	return "", false
}

// mergeUniqueTags joins two comma-separated strings ensuring no duplicate tags exist.
func mergeUniqueTags(current, parent string) string {
	tagMap := make(map[string]bool)
	var result []string

	// Helper to add unique tags to the slice
	addTags := func(s string) {
		parts := strings.Split(s, ",")
		for _, t := range parts {
			t = strings.TrimSpace(t)
			if t != "" && !tagMap[t] {
				tagMap[t] = true
				result = append(result, t)
			}
		}
	}

	addTags(current)
	addTags(parent)

	return strings.Join(result, ",")
}
