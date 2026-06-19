/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                          goVault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Filename templates support.
 *-----------------------------------------------------------------*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	vault "github.com/lordofscripts/govault"
	vfs "github.com/lordofscripts/govault/internal/fs"

	"github.com/lordofscripts/goapp/app"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	ListTemplateRules TemplateAction = iota
	ListTemplateTokens
	CheckFilenames
)

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

type TemplateAction byte

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Show help about the 'tags' command
func HelpTemplates() {
	fmt.Println(os.Args[0], " template {Option}")
	fmt.Println("Option:")
	fmt.Println("\t-rules   List defined Template Rules")
	fmt.Println("\t-tokens  List defined Template Tokens")
	fmt.Println("\t-check	Check filenames for template conformance")
	fmt.Println("\t-r		Recursive verification (-check)")

	fmt.Println()
	Version.BuyMeCoffee(vault.Reverse(vault.CO3))
}

// Main entry point for handling Template-related commands
func HandleTemplateCommands(action TemplateAction) (err error, exitCode int) {
	exitCode = 0
	err = nil

	switch action {
	case ListTemplateRules:
		fallthrough
	case ListTemplateTokens:
		err, exitCode = List(VaultConfig.TemplateConfigFile, action)

	case CheckFilenames:
		var cwd string
		/*
		 * Embedded function: Format directory summary output
		 */
		outputFormatter := func(count, percent int, dirname string) {
			fmt.Printf("%10d %3d%% %s\n", count, percent, dirname)
		}

		if cwd, err = os.Getwd(); err == nil {
			NameCheckOnDir(VaultConfig.TemplateConfigFile, cwd, OptRecursive, outputFormatter)
		}

	default:
		return fmt.Errorf("unknown TemplateAction value: %d", action), 50
	}

	return
}

// OPTION: -list tokens|rules
// Lists the rules or tokens defined in the Folder name template
func List(cfgFilename string, what TemplateAction) (err error, exitCode int) {
	exitCode = 0
	err = nil

	// templates read from configuration file
	th := vfs.NewTemplateHandler()
	if err = th.Load(cfgFilename); err != nil {
		exitCode = 51
	} else {
		switch what {
		case ListTemplateTokens:
			tokens := th.GetTokens()
			for _, token := range tokens {
				fmt.Printf("Token %s (%s)\n", token.Token, token.Description)
				if token.HasRule() {
					fmt.Printf("\tRule: %s\n", token.Rule)
				} else {
					//fmt.Printf("%*s Description:\n", 29, "Value")
					fmt.Printf("%*s Description\n", 12, "Value")
					for _, v := range token.Values {
						fmt.Println(v.String())
					}
				}
			}

		case ListTemplateRules:
			rules := th.GetRules()
			fmt.Printf("\tRules:\n")
			fmt.Printf("%*s (Token) Rule\n", 30, "Title")
			for _, r := range rules {
				fmt.Println(r.String())
			}

		default:
			err = fmt.Errorf("unknown TemplateAction in List()")
			exitCode = 52
		}
	}

	return
}

