/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package sync

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/lordofscripts/goapp/app"
	"github.com/lordofscripts/govault/internal/fs"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	_SKIP_HIDDEN bool = true
)

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

// Compares to directory (folder) file slices, the left one (virtual)
// comes from the JSON configuration file, and the right one (actual)
// is built from the real filesystem hieararchy in rootDir. It generates
// a slice of SyncResult indicating which synchronization actions should
// take place for both to be equal.
func CompareFolders(virtual []fs.Folder, rootDir string, updateAction UpdateAction) []SyncResult {
	var results []SyncResult

	// (a) Create slice with names of top-level virtual folders (filter criteria)
	allowed := make([]string, 0)
	for _, name := range virtual {
		allowed = append(allowed, name.Name)
	}

	// (b) Scan physical directory at rootDir skipping hidden files/dirs and
	//	   paths that do not begin with any in allowed
	//actual := scanDiskFiltered(rootDir, "", _SKIP_HIDDEN, allowed)
	actual := scanDisk(rootDir, "", _SKIP_HIDDEN)

	// (b) Filter out first-level physical directories not present in top-level virtual

	var filtered []fs.Folder
	for _, actualF := range actual {
		for _, prefix := range allowed {
			prefix := strings.Trim(prefix, "/")
			if strings.HasPrefix(actualF.Name, prefix) {
				filtered = append(filtered, actualF)
			}
		}
	}
	actual = nil

	// (c) Cross-compare both trees
	compareTrees(virtual, filtered, rootDir, nil, updateAction, &results)

	return results
}

// modifies the in-memory slice (virtual filesystem) based on the
// results of CompareFolders.
func SynchronizeVirtualFolders(results []SyncResult, virtual *[]fs.Folder) error {
	for _, res := range results {
		switch res.Action {
		case AddLeft:
			if res.Parent == nil {
				*virtual = append(*virtual, res.Target)
			} else {
				res.Parent.Children = append(res.Parent.Children, res.Target)
			}
		case DeleteLeft:
			if res.Parent == nil {
				*virtual = removeFolder(*virtual, res.Target.Name, res.Target.Path) // @audit
			} else {
				res.Parent.Children = removeFolder(res.Parent.Children, res.Target.Name, res.Target.Path) // @audit
			}
		}
	}
	return nil
}

func removeFolder(folders []fs.Folder, name, dirPath string) []fs.Folder {
	for i, f := range folders {
		if f.Name == name && f.Path == dirPath /*f.Name == name*/ {
			return append(folders[:i], folders[i+1:]...)
		}
	}
	return folders
}

// used by CompareFolders to scan the real filesystem and build-up
// a Folder slice that represents the actual filesystem.
// BROKEN
func scanDiskFiltered(currentPath string, logicalPath string, skipHidden bool, allowedPrefixes []string) []fs.Folder {
	var folders []fs.Folder
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		subPath := filepath.Join(currentPath, name)

		// 1. Hidden Filter
		if skipHidden {
			// Using the app helper as in your context
			if hidden, err := app.IsHiddenFile(subPath); err == nil && hidden {
				continue
			}
		}

		// 2. Build Logical Path (Always use forward slashes for consistency with allowedPrefixes)
		newLogical := name
		if logicalPath != "" {
			newLogical = logicalPath + "/" + name
		}

		// 3. Logic: Determine if we should enter or keep this folder
		isAllowed := false   // We are at or below the prefix
		shouldEnter := false // We are above the prefix and need to keep looking

		if len(allowedPrefixes) == 0 {
			isAllowed = true
		} else {
			for _, prefix := range allowedPrefixes {
				// Clean prefix of any leading/trailing slashes for comparison
				cleanPrefix := strings.Trim(prefix, "/")

				// Case A: Exact match or a child of the prefix
				// e.g., newLogical: "A/B/C", prefix: "A/B"
				if newLogical == cleanPrefix || strings.HasPrefix(newLogical, cleanPrefix+"/") {
					isAllowed = true
					break
				}

				// Case B: Ancestor of the prefix (Need to go deeper)
				// e.g., newLogical: "A", prefix: "A/B/C"
				if strings.HasPrefix(cleanPrefix, newLogical+"/") {
					shouldEnter = true
					break
				}
			}
		}

		// 4. Recurse and Prune
		if isAllowed || shouldEnter {
			f := fs.Folder{
				Name: name,
				Path: newLogical,
			}

			// Recurse to find children
			f.Children = scanDiskFiltered(subPath, newLogical, skipHidden, allowedPrefixes)

			// Add to results if:
			// - This folder itself is allowed
			// - OR it contains allowed children found during recursion
			if isAllowed || len(f.Children) > 0 {
				folders = append(folders, f)
			}
		}
	}
	// (DEGT) Populate derived properties.
	// NOTE 1: The Id (Renumber) will differ between virtual and real if
	//		   they are not the same directory structure! Thus, cannot use
	//		   Id for validating position.
	// NOTE 2: The Path is calculated from Name so it is contextualized
	// 	 	   to the actual hiearchy.
	fs.RenumberFolders(folders)
	fs.RecalculateAllPaths(folders) // @audit may rewrite logicalPath from above!

	return folders
}

