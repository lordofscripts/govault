/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           goVault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Token Entry object in JSON template configuration file.
 *-----------------------------------------------------------------*/
package fs

import (
	"fmt"
	"regexp"
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

// An entry for a token used to describe part of a valid filename
type TokenEntry struct {
	Token       string         `json:"token"`
	Description string         `json:"description"`
	Values      []KeyValuePair `json:"values,omitempty"` // empty if using a rule
	Rule        *TokenRule     `json:"rule,omitempty"`   // empty if using here-values @audit deprecate with Validator
}

// Key-value pairs used in a TokenEntry
type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Does the token entry has a listed value valStr?
func (te TokenEntry) HasValue(valStr string, caseSensitive bool) bool {
	found := false
	for _, kvp := range te.Values {
		if caseSensitive {
			found = strings.EqualFold(kvp.Key, valStr)
		} else {
			found = kvp.Key == valStr
		}

		if found {
			break
		}
	}
	return found
}

func (te TokenEntry) HasRule() bool {
	return te.Rule != nil
}

func (te TokenEntry) GetRule() *TokenRule {
	var result *TokenRule = nil
	if te.HasRule() {
		return te.Rule
	}
	return result
}

func (te TokenEntry) ValidateRE(reStr string, valueStr string) bool {
	/* Validator for %RFincaDGI
		tests := []string{
		"F123456-8715",
		"F123456",
		"F12345678-0000",
		"f123456-8715",   // invalid: lowercase f
		"F12345-8715",    // invalid: only 5 digits
		"F123456789-0000",// invalid: 9 digits
		"F123456-871",    // invalid: only 3 digits after dash
	}

	for _, t := range tests {
		fmt.Printf("%s -> %v\n", t, IsValid(t))
	}
	*/
	var re = regexp.MustCompile(reStr)
	return re.MatchString(valueStr)
}

// Implements fmt.Stringer for key-value pair
func (kv KeyValuePair) String() string {
	return fmt.Sprintf("%12s %s", kv.Key, kv.Value)
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
