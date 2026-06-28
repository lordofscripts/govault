/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lordofscripts/govault/cmd"
	"github.com/lordofscripts/govault/internal/fs"

	"github.com/lordofscripts/goapp/app"
	"github.com/lordofscripts/goapp/flagx"
	vault "github.com/lordofscripts/govault"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

var (
	OptVerbose   bool = false
	OptRecursive bool = false
)

var (
	ErrNoJsonPicsDocs  error = fmt.Errorf("please specify configuration JSON file path or app:pics or app:docs")
	ErrNoJsonTemplates error = fmt.Errorf("please specify configuration JSON file path or app:templates")
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Help command (-help)
func Help() {
	flag.PrintDefaults()
	fmt.Println()

	vault.ModuleVersion.BuyMeCoffee(vault.Reverse(vault.CO2))
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

func main() {
	// I. Command-line Flags
	/*
		$govault COMMAND {OPTIONS}

		Flag			Command					 Meaning
		--------------	------------------		 ---------------------
		-help			*						 Help
		-v				*						 Verbose output
		-list			tags					 Display Tag cloud
		-query string	tags					 Query folder tree for CSV tags
		-tokens			template				 Display Tag cloud
		-rules			template				 Query folder tree for CSV tags
		-r				template				 Recursive
		-check			template				 Template conformance of physical tree
		-root string	create 					 Physical path root
		-overwrite		create					 Overwite .foldermeta.json
		-perm octal		create			 		 Directory permissions
	*/
	const (
		CMD_TAGS     string = "tags"
		CMD_TEMPLATE string = "template"
		CMD_CREATE   string = "create"
		CMD_UPDATE   string = "update"
		CMD_HELP     string = "help"
	)

	// (1.a) CLI flag definition
	// (1.a.1) Common options
	var jsonPath string
	var optHelp bool

	// (1.a.2) Tags sub-command
	subCom := flagx.NewFlagSubCommander()
	var optList bool
	var optQuery string
	tagsCmd := subCom.Define(CMD_TAGS, flag.ExitOnError)
	tagsCmd.Usage = HelpTags
	tagsCmd.BoolVar(&optList, "list", false, "List tag cloud")
	tagsCmd.StringVar(&optQuery, "query", "", "Query folders matching 1 or more tags")
	tagsCmd.StringVar(&jsonPath, "json", "", "path to JSON folder tree file or app:pics or app:docs")
	tagsCmd.BoolVar(&optHelp, "help", false, "Help with options")

	// (1.a.3) Template sub-command
	var optRules, optTokens, optCheck bool
	templateCmd := subCom.Define(CMD_TEMPLATE, flag.ExitOnError)
	templateCmd.Usage = HelpTemplates
	templateCmd.BoolVar(&optRules, "rules", false, "Display defined template Rules")
	templateCmd.BoolVar(&optTokens, "tokens", false, "Display defined template Tokens")
	templateCmd.BoolVar(&OptRecursive, "r", false, "Recursive search on folders")
	templateCmd.BoolVar(&optCheck, "check", false, "Check template conformance on physical tree")
	templateCmd.StringVar(&jsonPath, "json", cmd.GOVAULT_TEMPLATES, "path to JSON folder tree file or app:pics or app:docs")
	templateCmd.BoolVar(&optHelp, "help", false, "Help with options")

	// (1.a.4) Create sub-command
	var optPerm, optRoot string
	var optOverwrite bool
	createCmd := subCom.Define(CMD_CREATE, flag.ExitOnError)
	createCmd.Usage = HelpCreate
	createCmd.StringVar(&optPerm, "perm", "0755", "Directory permissions in Octal")
	createCmd.StringVar(&optRoot, "root", ".", "Root path where folders should be created")
	createCmd.BoolVar(&optOverwrite, "overwrite", false, "Overwrite existing .foldermeta.json")
	createCmd.StringVar(&jsonPath, "json", "", "path to JSON folder tree file or app:pics or app:docs")
	createCmd.BoolVar(&optHelp, "help", false, "Help with options")

	// (1.a.5) Update sub-command
	updateCmd := subCom.Define(CMD_UPDATE, flag.ExitOnError)
	updateCmd.Usage = HelpUpdate
	updateCmd.StringVar(&optPerm, "perm", "0755", "Directory permissions in Octal")
	updateCmd.StringVar(&optRoot, "root", ".", "Root path where folders should be created")
	updateCmd.StringVar(&jsonPath, "json", "", "path to JSON folder tree file or app:pics or app:docs")
	updateCmd.BoolVar(&optHelp, "help", false, "Help with options")

	// (1.a.6) Help sub-command
	helpCmd := subCom.Define(CMD_HELP, flag.ExitOnError)
	helpCmd.Usage = Help

	// II. Application Prelude
	vault.ModuleVersion.Copyright(vault.Reverse(vault.CO1), vault.CHR_TRIDENT)
	fmt.Println()

	if err := subCom.Parse(); err != nil {
		app.DieWithError(err, 1)
	} else {
		var exitCode int = 0
		var err error = nil

		// (2.a) -help used on a valid sub-command
		if optHelp {
			if flagset := subCom.ChosenCommand(); flagset != nil {
				flagset.Usage()
			}
			os.Exit(0)
		}

		// (2.b) pre-process the JSON configuration filename.
		// app:docs|app:pics|app:templates get converted
		// pics & docs apply to tags, create & update commands, temps to templates.
		jsonPath = cmd.VaultConfig.ToFilename(jsonPath)

		// (2.c) Process sub-command
		switch subCom.ChosenCommandName() {
		case CMD_TAGS:
			// $govault tags [-list|-query 'TAGS']
			if len(jsonPath) == 0 {
				err = ErrNoJsonPicsDocs
				exitCode = 1
				break
			}

			var folders []fs.Folder
			folders, err = fs.LoadFolderTable(jsonPath)
			if err != nil {
				exitCode = 10
				break
			}
			if optList && len(optQuery) != 0 {
				exitCode = 11
				err = fmt.Errorf("%s -list and -query are mutually exclusive", CMD_TAGS)
				break
			}

			switch {
			case optList:
				ListTags(folders, true)
			case len(optQuery) != 0:
				Query(folders, optQuery)
			default:
				app.Die("tags requires [-list|-query]", 5)
			}

		case CMD_TEMPLATE:
			// $govault template [-check|-rules|-tokens] [-r]
			switch {
			case optCheck:
				var cwd string
				/*
				* Embedded function: Format directory summary output
				 */
				outputFormatter := func(count, percent int, dirname string) {
					fmt.Printf("%10d %3d%% %s\n", count, percent, dirname)
				}

				if cwd, err = os.Getwd(); err == nil {
					exitCode, err = NameCheckOnDir(cmd.VaultConfig.TemplateConfigFile, cwd, OptRecursive, outputFormatter)
				}
			case optRules:
				err, exitCode = List(cmd.VaultConfig.TemplateConfigFile, ListTemplateRules)
			case optTokens:
				err, exitCode = List(cmd.VaultConfig.TemplateConfigFile, ListTemplateTokens)
			default:
				app.Die("template requires [-tokens|-rules|-check]", 5)
			}

		case CMD_CREATE:
			// $govault create [-perm OCTAL][-overwrite] -root PATH
			if len(jsonPath) == 0 {
				err = ErrNoJsonPicsDocs
				exitCode = 1
				break
			}
			var folders []fs.Folder
			folders, err = fs.LoadFolderTable(jsonPath)
			err = Create(folders, optPerm, optRoot, optOverwrite)

		case CMD_UPDATE:
			// $govault update [-perm OCTAL] -root PATH
			if len(jsonPath) == 0 {
				err = ErrNoJsonPicsDocs
				exitCode = 1
				break
			}
			var folders []fs.Folder
			folders, err = fs.LoadFolderTable(jsonPath)
			err = Update(folders, optPerm, optRoot)

		case CMD_HELP:
			// $govault help {tags|template|create|update}
			if len(helpCmd.Args()) != 1 || !subCom.IsValidSubcommand(helpCmd.Arg(0)) {
				fmt.Printf("Usage:\n\t%s help {tags|template|create|update}\n", os.Args[0])
				os.Exit(1)
			} else {
				switch helpCmd.Arg(0) {
				case CMD_CREATE:
					HelpCreate()
				case CMD_UPDATE:
					HelpUpdate()
				case CMD_TAGS:
					HelpTags()
				case CMD_TEMPLATE:
					HelpTemplates()
				}
				os.Exit(0)
			}
		}

		if err != nil || exitCode != 0 {
			app.DieWithError(err, exitCode)
		}
	}

	fmt.Println()
	vault.ModuleVersion.BuyMeCoffee(vault.Reverse(vault.CO3))
}
