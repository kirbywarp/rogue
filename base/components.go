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

// POSITION ============================================================================ //
type Position struct {
    R engine.Entity
    X, Y, Z int64
}
func NewPosition(r engine.Entity, x, y, z int64) Position {
    return Position{R: r, X: x, Y: y, Z: z}
}

func CreatePosition() interface{} { return &Position{} }
func ClonePosition(val interface{}) interface{} { tmp := *(val.(*Position)); return &tmp }

// MOVEMENT ============================================================================ //
type Movement struct {
    Dx, Dy, Dz int64
}
func NewMovement(dx, dy, dz int64) Movement {
    return Movement{Dx: dx, Dy: dy, Dz: dz}
}

func CreateMovement() interface{} { return &Movement{} }
func CloneMovement(val interface{}) interface{} { tmp := *(val.(*Movement)); return &tmp }

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
    db.Register("position", CreatePosition, ClonePosition)
    db.Register("movement", CreateMovement, CloneMovement)
    db.Register("map", CreateEntityMap, CloneEntityMap)
}
