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
    X, Y, Z uint64
}
func NewLocation(x, y, z uint64) Location {
    return Location{X: x, Y: y, Z: z}
}

func CreateLocation() interface{} { return &Location{} }
func CloneLocation(val interface{}) interface{} { tmp := *(val.(*Location)); return &tmp }

// ============================================================================ //



func RegisterTypes(db engine.EntityDB) {
    db.Register("art", CreateArt, CloneArt)
    db.Register("location", CreateLocation, CloneLocation)
}
