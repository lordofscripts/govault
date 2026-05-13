/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package fs

import (
	"fmt"
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
	Name      string   `json:"name"`
	Encrypted bool     `json:"encrypted"`
	Tags      string   `json:"tags"` // comma-separated
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
type Meta struct {
	Encrypted bool     `json:"encrypted"`
	Tags      []string `json:"tags"`
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

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
