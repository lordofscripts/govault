//go:build ignore

/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Utility object to handle command-line sub-commands using the
 * standard flag package.
 * Superseeded by github.com/lordofscripts/goapp/flagx/
 *-----------------------------------------------------------------*/
package cmd

import (
	"flag"
	"fmt"
	"os"
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

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

type FlagSubCommander struct {
	subCommands      map[string]*flag.FlagSet
	chosenSet        *flag.FlagSet
	chosenSubCommand string
}

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (ctor) utility for dealing with flag subcommands
func NewFlagSubCommander() *FlagSubCommander {
	return &FlagSubCommander{
		subCommands: make(map[string]*flag.FlagSet, 0),
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Defines a new flag set (subcommand) with the subcommand name and
// error handling directive (ExitOnError, ContinueOnError, PanicOnError).
// if already defined a warning message is produced but returns the
// previously-defined flag set.
func (fsc *FlagSubCommander) Define(subCommandName string, eh flag.ErrorHandling) *flag.FlagSet {
	if fset, okay := fsc.subCommands[subCommandName]; okay {
		fmt.Fprintf(os.Stderr, "warning: subcommand '%s' already defined", subCommandName)
		return fset
	} else {
		subCommandName = strings.ToLower(subCommandName)
		fset := flag.NewFlagSet(subCommandName, eh)
		fsc.subCommands[subCommandName] = fset
		return fset
	}
}

// Returns true if name represents a valid/registered sub-command. It
// is normalized to lowercase.
func (fsc *FlagSubCommander) IsValidSubcommand(name string) bool {
	_, ok := fsc.subCommands[strings.ToLower(name)]
	return ok
}

// Verifies that the subcommand name is given and parses the flags after
// the subcommand, returning nil or an error if any.
func (fsc *FlagSubCommander) Parse() error {
	var err error = nil
	fsc.chosenSet = nil
	fsc.chosenSubCommand = ""

	if len(os.Args) < 2 {
		err = fmt.Errorf("expected sub-commands: %s", fsc.String())
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		// retrieve the sub-command from the CLI
		subCommand := strings.ToLower(os.Args[1])
		// verify sub-command is known
		if subFlags, exists := fsc.subCommands[subCommand]; !exists {
			err = fmt.Errorf("unknown sub-command '%s'", subCommand)
			fmt.Fprintln(os.Stderr, err.Error())
		} else {
			// parse all flags after the subcommand
			err = subFlags.Parse(os.Args[2:])
			// prepare for SubCommandName() && SubCommand()
			fsc.chosenSet = fsc.subCommands[subCommand]
			fsc.chosenSubCommand = subCommand
		}
	}

	return err
}

// Returns the list of defined subcommands. i.e. "foo,bar"
func (fsc *FlagSubCommander) String() string {
	var sb strings.Builder
	count := 0
	for k := range fsc.subCommands {
		count++
		sb.WriteString(k)
		if count < len(fsc.subCommands) {
			sb.WriteRune(',')
		}
	}
	return sb.String()
}

// After Parse() without error this method returns the
// name of the sub-command.
func (fsc *FlagSubCommander) ChosenCommandName() string {
	return fsc.chosenSubCommand
}

// After Parse() without error this method returns the
// flagset corresponding to the selected sub-command.
func (fsc *FlagSubCommander) ChosenCommand() *flag.FlagSet {
	return fsc.chosenSet
}

// Returns the free arguments after the flags of the sub-command.
func (fsc *FlagSubCommander) FreeArguments() []string {
	return flag.Args()
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

/*
func DemoMain() {
	const (
		SUBCMD_ENCODE string = "encrypt"
		SUBCMD_DECODE string = "decrypt"
	)
	subCom := NewFlagSubCommander()

	HelpEncrypt := func() {
		fmt.Println("aes encrypt -text 'PLAIN'")
	}
	encodeCmd := subCom.Define(SUBCMD_ENCODE, flag.ExitOnError)
	encodeCmd.Usage = HelpEncrypt
	plainTxt := encodeCmd.String("plain", "", "Plain text")

	HelpDecrypt := func() {
		fmt.Println("aes decrypt -cipher 'SECRET'")
	}
	decodeCmd := subCom.Define(SUBCMD_DECODE, flag.ExitOnError)
	decodeCmd.Usage = HelpDecrypt
	cipherTxt := decodeCmd.String("cipher", "", "Ciphered text to decrypt")

	if err := subCom.Parse(); err != nil {
		fmt.Println(err.Error())
	} else {
		// do something
		switch subCom.SubCommandName() {
		case SUBCMD_ENCODE:
			// do something
			AesCipher(plainTxt)
		case SUBCMD_DECODE:
			// do something
			AesCipher(cipherTxt)
		}
	}
}
*/
