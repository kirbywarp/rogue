package engine


import (
    "fmt"
)


/*
Entities are unique names for collections of components
*/
type Entity uint64



/*
Managers manage a database of components index by entity id
*/
type manager struct {
    create func() interface{}
    clone func(interface{}) interface{}
    comps map[Entity]interface{}
}
func (manager manager) Create(eid Entity) interface{} {
    comp := manager.create()
    manager.comps[eid] = comp
    return comp
}
func (manager manager) Clone(src, dst Entity) interface{} {
    oldc, ok := manager.comps[src]
    if !ok { return nil }

    newc := manager.clone(oldc)
    manager.comps[dst] = newc
    return newc
}
func (manager manager) Get(eid Entity) interface{} {
    return manager.comps[eid]
}
func (manager manager) Remove(eid Entity) {
    delete(manager.comps, eid)
}
func (manager manager) Has(eid Entity) bool {
    _, ok := manager.comps[eid]
    return ok
}
func (manager manager) Entities() []Entity {
    list := make([]Entity, len(manager.comps))
    i := 0
    for key, _ := range manager.comps {
        list[i] = key
        i++
    }
    return list
}



/*
EntityDBs hold entity data and provide convenience methods for accessing
components or component managers.
*/
type EntityDB struct {
    nextid Entity
    managers map[string]manager
}
func NewEntityDB() *EntityDB {
    return &EntityDB{nextid: 1, managers: make(map[string]manager)}
}

/*
Register registers a new component with the database under the passed name
*/
func (db *EntityDB) Register(name string, create func() interface{}, clone func(interface{}) interface{}) {
    db.managers[name] = manager{create: create, clone: clone, comps: make(map[Entity]interface{})}
}

/*
New creates a new entity, optionally with the specified empty components
*/
func (db *EntityDB) New(components ...string) Entity {
    eid := db.nextid
    db.nextid++

    for _, name := range components { db.managers[name].Create(eid) }
    return eid
}

/*
Instance creates a new entity using the passed entity as a template.
*/
func (db *EntityDB) Instance(template Entity) Entity {
    eid := db.nextid
    db.nextid++

    for _, manager := range db.managers { manager.Clone(template, eid) }
    return eid
}

/*
Delete removes an entity from the database
*/
func (db *EntityDB) Delete(eid Entity) {
    for _, manager := range db.managers { manager.Remove(eid) }
}

/*
Manager retrieves the appropriate manager for a component type
*/
func (db *EntityDB) Manager(name string) manager {
    manager, ok := db.managers[name]
    if !ok { panic(fmt.Sprintf("EntityDB: No manager registered for component '%s'", name)) }
    return manager
}

/*
Create creates a new component for the given entity and returns it
*/
func (db *EntityDB) Create(eid Entity, name string) interface{} {
    return db.Manager(name).Create(eid)
}

/*
Get retrieves a component for the given entity
*/
func (db *EntityDB) Get(eid Entity, name string) interface{} {
    return db.Manager(name).Get(eid)
}

/*
Remove removes a component from the given entity
*/
func (db *EntityDB) Remove(eid Entity, name string) {
    db.Manager(name).Remove(eid)
}

/*
Has returns true if the passed entity has every given component
*/
func (db *EntityDB) Has(eid Entity, components ...string) bool {
    for _, name := range components {
        if !db.Manager(name).Has(eid) { return false }
    }
    return true
}

/*
Search returns a list of Entities that have the given list of components
*/
func (db *EntityDB) Search(name string, components ...string) []Entity {
    base := db.Manager(name).Entities()
    retval := make([]Entity, 0, len(base))
    for _, eid := range base {
        if db.Has(eid, components...) {
            retval = append(retval, eid)
        }
    }
    return retval
}