// OPTION: -check
// Non-recursive check of all files in directory to verify whether their
// filenames comply with the name template (if defined).
// Returns the percentage of complying files.
func NameCheckOnDir(cfgFilename, directory string, isRecursive bool, formatterCB func(int, int, string)) (int, error) {
	// I. read templates from configuration file
	th := vfs.NewTemplateHandler()
	if err := th.Load(cfgFilename); err != nil {
		return 61, err
	}

	/*
	 * Embedded function: Pretty Print Warning
	 */
	/*
		warningOutOld := func(format string, args ...any) {
			fmt.Fprintln(os.Stderr, "\t\t\t", "⚠️ Warning ⚠️")
			fmt.Fprintf(os.Stderr, "\t·"+format, args...)
		}*/

	// temporary until next patch of lordofscripts/goapp v1.4.1
	warnCode := func(w *app.Warning) string { // @audit deprecate from goapp v1.4.1
		return fmt.Sprintf("WARN-%03d", w.Code)
	}

	warningOut := func(warn *app.Warning, dirName string, full bool) {
		// out-of-sync if we use Stderr :()
		fmt.Fprintf(os.Stdout, "\t%c %s: %s\n", vault.CHR_HIGHVOLTAGE, warnCode(warn), dirName)
		if full {
			fmt.Fprintf(os.Stdout, "\t· %s\n", warn.Error())
		}
	}

	/*
	 * Embedded function: Get List of Directories at This Level
	 */
	getDirectoriesHere := func(dirname string) ([]string, error) {
		entries, err := os.ReadDir(dirname)
		if err != nil {
			return nil, err
		}

		dirs := make([]string, 0)
		for _, e := range entries {
			if e.IsDir() { // faster but does NOT follow symbolic links, else use os.Stat() IsDir()
				fullname := filepath.Join(dirname, e.Name())
				// add fully-qualified directory name
				dirs = append(dirs, fullname)
			}
		}
		return dirs, nil
	}

	/*
	 * Embedded function: Process Directory
	 */
	processDirectory := func(dirname string) ([]string, *app.Warning, error) {
		var err error = nil
		dirs := make([]string, 0)
		if isRecursive {
			dirs, err = getDirectoriesHere(dirname)
			if err != nil {
				return dirs, nil, err
			}
		}

		totalCnt := 0
		// (a) Collect filenames in this directory (non-recursive)
		filesL, err := os.ReadDir(dirname)
		if err != nil {
			return nil, nil, err
		}

		// (b) is there a .foldermeta.json? if none there isn't anything to do
		metaFilename := filepath.Join(dirname, FOLDER_META_FILENAME)
		if app.CheckFileExistsAndReadable(metaFilename) != nil {
			w := app.NewWarning(WarnMissingMetadata, "Missing metadata file")
			//warningOut("This probably isn't a vaulted directory.")
			return dirs, w, nil
		}

		// (c) process .foldermeta.json
		var folderMeta *vfs.Meta
		if folderMeta, err = vfs.DeserializeFolderMeta(metaFilename); err != nil {
			w := app.NewWarning(WarnIgnoredFile, "NameCheck deserialize error %v\n", err)
			//fmt.Fprintf(os.Stderr, "NameCheck deserialize error %v\n", err)
			return dirs, w, nil
		}

		// (d) validate folder metadata
		if folderMeta == nil {
			//warningOut("This isn't a vaulted directory.")
			return dirs, app.NewWarning(WarnMissingMetadata, "Not a vaulted directory"), nil
		}
		if len(folderMeta.Template) == 0 {
			//warningOut("directory is not bound to filename templates: %s", dirname)
			return dirs, app.NewWarning(WarnEmptyTemplate, "dir not bound to templates"), nil
		}
		if folderMeta.Encrypted {
			//warningOut("This folder COULD contain encrypted/encoded files.")
		}
		//fmt.Printf("Tags: %s\n", folderMeta.Tags)

		// (e) Output
		validCnt := 0
		for _, f := range filesL {
			if !f.IsDir() {
				if f.Name() == FOLDER_META_FILENAME {
					continue
				}
				totalCnt++
				if passed, _ := th.ValidateTemplate(folderMeta.Template, f.Name()); passed {
					validCnt++
				}
				if OptVerbose {
					fmt.Printf("%3d %s\n", totalCnt, f.Name())
				}
			}
		}

		percentValid := (validCnt * 100 / totalCnt)
		formatterCB(validCnt, percentValid, dirname)

		return dirs, nil, err
	}

	// II. Processing adding other directories if recursive execution
	var allDirs, extraDirs []string
	var warn *app.Warning
	var err error
	var exitCode int = 0
	allDirs, warn, err = processDirectory(directory)
	if warn != nil {
		warningOut(warn, directory, OptVerbose)
	}

	for len(allDirs) > 0 {
		// (2.1) select next working directory, the 1st in the Queue
		cwd := allDirs[0]

		// (2.2) Process and get list of sub-directories. Show short warnings
		extraDirs, warn, err = processDirectory(cwd)
		if warn != nil {
			warningOut(warn, cwd, OptVerbose)
		}

		// (2.3) Update sub-directory queue
		if len(extraDirs) != 0 {
			allDirs = append(allDirs[1:], extraDirs...) // remove 1st add others
		} else {
			allDirs = allDirs[1:] // remove 1st, none more to add
		}
	}

	return exitCode, err
}
