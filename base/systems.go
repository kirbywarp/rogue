package base


import (
    "github.com/kirbywarp/rogue/engine"
)


/////////////////////
// MOVEMENT SYSTEM //
/////////////////////

/*
Move applies makes every entity with movement try to move in the map
*/
func SystemMove(db *engine.EntityDB) {
    for _, eid := range db.Search("movement", "position") {
        pos := db.Get(eid, "position").(*Position)
        mov := db.Get(eid, "movement").(*Movement)
        emap := db.Get(pos.R, "map").(*EntityMap)

        // Prevent entities from moving on top of each other, temporarily
        if emap.Get(pos.X+mov.Dx, pos.Y+mov.Dy, pos.Z+mov.Dz) != 0 { continue }

        // Move and update the map
        emap.Set(pos.X, pos.Y, pos.Z, 0)
        pos.X += mov.Dx; pos.Y += mov.Dy; pos.Z += mov.Dz
        emap.Set(pos.X, pos.Y, pos.Z, eid)
    }
}