// used by CompareFolders to scan the real filesystem and build-up
// a Folder slice that represents the actual filesystem.
func scanDisk(currentPath string, logicalPath string, skipHidden bool) []fs.Folder {
	var folders []fs.Folder
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()

			subPath := filepath.Join(currentPath, name)
			if skipHidden {
				// additional check for hidden directories
				hidden := false
				if hidden, err = app.IsHiddenFile(subPath); err == nil && hidden {
					continue // skip hidden directory
				} else if err != nil {
					println("hide check failed", err.Error(), subPath)
				}
			}

			newLogical := name
			if logicalPath != "" {
				newLogical = logicalPath + "/" + name
			}

			f := fs.Folder{
				Name: name,
				Path: newLogical,
			}
			f.Children = scanDisk(subPath, newLogical, _SKIP_HIDDEN)
			folders = append(folders, f)
		}
	}

	// (DEGT) Populate derived properties.
	// NOTE 1: The Id (Renumber) will differ between virtual and real if
	//		   they are not the same directory structure! Thus, cannot use
	//		   Id for validating position.
	// NOTE 2: The Path is calculated from Name so it is contextualized
	// 	 	   to the actual hiearchy.
	fs.RenumberFolders(folders)
	fs.RecalculateAllPaths(folders) // @audit may rewrite logicalPath from above!

	return folders
}

func compareTrees(virtual []fs.Folder, actual []fs.Folder, physicalRoot string, vParent *fs.Folder, updateAction UpdateAction, results *[]SyncResult) {
	// @note we track Path rather than Name because Name can be the same on sub-items.
	// However, for that the RecalculatePath has to have been called earlier.
	vMap := make(map[string]*fs.Folder)
	for i := range virtual {
		//vMap[virtual[i].Name] = &virtual[i]
		vMap[virtual[i].Path] = &virtual[i]
	}

	aMap := make(map[string]*fs.Folder)
	for i := range actual {
		//aMap[actual[i].Name] = &actual[i]
		aMap[actual[i].Path] = &actual[i]
	}

	doUpdateLeft := func(apMap map[string]*fs.Folder, physRoot string, updResults *[]SyncResult) {
		for name, aFolder := range apMap {
			physPath := filepath.Join(physRoot, aFolder.Name)
			if vFolder, exists := vMap[name]; !exists {
				// Folder is on Disk but not in Virtual -> add to Virtual
				newResult := SyncResult{
					Action:       AddLeft,
					Target:       *aFolder,
					Parent:       vParent,
					PhysicalPath: physPath}
				*updResults = append(*updResults, newResult)
			} else {
				// Recurse into common folders that exist on both. The parent reference
				// must be from the Virtual side.
				compareTrees(vFolder.Children, aFolder.Children, physPath, vFolder, updateAction, updResults)
			}
		}
	}

	doUpdateRight := func(vpMap map[string]*fs.Folder, physRoot string, updResults *[]SyncResult) {
		for name, vFolder := range vpMap {
			physPath := filepath.Join(physRoot, vFolder.Name)
			if aFolder, exists := aMap[name]; !exists {
				// add it to Virtual
				*updResults = append(*results, SyncResult{
					Action:       AddRight,
					Target:       *vFolder,
					Parent:       nil,
					PhysicalPath: physPath})
			} else {
				// Recurse into common folders that exist on both
				compareTrees(vFolder.Children, aFolder.Children, physPath, vFolder, updateAction, updResults)
			}
		}
	}

	// 1. Check Virtual vs Actual (Identify what to Delete from Virtual or Add to Disk)
	//for name, vFolder := range vMap {
	//physPath := filepath.Join(physicalRoot, vFolder.Name)
	switch updateAction {
	case UpdateLeft: // make left (Virtual) structure match right (Physical/Disk)
		doUpdateLeft(aMap, physicalRoot, results)
		// Left ends up with everything from Right and PERHAPS extra stuff that
		// was present on Left but not on right, however, we are UPDATING left.

	case UpdateRight: // make right (Physical) structure match left (Virtual)
		doUpdateRight(vMap, physicalRoot, results)
		// Right ends up with everything from Left and PERHAPS extra stuff that
		// was present on Right but not on left, however, we are UPDATING left.

	case FullSync: // make both sides look the same
		panic("Not Implemented")
	}
	/*
		if aFolder, exists := aMap[name]; !exists {
			*results = append(*results, SyncResult{Action: AddRight, Target: *vFolder, PhysicalPath: physPath})
			*results = append(*results, SyncResult{Action: DeleteLeft, Target: *vFolder, Parent: vParent})
		} else {
			// Recurse into common folders
			compareTrees(vFolder.Children, aFolder.Children, physPath, vFolder, results)
		}
	}*/

	// 2. Check Actual vs Virtual (Identify what to Add to Virtual or Delete from Disk)
	/*
		for name, aFolder := range aMap {
			if _, exists := vMap[name]; !exists {
				physPath := filepath.Join(physicalRoot, aFolder.Name)
				*results = append(*results, SyncResult{Action: AddLeft, Target: *aFolder, Parent: vParent})
				*results = append(*results, SyncResult{Action: DeleteRight, Target: *aFolder, PhysicalPath: physPath})
			}
		}
	*/
}
