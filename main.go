/*
Every file in Go defines a part of a package.  All packages except "main"
can be split across files; "main" defines the entrypoint for the application.
Packages other than "main" should be in an appropriate folder in the source
directory, so an "engine" package would be in a folder called "engine".
*/
package main

/*
Each import is a path to a package (or just a name for the built-in things
like math).  The last element of the name is usually the actual name referring
to the packge, so you'd use the random number generator package like
"rand.FunctionName()".

In termbox's case, termbox defines the name as "termbox" rather than
"termbox-go".

For packages, all public names start with a capital letter.  You'll see a lot
of structs and functions all start with caps, including struct members.
*/
import (
    "github.com/nsf/termbox-go"
    "math/rand"
    "time"
    "fmt"
)


/*
This is a struct with three members R, G, and B all of type float64.  They're
capitalized so they're publically accessable.  Go puts the type information at
the end, and if you have a string of variables of the same type, you can collapse
the list and just put one type call.
*/
/*
Color represents an RGB value for use on a 256-color terminal.
*/
type Color struct {
    R, G, B float64
}
/*
A common pattern is to create "constructor" functions for structs; this is just
so you don't have to use the struct literal definition everywhere.  Note that
even the function arguments can have their type definitions condensed, and the type
always goes at the end.

This function takes three parameters of type float64 and returns a Color struct,
which is the actual struct itself rather than a reference or a pointer.
*/
func RGB(r, g, b float64) Color {
    // This is a struct literal; you're creating a new struct and setting the
    //  values for the named members with almost a dictionary syntax
    return Color{R: r, G: g, B: b}
}


/*
Another struct.  "rune" is sort of Go's name for a character, except it gets
a bit tricky with unicode characters which are not all the same size in memory.
Don't have to worry too much about this for now.
*/
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
Here's a type definition for something other than a struct.  In this case,
a MapChunk is going to be an alias for an array of length 256.

In Go, arrays are typed according to their contents and accoring to their
length, so an array of 5 ints is a different type than an array of 10 ints.
There are ways to have variable sized arrays with slices, but it's useful
to use a constant size here.

What this means is that arrays behave almost like a struct with different
syntax; the memory footprint of a MapChunk is literally a contiguous section
of 256*(size of Art struct in bytes) bytes.
*/
/*
MapChunks represent one 16x16 square of tiles in a map that can be
dynamically loaded in as needed.
*/
type MapChunk [256]Art

/*
Another type definition of something other than a struct.  This time, it
aliases the name Map to a Go built-in type called... "map."  Notice how,
as with all the type related stuff in go, the type information goes after
the name.  This is a 'map' that maps uint64 keys to *MapChunk values,
or pointers to a MapChunk.  Pointers are fundamentally just the memory
address of where a value is stored.  In Go, since you can't do pointer
arithmetic like you can in C, they behave most closely to Java references.
*/
/*
Maps hold tile data for a contiguous section of the world, addressable
via x,y,z coordinates.
*/
type Map map[uint64]*MapChunk
func NewMap() Map {
    /*
    This "make" call is how you initializes maps and slices in Go, rather
    than an object literal.  It's a built-in function that takes the type
    of thing to build (in this case, the type is "a map of uint64 to
    MapChunk pointers) and returns the thing.
    You pretty much never use pointers to maps or slices, since they are
    already something similar to pointers.  The actual value that's returned
    is only a couple bytes large, and the actual map data is stored elsewhere.
    */
    return make(map[uint64]*MapChunk)
}

