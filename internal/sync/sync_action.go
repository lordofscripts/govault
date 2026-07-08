/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package sync

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	NoSyncAction SyncAction = iota
	AddLeft                 // Virtual needs to add folder from disk
	AddRight                // Disk needs to add folder from virtual
	DeleteLeft              // Virtual needs to remove folder (not on disk)
	DeleteRight             // Disk needs to remove folder (not in virtual)
)

const (
	FullSync    UpdateAction = iota // Virtual & Disk updated so that they both match
	UpdateLeft                      // Virtual is updated to match Disk structure
	UpdateRight                     // Disk is updated to match Virtual
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// SyncAction enumerates which operation should be done on an item
// during an UpdateAction.
type SyncAction int

// UpdateAction enumerates the overall type of synchronization.
type UpdateAction int

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

func (e SyncAction) String() string {
	names := map[SyncAction]string{
		NoSyncAction: "-",
		AddLeft:      "◀(+)", // ⬅ ◀
		AddRight:     "(+)▶", // ➡ ▶
		DeleteLeft:   "◀(-)",
		DeleteRight:  "(-)▶",
	}
	return names[e]
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
