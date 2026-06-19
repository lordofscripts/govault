/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package cmd

import "fmt"

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	// Colored console ANSI colors control codes
	ansi_RESET         AnsiCode = "\033[0m"
	ansi_COLOR_RESET   AnsiCode = "\033[0m"
	ansi_RED           AnsiCode = "\033[31m"
	ansi_BRIGHT_RED    AnsiCode = "\033[91m"
	ansi_GREEN         AnsiCode = "\033[32m"
	ansi_BRIGHT_GREEN  AnsiCode = "\033[92m"
	ansi_YELLOW        AnsiCode = "\033[33m"
	ansi_BRIGHT_YELLOW AnsiCode = "\033[93m"
	ansi_PURPLE        AnsiCode = "\u001b[35m"
	ansi_BRIGHT_PURPLE AnsiCode = "\u001b[95m"
	ansi_CYAN          AnsiCode = "\u001b[36m"
	ansi_BRIGHT_CYAN   AnsiCode = "\u001b[96m"
	ansi_WHITE         AnsiCode = "\033[37m"
	ansi_BRIGHT_WHITE  AnsiCode = "\033[97m"
	ansi_HIDE_CURSOR   AnsiCode = "\033[?25l" // does not appear to work
	ansi_SHOW_CURSOR   AnsiCode = "\033[?25h" // idem
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// alias for ANSI control codes
type AnsiCode string

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

// Reset the terminal
func Reset() {
	fmt.Print("\033[39m\\033[49m")
}

// clear the screen and move cursor to (0,0)
func Clear() {
	fmt.Print("\033[2J")
}

// print bold text (no CR) and reset bold
func Bolded(str string) {
	fmt.Printf("\033[1m%s\033[21m", str)
}

// start using Bold
func Bold() {
	fmt.Print("\033[1m")
}

// terminate using Bold
func BoldOff() {
	fmt.Print("\033[22m")
}

// Print the args in color
func Color(color AnsiCode, args ...any) {
	fmt.Print(color)
	fmt.Print(args...)
	fmt.Print(ansi_COLOR_RESET)
}

// print in Red
func Red(args ...any) {
	Color(ansi_RED, args...)
}

// print in Green
func Green(args ...any) {
	Color(ansi_GREEN, args...)
}

// print in Yellow
func Yellow(args ...any) {
	Color(ansi_YELLOW, args...)
}

// Print in Magenta
func Purple(args ...any) {
	Color(ansi_PURPLE, args...)
}

// Print in Cyan
func Cyan(args ...any) {
	Color(ansi_CYAN, args...)
}

func BrightRed(args ...any) {
	Color(ansi_BRIGHT_RED, args...)
}

func BrightGreen(args ...any) {
	Color(ansi_BRIGHT_GREEN, args...)
}

func BrightYellow(args ...any) {
	Color(ansi_BRIGHT_YELLOW, args...)
}

func BrightPurple(args ...any) {
	Color(ansi_BRIGHT_PURPLE, args...)
}

func BrightCyan(args ...any) {
	Color(ansi_BRIGHT_CYAN, args...)
}

func BrightWhite(args ...any) {
	Color(ansi_BRIGHT_WHITE, args...)
}
