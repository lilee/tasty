package utils

import (
    "sort"
)

type IdSet map[uint64]struct{}

func NewIdSet() IdSet {
    return make(IdSet)
}

// Adds an item to the current set if it doesn't already exist in the set.
func (this IdSet) Add(id uint64) bool {
    _, found := this[id]
    this[id] = struct{}{}
    return !found //False if it existed already
}

// Determines if a given item is already in the set.
func (this IdSet) Contains(id uint64) bool {
    _, found := this[id]
    return found
}

// Allows the removal of a single item in the set.
func (this IdSet) Remove(id uint64) {
    delete(this, id)
}

func (this IdSet) AddArray(ids []uint64) {
    for _, id := range ids {
        this.Add(id)
    }
}

func (this IdSet) RemoveArray(ids []uint64) {
    for _, id := range ids {
        this.Remove(id)
    }
}

func (this IdSet) ToArray() (ids []uint64) {
    for id := range this {
        ids = append(ids, id)
    }
    sort.Sort(IdSlice(ids))
    return
}

