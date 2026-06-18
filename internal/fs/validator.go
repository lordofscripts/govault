/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A validator template processor for validating filenames against
 * allowed patterns/templates. It has the flexibility of allowing
 * predefined tokens (like country codes, image prefixes), lists
 * and regular-expression rules.
 *-----------------------------------------------------------------*/
package fs

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
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

// Array or slice of Key-Value pairs
type KeyValuePairArray = []KeyValuePair
type StringMap = map[string]string

// A validator object
type Validator struct {
	Lists      map[string]StringMap
	Regexes    StringMap
	Predefined map[string]KeyValuePairArray
}

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func NewValidator() *Validator {
	return &Validator{
		Lists:      make(map[string]StringMap),
		Regexes:    make(map[string]string),
		Predefined: make(map[string]KeyValuePairArray),
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// replace the validation lists. Omit the leading '$' in the List key.
// Each key must begin with 'L' followed by 1-2 digits.
func (v *Validator) WithLists(l map[string]StringMap) *Validator {
	v.Lists = l
	return v
}

// replace the validation Regular Expression rules. Omit leading '$'
// in the Rule key. Each key must begin with 'R' followed by 1-2 digits.
func (v *Validator) WithRules(l map[string]string) *Validator {
	v.Regexes = l
	return v
}

// replace the validation predefined tokens. Omit the leading '%' in the
// one-letter token.
func (v *Validator) WithTokens(l map[string]KeyValuePairArray) *Validator {
	v.Predefined = l
	return v
}

// Add a new list item. The corresponding list item token that can be
// referred in the templates is in the form $Ln or $Lnn where n is 0-9
func (v *Validator) AddListItem(id byte, sm StringMap) *Validator {
	key := fmt.Sprintf("L%d", id)
	v.Lists[key] = sm
	return v
}

// Add a new list item. The corresponding list item token that can be
// referred in the templates is in the form $Ln or $Lnn where n is 0-9
func (v *Validator) AddRuleItem(id byte, sm string) *Validator {
	key := fmt.Sprintf("R%d", id)
	v.Regexes[key] = sm
	return v
}

// Add a new list item. The corresponding list item token that can be
// referred in the templates is in the form $Ln or $Lnn where n is 0-9
func (v *Validator) AddPredefinedItem(id rune, kv KeyValuePairArray) *Validator {
	v.Predefined[string(id)] = kv
	return v
}

// validate a filename against the provided template.
func (v *Validator) ValidateFilename(template, filename string) (bool, error) {
	// 1. Escape the template. This turns '*' into '\*' and '$' into '\$'
	rePattern := regexp.QuoteMeta(template)

	// 2. Handle Wildcard (replace escaped '\*' with '.*')
	rePattern = strings.ReplaceAll(rePattern, `\*`, ".*")

	// 3. Replace Predefined tokens (% is not escaped by QuoteMeta)
	rePredefined := regexp.MustCompile(`%[A-Z]`)
	rePattern = rePredefined.ReplaceAllStringFunc(rePattern, func(m string) string {
		token := m[1:]
		switch token {
		case "D": // YYYYMMDD
			return `(\d{4}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01]))`
		case "T": // HHMMSS
			return `(([01]\d|2[0-3])[0-5]\d[0-5]\d)`
		default:
			if kvArray, ok := v.Predefined[token]; ok {
				keys := make([]string, 0, len(kvArray))
				for _, kv := range kvArray {
					keys = append(keys, regexp.QuoteMeta(kv.Key))
				}
				sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })
				return "(" + strings.Join(keys, "|") + ")"
			}
		}
		return m
	})

	// 4. Replace Regex tokens (Match the escaped '\$R')
	reRegexTokens := regexp.MustCompile(`\\\$R\d{1,2}`)
	rePattern = reRegexTokens.ReplaceAllStringFunc(rePattern, func(m string) string {
		// m is "\$R1", we need "R1"
		id := m[2:]
		if regStr, ok := v.Regexes[id]; ok {
			return "(" + regStr + ")"
		}
		return m
	})

	// 5. Replace List tokens (Match the escaped '\$L')
	reListTokens := regexp.MustCompile(`\\\$L\d{1,2}`)
	rePattern = reListTokens.ReplaceAllStringFunc(rePattern, func(m string) string {
		// m is "\$L1", we need "L1"
		id := m[2:]
		if list, ok := v.Lists[id]; ok {
			keys := make([]string, 0, len(list))
			for k := range list {
				keys = append(keys, regexp.QuoteMeta(k))
			}
			sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })
			return "(" + strings.Join(keys, "|") + ")"
		}
		return m
	})

	finalRegex, err := regexp.Compile("^" + rePattern + "$")
	if err != nil {
		return false, fmt.Errorf("invalid template regex: %w", err)
	}

	return finalRegex.MatchString(filename), nil
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

/*
func DemoValidator() {
	v := &Validator{
		Lists:   map[string]map[string]string{"L1": {"LOG": "Logfile", "ERR": "Errorfile"}},
		Regexes: map[string]string{"R1": `[a-z]{3}`},
		Predefined: map[string]KeyValuePairArray{
			"C": {
				{Key: "NL", Value: "Nederland"},
				{Key: "NLD", Value: "Nederland"},
				{Key: "PA", Value: "Panama"},
				{Key: "PAN", Value: "Panama"},
				{Key: "US", Value: "U.S.A."},
				{Key: "USA", Value: "U.S.A."},
			},
		},
	}

	// Test case: List, Date, Time, Regex, Wildcard
	template := "app_$L1_%D_%T_$R1_*.txt"
	filename := "app_LOG_20231025_143005_abc_session99.txt"

	isValid, _ := v.ValidateFilename(template, filename)
	fmt.Printf("Filename: %s\nValid: %v\n", filename, isValid)
}
*/
