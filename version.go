/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2026 Dídimo Grimaldo T.
 *							 go-vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A tagged Picture & Document Vault maker with SQL-like query
 * capabilities.
 *-----------------------------------------------------------------*/
package vault

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/lordofscripts/goapp/app"
	//_ "embed"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	_NAME          string = "goVault"
	_DESC          string = "File vault maker with SQL-like query capability"
	MANUAL_VERSION string = "1.1.0"
)

const (
	// Useful Unicode Characters
	CHR_COPYRIGHT       = '\u00a9'      // ©
	CHR_REGISTERED      = '\u00ae'      // ®
	CHR_GUILLEMET_L     = '\u00ab'      // «
	CHR_GUILLEMET_R     = '\u00bb'      // »
	CHR_TRADEMARK       = '\u2122'      // ™
	CHR_SAMARITAN       = '\u214f'      // ⅏
	CHR_PLACEOFINTEREST = '\u2318'      // ⌘
	CHR_HIGHVOLTAGE     = '\u26a1'      // ⚡
	CHR_TRIDENT         = rune(0x1f531) // 🔱
	CHR_SPLATTER        = rune(0x1fadf)
	CHR_WARNING         = '\u26a0' // ⚠
	CHR_EXCLAMATION     = '\u2757'
	CHR_SKULL           = '\u2620' // ☠

	CO1 = "odlamirG omidiD 6202-5202)C("
	CO2 = "stpircS fO droL 6202-5202)C("
	CO3 = "gnitirwnitsol"
)

var (
	ModuleVersion app.PackageVersion = app.NewReleaseCandidateVersion(_NAME, _DESC, MANUAL_VERSION, 3)
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Funny LordOfScripts logo
func Logo() string {
	const (
		whiteStar rune = '\u269d' // ⚝
		unisex    rune = '\u26a5' // ⚥
		hotSpring rune = '\u2668' // ♨
		leftConv  rune = '\u269e' // ⚞
		rightConv rune = '\u269f' // ⚟
		eye       rune = '\u25d5' // ◕
		mouth     rune = '\u035c' // ͜	‿ \u203f
		skull     rune = '\u2620' // ☠
	)
	return fmt.Sprintf("%c%c%c %c%c", leftConv, eye, mouth, eye, rightConv)
	//fmt.Sprintf("(%c%c %c)", eye, mouth, eye)
}

// Hey! My time costs money too!
func BuyMeCoffee(coffee4 ...string) {
	const (
		coffee rune = '\u2615' // ☕
	)

	var recipient string
	if len(coffee4) == 0 {
		recipient = Reverse(CO3)
	} else {
		recipient = coffee4[0]
	}

	fmt.Printf("\t%c Buy me a Coffee? https://www.buymeacoffee/%s\n", coffee, recipient)
}

func Copyright(owner string, withLogo bool) {
	//fmt.Printf("\t\u2720 %s %s \u269d\n", Version, Reverse(owner))
	fmt.Printf("\t%c %s %s %c\n", CHR_TRIDENT, ModuleVersion, Reverse(owner), CHR_TRIDENT)
	fmt.Println("\t\t\t\t", Logo())
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// get the current GO language version
func GoVersion() string {
	ver := strings.Replace(runtime.Version(), "go", "", -1)
	return ver
}

// retrieve the current GO language version and compare it
// to the minimum required. It returns the current version
// and whether the condition current >= min is fulfilled or not.
func GoVersionMin(min string) (string, bool) {
	current := GoVersion()
	ok := current >= min
	return current, ok
}
