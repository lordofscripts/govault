/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           goVault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Warning codes produced by GoVault application
 *-----------------------------------------------------------------*/
package main

import "github.com/lordofscripts/goapp/app"

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	WarnMissingMetadata app.WarningCode = iota // no .foldermeta.json file
	WarnIgnoredFile
	WarnEmptyTemplate
)
