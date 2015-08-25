package base


import (
    "../engine"
)


//////////////////////
// MOVEMENT AND MAP //
//////////////////////

/*
Place puts an entity on to the map at a particular position
*/
func Place(db *engine.EntityDB, eid engine.Entity, r engine.Entity, x, y, z int64) {
    if !db.Has(r, "map") { return }

    db.Get(r, "map").(EntityMap).Set(x, y, z, eid)
    pos := db.Create(eid, "position").(*Position)
    pos.R, pos.X, pos.Y, pos.Z = r, x, y, z
}
