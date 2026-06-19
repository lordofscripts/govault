/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Tests for Folder struct and its methods.
 *-----------------------------------------------------------------*/
package tests

import (
	"testing"

	vfs "github.com/lordofscripts/govault/internal/fs"
)

// After LoadFolderTable both Path and Id properties are calculated.
// We test by doing a deep search by Logical Path
func TestFolder_SearchFolders(t *testing.T) {
	folders := loadLogicalStore(t, SAMPLE_DOCS_TABLE)

	const PATH = "Documents/Legal/Citizenship"
	target := vfs.SearchFolders(folders, PATH)
	if target == nil {
		t.Errorf("expected a result in SearchFolders")
	}
	if target.Path != PATH {
		t.Errorf("expected '%s' but got '%s'", PATH, target.Path)
	}
}

// After LoadFolderTable both Path and Id properties are calculated.
// We test by doing a deep search by Pysical Path which internally
// requires SearchFolders to work properly. In return we get the
// object with corresponding Logical Path.
func TestFolder_FindByPhysicalPath(t *testing.T) {
	folders := loadLogicalStore(t, SAMPLE_DOCS_TABLE)

	const PHYSICAL_PATH = "/home/lordofscripts/Documents/Legal/Citizenship"
	const LOGICAL_PATH = "Documents/Legal/Citizenship"
	physicalMapping := map[string]string{
		"Documents": "/home/lordofscripts/Documents",
		"Pictures":  "/home/lordofscripts/Pictures",
	}
	target := vfs.FindByPhysicalPath(folders, physicalMapping, PHYSICAL_PATH)
	if target == nil {
		t.Errorf("expected a result in FindByPhysicalPath")
	} else if target.Path != LOGICAL_PATH {
		t.Errorf("expected '%s' but got '%s'", PHYSICAL_PATH, target.Path)
	}
}

// Using LIKE function to search a Folder slice on their Id, Path or Name.
func TestFolder_SearchLike(t *testing.T) {
	t.Parallel()

	folders := loadLogicalStore(t, SAMPLE_DOCS_TABLE)

	tcs := []struct {
		By            string
		Description   string
		Search        string
		ExpectedCount int
	}{
		{"Id", "Valid ID (specific)", "1.10.3", 1},
		{"Id", "Valid ID (char wildcard)", "1.10._", 3},
		{"Id", "Invalid ID (specific)", "15", 0},
		{"Name", "Valid Name (specific)", "People", 1},
		{"Name", "Valid Name (char wildcard)", "Pe_ple", 1},
		{"Name", "Valid Name (specific)", "Housing & Property", 1},
		{"Name", "Valid Name (words wildcard)", "Housing%", 1},
		{"Name", "Invalid Name (specific)", "Carajo", 0},
		{"Path", "Valid Path (specific)", "Documents/Legal/Citizenship", 1},
		{"Path", "Valid Path (char wildcard)", "Docu_ents/Le_al/Citizenship", 1},
		{"Path", "Valid Path (words wildcard)", "%/Legal/%", 3},
		{"Path", "Invalid Path (specific)", "Pictures", 0},
	}

	for n, tc := range tcs {
		t.Run(tc.Description, func(t *testing.T) {
			results := vfs.SearchLike(folders, tc.Search)
			if len(results) != tc.ExpectedCount {
				t.Errorf("#%d expected %d results, got %d", n+1, tc.ExpectedCount, len(results))
			}
		})
	}
}

// Helps loading any of the Docs, Pics or custom JSON files containing Folders
func loadLogicalStore(t testing.TB, storeJSON string) []vfs.Folder {
	t.Helper()

	folders, err := vfs.LoadFolderTable(storeJSON)
	if err != nil {
		t.Errorf("unable to open Folder table '%s': %v", storeJSON, err)
	}

	return folders
}
