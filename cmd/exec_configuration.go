/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           goVault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Application configuration files found on ~/{CONFIG_DIR}/coralys/govault/
 * where we may find one or more of the following:
 * - govault.json  general application configuration (not used yet)
 * - govault_templates.json Filename templates
 * - govault_pics.json  Directory hierarchy for Pictures folder
 * - govault_docs.json  Directory hierarchy for Documents folder
 *-----------------------------------------------------------------*/
package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/lordofscripts/goapp/app"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/
const (
	FN_CONFIG_MAIN      string = "govault.json"           // application configuration (not used)
	FN_CONFIG_TEMPLATES string = "govault_templates.json" // Filename template config
	FN_CONFIG_TREE_DOC  string = "govault_docs.json"      // Documents hierarchy config
	FN_CONFIG_TREE_PIC  string = "govault_pics.json"      // Pictures hierarchy config

	ORG_NAME       string      = "coralys"
	APP_NAME       string      = "govault"
	META_FILE_MODE os.FileMode = 0644
)

const (
	GOVAULT_PICS      string = "app:pics"      // when used for -json use {CONFIG}/govault_pics.json
	GOVAULT_DOCS      string = "app:docs"      // when used for -json use {CONFIG}/govault_docs.json
	GOVAULT_TEMPLATES string = "app:templates" // when used for -json use {CONFIG}/govault_templates.json
)

var (
	VaultConfig *GoVaultConfig
)

/* ----------------------------------------------------------------
 *             E M B E D D E D    R E S O U R C E S
 *-----------------------------------------------------------------*/

// Contains a sample Filename Template configuration file to start with
//
//go:embed json/empty_template.json
var jsonTemplateData []byte

// Contains a sample Documents tree configuration file to start with
//
//go:embed json/tree_sample_docs.json
var jsonDocsData []byte

// Contains a sample Pictures tree configuration file to start with
//
//go:embed json/tree_sample_pics.json
var jsonPicsData []byte

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// goVault Global Configuration
type GoVaultConfig struct {
	MainConfigFile     string // govault.json
	TemplateConfigFile string // govault_templates.json
	PicturesTreeFile   string // govault_pics.json
	DocumentsTreeFile  string // govault_docs.json
}

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                     I N I T I A L I Z E R
 *-----------------------------------------------------------------*/

// Prepare application's global configuration prior to execution of
// main and ensure configuration directory and (sample) configuration
// files exist or create them.
func init() {
	// (a) instantiate global configuration
	VaultConfig = &GoVaultConfig{}
	// (b) setup filename template validation
	VaultConfig.TemplateConfigFile = initConfig("_templates", jsonTemplateData)
	// (c) setup Pictures directory tree sample configuration
	VaultConfig.PicturesTreeFile = initConfig("_pics", jsonPicsData)
	// (d) setup Documents directory tree sample configuration
	VaultConfig.DocumentsTreeFile = initConfig("_docs", jsonDocsData)
}

// Ensure the configuration file exists, create a dummy/sample
// file if necessary with the supplied initial data. Request the user
// to customize it.
func initConfig(suffix string, initialData []byte) string {
	var err error
	var cfgJSON string
	if cfgJSON, err = app.EnsureConfigWithSuffix(ORG_NAME, APP_NAME, suffix, ".json"); err != nil {
		// did not exist, attempt to create one with sample data
		if err := os.WriteFile(cfgJSON, initialData, META_FILE_MODE); err != nil {
			app.DieWithError(err, 1)
		} else {
			fmt.Fprintf(os.Stderr, "Please customize your config file %s\n", cfgJSON)
		}
	}

	return cfgJSON
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// If nameOrToken is any of app:docs|app:pics|app:templates it returns
// the corresponding fully-qualified filename of the corresponding
// JSON configuration file. Else, nameOrToken is returned as-is.
func (gvc *GoVaultConfig) ToFilename(nameOrToken string) string {
	switch nameOrToken {
	case GOVAULT_DOCS:
		nameOrToken = gvc.DocumentsTreeFile
	case GOVAULT_PICS:
		nameOrToken = gvc.PicturesTreeFile
	case GOVAULT_TEMPLATES:
		nameOrToken = gvc.TemplateConfigFile
	default:
	}
	return nameOrToken
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
