/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"lordofscripts/vault"
	vfs "lordofscripts/vault/internal/fs"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const META_FILE_MODE os.FileMode = 0644

/* ----------------------------------------------------------------
 *                    A U X   F U N C T I O N S
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

// Recursive creation of vault's directory/folder structure
func createFolderHierarchy(root string, f vfs.Folder, dirMode fs.FileMode, overwrite, update bool) error {
	dirPath := filepath.Join(root, f.Name)
	if err := os.MkdirAll(dirPath, dirMode); err != nil {
		return fmt.Errorf("mkdir %s: %w", dirPath, err)
	}
	// prepare metadata
	tm := vfs.NewTagManager().FromCsvString(f.Tags) // normalize and remove duplicates
	m := vfs.Meta{Encrypted: f.Encrypted, Tags: tm.Tags()}
	metaPath := filepath.Join(dirPath, ".foldermeta.json")

	if _, err := os.Stat(metaPath); err == nil && !overwrite {
		fmt.Printf("exists: %s (metadata exists, use -overwrite to replace)\n", dirPath)
	} else {
		j, _ := json.MarshalIndent(m, "", "  ")
		if err := os.WriteFile(metaPath, j, META_FILE_MODE); err != nil {
			return fmt.Errorf("write meta %s: %w", metaPath, err)
		}
		label := "created"
		if update {
			label = "updated"
		}
		fmt.Printf("%s: %s\n", label, dirPath)
	}

	// recurse into children
	for _, c := range f.Children {
		child := c
		// if child.Name is relative (no slashes) join automatically
		// otherwise treat as nested path
		if err := createFolderHierarchy(dirPath, child, dirMode, overwrite, update); err != nil {
			return err
		}
	}
	return nil
}

func parseFileMode(s string) (fs.FileMode, error) {
	// expects octal like "0755"
	if len(s) == 0 {
		return 0755, nil
	}
	// strconv.ParseUint
	val, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return 0, err
	}
	return fs.FileMode(val), nil
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Help command (-help)
func Help() {
	flag.PrintDefaults()
	fmt.Println()
	vault.BuyMeCoffee()
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

// Create folder hieararchy command (-create)
func Create(folders []vfs.Folder, permStr, rootDir string, overwrite bool) {
	mode, err := parseFileMode(permStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid perm: %v\n", err)
		os.Exit(1)
	}

	for _, f := range folders {
		if err := createFolderHierarchy(rootDir, f, mode, overwrite, false); err != nil {
			fmt.Fprintf(os.Stderr, "error creating folders: %v\n", err)
		}
	}
}

// Create folder hieararchy command (-create)
func Update(folders []vfs.Folder, permStr, rootDir string) {
	const OVERWRITE bool = true
	mode, err := parseFileMode(permStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid perm: %v\n", err)
		os.Exit(1)
	}

	for _, f := range folders {
		if err := createFolderHierarchy(rootDir, f, mode, OVERWRITE, true); err != nil {
			fmt.Fprintf(os.Stderr, "error updating folders: %v\n", err)
		}
	}
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

func main() {
	jsonPath := flag.String("json", "tree.json", "path to JSON folder tree")
	root := flag.String("root", ".", "root path where folders will be created")
	overwrite := flag.Bool("overwrite", false, "overwrite existing folder metadata")
	perm := flag.String("perm", "0755", "directory permissions (octal)")
	cmdUpdate := flag.Bool("update", false, "Update folder hierarchy from JSON file")
	cmdCreate := flag.Bool("create", false, "Create directory structure")
	cmdQuery := flag.String("query", "", "Tags to query to know which folders apply")
	cmdListTags := flag.Bool("tags", false, "Show tag cloud")
	cmdHelp := flag.Bool("help", false, "Help with options")
	flag.Parse()

	vault.Copyright(vault.CO1, true)

	if !*cmdHelp && !*cmdCreate && !*cmdUpdate && !*cmdListTags && len(*cmdQuery) == 0 {
		fmt.Println("Please specify either of -help -create -update -tags or -query TAGS")
		Help()
		os.Exit(1)
	}

	if *cmdHelp {
		Help()
		os.Exit(0)
	}

	if *cmdCreate && *cmdUpdate {
		fmt.Println("-update and -create are mutually exclusive")
		Help()
		os.Exit(2)
	}

	if (len(*cmdQuery) != 0) && *cmdListTags {
		fmt.Println("-query and -tags are mutually exclusive")
		Help()
		os.Exit(3)
	}

	// for Create or Query we need the JSON file with metadata
	data, err := os.ReadFile(*jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading JSON: %v\n", err)
		os.Exit(4)
	}

	var folders []vfs.Folder
	if err := json.Unmarshal(data, &folders); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing JSON: %v\n", err)
		os.Exit(5)
	}

	if *cmdCreate {
		Create(folders, *perm, *root, *overwrite)
	}
	if *cmdUpdate {
		Update(folders, *perm, *root)
	}
	if len(*cmdQuery) != 0 {
		Query(folders, *cmdQuery)
	}
	if *cmdListTags {
		ListTags(folders, true)
	}

	vault.BuyMeCoffee()
}
