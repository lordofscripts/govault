/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * QUERY ENGINE :: SELECT
 * Parses the SQL-like string and executes the search.
 *-----------------------------------------------------------------*/
package fql

import (
	"regexp"
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

type OperatorFunc func(any, any) bool

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

// Dynamic Projections: The returned []map[string]any ensures you only get
// the keys you asked for in the SELECT part.
func QueryFolders(
	folders []fs.Folder,
	query string,
	operators map[string]OperatorFunc,
) []FQLXResult {
	// Updated Regex: Handles SELECT and FROM keywords, captures remainder for WHERE
	re := regexp.MustCompile(`(?i)SELECT\s+(.*?)\s+FROM\s+(.*)`)
	matches := re.FindStringSubmatch(query)
	if len(matches) < 3 {
		return nil
	}

	selectPart := strings.TrimSpace(matches[1])
	fromAndWhere := strings.TrimSpace(matches[2])

	// 1. Extract Source (handling quotes)
	sourcePart, remainder := extractQuotedString(fromAndWhere)

	// Extract the WHERE clause from the remainder (remove "WHERE " prefix)
	wherePart := ""
	if idx := strings.Index(strings.ToUpper(remainder), "WHERE"); idx != -1 {
		wherePart = strings.TrimSpace(remainder[idx+5:])
	}

	// 2. Determine search scope
	var searchScope []fs.Folder
	if sourcePart == "*" || strings.ToLower(sourcePart) == "any" {
		searchScope = folders
	} else {
		for i := range folders {
			if folders[i].Name == sourcePart {
				searchScope = []fs.Folder{folders[i]}
				break
			}
		}
	}

	if len(searchScope) == 0 {
		return nil
	}

	// 3. Process results
	var results []FQLXResult
	var requestedFields []string
	if selectPart != "*" {
		for _, f := range strings.Split(selectPart, ",") {
			requestedFields = append(requestedFields, strings.TrimSpace(f))
		}
	}

	var walk func([]fs.Folder)
	walk = func(current []fs.Folder) {
		for i := range current {
			if evaluateWhere(current[i], wherePart, operators) {
				//row := make(map[string]any)
				var row FQLXResult
				if selectPart == "*" {
					//row["Object"] = current[i]
					//row = getObjectValue(current[i])
					row = NewFQLXResultFrom(current[i])
				} else {
					row = NewFQLXResult()
					for _, fName := range requestedFields {
						if enum, ok := fieldMap[fName]; ok {
							row[fName] = getFieldValue(current[i], enum)
						}
					}
				}
				results = append(results, row)
			}
			walk(current[i].Children)
		}
	}

	walk(searchScope)
	return results
}

// Helper to handle quoted or unquoted string extraction in FROM & WHERE clauses.
// Handles the following forms of FROM and WHERE clauses:
// 1. CLAUSE 'this or that'
// 2. CLAUSE "this or that"
// 3. CLAUSE Documents
func extractQuotedString(input string) (value string, remainder string) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return "", ""
	}

	// Check for quotes
	if input[0] == '"' || input[0] == '\'' {
		quote := input[0]
		end := strings.IndexByte(input[1:], quote)
		if end != -1 {
			// Found the closing quote
			return input[1 : end+1], strings.TrimSpace(input[end+2:])
		}
	}

	// No quotes: take the first word
	parts := strings.Fields(input)
	if len(parts) > 0 {
		return parts[0], strings.TrimSpace(input[len(parts[0]):])
	}
	return "", ""
}
