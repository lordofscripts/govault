/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Represents a directory or Folder in disk with metadata such as
 * tags and security info. It allows for queries, especially with
 * the FQL object.
 *   When the directory tree specifier is read from
 * the JSON file, we build an in-memory tree and each element gets
 * assigned a hiearchical ID using 'Materialized Path'. Renumbering
 * can (and should) be done if elements are added or deleted.
 *   Similarly we recalculate the 'Path' property recursively.
 *-----------------------------------------------------------------*/
package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
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

// Encapsulates directory (folder) information for creating and
// maintaining the document container.
type Folder struct {
	Id        string   `json:"-"` // assigned on the go, not a real ident
	Name      string   `json:"name"`
	Encrypted bool     `json:"encrypted"`
	Tags      string   `json:"tags"` // comma-separated
	Template  string   `json:"template,omitempty"`
	Children  []Folder `json:"children,omitempty"`
	Path      string   `json:"-"`

	nameWidth      int
	showFlags      bool // controls how fmt.Stringer behaves
	showChildCount bool // idem
	showPath       bool // idem
}

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

// Directory metadata that is stored in the .foldermeta.json file
// found in each directory
type FolderMeta struct {
	Encrypted bool     `json:"encrypted,omitempty"` // whether it is supposedly encrypted
	Tags      []string `json:"tags"`                // category tags
	Template  string   `json:"template,omitempty"`  // suggested filename template
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func NewMeta(tags []string, encrypted bool, nameFormat string) FolderMeta {
	return FolderMeta{
		Encrypted: encrypted,
		Tags:      tags,
		Template:  nameFormat,
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// implements fmt.Stringer. The output format can be modified using
// the WithFlag, WithChildCount, WithPath and WithNameWidth constructor
// fluent API.
func (f Folder) String() string {
	if f.nameWidth <= 0 {
		f.nameWidth = 15
	}
	var sb strings.Builder
	if f.showPath && len(f.Path) > 0 {
		fmt.Fprintf(&sb, "%-*s ", f.nameWidth, f.Path)
	} else {
		fmt.Fprintf(&sb, "%-*s ", f.nameWidth, f.Name)
	}

	if f.showChildCount {
		fmt.Fprintf(&sb, "%02d ", len(f.Children))
	}
	if f.showFlags {
		flags := ""
		if f.Encrypted {
			flags += "E"
		}
		if len(f.Tags) > 0 {
			flags += "T"
		}
		if len(f.Children) > 0 {
			flags += "C"
		}
		fmt.Fprintf(&sb, "%3s ", flags)
	}
	fmt.Fprintf(&sb, "tags: %s", f.Tags)
	return sb.String()
}

// Outputs FOLDER_NAME:15 [FLAGS] TAG_LIST
func (f Folder) StringPlus() string {
	flags := ""
	if f.Encrypted {
		flags += "E"
	}
	flags += fmt.Sprintf("%d", len(f.Children))

	return fmt.Sprintf("%-15s [%s] %s", f.Name, flags, strings.ToLower(f.Tags))
}

// Let String() output the folder's children count
func (f *Folder) WithChildCount() *Folder {
	f.showChildCount = true
	return f
}

// Let String() output the folder flags
func (f *Folder) WithFlags() *Folder {
	f.showFlags = true
	return f
}

// Let String() output Path instead of Name
func (f *Folder) WithPath() *Folder {
	f.showPath = true
	return f
}

// Let String() output the Name/Path in this width (defaults to 15)
func (f *Folder) WithNameWidth(length int) *Folder {
	f.nameWidth = length
	return f
}

// Normalize tags, use before or after serialization
func (f *Folder) Normalize() {
	tm := NewTagManager().FromCsvString(f.Tags)
	f.Tags = tm.String()
}

// Generates a tag cloud
func (f *Folder) TagCloud() map[string]int {
	tagCloud := make(map[string]int)
	f.summarizeTags(tagCloud, f.Tags)
	for _, child := range f.Children {
		f.summarizeTags(tagCloud, child.Tags)
	}
	return tagCloud
}

// UpdatePaths recalculates the Path field for the folder and all its children.
func (f *Folder) UpdatePaths(parentPath string) {
	if parentPath == "" {
		f.Path = f.Name
	} else {
		f.Path = parentPath + "/" + f.Name
	}

	for i := range f.Children {
		// Pass the current folder's path as the parent path for children
		f.Children[i].UpdatePaths(f.Path)
	}
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

func (f *Folder) summarizeTags(tagCloud map[string]int, tags string) {
	tm := NewTagManager().FromCsvString(tags)
	for _, tag := range tm.Tags() {
		tagCloud[tag]++
	}
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

//				((( JSON Functions )))

// Serialize folder metadata to filename.
func SerializeFolderMeta(metadata FolderMeta, filename string) error {
	if dataBin, err := json.MarshalIndent(metadata, "", "  "); err == nil {
		if err = os.WriteFile(filename, dataBin, 0644); err != nil {
			return fmt.Errorf("write meta %s: %w", filename, err)
		}
	} else {
		return err
	}

	return nil
}

// Deserialize folder metadata from filename.
func DeserializeFolderMeta(filename string) (*FolderMeta, error) {
	var metadata FolderMeta
	if dataBin, err := os.ReadFile(filename); err == nil {
		if err = json.Unmarshal(dataBin, &metadata); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("read meta %s: %w", filename, err)
	}

	return &metadata, nil
}

// Load a JSON file containing the tree of folders. Upon return the
// derived fields Id and Path are recalculated and populated.
func LoadFolderTable(fromPath string) ([]Folder, error) {
	// (a) Read reference metadata from govault_pics.json OR govault_docs.json
	data, err := os.ReadFile(fromPath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON: %v\n", err)
	}

	// (b) Deserialization into a model
	var folders []Folder
	if err = json.Unmarshal(data, &folders); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v\n", err)
	}
	RenumberFolders(folders)
	RecalculateAllPaths(folders)

	return folders, nil
}

// Writes a []Folder slice to a JSON file.
func SaveFolderTable(folders []Folder, toPath string) error {
	RenumberFolders(folders)
	RecalculateAllPaths(folders)

	if bytes, err := json.MarshalIndent(folders, "", "  "); err == nil {
		err = os.WriteFile(toPath, bytes, 0755)
		return err
	} else {
		return err
	}
}

// 				((( Folder.Id Recalculation )))

// RenumberFolders resets and populates the Id field for a slice of folders
// and all descendants using a 'Materialized Path' method.
// Ref.: https://medium.com/@rishabhdevmanu/from-trees-to-tables-storing-hierarchical-data-in-relational-databases-a5e5e6e1bd64
func RenumberFolders(folders []Folder) {
	assignIDs(folders, "")
}

// Renumbering (internal) assigns a hiearchical ID to a Folder item.
func assignIDs(folders []Folder, prefix string) {
	// @note uses for i := range folders and accesses elements via folders[i].
	// This is critical because Folder.Children is a slice of values ([]Folder),
	// not pointers. Using a value-based for _, f := range would modify a
	// copy of the folder instead of the one in the slice.
	for i := range folders {
		// Generate the current ID based on 1-based indexing
		currentID := strconv.Itoa(i + 1)
		if prefix != "" {
			currentID = prefix + "." + currentID
		}

		// Assign the ID to the current folder
		folders[i].Id = currentID

		// Recursively process children if they exist
		if len(folders[i].Children) > 0 {
			assignIDs(folders[i].Children, currentID)
		}
	}
}

//				((( Folder.Path recalculation )))

// RecalculateAllPaths is a helper to update paths for a slice of root folders.
func RecalculateAllPaths(folders []Folder) {
	for i := range folders {
		folders[i].UpdatePaths("")
	}
}

//				((( Physical to Logical Mapping with Path & Name )))

// FindByPhysicalPath attempts to find a Folder given a physical path and a mapping
// of root Folder names to their base physical directories.
func FindByPhysicalPath(folders []Folder, rootMappings map[string]string, physicalPath string) *Folder {
	// 1. Identify which root the physical path belongs to and get the logical path
	var targetLogicalPath string

	for rootName, baseDir := range rootMappings {
		if strings.HasPrefix(physicalPath, baseDir) {
			// Example:
			// physicalPath: /home/dev/test/Documents/Legal
			// baseDir:      /home/dev/test/Documents
			// result:       Documents/Legal

			// Extract the part relative to the parent of the baseDir
			// We find where the root name starts in the physical path
			idx := strings.LastIndex(baseDir, rootName)
			if idx != -1 {
				targetLogicalPath = physicalPath[idx:]
				break
			}
		}
	}

	if targetLogicalPath == "" {
		return nil
	}

	// 2. Search the tree for the calculated logical path
	return SearchFolders(folders, targetLogicalPath)
}

// SearchFolders is a recursive helper that finds a folder by its internal Path member.
func SearchFolders(folders []Folder, targetPath string) *Folder {
	for i := range folders {
		if folders[i].Path == targetPath {
			return &folders[i]
		}
		// If the target path starts with this folder's path, look deeper
		if strings.HasPrefix(targetPath, folders[i].Path+"/") {
			if found := SearchFolders(folders[i].Children, targetPath); found != nil {
				return found
			}
		}
	}
	return nil
}

//				((( Search using a LIKE function )))

// SearchLike finds all folders where either Id or Path matches the SQL-like pattern.
// The search automatically targets Id, Path and Name.
func SearchLike(folders []Folder, pattern string) []*Folder {
	// 1. Convert SQL LIKE syntax to Regex
	// Escape special regex chars like +,(or [, then swap % with .* and _ with .
	regexPattern := regexp.QuoteMeta(pattern)

	regexPattern = strings.ReplaceAll(regexPattern, "%", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "_", ".")
	regexPattern = "^" + regexPattern + "$"

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil
	}

	// Pointer slice so that caller can modify the actual tree elements
	var results []*Folder
	recursiveSearchLike(folders, re, &results)
	return results
}

// a global recursive search regardless of depth in the logical tree.
// The recursive search automatically targets Id, Path and Name fields.
func recursiveSearchLike(folders []Folder, re *regexp.Regexp, results *[]*Folder) {
	for i := range folders {
		// Check both Id and Path
		if re.MatchString(folders[i].Id) ||
			re.MatchString(folders[i].Path) ||
			re.MatchString(folders[i].Name) {
			*results = append(*results, &folders[i])
		}

		if len(folders[i].Children) > 0 {
			recursiveSearchLike(folders[i].Children, re, results)
		}
	}
}

//				((( Search by Predicate function )))

// SearchCustom performs a recursive search using a caller-provided predicate function.
func SearchCustom(folders []Folder, predicate func(Folder) bool) []*Folder {
	var matches []*Folder

	// Define the recursive worker
	var walk func([]Folder)
	walk = func(currentLevel []Folder) {
		for i := range currentLevel {
			// Run the custom comparison
			if predicate(currentLevel[i]) {
				matches = append(matches, &currentLevel[i])
			}

			// Recurse into children
			if len(currentLevel[i].Children) > 0 {
				walk(currentLevel[i].Children)
			}
		}
	}

	walk(folders)
	return matches
}

//				((( Utility )))

// Dump a []Folder as an indented hiearchical tree
func DumpVirtual(virtualFS []Folder, showLevel bool) {
	for _, f := range virtualFS {
		prefixRepeat := strings.Count(f.Id, ".")
		prefix := strings.Repeat("    ", prefixRepeat)
		if showLevel {
			fmt.Printf("%03d %s%s\n", prefixRepeat, prefix, path.Base(f.Path))
		} else {
			fmt.Println(prefix, path.Base(f.Path))
		}
		if len(f.Children) > 0 {
			DumpVirtual(f.Children, showLevel)
		}
	}
}
