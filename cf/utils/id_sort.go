package utils

import "sort"

// IdSlice attaches the methods of Interface to []uint64, sorting in increasing order.
type IdSlice []uint64
func (p IdSlice) Len() int           { return len(p) }
func (p IdSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IdSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p IdSlice) Sort() { sort.Sort(p) }
