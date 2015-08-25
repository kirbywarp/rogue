package base


import (
    "../engine"
)


///////////////////////
// MAP MANIPULATIONS //
///////////////////////

/*
Place puts an entity on to the map at a particular location
*/
func Place(db *engine.EntityDB, eid engine.Entity, r engine.Entity, x, y, z uint64) {
    if !db.Has(r, "map") { return }

    db.Get(r, "map").(EntityMap).Set(x, y, z, eid)
    loc := db.Create(eid, "location").(*Location)
    loc.R, loc.X, loc.Y, loc.Z = r, x, y, z
}
