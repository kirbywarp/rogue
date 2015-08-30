package main

import (
    "fmt"
    "github.com/nsf/termbox-go"
    "math"
    "math/rand"
    "os"
    "strconv"
    "time"

    "github.com/kirbywarp/rogue/engine"
    "github.com/kirbywarp/rogue/base"
)



/*
Draw draws a bit of tileart to the given x,y coordinates on the terminal.
*/
func Draw(x, y int, symbol rune, fg, bg base.Color) {
    // Convert colors into termbox attributes.
    // To get a nice scale, round each channel.
    fga := termbox.Attribute(uint16(fg.R*5+.5)*36 + uint16(fg.G*5+.5)*6 + uint16(fg.B*5+.5) + 1)
    bga := termbox.Attribute(uint16(bg.R*5+.5)*36 + uint16(bg.G*5+.5)*6 + uint16(bg.B*5+.5) + 1)

    termbox.SetCell(x, y, symbol, fga, bga)
}
/*
DrawString prints text to the termbox on row y starting at column x
*/
func DrawString(x, y int, text string, fg, bg base.Color) {
    for _, char := range text {
        Draw(x, y, char, fg, bg)
        x++
    }
}
/*
DrawPaddedString prints text to termbox with size cells filled up. Any remaining cells not covered by text
will be drawn empty. If text is longer than size, it will be truncated.
*/

func DrawPaddedString(x, y int, text string, fg, bg base.Color, size int) {
    var i int
    var char rune
    for _, char = range text {
        if i >= size {
            break
        }
        Draw(x+i, y, char, fg, bg)
        i++
    }
    for i < size {
        Draw(x+i, y, '.', bg, bg)
        i++
    }
}

/*
DrawTextBox creates a textbox of length size at coord x, y on the termbox and blocks while the user
inputs a string. It stops blocking when the enter key is pressed
*/

func DrawTextBox(x, y int, fg, bg base.Color, size int) string {
    //Draw the box
    x1 := x
    for i := 0; i < size+2; i++ {
        Draw(x1, y, '-', fg, bg)
        Draw(x1, y-2, '-', fg, bg)
        x1++
    }
    Draw(x, y-1, '|', fg, bg)
    Draw(x+size+1, y-1, '|', fg, bg)
    termbox.Flush()

    //check input
    buffer := make([]rune, 0)
    done := false
    for !done {
        event := termbox.PollEvent()
        switch event.Ch {
        default:
            buffer = append(buffer, event.Ch)
            DrawPaddedString(x+1, y-1, string(buffer), fg, bg, size)
            termbox.Flush()
        case 0:
            switch event.Key {
            case termbox.KeyEnter : done = true
            case termbox.KeyCtrlQ:
                termbox.Close()
                os.Exit(0)
            case termbox.KeyBackspace :
                if len(buffer) > 0{
                    buffer = buffer[:len(buffer)-1]
                    DrawPaddedString(x+1, y-1, string(buffer), fg, bg, size)
                    termbox.Flush()
                }
            }
        }
    }
    DrawPaddedString(x+1, y-1, "", fg, bg, size)
    return string(buffer)
}






/*
MAP GENERATION
*/
type StoneFieldGenerator struct {
    stone, grass engine.Entity
    fill float64
}
func NewStoneFieldGenerator(db *engine.EntityDB, fill float64) *StoneFieldGenerator {
    grass := db.New(); db.Set(grass, "art", base.NewArt('.', 0, 1, 0, 0, 0, 0))
    stone := db.New(); db.Set(stone, "art", base.NewArt('#',  .7,  .7 , .7, 0, 0, 0))
    return &StoneFieldGenerator{stone: stone, grass: grass, fill: fill}
}
func (g *StoneFieldGenerator) GenerateChunk(emap *base.EntityMap, x, y, z int64) {
    chunk := emap.CreateChunk(x, y, z)
    for x := int64(0); x < 16; x++ {
        for y := int64(0); y < 16; y++ {
            if rand.Float64() < g.fill {
                chunk.Set(x, y, 0, g.stone)
            } else {
                chunk.Set(x, y, 0, g.grass)
            }
        }
    }
}

/*
CreateMap creates a new map region entity and returns it
*/
func CreateMap(db *engine.EntityDB) engine.Entity {
    retval := db.New("map")

    // Register a chunk generator on the map
    emap := db.Get(retval, "map").(*base.EntityMap)
    emap.RegisterChunkGenerator(NewStoneFieldGenerator(db, .05))

    return retval
}

