/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Tests for Folder Query with Fluent Interface.
 *-----------------------------------------------------------------*/
package tests

import (
	"fmt"
	"testing"

	"github.com/lordofscripts/govault/internal/fs"

	"github.com/lordofscripts/govault/internal/fql"
)

const (
	SAMPLE_DOCS_TABLE = "../cmd/json/tree_sample_docs.json"
	SAMPLE_PICS_TABLE = "../cmd/json/tree_sample_pics.json"
)

/*
 * Test Comparer operators for use in Where statement
 */

// FQL.Is case-sensitive string comparison
func TestFQL_IsEqual(t *testing.T) {
	testCases := []struct {
		Target  string
		Value   string
		IsMatch bool
	}{
		{"John", "John", true},
		{"John", "john", false},
		{"John", "Jane", false},
	}

	for i, tcase := range testCases {
		got := fql.IsEqual(tcase.Target, tcase.Value)
		if got != tcase.IsMatch {
			t.Errorf("#%d expected:%t got:%t", i+1, tcase.IsMatch, got)
		}
	}
}

// FQL.IsInsensitive case-insensitive string comparison
func TestFQL_Is(t *testing.T) {
	testCases := []struct {
		Target  string
		Value   string
		IsMatch bool
	}{
		{"John", "John", true},
		{"John", "john", true},
		{"John", "Jane", false},
	}

	for i, tcase := range testCases {
		got := fql.Is(tcase.Target, tcase.Value)
		if got != tcase.IsMatch {
			t.Errorf("#%d expected:%t got:%t", i+1, tcase.IsMatch, got)
		}
	}
}

// FQL.Contains case-sensitive string containment
func TestFQL_Contains(t *testing.T) {
	testCases := []struct {
		Target  string
		Value   string
		IsMatch bool
	}{
		{"John", "John", true},
		{"John is a Programmer", "john", false},
		{"John is a Programmer", "Programmer", true},
		{"John is a Designer", "is a Designer", true},
	}

	for i, tcase := range testCases {
		got := fql.Contains(tcase.Target, tcase.Value)
		if got != tcase.IsMatch {
			t.Errorf("#%d expected:%t got:%t", i+1, tcase.IsMatch, got)
		}
	}
}

// FQL.Contains case-insensitive string containment
func TestFQL_ContainsInsensitive(t *testing.T) {
	testCases := []struct {
		Target  string
		Value   string
		IsMatch bool
	}{
		{"John", "John", true},
		{"John is a Programmer", "john", true},
		{"John is a Programmer", "programmer", true},
		{"John is a Designer", "is a designer", true},
		{"John is a Designer", "is a Programmer", false},
	}

	for i, tcase := range testCases {
		got := fql.ContainsInsensitive(tcase.Target, tcase.Value)
		if got != tcase.IsMatch {
			t.Errorf("#%d expected:%t got:%t", i+1, tcase.IsMatch, got)
		}
	}
}

/*
 * Test Select() & Selected()
 */

func TestFQL_Select(t *testing.T) {
	stmt := fql.NewFolderQuery()

	if stmt.FieldCount() != 0 {
		t.Error("Expected field count = 0")
	}

	stmt.Select(fql.FolderAll)
	if stmt.FieldCount() != 1 {
		t.Error("Expected field count = 1 for *")
	}
}

func TestFQL_Selected(t *testing.T) {
	stmt := fql.NewFolderQuery()
	stmt.Select(fql.FolderTags, fql.FolderPath, fql.FolderName) // any order
	if stmt.FieldCount() != 3 {
		t.Error("Expected field count = 3 for *")
	}
	expStr := "Name,Path,Tags"
	gotStr := stmt.Selected()
	if gotStr != expStr { // always returned in a sorted fashion
		t.Errorf("didn't get expected '%s' selected field names: '%s'", expStr, gotStr)
	}
}

/*
 * Test Where()
 */

// Test WhereX() which uses a user-provided comparator. Use the sample
// Docs JSON file and produce a deep match. @note when source file gets
// elements added/removed the test case(s) should be re-evaluated.
func TestFQL_WhereX(t *testing.T) {
	table := loadLogicalStore(t, SAMPLE_DOCS_TABLE)

	// a custom evaluator
	tableProcedure := func(f fs.Folder) bool {
		return fql.IsEqual(f.Name, "Citizenship") && fql.ContainsInsensitive(f.Tags, "foreign")
	}

	stmt := fql.NewFolderQuery().
		Select(fql.FolderId, fql.FolderName, fql.FolderPath, fql.FolderTags).
		From(table).
		WhereX(tableProcedure)
	totalResultRows := stmt.Count()
	if totalResultRows != 1 {
		t.Errorf("expected %d result in sample Docs, got %d", 1, totalResultRows)
	} else {
		record := stmt.Fetch()[0]
		fmt.Println(record.String())
	}
}

