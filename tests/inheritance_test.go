/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Tests for Tag & Template inheritance.
 *-----------------------------------------------------------------*/
package tests

import (
	"testing"

	"github.com/lordofscripts/govault/internal/fql"
	"github.com/lordofscripts/govault/internal/fs"
)

func TestInheritFromParent(t *testing.T) {
	// 1. Setup a test hierarchy
	// Root (Tags: "work", Template: "base")
	//   -> Sub1 (Tags: "", Template: "")
	//     -> Sub2 (Tags: "urgent,work", Template: "")
	folders := []fs.Folder{
		{
			Id: "1", Name: "Root", Path: "Root", Tags: "work", Template: "base",
			Children: []fs.Folder{
				{
					Id: "1.1", Name: "Sub1", Path: "Root/Sub1", Tags: "", Template: "",
					Children: []fs.Folder{
						{
							Id: "1.1.1", Name: "Sub2", Path: "Root/Sub1/Sub2", Tags: "urgent,work", Template: "",
						},
					},
				},
			},
		},
	}

	t.Run("Template Inheritance", func(t *testing.T) {
		// Sub1 should inherit "base" from Root
		res := fql.InheritFromParent(fql.FolderTemplate, folders, "1.1")
		if res != "base" {
			t.Errorf("Expected template 'base', got '%s'", res)
		}
		if folders[0].Children[0].Template != "base" {
			t.Error("Folder struct was not updated in-place")
		}

		// Sub2 should inherit "base" from Root (nearest non-empty ancestor)
		resSub2 := fql.InheritFromParent(fql.FolderTemplate, folders, "Root/Sub1/Sub2")
		if resSub2 != "base" {
			t.Errorf("Expected Sub2 to inherit 'base' from grandparent, got '%s'", resSub2)
		}
	})

	t.Run("Tag Merging Unique", func(t *testing.T) {
		// Sub2 has "urgent,work". Root has "work".
		// Merge should result in "urgent,work" (no duplicates)
		res := fql.InheritFromParent(fql.FolderTags, folders, "1.1.1")

		// Note: order depends on implementation, but unique is required
		expected := "urgent,work"
		if res != expected {
			t.Errorf("Expected tags '%s', got '%s'", expected, res)
		}
	})

	t.Run("Tag Simple Inheritance", func(t *testing.T) {
		// Sub1 is empty, should just take Root's "work"
		res := fql.InheritFromParent(fql.FolderTags, folders, "1.1")
		if res != "work" {
			t.Errorf("Expected inherited tags 'work', got '%s'", res)
		}
	})

	t.Run("Non-existent target", func(t *testing.T) {
		res := fql.InheritFromParent(fql.FolderTags, folders, "9.9.9")
		if res != "" {
			t.Error("Expected empty string for non-existent target")
		}
	})
}
