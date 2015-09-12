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
FollowAI makes an entity doggedly follow a target
*/
type FollowAI struct {
    Target engine.Entity
}
func NewFollowAI(target engine.Entity) *FollowAI {
    return &FollowAI{Target: target}
}
func (ai *FollowAI) Clone() base.AIController {
    return &FollowAI{Target: ai.Target}
}
func (ai *FollowAI) Act(db *engine.EntityDB, eid engine.Entity) {
    epos := db.Get(eid, "position").(*base.Position)
    tpos := db.Get(ai.Target, "position").(*base.Position)

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


var (
    done = false
)

/*
PlayerAI makes an entity respond to player controls.
*/
type PlayerAI struct {}
func NewPlayerAI() *PlayerAI {
    return &PlayerAI{}
}
func (ai *PlayerAI) Clone() base.AIController {
    return &PlayerAI{}
}
func (ai *PlayerAI) Act(db *engine.EntityDB, eid engine.Entity) {
        // RENDERING
        RenderMapAt(db, eid)

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

        base.HelperMove(db, eid, dx, dy, dz)
}



/************************
 ** State Machine Code **
 ************************/

/*
States are UI states that respond to time ticks and input, and can
transition to other states in various ways.
*/
type State interface {
    Enter(*UI)                  // Called when this state is entered from a different state
    Exit(*UI)                   // Called when the machine transfers away from this state
    Update(*UI, float64)        // Called every frame of animation
}

/*
UI holds state information and provides convenience methods for
things like dialogue boxes and input areas.
*/
type UI struct {
    States map[string]State
    StateStack []string
    Properties map[string]interface{}
}
func NewUI() *UI {
    return &UI{States: make(map[string]State), StateStack: make([]string, 0), Properties: make(map[string]interface{})}
}

func (ui *UI) RegisterState(name string, state State) {
    ui.States[name] = state
}

func (ui *UI) Transition(name string) {
    l := len(ui.StateStack)
    if l > 0 {
        ui.StateStack[l-1] = name
    }
}
func (ui *UI) Push(name string) {
    ui.StateStack = append(ui.StateStack, name)
}
func (ui *UI) Pop() {
    l := len(ui.StateStack)
    if l > 0 {
        ui.StateStack = ui.StateStack[:l-1]
    }
}
func (ui *UI) Peek() string {
    l := len(ui.StateStack)
    if l > 0 {
        return ui.StateStack[l-1]
    }
    return ""
}

func (ui *UI) Run() {
    name := ui.Peek()
    state := ui.States[name]
    state.Enter(ui)

    for {
        state.Update(ui, 1)

        next := ui.Peek()
        if next == "" {
            return
        } else if next != name {
            state.Exit(ui)
            name = next
            state = ui.States[name]
            state.Enter(ui)
        }
    }
}





/************************
 * State implementation *
 ************************/

/*
TitleState
*/
type TitleState struct {}
func NewTitleState() *TitleState {
    return &TitleState{}
}

func (title *TitleState) Enter(ui *UI) {}
func (title *TitleState) Exit(ui *UI) {}
func (title *TitleState) Update(ui *UI, dt float64) {
    width, height := termbox.Size()

    text := "Press any key to play. Press 'y' to face an bat at your own risk!"
    DrawString(width/8, height/4, text, base.RGB(0, 0, 1), base.RGB(0, 0, 0))
    termbox.Flush()

    event := termbox.PollEvent()
    for event.Type != termbox.EventKey {
        event = termbox.PollEvent()
    }

    switch event.Ch {
    case 'y':
        ui.Transition("batmenu")
        return
    case 0:
        switch event.Key {
        case termbox.KeyCtrlQ:
            ui.Pop()
            return
        }
    }

    // Create a game without bats
    ui.RegisterState("game", NewGameState(0))
    ui.Transition("game")
}

/*
BatMenuState
*/
type BatMenuState struct {}
func NewBatMenuState() *BatMenuState {
    return &BatMenuState{}
}

func (menu *BatMenuState) Enter(ui *UI) {}
func (menu *BatMenuState) Exit(ui *UI) {}
func (menu *BatMenuState) Update(ui *UI, dt float64) {
    width, height := termbox.Size()

    label := "How many bats?"
    tryAgain := "Please enter an integer"

    DrawString(width/2, height/2-4, label, base.RGB(0, 0, 1), base.RGB(0, 0, 0))
    numBatstr := DrawTextBox(width/2, height/2, base.RGB(0, 0, 1), base.RGB(0, 0, 0), 4)
    numBats, err := strconv.ParseInt(numBatstr, 10, 64)
    for err != nil {
        DrawString(width/2, height/2-3, tryAgain, base.RGB(0, 0, 1), base.RGB(0, 0, 0))
        numBatstr := DrawTextBox(width/2, height/2, base.RGB(0, 0, 1), base.RGB(0, 0, 0), 4)
        numBats, err = strconv.ParseInt(numBatstr, 10, 64)
    }

    ui.RegisterState("game", NewGameState(numBats))
    ui.Transition("game")
}


/*
GameState
*/
type GameState struct {
    DB *engine.EntityDB
}
func NewGameState(numbats int64) *GameState {
    // Game Data Initialization
    db := engine.NewEntityDB()
    base.RegisterTypes(db)

    tilemap := CreateMap(db)

    player := db.New("movement")
    db.Set(player, "ai", base.NewAI(NewPlayerAI()))
    db.Set(player, "art", base.NewArt('@', 1, 0, 0, 0, 0, 0))
    base.HelperPlace(db, player, tilemap, 0, 0, 1)

    bat := db.New("movement")
    db.Set(bat, "ai", base.NewAI(NewFollowAI(player)))
    db.Set(bat, "art", base.NewArt('b', 0, 0, 1, 0, 0, 0))

    // Create bats from the template entity
    for i := int64(0); i < numbats; i++ {
        newBat := db.Instance(bat)
        base.HelperPlace(db, newBat, tilemap, rand.Int63n(numbats)-numbats/2, rand.Int63n(numbats)-numbats/2, 2)
    }

    return &GameState{DB: db}
}

func (game *GameState) Enter(ui *UI) {}
func (game *GameState) Exit(ui *UI) {}
func (game *GameState) Update(ui *UI, dt float64) {
    base.SystemAct(game.DB)
    base.SystemMove(game.DB)

    if done {
        ui.Pop()
    }
}



func main() {
    // Seed the random number generator!
    rand.Seed(time.Now().UTC().UnixNano())

    // GUI and input initialization
    err := termbox.Init()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer termbox.Close()
    termbox.SetOutputMode(termbox.Output216)
    termbox.Clear(0, 0)

    // Initial states
    ui := NewUI()
    ui.RegisterState("title", NewTitleState())
    ui.RegisterState("batmenu", NewBatMenuState())

    ui.Push("title")
    ui.Run()
}
