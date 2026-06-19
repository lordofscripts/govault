/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Tests for FolderField enumeration.
 *-----------------------------------------------------------------*/
package tests

import (
	"testing"

	"github.com/lordofscripts/govault/internal/fql"
)

func TestFolderField_String(t *testing.T) {
	tCases := []struct {
		Value   fql.FolderField
		ExpName string
	}{
		{fql.FolderAll, "All"},
		{fql.FolderId, "Id"},
		{fql.FolderName, "Name"},
		{fql.FolderPath, "Path"},
		{fql.FolderTags, "Tags"},
		{fql.FolderTemplate, "Template"},
		{fql.FolderEncrypted, "Encrypted"},
	}

	for _, tCase := range tCases {
		got := tCase.Value.String()
		if got != tCase.ExpName {
			t.Errorf("expected %s got %s", tCase.ExpName, got)
		}
	}
}

func TestFolderField_Parse_NoPanic(t *testing.T) {
	tCases := []struct {
		Text    string
		IsValid bool
		Expect  fql.FolderField
	}{
		{"All", true, fql.FolderAll},
		{"FolderAll", true, fql.FolderAll},
		{"Id", true, fql.FolderId},
		{"FolderId", true, fql.FolderId},
		{"Name", true, fql.FolderName},
		{"FolderName", true, fql.FolderName},
		{"Path", true, fql.FolderPath},
		{"FolderPath", true, fql.FolderPath},
		{"Tags", true, fql.FolderTags},
		{"FolderTags", true, fql.FolderTags},
		{"Template", true, fql.FolderTemplate},
		{"FolderTemplate", true, fql.FolderTemplate},
		{"Encrypted", true, fql.FolderEncrypted},
		{"FolderEncrypted", true, fql.FolderEncrypted},
	}

	for _, tCase := range tCases {
		const OUT_OF_RANGE int = 15
		value := fql.FolderField(OUT_OF_RANGE)
		// instruct Parse not to panic. Value remains unaltered.
		value.Parse(tCase.Text, false)
		if tCase.IsValid && value != tCase.Expect {
			t.Errorf("for %s expected %s (%d)", tCase.Text, tCase.Expect, tCase.Expect)
		}
		if !tCase.IsValid && value != fql.FolderField(OUT_OF_RANGE) {
			t.Errorf("for '%s' expected Parse to fail, got %s", tCase.Text, value)
		}
	}
}

func TestFolderField_Parse_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	tCases := []struct {
		Text    string
		IsValid bool
		Expect  fql.FolderField
	}{
		{"Allx", true, fql.FolderAll},
		{"Idx", true, fql.FolderId},
		{"Namex", true, fql.FolderName},
		{"Pathx", true, fql.FolderPath},
		{"Tagsx", true, fql.FolderTags},
		{"Templatex", true, fql.FolderTemplate},
		{"Encryptedx", true, fql.FolderEncrypted},
	}

	for _, tCase := range tCases {
		var value fql.FolderField = fql.FolderAll
		// instruct Parse to panic when text is not recognized
		value.Parse(tCase.Text, true)
	}
}
