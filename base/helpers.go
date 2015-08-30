package base


import (
    "github.com/kirbywarp/rogue/engine"
)


//////////////////////
// MOVEMENT AND MAP //
//////////////////////

/*
HelperPlace puts an entity on to the map at a particular position
*/
func HelperPlace(db *engine.EntityDB, eid engine.Entity, r engine.Entity, x, y, z int64) {
    if !db.Has(r, "map") { return }

    db.Get(r, "map").(*EntityMap).Set(x, y, z, eid)
    pos := db.Create(eid, "position").(*Position)
    pos.R, pos.X, pos.Y, pos.Z = r, x, y, z
}

/*
HelperMove validates an entity's attempt to move and sets the entity's
movement component appropriately.  Returns true if the entity could move
*/
func HelperMove(db *engine.EntityDB, eid engine.Entity, dx, dy, dz int64) bool {
    if !db.Has(eid, "movement", "position") { return false }

    pos := db.Get(eid, "position").(*Position)
    emap := db.Get(pos.R, "map").(*EntityMap)
    mov := db.Get(eid, "movement").(*Movement)

    // TODO: fix layer system and make it real.
    // Check the entities in the layer beneath the move target
    // to see if the entity is a wall (with art symbol '#')
    target := emap.Get(pos.X+dx, pos.Y+dy, pos.Z+dz-1)
    if db.Has(target, "art") && db.Get(target, "art").(*Art).Symbol == '#' {
        mov.Dx = 0
        mov.Dy = 0
        mov.Dz = 0
        return false
    }

    // Move can be done!
    mov.Dx = dx
    mov.Dy = dy
    mov.Dz = dz
    return true
}
