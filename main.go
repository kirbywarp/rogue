package main

import (
    "fmt"
    "github.com/nsf/termbox-go"
    "math"
    "math/rand"
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
func DrawString(text string, x, y int, fg, bg base.Color) {
    for _, char := range text {
        Draw(x, y, char, fg, bg)
        x++
    }
}


/*
CreateMap creates a new map region entity and returns it

This map sort of has 3 "layers".
    Z == 0: The layer for tiles, which includes walls
    Z == 1: The layer for normal entities, mainly the player
    Z == 2: The layer for flying things, like the bat

    This ensures that the player, the bat, and a tile can
    all exist at the same x/y position in the map.  When
    moving an entity, the move function will look at the layer
    one level below the entity moving to see if a wall exists.
    For the player, it will look at Z == 0, and see the tiles.
    The player won't be able to go through walls.
    For the bat, it will look at Z == 1, and usually see nothing.
    The bat will be able to go through walls therefore.
*/
func CreateMap(db *engine.EntityDB) engine.Entity {
    retval := db.New("map")

    // Tiles for the map
    floor := db.New(); db.Set(floor, "art", base.NewArt('.', .5, .5, .5, 0, 0, 0))
    wall  := db.New(); db.Set(wall,  "art", base.NewArt('#',  1,  1 , 1, 0, 0, 0))

    // Filling out the map with a dirt-simple structure
    emap := db.Get(retval, "map").(base.EntityMap)
    for x := int64(-30); x <= 30; x++ {
        for y := int64(-30); y <= 30; y++ {
            if x == 30 || x == -30 || y == 30 || y == -30 ||
               rand.Float64() < .05 {
                emap.Set(x, y, 0, wall)
            } else {
                emap.Set(x, y, 0, floor)
            }
        }
    }

    return retval
}

/*
RenderMapAt draws a portion of the map centered at the given entity
*/
func RenderMapAt(db *engine.EntityDB, eid engine.Entity) {
    width, height := termbox.Size()

    pos := db.Get(eid, "position").(*base.Position)
    emap := db.Get(pos.R, "map").(base.EntityMap)

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


    // menu
    DrawString(title, width/8, height/2, base.RGB(0, 0, 1), base.RGB(0, 0, 0))
    termbox.Flush()

    event1 := termbox.PollEvent()
    for event1.Type != termbox.EventKey {
        event1 = termbox.PollEvent()
    }

    showbat := false
    switch event1.Ch {
    case 'y':
        showbat = true;
        base.HelperPlace(db, bat, tilemap, 5, 5, 2)
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
        case '>': dz = -3
        case '<': dz =  3
        case 0:
            switch event.Key {
            case termbox.KeyCtrlQ:
                done = true
            }
        }

        base.HelperMove(db, player, dx, dy, dz)


        // UPDATING
        if showbat { Follow(db, bat, player) }

        base.SystemMove(db)
    }
}
