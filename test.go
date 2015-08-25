package main


import (
    "./engine"
    "./base"

    "fmt"
)


func main() {
    db := engine.NewEntityDB()
    base.RegisterTypes(db)

    entity := db.New("art", "movement")
    db.Get(entity, "movement").(*base.Movement).Dx = 1

    world := db.New("map")
    emap := db.Get(world, "map").(base.EntityMap)
    base.Place(db, entity, world, 10, 20, 0)

    fmt.Println(db.Get(entity, "position"))
    fmt.Println(emap.Get(10, 20, 0), emap.Get(11, 20, 0))
    base.SystemMove(db);
    fmt.Println(db.Get(entity, "position"))
    fmt.Println(emap.Get(10, 20, 0), emap.Get(11, 20, 0))
}
