/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           goVault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package fs

import "fmt"

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

// A Rule descriptor is used to validate custom tokens. It
// can be a "|" separated list of values OR a Regular Expression validator
type TokenRule struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
	Title string `json:"title"`
	Regex string `json:"regex,omitempty"`
	List  string `json:"list,omitempty"` // @audit deprecate with Validator
}

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Implements fmt.Stringer for either the expression rule
// or the list rule but cannot be both.
func (tr TokenRule) String() string {
	if tr.IsExpression() {
		return fmt.Sprintf("%30s (%s) %s", tr.Title, tr.Token, tr.Regex)
	} else {
		return fmt.Sprintf("%30s (%s) %s", tr.Title, tr.Token, tr.List)
	}
}

// Returns true if the rule is a Regular Expression,
// false if it is a list.
func (tr TokenRule) IsExpression() bool {
	return len(tr.Regex) != 0
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
