package base


import (
    "../engine"
)


// ART ============================================================================ //
type Color struct {
    R, G, B float64
}
func RGB(r, g, b float64) Color {
    return Color{R: r, G: g, B: b}
}

type Art struct {
    Symbol rune
    Fg, Bg Color
}

func CreateArt() interface{} { return &Art{} }
func CloneArt(val interface{}) interface{} { tmp := *(val.(*Art)); return &tmp }

// LOCATION ============================================================================ //
type Location struct {
    R engine.Entity
    X, Y, Z uint64
}
func NewLocation(r engine.Entity, x, y, z uint64) Location {
    return Location{R: r, X: x, Y: y, Z: z}
}

func CreateLocation() interface{} { return &Location{} }
func CloneLocation(val interface{}) interface{} { tmp := *(val.(*Location)); return &tmp }

// ENTITY MAP ============================================================================ //
/* See map.go for type definition */

func CreateEntityMap() interface{} { return NewEntityMap() }
func CloneEntityMap(val interface{}) interface{} {
    newmap := NewEntityMap()
    for k, v := range val.(EntityMap) {
        newmap[k] = v
    }
    return newmap
}

// ============================================================================ //



func RegisterTypes(db *engine.EntityDB) {
    db.Register("art", CreateArt, CloneArt)
    db.Register("location", CreateLocation, CloneLocation)
    db.Register("map", CreateEntityMap, CloneEntityMap)
}
