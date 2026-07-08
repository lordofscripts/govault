/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package sync

import (
	"errors"
	"fmt"
	"os"
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

// performs physical OS filesystem synchronization operations
func SynchronizeRealFolders(results []SyncResult, dryRun bool) error {
	const ANSI_BRIGHT_YELLOW string = "\033[93m"
	const ANSI_RESET_COLOR string = "\033[0m"
	var errs []error
	for _, res := range results {
		var thisErr error = nil
		switch res.Action {
		case AddRight:
			if dryRun {
				fmt.Printf("\t%s[Dry-Run] mkdir: %s%s\n", ANSI_BRIGHT_YELLOW, res.PhysicalPath, ANSI_RESET_COLOR)
			} else {
				thisErr = os.MkdirAll(res.PhysicalPath, 0755)
			}
		case DeleteRight:
			if dryRun {
				fmt.Printf("\t%s[Dry-Run] remove: %s%s\n", ANSI_BRIGHT_YELLOW, res.PhysicalPath, ANSI_RESET_COLOR)
			} else {
				thisErr = os.RemoveAll(res.PhysicalPath)
			}
		}

		if thisErr != nil {
			errs = append(errs, thisErr)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
