/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"fmt"
	"os"
	"sort"

	vault "github.com/lordofscripts/govault"
	"github.com/lordofscripts/govault/cmd"
	"github.com/lordofscripts/govault/internal/fs"
	"github.com/lordofscripts/govault/internal/sync"
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

func HelpSync() {
	fmt.Println(os.Args[0], " sync -json PATH {Option} [-virtual|-system]")
	fmt.Println("Option:")
	fmt.Println("\t-json string           The path to JSON config or app:docs/app:pics")
	fmt.Println("\t-root string           The root directory to examine")
	fmt.Println("\t-virtual               Update the configuration (virtual)")
	fmt.Println("\t-system                Create/Delete subdirectories (filesystem)")
	fmt.Println("\t-dry                   Dry run, only lists actions without altering filesystem")
	fmt.Println("WARNING!!!")
	fmt.Println("\tMake sure to first use it with -dry. Erased folders cannot be recovered!")
	fmt.Println()

	vault.ModuleVersion.BuyMeCoffee(vault.Reverse(vault.CO3))
}

// Update the left tree (Virtual Directory Structure) from the JSON file
// to match that of the right tree (Physical Directory Structure)
func SynchronizeVirtual(foldersV []fs.Folder, cwd string, dryRun, showOperation bool, toFilename string) error {
	var err error = nil
	changes := sync.CompareFolders(foldersV, cwd, sync.UpdateLeft)
	if len(changes) != 0 {
		fmt.Println("Deltas : ", len(changes))

		if showOperation {
			// List them in ascending order
			sort.Slice(changes, func(i, j int) bool {
				return changes[i].Target.Path < changes[j].Target.Path
			})

			for _, chg := range changes {
				cmd.BrightPurple("\t" + chg.String() + "\n")
			}
			fmt.Println()
		}

		err = sync.SynchronizeVirtualFolders(changes, &foldersV) //  always returns nil

		if !dryRun {
			// @todo Update JSON file
			fs.SaveFolderTable(foldersV, toFilename)
		}
		fs.DumpVirtual(foldersV, true) // @note debug only
	} else {
		fmt.Println("· The real filesystem matches 100% with configuration (virtual)")
	}

	return err
}

// Update the right tree (Physical Directory Structure) to match that
// of the left tree (Virtual Directory Structure) from the JSON file
func SynchronizeReal(foldersV []fs.Folder, cwd string, dryRun, showOperation bool) error {
	var err error = nil
	changes := sync.CompareFolders(foldersV, cwd, sync.UpdateRight)
	if len(changes) != 0 {
		fmt.Println("Deltas: ", len(changes))

		if showOperation && len(changes) > 0 {
			for _, chg := range changes {
				cmd.BrightPurple(fmt.Sprintf("\t%s %s\n", chg.Action, chg.Target.Path))
			}
			fmt.Println()
		}

		err = sync.SynchronizeRealFolders(changes, dryRun)
	} else {
		fmt.Println("· The real filesystem matches 100% with configuration (virtual)")
	}

	return err
}
