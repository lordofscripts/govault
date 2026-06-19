package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lordofscripts/goapp/app"
	vault "github.com/lordofscripts/govault"
	"github.com/lordofscripts/govault/cmd"
	"github.com/lordofscripts/govault/internal/fql"
)

const (
	MANUAL_VERSION string = "1.0.0"
)

var (
	Version app.PackageVersion = app.NewReleaseVersion("qVault", "Query Document Vault using SQL", MANUAL_VERSION)
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

func gatherFreeArguments() string {
	var unusedArgs []string
	//for _, arg := range os.Args[1:] { // Skip the program name
	for _, arg := range flag.Args() {
		unusedArgs = append(unusedArgs, arg)
	}
	unusedArgsStr := strings.Join(unusedArgs, " ")
	return unusedArgsStr
}

// colored output of a Row Result from FQLX
func displayRow(row fql.FQLXResult) {
	for k, v := range row {
		cmd.BrightGreen(fmt.Sprintf("%12s: ", k))
		cmd.BrightPurple(fmt.Sprintf("%v", v))
		fmt.Println()
	}
}

// Colored output of Result Query header
func displayQueryHeader(query string, resultQty uint) {
	cmd.BrightWhite(fmt.Sprintf("· %s\n", query))
	cmd.BrightWhite(fmt.Sprintf("· Produced %d records/modifications\n", resultQty))
}

// Colored output of Result Row header
func displayRowHeader(rowNr int) {
	cmd.BrightGreen(fmt.Sprintf("Row #%d\n", rowNr))
}

func SampleRun(jsonFile, container string) {
	// SELECT
	// 1. Form #1: Custom Operator
	//q1 := fmt.Sprintf("SELECT Name,Path FROM %s WHERE Encrypted = true", container)
	q1 := fmt.Sprintf("SELECT Name,Path FROM %s WHERE Encrypted = true", container)

	// 2. Form #2: Specific Field LIKE
	q2 := fmt.Sprintf("SELECT Id,Name,Path FROM %s WHERE Path LIKE 'Documents/%%'", container)

	// 3. Form #3: Global LIKE
	q3 := fmt.Sprintf("SELECT Name,Tags FROM %s WHERE * LIKE '%%Project%%'", container)

	// 4. Form #4: Entire object
	q4 := fmt.Sprintf("SELECT * FROM %s WHERE Name LIKE '%%Lisbon%%'", container)

	// 5. Update
	q5 := "UPDATE folders SET Name = 'Archived' WHERE Id = '1.2'"

	q6 := "UPDATE folders SET Encrypted = false, Tags = 'public' WHERE Name = 'Public Docs'"

	db := fql.NewFQLX().
		Use(jsonFile).
		Connect()

	for QNum, queryStr := range []string{q1, q2, q3, q4, q5, q6} {
		cnt := db.Query(queryStr) // handles both SELECT & UPDATE
		fmt.Println("\nQuery #", QNum+1, " ", queryStr)
		fmt.Printf("· Produced %d records or modifications\n", cnt)

		if db.RowCount() != 0 {
			rows := db.Fetch()
			for n, row := range rows {
				fmt.Println("Row #", n+1)
				fmt.Println(row.String())
			}
			db.Clear()
		}
	}

	/*
		// (Inheritance of Template and Tags)
		// Inherit Template from nearest parent for the "Legal" folder
		inheritedTemplate := fql.InheritFromParent(fql.FolderTemplate, myFolders, "Documents/Legal")
		fmt.Println("Inherited Template: ", inheritedTemplate)

		// Merge Tags for Folder with ID 2.1
		inheritedTags := fql.InheritFromParent(fql.FolderTags, myFolders, "2.1")
		fmt.Println("Inherited Tags: ", inheritedTags)
	*/
}

func CustomRun(jsonFile, queryStr string) {
	db := fql.NewFQLX().
		Use(jsonFile).
		Connect()

	cnt := db.Query(queryStr) // handles both SELECT & UPDATE
	//fmt.Println("\nQuery: ", queryStr)
	//fmt.Printf("\t· Produced %d records or modifications\n", cnt)
	displayQueryHeader(queryStr, cnt)

	if db.RowCount() != 0 {
		rows := db.Fetch()
		for n, row := range rows {
			//fmt.Println("Row #", n+1)
			//fmt.Println(row.String())
			displayRowHeader(n + 1)
			displayRow(row)
		}
		db.Clear()
	}
}

func Help() {
	fmt.Println("Usage:")
	fmt.Println("\tqvault -json /Path/To/Config/folders.json -root TopFolder")
	fmt.Println("\tqvault -json /Path/To/Config/folders.json 'SELECT ...'")
	fmt.Println("\tqvault -json /Path/To/Config/folders.json 'UPDATE ...'")
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

func main() {
	var reqJSON, optRootNode string
	var optHelp bool
	flag.BoolVar(&optHelp, "help", false, "Help me")
	flag.StringVar(&reqJSON, "json", "", "JSON file with Folder logical hieararchy")
	flag.StringVar(&optRootNode, "root", "Documents", "Root node name for search start")
	flag.Usage = Help
	flag.Parse()

	Version.Copyright(vault.Reverse(vault.CO1), vault.CHR_TRIDENT)
	fmt.Println()

	if optHelp {
		Help()
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		SampleRun(reqJSON, optRootNode)
	} else if flag.NArg() > 1 {
		CustomRun(reqJSON, gatherFreeArguments())
	} else {
		CustomRun(reqJSON, flag.Arg(0))
	}

	fmt.Println()
	Version.BuyMeCoffee(vault.Reverse(vault.CO3))
}
