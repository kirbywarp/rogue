package base


import (
    "github.com/kirbywarp/rogue/engine"
)


/*
MapChunks represent a square of tiles in a map that can be
dynamically loaded in as needed.
*/
type MapChunk [16*16]engine.Entity

/*
EntityMaps hold tile data for a contiguous section of the world, addressable
via x,y,z coordinates.
*/
type EntityMap map[int64]*MapChunk
func NewEntityMap() EntityMap {
    return make(EntityMap)
}

/*
Get returns the tile at the given x,y,z coordinates.  Each coordinate is 3 bytes wide.
*/
func (m EntityMap) Get(x, y, z int64) engine.Entity {
    key := (x&0xFFFFF0)<<36 + (y&0xFFFFF0)<<16 + (z & 0xFFFFF)
    chunk, ok := m[key]
    if ok {
        return chunk[(x&0xF)<<4+(y&0xF)]
    }
    return 0
}

/*
Set sets the value of a tile location on the EntityMap.
*/
func (m EntityMap) Set(x, y, z int64, eid engine.Entity) {
    key := (x&0xFFFFF0)<<36 + (y&0xFFFFF0)<<16 + (z & 0xFFFFF)
    chunk, ok := m[key]
    if !ok {
        chunk = &MapChunk{}
        m[key] = chunk
    }
    chunk[(x&0xF)<<4+(y&0xF)] = eid
}
