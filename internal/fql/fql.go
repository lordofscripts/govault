/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Simple Folder Query Language resembling SQL or LINQ.
 * @audit Perhaps consider using https://github.com/ahmetb/go-linq
 * Alternatively, if the query is in the form of a string, use the
 * FQLX object instead of this.
 *-----------------------------------------------------------------*/
package fql

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/lordofscripts/govault/internal/fs"
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

type FQL struct {
	target []fs.Folder
	//result map[FolderField]fqlResult
	selector []FolderField
	result   []FqlResult
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

// (ctor) creates an instance of a Folder Query Language object.
// The caller can use the object to perform some SQL-like queries.
func NewFolderQuery() *FQL {
	return &FQL{
		target:   nil,
		selector: make([]FolderField, 0),
		result:   make([]FqlResult, 0),
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Select the fields we want to extract
func (fql *FQL) Select(fields ...FolderField) *FQL {
	fql.result = make([]FqlResult, 0)
	for _, field := range fields {
		if field == FolderAll {
			fql.selector = []FolderField{FolderAll}
		} else {
			if fql.isValidField(field) {
				fql.selector = append(fql.selector, field)
			}
		}
	}

	return fql
}

// Specify the target of the structured query
func (fql *FQL) From(folders []fs.Folder) *FQL {
	fql.target = folders
	return fql
}

func (fql *FQL) WhereFlex(args ...any) *FQL {
	if len(args) == 3 && IsType[FolderField](args[0]) &&
		IsType[Comparer](args[1]) &&
		IsType[string](args[2]) {
		field, _ := (args[0]).(FolderField)
		op, _ := (args[1]).(Comparer)
		valueStr, _ := (args[2]).(string)
		return fql.Where(field, op, valueStr)
	}
	if len(args) == 1 && IsType[FolderComparer](args[1]) {
		op := (args[0]).(FolderComparer)
		return fql.WhereX(op)
	}
	if len(args) == 2 && IsType[LikeAnyComparer](args[0]) &&
		IsType[string](args[1]) {
		op := (args[0]).(LikeAnyComparer)
		pattern := (args[1]).(string)
		subset := op(fql.target, pattern)
		for _, itemF := range subset {
			result := fql.newResult(itemF)
			fql.result = append(fql.result, *result)
		}
	}

	return fql
}

func (fql *FQL) Where(field FolderField, operatorCB Comparer, value string) *FQL {
	if !fql.isValidField(field) {
		fmt.Fprintf(os.Stderr, "invalid selector #'%d' in Where()\n", field)
		return fql
	}
	if field == FolderAll {
		fmt.Fprintf(os.Stderr, "cannot use '*' in Where()\n")
		return fql
	}

	qualified := fql.seekTable(fql.target, field, operatorCB, value)
	for _, itemF := range qualified {
		result := fql.newResult(itemF)
		fql.result = append(fql.result, *result)
	}

	return fql
}

func (fql *FQL) WhereX(customCompare FolderComparer) *FQL {
	if customCompare == nil {
		return fql
	}

	matches := fs.SearchCustom(fql.target, customCompare)
	for _, itemF := range matches {
		result := fql.newResult(itemF)
		fql.result = append(fql.result, *result)
	}

	return fql
}

// Uses the LIKE operator on any of Folder.Id, Folder.Path or
// Folder.Name
func (fql *FQL) WhereLike(pattern string) *FQL {
	some := fs.SearchLike(fql.target, pattern)
	for _, itemF := range some {
		result := fql.newResult(itemF) // @audit make Result a Pointer to Folder
		fql.result = append(fql.result, *result)
	}

	return fql
}

func (fql *FQL) String() string {
	var sb strings.Builder
	sb.WriteString("SELECT ")
	sb.WriteString("[" + fql.Selected() + "] ")
	sb.WriteString("FROM Folders ")
	//sb.WriteString("WHERE ")
	return sb.String()
}

// End of a query. Returns the number of rows in the result.
func (fql *FQL) Count() int {
	return len(fql.result)
}

// End of a query. Returns the number of Select fields
func (fql *FQL) FieldCount() int {
	return len(fql.selector)
}

// return the total number of hierarchical elements in the table
func (fql *FQL) RowCount() int {
	// tree traversal
	var traverseFunc func(vCol []fs.Folder) int
	traverseFunc = func(vCol []fs.Folder) int {
		localCount := 0
		for _, cwf := range vCol {
			localCount++
			if len(cwf.Children) != 0 {
				moreCnt := traverseFunc(cwf.Children)
				localCount += moreCnt
			}
		}

		return localCount
	}

	return traverseFunc(fql.target)
}

// End of a query. Returns the names of fields chosen for Query
func (fql *FQL) Selected() string {
	if len(fql.selector) == 1 && fql.selector[0] == FolderAll {
		return ALL_COLUMNS
	} else {
		names := make([]string, 0)
		for _, field := range fql.selector {
			names = append(names, field.String())
		}
		sort.Strings(names)
		return strings.Join(names, ",")
	}
}

// Fetch all results from the previous query
func (fql *FQL) Fetch() []FqlResult {
	all := make([]FqlResult, 0)
	all = append(all, fql.result...)
	return all
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// implements FolderFieldPicker
// given the field handle, it retrieves its value
func (fql *FQL) getField(fld FolderField, obj fs.Folder) any {
	var targetValue any // @audit getField func but it ignores *
	switch fld {
	case FolderId:
		targetValue = obj.Id
	case FolderName:
		targetValue = obj.Name
	case FolderTags:
		targetValue = obj.Tags
	case FolderTemplate:
		targetValue = obj.Template
	case FolderPath:
		targetValue = obj.Path
	case FolderEncrypted:
		targetValue = fmt.Sprintf("%t", obj.Encrypted)
	case FolderAll:
		targetValue = obj
	}
	return targetValue
}

func (fql *FQL) newResult(f *fs.Folder) *FqlResult {
	var r *FqlResult = nil
	if len(fql.selector) > 0 {
		r = newFQLResult()
		if slices.Contains(fql.selector, FolderAll) {
			r.Whole = f
		} else {
			for _, fld := range fql.selector { // All already processed, next is always string
				r.Fields[fld] = (getFieldFrom(fld, *f)).(string) // @audit revise in the future
			}
		}
	}

	return r
}

// Convert the hieararchical tree into a flattened version without children
func (fql *FQL) flattenTable(vCol []fs.Folder) []fs.Folder {
	flattened := make([]fs.Folder, 0)
	branch := make([]fs.Folder, 0)

	// start with the top branches on the trunk
	branch = append(branch, vCol...)

	for len(branch) > 0 {
		// trim
		leaf := branch[0]
		// harvest
		flattened = append(flattened, getShallowFolder(leaf))
		// eliminate the head which we just processed
		branch = branch[1:]
		// to be pruned, put on the (shortened) queue
		branch = append(branch, leaf.Children...)
	}
	return flattened
}

// search the tree/table for an item/Folder with a field containing that value
func (fql *FQL) seekTable(vCol []fs.Folder, fld FolderField, operatorCB Comparer, value string) []*fs.Folder {
	_, items := seekTableInternal(vCol, fld, fql.getField, operatorCB, value, true)
	return items
}

func (fql *FQL) isValidFieldName(name string) bool {
	return strings.EqualFold(name, "id") || // @audit update when fs.Folder changed!
		strings.EqualFold(name, "name") ||
		strings.EqualFold(name, "tag") ||
		strings.EqualFold(name, "template") ||
		strings.EqualFold(name, "path") ||
		strings.EqualFold(name, "encrypted") ||
		name == ALL_COLUMNS
}

func (fql *FQL) isValidField(name FolderField) bool {
	return name == FolderId || // @audit update when fs.Folder changed!
		name == FolderName ||
		name == FolderTags ||
		name == FolderTemplate ||
		name == FolderPath ||
		name == FolderEncrypted ||
		name == FolderAll
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

func seekTableInternal(vCol []fs.Folder, fld FolderField, getterCB FieldGetter, operatorCB Comparer, value string, deep bool) (count int, chosen []*fs.Folder) {
	count = 0
	chosen = make([]*fs.Folder, 0)

	for _, cwf := range vCol {
		thisVal := (getterCB(fld, cwf)).(string)
		if operatorCB(thisVal, value) { // this element matches
			count++
			chosen = append(chosen, &cwf) // add to current results
		}
		// do the children if depth search
		// now seek its children as well
		if len(cwf.Children) != 0 && deep {
			moreCnt, moreFolders := seekTableInternal(cwf.Children, fld, getterCB, operatorCB, value, deep)
			count += moreCnt
			if len(moreFolders) > 0 {
				chosen = append(chosen, moreFolders...)
			}
		}
	}

	return count, chosen
}

// a shallow copy of f that does NOT include its children.
func getShallowFolder(f fs.Folder) fs.Folder {
	return fs.Folder{
		Id:        f.Id,
		Name:      f.Name,
		Encrypted: f.Encrypted,
		Tags:      f.Tags,
		Template:  f.Template,
		Path:      f.Path,
		Children:  nil,
	}
}

func getField(fld FolderField, obj fs.Folder) string {
	var targetValue string
	switch fld {
	case FolderId:
		targetValue = obj.Id
	case FolderName:
		targetValue = obj.Name
	case FolderTags:
		targetValue = obj.Tags
	case FolderTemplate:
		targetValue = obj.Template
	case FolderPath:
		targetValue = obj.Path
	case FolderEncrypted:
		targetValue = fmt.Sprintf("%t", obj.Encrypted)
	default:
		fmt.Fprintf(os.Stderr, "ignoring getField(*)")
		targetValue = ""
	}

	return targetValue
}
