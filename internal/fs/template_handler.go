/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	TOKEN_COUNTRY     string = "%C" // 2-3 letter ISO country code
	TOKEN_PERSON      string = "%P" // 3-4 letras para personas
	TOKEN_DATE        string = "%D" // para fechas en YYYYMMDD
	TOKEN_INSTITUTION string = "%I" // para instituciones
	TOKEN_MEDICAL     string = "%M" // para especialidades medicas
	TOKEN_FREE        string = "*"  // free-form (libre)
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

type TemplateHandler struct {
	templates *Template
	validator *Validator
}

// JSON name template file templates.json
type Template struct {
	Tokens []TokenEntry `json:"tokens"`
	Rules  []TokenRule  `json:"rules"`
}

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// Ctor. Creates a new instance of the filename template validator
func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{
		templates: nil,
		validator: NewValidator(),
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

func (th *TemplateHandler) Load(jsonFilename string) error {
	// for Create or Query we need the JSON file with metadata
	data, err := os.ReadFile(jsonFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading JSON: %v\n", err)
		return err
	}

	var templates Template
	if err = json.Unmarshal(data, &templates); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing JSON: %v\n", err)
		return err
	}

	th.templates = &templates
	th.setupValidator()
	return nil
}

// validates against a regular expression
func (th *TemplateHandler) ValidateRE(reStr string, valueStr string) bool {
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

/*
func (th *TemplateHandler) ValidateTemplate(templateStr string, filename string) error {
	// (a) parameter validation
	if len(templateStr) == 0 {
		return fmt.Errorf("cannot verify nil or empty template")
	}
	if len(filename) == 0 {
		return fmt.Errorf("cannot verify nil or empty filename")
	}

	// (b)
}*/

// Validates a filename against the provided template. The template is
// composed of tokens and rules from the template configuration file (govault_templates.json)
// Returns true if it conforms to the template, otherwise false. If an error is
// encountered (other than validation) an error is returned.
func (th *TemplateHandler) ValidateTemplate(templateStr string, filename string) (bool, error) {
	passed, err := th.validator.ValidateFilename(templateStr, filename)
	return passed, err
}

func (th *TemplateHandler) GetTokens() []TokenEntry {
	return th.templates.Tokens
}

func (th *TemplateHandler) GetRules() []TokenRule {
	return th.templates.Rules
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// transform the read-in JSON catalog to descriptors suitable
// for the Validator object.
func (th *TemplateHandler) setupValidator() {
	for _, utoken := range th.templates.Tokens {
		if !utoken.HasRule() {
			// a list of values
			tokenID := []rune(utoken.Token)[1] // skip leading %
			th.validator.AddPredefinedItem(tokenID, utoken.Values)
		} else {
			// deprecated
		}
	}

	for _, urule := range th.templates.Rules {
		if urule.IsExpression() {
			th.validator.AddRuleItem(byte(urule.Id), urule.Regex)
		} else {
			th.validator.AddRuleItem(byte(urule.Id), urule.List) // The | list is a valid RegEx too!
		}
	}
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
