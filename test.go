package main


import (
    "github.com/kirbywarp/rogue/engine"
    "github.com/kirbywarp/rogue/base"

    "fmt"
)


func main() {
    db := engine.NewEntityDB()
    base.RegisterTypes(db)

    entity := db.New("health","attack")

    health := db.Get(entity, "health").(*base.Health)
    attack := db.Get(entity, "attack").(*base.Attack)

    health.SetMax(10)
    health.SetCurrent(5)
    health.Mod(-9)
    attack.Damage = 5


    fmt.Println(db.Get(entity, "health"))
    fmt.Println(db.Get(entity, "attack"))
}
