/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Extended Folder Query Language parsed from a string. This
 * handles the WHERE clause of the SQL-like statement.
 *-----------------------------------------------------------------*/
package fql

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/lordofscripts/govault/internal/fs"
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

/*
 *		CLAUSE EVALUATOR
 *		-------------------------------------------------
 *		The query looks like: "SELECT field1,field2,field3 FROM string WHERE clause" where the comma-separated list of fields names are strings mapped to an enumeration. the clause in the WHERE part can take several forms:
 *
 *		1. WHERE field_name Operator VALUE
 *		2. WHERE field_name LIKE VALUE
 * 		3. WHERE * LIKE VALUE
 *
 *		In form #1 The OPERATOR is a func(any,any) that gets passed the Folder member field
 *		corresponding to field_name as the 1st parameter, and VALUE as the 2nd parameter.
 *
 *		In form #2 with LIKE we use the same mechanism as the SearchLike function except that
 *		we apply LIKE on the selected field_name value, and
 *
 *		In form #3 we use the SearchLike function as earlier in this project where the "*"
 *		as field_name means we would use LIKE on the entire Folder object but applying the
 *		Like on either Tags, Path, Name or Id.
 */

func evaluateWhere(f fs.Folder, clause string, operators map[string]OperatorFunc) bool {
	if clause == "" {
		return true
	}

	// Form 3: * LIKE 'Value'
	if strings.HasPrefix(strings.ToUpper(clause), "* LIKE") {
		_, val := extractQuotedString(clause[6:])
		val, _ = extractQuotedString(clause[strings.Index(strings.ToUpper(clause), "LIKE")+4:])
		return matchLike(f.Name, val) || matchLike(f.Id, val) || matchLike(f.Path, val) || matchLike(f.Tags, val)
	}

	// General Form: FieldName Operator Value
	// 1. Get Field
	fieldName, rem := extractQuotedString(clause)
	fieldEnum, ok := fieldMap[fieldName]
	if !ok {
		return false
	}

	// 2. Get Operator
	opName, rem := extractQuotedString(rem)
	opName = strings.ToUpper(opName)

	// 3. Get Value
	rawValue, _ := extractQuotedString(rem)

	actualValue := getFieldValue(f, fieldEnum)

	// Handle Boolean Conversion (Fix for Requirement #1)
	var processedValue any = rawValue
	if b, ok := actualValue.(bool); ok {
		parsedBool, err := strconv.ParseBool(rawValue)
		if err == nil {
			processedValue = parsedBool
		} else {
			// If it's not a valid bool string, the comparison should probably be false
			processedValue = !b // Ensure mismatch
		}
	}

	if opName == "LIKE" {
		strVal, _ := actualValue.(string)
		return matchLike(strVal, rawValue)
	}

	if op, ok := operators[opName]; ok {
		return op(actualValue, processedValue)
	}

	return false
}

// Standard LIKE logic (reusing your SearchLike concept)
func matchLike(text, pattern string) bool {
	p := regexp.QuoteMeta(pattern)
	p = strings.ReplaceAll(p, "%", ".*")
	p = strings.ReplaceAll(p, "_", ".")
	matched, _ := regexp.MatchString("(?i)^"+p+"$", text)
	return matched
}

func extractValue(clause, op string) string {
	// Simplistic extraction: get everything after the operator and strip quotes
	idx := strings.Index(strings.ToUpper(clause), strings.ToUpper(op))
	val := strings.TrimSpace(clause[idx+len(op):])
	return strings.Trim(val, "'\"")
}
