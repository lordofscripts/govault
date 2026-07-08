/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           goVault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Application Folder handling.
 *-----------------------------------------------------------------*/
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	vault "github.com/lordofscripts/govault"
	vfs "github.com/lordofscripts/govault/internal/fs"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	FOLDER_META_FILENAME string = ".foldermeta.json"
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Show help about the 'create' command
func HelpCreate() {
	fmt.Println(os.Args[0], " create [-perm OCTAL]")
	fmt.Println("Option:")
	fmt.Println("\t-root PATH	Directories will be created at PATH")
	fmt.Println("\t-json PATH	Path to the JSON file with Virtual structure")
	fmt.Println("\t-perm OCTAL  Directory permissions in OCTAL")
	fmt.Println("\t-overwrite	Overwrite existing .foldermeta.json files")

	fmt.Println()
	vault.ModuleVersion.BuyMeCoffee(vault.Reverse(vault.CO3))
}

// Show help about the 'update' command
func HelpUpdate() {
	fmt.Println(os.Args[0], " update [-perm OCTAL]")
	fmt.Println("Option:")
	fmt.Println("\t-root PATH	Directories will be created at PATH")
	fmt.Println("\t-json PATH	Path to the JSON file with Virtual structure")
	fmt.Println("\t-perm OCTAL  Directory permissions in OCTAL")

	fmt.Println()
	vault.ModuleVersion.BuyMeCoffee(vault.Reverse(vault.CO3))
}

// Create folder hieararchy command (-create)
func Create(folders []vfs.Folder, permStr, rootDir string, overwrite bool) error {
	mode, err := parseFileMode(permStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid perm: %v\n", err)
	} else {
		for _, f := range folders {
			if err := createFolderHierarchy(rootDir, f, mode, overwrite, false); err != nil {
				fmt.Fprintf(os.Stderr, "error creating folders: %v\n", err)
			}
		}
	}

	return err
}

// Create folder hierarchy command (-update)
func Update(folders []vfs.Folder, permStr, rootDir string) error {
	const OVERWRITE bool = true
	mode, err := parseFileMode(permStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid perm: %v\n", err)
	} else {
		for _, f := range folders {
			if err := createFolderHierarchy(rootDir, f, mode, OVERWRITE, true); err != nil {
				fmt.Fprintf(os.Stderr, "error updating folders: %v\n", err)
			}
		}
	}

	return err
}

/* ----------------------------------------------------------------
 *                       I N T E R N A L
 *-----------------------------------------------------------------*/

// Recursive creation of vault's directory/folder structure
func createFolderHierarchy(root string, f vfs.Folder, dirMode fs.FileMode, overwrite, update bool) error {
	dirPath := filepath.Join(root, f.Name)
	if err := os.MkdirAll(dirPath, dirMode); err != nil {
		return fmt.Errorf("mkdir %s: %w", dirPath, err)
	}
	// prepare metadata
	tm := vfs.NewTagManager().FromCsvString(f.Tags) // normalize and remove duplicates
	m := vfs.NewMeta(tm.Tags(), f.Encrypted, f.Template)
	metaPath := filepath.Join(dirPath, FOLDER_META_FILENAME)

	if _, err := os.Stat(metaPath); err == nil && !overwrite {
		fmt.Fprintf(os.Stderr, "exists: %s (metadata exists, use -overwrite to replace)\n", dirPath)
	} else {
		if err := vfs.SerializeFolderMeta(m, metaPath); err != nil {
			return err
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
