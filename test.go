package main


import (
    "github.com/kirbywarp/rogue/engine"
    "github.com/kirbywarp/rogue/base"

    "fmt"
)


func main() {
    db := engine.NewEntityDB()
    base.RegisterTypes(db)

    entity := db.New("health")

    health := db.Get(entity, "health").(*base.Health)

    health.SetMax(10)
    health.SetCurrent(5)
    health.Mod(-9)


    fmt.Println(db.Get(entity, "health"))
}