/*
RenderMapAt draws a portion of the map centered at the given entity
*/
func RenderMapAt(db *engine.EntityDB, eid engine.Entity) {
    width, height := termbox.Size()

    pos := db.Get(eid, "position").(*base.Position)
    emap := db.Get(pos.R, "map").(*base.EntityMap)

    for y := 0; y < height; y++ {
        py := pos.Y+int64(y-height/2)
        for x := 0; x < width; x++ {
            px := pos.X+int64(x-width/2)

            // Search for the highest entity on the map in the same general layer
            //  as the passed entity that can be drawn.  This will need to be more
            //  formal in the future (not hardcoded knowing the player is on z level
            //  1, tiles are on z level 0, and the bat at z level 2)
            topArt := &base.Art{}
            for i := int64(1); i >= -1; i-- {
                entity := emap.Get(px, py, pos.Z+i)
                if db.Has(entity, "art") {
                    topArt = db.Get(entity, "art").(*base.Art)
                    break
                }
            }

            // And draw the found art, which will be an empty black square
            //  if nothing was found
            Draw(x, height-1-y, topArt.Symbol, topArt.Fg, topArt.Bg)
        }
    }
    termbox.Flush()
}



/*
Follow attempts to move an entity to chase another entity (in this case, bat chases player)
*/
func Follow(db *engine.EntityDB, eid, target engine.Entity) {
    epos := db.Get(eid, "position").(*base.Position)
    tpos := db.Get(target, "position").(*base.Position)

    dx, dy := int64(0), int64(0)
    if epos.X > tpos.X+1 || epos.X < tpos.X-1 {
        dx = int64(math.Copysign(1, float64(tpos.X-epos.X)))
    }
    if epos.Y > tpos.Y+1 || epos.Y < tpos.Y-1 {
        dy = int64(math.Copysign(1, float64(tpos.Y-epos.Y)))
    }

    // Since the bat is situated at z-level 2, it will never find
    // a wall with art symbol '#' at one level below, so this move
    // function still allows the bat to fly over stones! ^_^
    // Hence the "+1" for the z; so the bat stays one level above
    //  the player.
    base.HelperMove(db, eid, dx, dy, tpos.Z-epos.Z+1)
}





func main() {
    // Seed the random number generator!
    rand.Seed(time.Now().UTC().UnixNano())



    // Game Data Initialization
    db := engine.NewEntityDB()
    base.RegisterTypes(db)

    tilemap := CreateMap(db)

    player := db.New("movement")
    db.Set(player, "art", base.NewArt('@', 1, 0, 0, 0, 0, 0))
    base.HelperPlace(db, player, tilemap, 0, 0, 1)

    bat := db.New("movement")
    db.Set(bat, "art", base.NewArt('b', 0, 0, 1, 0, 0, 0))
    bats := make([]engine.Entity, 0)



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
    batTextBox := "How many bats?"
    tryAgain := "Please enter an integer"


    // menu
    DrawString(width/8, height/4, title, base.RGB(0, 0, 1), base.RGB(0, 0, 0))
    termbox.Flush()

    event1 := termbox.PollEvent()
    for event1.Type != termbox.EventKey {
        event1 = termbox.PollEvent()
    }

    showbat := false
    var numBats int64

    switch event1.Ch {
    case 'y':
        showbat = true;
    }

    if showbat{
        DrawString(width/2, height/2-4, batTextBox, base.RGB(0, 0, 1), base.RGB(0, 0, 0))
        numBatstr := DrawTextBox(width/2, height/2, base.RGB(0, 0, 1), base.RGB(0, 0, 0), 4)
        var err error
        numBats, err = strconv.ParseInt(numBatstr, 10, 64)
        for err != nil {
            DrawString(width/2, height/2-3, tryAgain, base.RGB(0, 0, 1), base.RGB(0, 0, 0))
            numBatstr := DrawTextBox(width/2, height/2, base.RGB(0, 0, 1), base.RGB(0, 0, 0), 4)
            numBats, err = strconv.ParseInt(numBatstr, 10, 64)
        }

        //TODO: Copy bat numBats times (lol numBats) and remove print statement
        for i := int64(0); i < numBats; i++ {
            newBat := db.Instance(bat)
            bats= append(bats, newBat)
            base.HelperPlace(db, newBat, tilemap, rand.Int63n(numBats)- numBats/2, rand.Int63n(numBats)- numBats/2, 2)
        }
    }



    // game loop
    done := false
    for !done {
        // RENDERING
        RenderMapAt(db, player)


        // INPUT HANDLING
        event := termbox.PollEvent()

        var dx, dy, dz int64

        switch event.Ch {
        case 'h': dx = -1
        case 'j': dy = -1
        case 'k': dy =  1
        case 'l': dx =  1
        case 'y': dx = -1; dy =  1
        case 'u': dx =  1; dy =  1
        case 'b': dx = -1; dy = -1
        case 'n': dx =  1; dy = -1
        case '>': dz = -4
        case '<': dz =  4
        case 0:
            switch event.Key {
            case termbox.KeyCtrlQ:
                done = true
            }
        }

        base.HelperMove(db, player, dx, dy, dz)


        // UPDATING
        if showbat { 
            for _ , bat := range bats {
                Follow(db, bat, player) 
            }
        }

        base.SystemMove(db)
    }
}
