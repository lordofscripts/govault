/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           go-Vault
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package fs

import (
	"slices"
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

type TagManager struct {
	tags []string
}

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func NewTagManager() *TagManager {
	return &TagManager{tags: nil}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// Sets the tag slice to tags after normalizing them and removing duplicates
func (tm *TagManager) FromCsvString(tags string) *TagManager {
	tm.tags = tm.normalizeTags(tags)
	slices.Sort(tm.tags)
	tm.tags = slices.Compact(tm.tags)
	return tm
}

// Sets the tag slice to tags after normalizing them and removing duplicates
func (tm *TagManager) FromSlice(tags []string) *TagManager {
	flat := strings.Join(tags, ",")
	tm.tags = tm.normalizeTags(flat)
	slices.Sort(tm.tags)
	tm.tags = slices.Compact(tm.tags)
	return tm
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Implements fmt.Stringer and returns a comma-separated list of tags
func (tm *TagManager) String() string {
	return strings.Join(tm.tags, ",")
}

// Clears all managed tags
func (tm *TagManager) Clear() {
	tm.tags = make([]string, 0)
}

// Check if tag exists
func (tm *TagManager) HasTag(tag string) bool {
	ok, _, _ := tm.HasTags([]string{tag})
	return ok
}

// how many tags from Search are found in Set
// Returns:
//   - (bool) true if at least one tag from Search in Set
//   - (int) How many tags from Search found in Set
//   - (string) Comma-separated list of tags found
func (tm *TagManager) HasTags(search []string) (bool, int, string) {
	found := 0
	tagsFound := make([]string, 0)
	for _, tag := range search {
		if slices.Contains(tm.tags, tag) {
			found++
			tagsFound = append(tagsFound, tag)
		}
	}
	return (found > 0), found, strings.Join(tagsFound, ",")
}

// Join tags to the current managed tags, normalize and remove duplicates
func (tm *TagManager) Join(tags []string) {
	// join them, it may contain repeated elements
	tagsNew := append(tm.tags, tags...)
	// normalize
	tagsNew = tm.normalizeTagSlice(tagsNew)
	// sort them to put repeated elements consecutively
	slices.Sort(tagsNew)
	// remove consecutive repeated elements
	tm.tags = slices.Compact(tagsNew)
}

func (tm *TagManager) JoinFromString(tags string) {
	parts := strings.Split(tags, ",")
	tm.Join(parts)
}

// get all managed tags (normalized and sorted)
func (tm *TagManager) Tags() []string {
	return tm.tags
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

func (tm *TagManager) normalizeTags(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []string
	seen := map[string]bool{}
	for _, p := range parts {
		t := strings.ToLower(strings.TrimSpace(p))
		if t == "" {
			continue
		}
		if !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}

func (tm *TagManager) normalizeTagSlice(tags []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, p := range tags {
		t := strings.ToLower(strings.TrimSpace(p))
		if t == "" {
			continue
		}
		if !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                          T E S T S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/
