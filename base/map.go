package base


import (
    "github.com/kirbywarp/rogue/engine"
)


/*
ChunkGenerators are tasked with generating new chunks on-demand when something
tries to get/set a chunk that doesn't exist.
*/
type ChunkGenerator interface {
    GenerateChunk(*EntityMap, int64, int64, int64)
}





/*
MapChunks represent a square of tiles in a map that can be
dynamically loaded in as needed.
*/
type MapChunk [16*16*4]engine.Entity

/*
EntityMaps hold tile data for a contiguous section of the world, addressable
via x,y,z coordinates.
*/
type EntityMap struct{
    chunks map[int64]*MapChunk
    generator ChunkGenerator
}
func NewEntityMap() *EntityMap {
    return &EntityMap{chunks: make(map[int64]*MapChunk)}
}

/*
Get returns the tile at the given x,y,z coordinates.  Each coordinate is 3 bytes wide.
*/
func (m *EntityMap) Get(x, y, z int64) engine.Entity {
    key := (x&0xFFFFF0)<<36 + (y&0xFFFFF0)<<16 + (z&0x3FFFFC)>>2
    chunk, ok := m.chunks[key]
    if chunk != nil {
        return chunk[(x&0xF)<<6+(y&0xF)<<2+(z&0x3)]
    } else if !ok && m.generator != nil{
        m.generator.GenerateChunk(m, (x&0xFFFFF0)>>4, (y&0xFFFFF0)>>4, (z&0x3FFFFC)>>2)
        return m.Get(x, y, z)
    }
    return 0
}

/*
Set sets the value of a tile location on the EntityMap.
*/
func (m *EntityMap) Set(x, y, z int64, eid engine.Entity) {
    key := (x&0xFFFFF0)<<36 + (y&0xFFFFF0)<<16 + (z&0x3FFFFC)>>2
    chunk, ok := m.chunks[key]
    if chunk != nil {
        chunk[(x&0xF)<<6+(y&0xF)<<2+(z&0x3)] = eid
    } else if !ok && m.generator != nil{
        m.generator.GenerateChunk(m, (x&0xFFFFF0)>>4, (y&0xFFFFF0)>>4, (z&0x3FFFFC)>>2)
        m.Set(x, y, z, eid)
    }
}

/*
ChunkGenerated returns true if a particular chunk has already been generated.  This
doesn't mean the chunk won't be nil internally, just that it shouldn't be passed
to the chunk generator again.
*/
func (m *EntityMap) ChunkGenerated(x, y, z int64) bool {
    _, ok := m.chunks[(x&0xFFFFF)<<40 + (y&0xFFFFF)<<20 + z&0xFFFFF]
    return ok
}

/*
CreateChunk creates a new chunk at the passed coordinates, overwriting any existing
data, and returns a pointer to the new chunk.
*/
func (m *EntityMap) CreateChunk(x, y, z int64) *MapChunk {
    chunk := &MapChunk{}
    m.chunks[(x&0xFFFFF)<<40 + (y&0xFFFFF)<<20 + z&0xFFFFF] = chunk
    return chunk
}