/*
This function definition is analagous to a method definition in Java.  Notice
the "(m Map) in the beginning, before the function mame; that's the "acceptor."
Basically, this function is part of the "method set" of Map types (not *Map, just
Map, although Go will do casting for you if necessary), it is called "Get",
it has three parameters and returns an Art struct.

There is no "this" in these functions; instead, the acceptor gives you the object
reference, which is passed by value.  Because of this, it's common to see the acceptors
be pointers so that you don't duplicate the object's memory on each function call,
but it doesn't matter for small objects like a small struct or maps/slices.
*/
/*
Get returns the tile at the given x,y,z coordinates.  Each coordinate is 3 bytes wide.
*/
func (m Map) Get(x, y, z uint64) Art {
    // Pack the three coordinates into a key for the cunk.  The key itself
    // could be defined as "var key uint64; key = ...", but this shortcut
    // lets you skip declaring the type and having the type just be inferred.
    // Once this define statement is done, key will be of type uint64 in this
    // scope, so things are still strongly typed. It's just nice.
    key := (x&0xFFFFF0)<<36 + (y&0xFFFFF0)<<16 + (z&0xFFFFF)

    // Multiple definitions!  Go let's you return multiple values from functions
    // and assign multiple values in one statement.  It makes things like swapping
    // easy: "a, b = b, a".  In this case, we're reading "key" from the Map m (which
    // if you remember is of type map[uint64]Art and therefore can be used like
    // a builtin map).  The first value is the actual value in the map, or a zeroed
    // version of the return type if the key isn't mapped, and a boolean saying if
    // the read was successful.
    chunk, ok := m[key]
    if ok {
        return chunk[(x&0xF)<<4 + (y&0xF)]
    } else {
        // Create a new map chunk dynamically; remember a MapChunk is just another
        // name for a 256 element array.  Since Go has the length as part of the
        // type, you don't need to specify the length when creating an array value.
        // So you can use the normal construction literal MapChunk{}.  The & operator
        // takes a valure and returns a pointer to it, so here chunk is assigned a
        // pointer to a new MapChunk.
        chunk = &MapChunk{}

        // Good place to talk about syntax.  Go tends to remove a bunch of extraneous
        // syntax stuff.  You won't find parentheses in for and if statements, and you
        // never need ';' at the end of lines unless you actually want two statements
        // in one line.  Otherwise, this should be straightforward; a for loop with
        // a new variable 'i' defined as 0, loop while 'i < 256', and increment on
        // each loop.
        for i := 0; i < 256; i++ {
            // Here you see a cast.  Go is extremely picky about types; it won't even
            // auto-cast numbers, and not even for things like comparisons.  You have
            // to manually cast things using the form <typename>(<value>).  Numeric
            // constants like the 1000 here are stored as the right type by the compiler,
            // so even though it isn't 1000.0 (a float64), the typing works out.
            if rand.Float64() < float64(z%100)/1000 {
                chunk[i] = NewArt('#', RGB(1, 1, 1), RGB(0, 0, 0))
            } else {
                chunk[i] = NewArt('.', RGB(.5, .5, .5), RGB(0,0,0))
            }
        }

        // This is how you store a value into a map, pretty straightforward
        m[key] = chunk

        // Remember, MapChunks are arrays, so this is an array index using
        // a packed key of the coordinates
        return chunk[(x&0xF)<<4 + (y&0xF)]
    }
}

/*
Set sets the value of a tile location on the Map.
*/
func (m Map) Set(x, y, z uint64, a Art) {
    // It should be pretty clear what's going on here.  The & statement is
    // a bitwise 'and', so this is basically masking most of the bottom bits
    // and shifting them to combine into a single uint64 value.
    key := (x&0xFFFFF0)<<36 + (y&0xFFFFF0)<<16 + (z&0xFFFFF)
    chunk, ok := m[key]
    if !ok {
        chunk = &MapChunk{}
        m[key] = chunk
    }
    chunk[(x&0xF)<<4 + (y&0xF)] = a
}