// Test WhereX() which uses a user-provided comparator. Use the sample
// Docs JSON file and produce a deep match. @note when source file gets
// elements added/removed the test case(s) should be re-evaluated.
func TestFQL_Id(t *testing.T) {
	table := loadLogicalStore(t, SAMPLE_DOCS_TABLE)

	// a custom evaluator
	tableProcedure := func(f fs.Folder) bool {
		return fql.IsEqual(f.Name, "Citizenship") && fql.ContainsInsensitive(f.Tags, "foreign")
	}

	stmt := fql.NewFolderQuery().
		Select(fql.FolderId, fql.FolderName, fql.FolderPath, fql.FolderTags).
		From(table).
		WhereX(tableProcedure)
	const EXPECT = "1.10.3"
	got := stmt.Fetch()[0].GetField(fql.FolderId)
	if got != EXPECT {
		t.Errorf("expected %s result in sample Docs, got %s", EXPECT, got)
	} else {
		record := stmt.Fetch()[0]
		fmt.Println(record.String())
	}
}

// Test WhereX() which uses a user-provided comparator. Use the sample
// Docs JSON file and produce a deep match. @note when source file gets
// elements added/removed the test case(s) should be re-evaluated.
func TestFQL_Path(t *testing.T) {
	table := loadLogicalStore(t, SAMPLE_DOCS_TABLE)

	// a custom evaluator
	tableProcedure := func(f fs.Folder) bool {
		return fql.IsEqual(f.Name, "Citizenship") && fql.ContainsInsensitive(f.Tags, "foreign")
	}

	stmt := fql.NewFolderQuery().
		Select(fql.FolderId, fql.FolderName, fql.FolderPath, fql.FolderTags).
		From(table).
		WhereX(tableProcedure)
	const EXPECT = "Documents/Legal/Citizenship"
	got := stmt.Fetch()[0].GetField(fql.FolderPath)
	if got != EXPECT {
		t.Errorf("expected %s result in sample Docs, got %s", EXPECT, got)
	} else {
		record := stmt.Fetch()[0]
		fmt.Println(record.String())
	}
}

func TestFQL_WhereFetch(t *testing.T) {
	table := loadLogicalStore(t, SAMPLE_DOCS_TABLE)

	testCases := []struct {
		Field fql.FolderField
		Oper  fql.Comparer
		Value string
		Qty   int
	}{
		{fql.FolderName, fql.IsEqual, "Documents", 1}, // Level 0
		{fql.FolderName, fql.Is, "documents", 1},      // Level 0 (case-insensitive)
		{fql.FolderName, fql.IsEqual, "documents", 0}, // Level 0 (case-sensitive)
		{fql.FolderName, fql.IsEqual, "Admin", 1},     // Level 1
		{fql.FolderName, fql.IsEqual, "Lisbon", 1},
		{fql.FolderName, fql.IsEqual, "Certificates", 1},
		// now try Tags which contains CSV therefore use Contain* instead of Is*
		{fql.FolderTags, fql.ContainsInsensitive, "projects", 1},
	}

	for n, tCase := range testCases {
		stmt := fql.NewFolderQuery().
			Select(tCase.Field, fql.FolderPath, fql.FolderTags).
			From(table).
			Where(tCase.Field, tCase.Oper, tCase.Value)
		totalResultRows := stmt.Count()
		if totalResultRows != tCase.Qty {
			t.Errorf("#%02d expected %d result in sample Docs, got %d", n+1, tCase.Qty, totalResultRows)
		} else if tCase.Qty > 0 {
			record := stmt.Fetch()[0]
			fmt.Printf("%d Path:%s Tags:%s\n", n+1, record.AsString(fql.FolderPath), record.AsString(fql.FolderTags))
		}
	}
}

/*
 * Test RowCount() the total number of records in a virtual table.
 * @note If you add elements or children to tree_sample_docs/pics.json,
 *		 you must update the SAMPLE_COUNT constants in this test!
 */
func TestFQL_RowCount(t *testing.T) {
	const DOCS_SAMPLE_COUNT int = 28
	const PICS_SAMPLE_COUNT int = 26
	table := loadLogicalStore(t, SAMPLE_DOCS_TABLE)

	stmt := fql.NewFolderQuery().
		Select(fql.FolderName, fql.FolderPath).
		From(table) // we don't need Where because we are counting table size
	totalRows := stmt.RowCount()
	if totalRows != DOCS_SAMPLE_COUNT {
		t.Errorf("expected %d elements in sample Docs, got %d", DOCS_SAMPLE_COUNT, totalRows)
	}

	table = loadLogicalStore(t, SAMPLE_PICS_TABLE)

	stmt = fql.NewFolderQuery().
		Select(fql.FolderName, fql.FolderPath).
		From(table).
		Where(fql.FolderName, fql.IsEqual, "Pictures") // irrelevant because we count rows not result
	totalRows = stmt.RowCount()
	if totalRows != PICS_SAMPLE_COUNT {
		t.Errorf("expected %d elements in sample Pics, got %d", PICS_SAMPLE_COUNT, totalRows)
	}
}
