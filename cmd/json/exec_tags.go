/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Application Tag handling for Folder Metadata.
 *-----------------------------------------------------------------*/
package main

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	vault "github.com/lordofscripts/govault"
	vfs "github.com/lordofscripts/govault/internal/fs"
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Show help about the 'tags' command
func HelpTags() {
	fmt.Println(os.Args[0], " tags -json PATH {Option}")
	fmt.Println("Option:")
	fmt.Println("\t-list                  List tag cloud")
	fmt.Println("\t-query 'tag1,tag2...'  Query folders matching tag(s)")

	fmt.Println()
	Version.BuyMeCoffee(vault.Reverse(vault.CO3))
}

// List Tag Cloud command (-tags)
func ListTags(folders []vfs.Folder, byColumns bool) {
	// generate Tag Cloud recursively
	tagCloud := make(map[string]int)

	for _, f := range folders {
		fmap := f.TagCloud()
		for k, v := range fmap {
			tagCloud[k] += v
		}
	}

	// output tagName:  count  in one column.
	outputFormatA := func() {
		for _, k := range slices.Sorted(maps.Keys(tagCloud)) {
			fmt.Printf("\t%-15s: %d\n", k, tagCloud[k])
		}
	}

	// output tagName(count) in 3 columns
	outputFormatB := func() {
		const MAX_COL = 3
		column := 0
		for _, k := range slices.Sorted(maps.Keys(tagCloud)) {
			switch column % MAX_COL {
			case 0:
				fmt.Printf("%20s (%d)", k, tagCloud[k])
			case 1:
				fmt.Printf("%20s (%d)", k, tagCloud[k])
			case 2:
				fmt.Printf("%20s (%d)\n", k, tagCloud[k])
			}
			column++
		}
		if column%MAX_COL != 2 {
			fmt.Println()
		}
	}

	fmt.Println("Tag Cloud")
	if byColumns {
		outputFormatB()
	} else {
		outputFormatA()
	}
}

// Query vault hierarchy for tags (-query "tag1,tag2,...")
// It is useful when you want to store a new item and want
// to query the structure for the best place to store. It
// it advisable to try -tags first to know what to query.
func Query(folders []vfs.Folder, tags string) {
	fmt.Printf("Query Tags: %s\n", tags)
	tm := vfs.NewTagManager().FromCsvString(tags)
	tagsQ := tm.Tags() // query Tags

	pathifyChildren := func(folder vfs.Folder, parent string) vfs.Folder {
		children := make([]vfs.Folder, 0)
		// recurse into children
		for _, child := range folder.Children {
			tm := vfs.NewTagManager().FromCsvString(folder.Tags)
			tm.JoinFromString(child.Tags) // inherit tags from parent
			dirPath := filepath.Join(parent, child.Name)
			children = append(children, vfs.Folder{
				Name:      child.Name,
				Encrypted: child.Encrypted,
				Tags:      tm.String(),
				Path:      dirPath,
			})
		}
		folder.Children = children
		return folder
	}

	// calculate full paths in hierarchy of folders
	var foldersOut []vfs.Folder = make([]vfs.Folder, 0)
	for _, f := range folders {
		folderNew := pathifyChildren(f, f.Name)
		foldersOut = append(foldersOut, folderNew)
	}

	for _, f := range foldersOut {
		inspect(f, tagsQ)
	}
}

/* ----------------------------------------------------------------
 *                       I N T E R N A L
 *-----------------------------------------------------------------*/

// hierarchical inspection of a folder
func inspect(folder vfs.Folder, tagsQ []string) {
	//log.Printf("%s", folder.StringPlus())
	tm := vfs.NewTagManager().FromCsvString(folder.Tags)
	if ok, _, _ := tm.HasTags(tagsQ); ok {
		/*secure := ""
		if folder.Encrypted {
			secure = "(secured)"
		}*/

		folder.WithChildCount().WithFlags().WithNameWidth(20).WithPath()
		fmt.Printf("\t%s\n", folder)
		//fmt.Printf("Folder: %s %s %d tags: %s\n", folder.Path, secure, count, list)
	}

	for _, child := range folder.Children {
		inspect(child, tagsQ)
	}
}
