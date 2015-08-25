package main


import (
    "./engine"
    "./base"

    "fmt"
)


func main() {
    db := engine.NewEntityDB()
    base.RegisterTypes(db)

    entity := db.New("art")
    world := db.New("map")

    base.Place(db, entity, world, 10, 20, 0)

    emap := db.Get(world, "map").(base.EntityMap)
    fmt.Println(entity, world, emap)
    fmt.Println(db.Get(entity, "location"))
}