/*
Here we have a function without an acceptor, so it's just a normal function
that can be called standalone or passed around.
*/
/*
Draw draws a bit of tileart to the given x,y coordinates on the terminal.
*/
func Draw(x, y int, a Art) {
    // Here we have termbox being used.  A termbox.Attribute is basically the color information
    // of a cell on the terminal.  This takes the "color" portions of the Art and converts them
    // into an 18-bit color (which is all we have for this color scheme).
    // Notice we don't have to know exactly what type comes from termbox.Attribute(), the compiler
    // handles it all implicitly.
    fg := termbox.Attribute(uint16(a.Fg.R*5+.5)*36+uint16(a.Fg.G*5+.5)*6+uint16(a.Fg.B*5+.5)+1)
    bg := termbox.Attribute(uint16(a.Bg.R*5+.5)*36+uint16(a.Bg.G*5+.5)*6+uint16(a.Bg.B*5+.5)+1)

    // And this part actually sets the cell, with a particular rune and foreground/background.
    termbox.SetCell(x, y, a.Symbol, fg, bg)
}



/*
Here's the main function definition, called when you actually run the application.
This function must exist in a main package file.  There are a couple other functions
that the compiler looks for as well, like an "init" function, but that doesn't matter
right now.
*/
func main() {
    // Seed the random number generator!
    rand.Seed(time.Now().UTC().UnixNano())



    // Create a new tilemap and player coordinates for the game
    tilemap := NewMap()
    px, py, pz := uint64(10), uint64(10), uint64(100)






    // GUI and input initialization
    // This is a common pattern; error handling in Go happens with passed values
    // (which some might consider antequated... I don't know...).  The "err" value
    // is of type error, and will be nil (null) if no error occured.  It's common
    // to see this style of error checking with an if statement right after the
    // call.
    err := termbox.Init()
    if err != nil {
        fmt.Println(err)
        return
    }

    // defer is an awesome thing.  What it does is take a function call, but only
    // actually evaluate just before the closing scope is about to exit.  In this
    // case, whenever control is about to leave the main function (through a return
    // or whatever), this will call termbox.Close() first.  It's great for cleaning
    // up open files and the like, and in this case unsetting some terminal mode
    // settings that make it more useful for drawing and less for regular typing
    // and command entering.
    // It also means that it doesn't matter where you exit; you could have return
    // statements all over the place and this will be executed at the right point.
    defer termbox.Close()

    // Some particular terminal stuff.  This sets the terminal into 216 color mode,
    // and blanks out the terminal.  It also gets the width and height in cells.
    termbox.SetOutputMode(termbox.Output216)
    termbox.Clear(0, 0)
    width, height := termbox.Size()


    // And start the game!
    // Go doesn't have a "while" loop, it just as a "for" loop with optional
    // portions.  In this case, a for loop with only the boolean check.
    done := false
    for !done {
        // RENDERING
        for y := 0; y < height; y++ {
            for x := 0; x < width; x++ {
                Draw(x, height-y-1, tilemap.Get(px+uint64(x-width/2), py+uint64(y-height/2), pz))
            }
        }
        Draw(width/2, height-1-height/2, NewArt('@', RGB(1, 0, 0), RGB(0, 0, 0)))
        // Cell sets will only display when you actually call termbox.Flush().  This
        // is handy because you don't want to show partial updates on screen, which
        // would result in ugliness.
        termbox.Flush()



        // INPUT HANDLING
        // Termbox lets you listen for keyboard input through PollEvent(), which will
        // block until input comes in.  Just realized this is badly named... because
        // Polling usually doesn't block.  Oh well.
        event := termbox.PollEvent()

        // Here is a standard variable declaration, used because there aren't any
        // good values for them yet.  New variables are initialized to 0.
        var dx, dy, dz uint64

        // Here's a switch statement; pretty straightforward, but notice the lack of
        // parentheses.  Also, each case statement doesn't need a "break."  Cases don't
        // "fall through" in Go like they do in other languages.
        // In termbox, you can get the printable character of the key pressed as
        // a rune, or if the character couldn't print (like an arrow key or something),
        // you can get a key code.
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

        // Stupid "collision testing."  Don't let the player move to the new
        // location if there's a wall there.
        if tilemap.Get(px+dx, py+dy, pz+dz).Symbol != '#' {
            px += dx
            py += dy
            pz += dz
        }
    }
}
