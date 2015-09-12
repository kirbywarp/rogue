package base


import (
    "github.com/kirbywarp/rogue/engine"
)


// AI ============================================================================ //
type AIController interface {
    Act(*engine.EntityDB, engine.Entity)
    Clone() AIController
}

type AI struct {
    Controller AIController
}
func NewAI(controller AIController) *AI {
    return &AI{Controller: controller}
}

func CreateAI() interface{} { return &AI{} }
func CloneAI(val interface{}) interface{} {
    tmp := *(val.(*AI))
    tmp.Controller = tmp.Controller.Clone()
    return &tmp
}

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
func NewArt(symbol rune, fgr, fgg, fgb, bgr, bgg, bgb float64) *Art {
    return &Art{Symbol: symbol, Fg: RGB(fgr, fgg, fgb), Bg: RGB(bgr, bgg, bgb)}
}

func CreateArt() interface{} { return &Art{} }
func CloneArt(val interface{}) interface{} { tmp := *(val.(*Art)); return &tmp }

// HEALTH ============================================================================== //

type Health struct {
    Max, Current float64
}

func NewHealth(max, current float64) *Health {
    return &Health{Max: max, Current: current}
}

func (health *Health) Mod(delta float64) {
    health.Current += delta
    if health.Current >= health.Max {
        health.Current = health.Max
    } else if health.Current <= 0 {
        health.Current = 0
    }
}

func (health *Health) SetMax(max float64) {
    if health.Current >= max {
        health.Current = max
    }
    health.Max = max
}

func (health *Health) SetCurrent(current float64) {
    health.Current = current
    if health.Current >= health.Max {
        health.Current = health.Max
    }
}

func CreateHealth() interface{} { return &Health{} }
func CloneHealth(val interface{}) interface{} { tmp := *(val.(*Health)); return &tmp }

// POSITION ============================================================================ //
type Position struct {
    R engine.Entity
    X, Y, Z int64
}
func NewPosition(r engine.Entity, x, y, z int64) *Position {
    return &Position{R: r, X: x, Y: y, Z: z}
}

func CreatePosition() interface{} { return &Position{} }
func ClonePosition(val interface{}) interface{} { tmp := *(val.(*Position)); return &tmp }

// MOVEMENT ============================================================================ //
type Movement struct {
    Dx, Dy, Dz int64
}
func NewMovement(dx, dy, dz int64) *Movement {
    return &Movement{Dx: dx, Dy: dy, Dz: dz}
}

func CreateMovement() interface{} { return &Movement{} }
func CloneMovement(val interface{}) interface{} { tmp := *(val.(*Movement)); return &tmp }

// ENTITY MAP ============================================================================ //
/* See map.go for type definition */

func CreateEntityMap() interface{} { return NewEntityMap() }
func CloneEntityMap(val interface{}) interface{} {
    newmap := NewEntityMap()
    for k, v := range val.(*EntityMap).chunks {
        newmap.chunks[k] = v
    }
    newmap.generator = val.(*EntityMap).generator
    return newmap
}

// ============================================================================ //



func RegisterTypes(db *engine.EntityDB) {
    db.Register("ai", CreateAI, CloneAI)
    db.Register("art", CreateArt, CloneArt)
    db.Register("position", CreatePosition, ClonePosition)
    db.Register("movement", CreateMovement, CloneMovement)
    db.Register("map", CreateEntityMap, CloneEntityMap)
    db.Register("health", CreateHealth, CloneHealth)
}
