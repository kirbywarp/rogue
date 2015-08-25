package main

import (
    "fmt"
    "github.com/nsf/termbox-go"
    "math"
    "math/rand"
    "time"
)

/*
Color represents an RGB value for use on a 256-color terminal.
*/
type Color struct {
    R, G, B float64
}

func RGB(r, g, b float64) Color {
    // This is a struct literal; you're creating a new struct and setting the
    //  values for the named members with almost a dictionary syntax
    return Color{R: r, G: g, B: b}
}

/*
Art holds information for a colored symbol on the terminal.
*/
type Art struct {
    Symbol rune
    Fg, Bg Color
}

func NewArt(r rune, fg Color, bg Color) Art {
    return Art{Symbol: r, Fg: fg, Bg: bg}
}

/*
MapChunks represent one 16x16 square of tiles in a map that can be
dynamically loaded in as needed.
*/
type MapChunk [256]Art

/*
Maps hold tile data for a contiguous section of the world, addressable
via x,y,z coordinates.
*/
type Map map[uint64]*MapChunk

func NewMap() Map {
    return make(Map)
}

/*
Get returns the tile at the given x,y,z coordinates.  Each coordinate is 3 bytes wide.
*/
func (m Map) Get(x, y, z uint64) Art {
    key := (x&0xFFFFF0)<<36 + (y&0xFFFFF0)<<16 + (z & 0xFFFFF)

    chunk, ok := m[key]
    if ok {
        return chunk[(x&0xF)<<4+(y&0xF)]
    } else {
        chunk = &MapChunk{}

        for i := 0; i < 256; i++ {

            if rand.Float64() < float64(z%100)/1000 {
                chunk[i] = NewArt('#', RGB(1, 1, 1), RGB(0, 0, 0))
            } else {
                chunk[i] = NewArt('.', RGB(.5, .5, .5), RGB(0, 0, 0))
            }
        }

        m[key] = chunk

        return chunk[(x&0xF)<<4+(y&0xF)]
    }
}

/*
Set sets the value of a tile location on the Map.
*/
func (m Map) Set(x, y, z uint64, a Art) {
    key := (x&0xFFFFF0)<<36 + (y&0xFFFFF0)<<16 + (z & 0xFFFFF)
    chunk, ok := m[key]
    if !ok {
        chunk = &MapChunk{}
        m[key] = chunk
    }
    chunk[(x&0xF)<<4+(y&0xF)] = a
}

/*
Draw draws a bit of tileart to the given x,y coordinates on the terminal.
*/
func Draw(x, y int, a Art) {
    fg := termbox.Attribute(uint16(a.Fg.R*5+.5)*36 + uint16(a.Fg.G*5+.5)*6 + uint16(a.Fg.B*5+.5) + 1)
    bg := termbox.Attribute(uint16(a.Bg.R*5+.5)*36 + uint16(a.Bg.G*5+.5)*6 + uint16(a.Bg.B*5+.5) + 1)

    termbox.SetCell(x, y, a.Symbol, fg, bg)
}
/*
BatMovement determines which direction the an bat should move in order to chase the player.
It returns dx and dy where dx,dy = 0 or ±1
An bats can fly over stones, and automatically follow the player in the z direction.

*/
func BatMovement(bx, by, px, py uint64) (float64, float64) {     
    dx, dy := float64(0), float64(0)
    
    if (float64(bx) > (float64(px)+1)) || (float64(bx) < (float64(px)-1)) {
        dx = math.Copysign(1,float64(px)-float64(bx))
    } 
    if (float64(by) > (float64(py)+1)) || (float64(by) < (float64(py)-1)) {
        dy = math.Copysign(1,float64(py)-float64(by))
    }    
    return dx, dy
}

/*
TermboxPrint prints text to the termbox on row y starting at column x
*/

func TermboxPrint(text string, x, y int, fg, bg termbox.Attribute) {
	for _, char := range text {
		termbox.SetCell(x, y, char, fg, bg)
		x++
	}
}

func main() {
    // Seed the random number generator!
    rand.Seed(time.Now().UTC().UnixNano())

    // Create a new tilemap and player coordinates for the game
    tilemap := NewMap()
    cx, cy, cz := 0x80000, 0x80000, 0x80000
    px, py, pz := uint64(cx+10), uint64(cy+10), uint64(cz+100)
    bx, by := uint64(cx+5), uint64(cy+5)

    // GUI and input initialization
    err := termbox.Init()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer termbox.Close()

    termbox.SetOutputMode(termbox.Output216)
    termbox.Clear(0, 0)
    width, height := termbox.Size()
    title := "Press any key to play. Press 'y' to face an bat at your own risk!"


    // And start the game!
    
    //menu??
	TermboxPrint(title, width/8, height/2, termbox.ColorBlue, termbox.ColorBlack)
	termbox.Flush()
	
	//TODO: blocking?
    event1 := termbox.PollEvent()

    bat := false
    
    switch event1.Ch {
        case 'y': bat = true
        }
    
    done := false

    for !done {
        
        // UPDATING
        bdx, bdy := BatMovement(bx,by,px,py)
        bx = uint64(float64(bx) + bdx)
        by = uint64(float64(by) + bdy)
             
        // RENDERING
        for y := 0; y < height; y++ {
            for x := 0; x < width; x++ {
                Draw(x, height-y-1, tilemap.Get(px+uint64(x-width/2), py+uint64(y-height/2), pz))
            }
        }
        Draw(width/2, height-1-height/2, NewArt('@', RGB(1, 0, 0), RGB(0, 0, 0)))
        if bat{
            Draw(width/2+int(bx)-int(px), height-1-(height/2+int(by)-int(py)), NewArt('b', RGB(0, 0, 1), RGB(0, 0, 0)))  
        }
        termbox.Flush()
        

        // INPUT HANDLING
        event := termbox.PollEvent()

        var dx, dy, dz uint64

        switch event.Ch {
        case 'h': dx -= 1
        case 'j': dy -= 1
        case 'k': dy += 1
        case 'l': dx += 1
        case 'y': dx -= 1; dy += 1
        case 'u': dx += 1; dy += 1
        case 'b': dx -= 1; dy -= 1
        case 'n': dx += 1; dy -= 1
        case '>': dz -= 1
        case '<': dz += 1
        case 0:
            switch event.Key {
            case termbox.KeyCtrlQ:
                done = true
            }
        }
        
        

        if tilemap.Get(px+dx, py+dy, pz+dz).Symbol != '#' {
            px += dx
            py += dy
            pz += dz
        }
        
        
    }
}
