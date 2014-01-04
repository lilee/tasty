package recommender

import (
    "container/heap"
    "sort"

    "github.com/lilee/tasty/cf"
    "github.com/lilee/tasty/cf/model"
)

type estimator interface {
    estimate(thing interface{}) (float64, error)
}

type estimatorFunc func(id uint64) (float64, error)

func (f estimatorFunc) estimate(thing interface{}) (float64, error) {
    id, ok := thing.(uint64)
    if !ok {
        return 0.0, cf.NaNError
    }
    return f(id)
}

func capper(v float64, m model.DataModel) float64 {
    if v > m.MaxPreferenceValue() {
        v = m.MaxPreferenceValue()
    } else if v < m.MinPreferenceValue() {
        v = m.MinPreferenceValue()
    }
    return v
}

type pqitem interface {
    Less(item interface{}) bool
}

// TopItems
type pq struct {
    items []pqitem
    n int
    lowest pqitem
    full bool
}

func newPQ(n int) *pq {
    a := &pq {
        items: []pqitem{},
        n: n,
        full: false,
    }
    heap.Init(a)
    return a
}

// interface for heap.Interface
func (this *pq) Len() int {
    return len(this.items)
}

func (this *pq) Less(i, j int) bool {
    return this.items[i].Less(this.items[j])
}

func (this *pq) Swap(i, j int) {
    this.items[i], this.items[j] = this.items[j], this.items[i]
}

func (this *pq) Push(x interface{}) {
    this.items = append(this.items, x.(pqitem))
}

func (this *pq) Pop() interface {} {
    old := this.items
    n := len(old)
    item := old[n-1]
    this.items = old[0:n-1]
    return item
}

func (this *pq) add(item pqitem) {
    if !this.full || this.lowest == nil || this.lowest.Less(item) {
        heap.Push(this, item)
        if this.full {
            heap.Pop(this)
        } else if this.Len() > this.n {
            this.full = true
            heap.Pop(this)
        }
        this.lowest = this.items[0]
    }
}

func (this *pq) get() []pqitem {
    a := this
    sort.Stable(sort.Reverse(a))
    return a.items
}

type pqRecommendedItem model.RecommendedItem

func (this pqRecommendedItem) Less(x interface{}) bool {
    other := x.(pqRecommendedItem)
    return this.Value < other.Value
}

func getTopItems(howMany int, possibleItemIds []uint64, e estimator) ([]model.RecommendedItem, error) {
    tops := newPQ(howMany)
    for _, itemId := range(possibleItemIds) {
        preference, err := e.estimate(itemId)
        if err != nil {
            continue
        }
        tops.add(pqRecommendedItem{
            ItemId: itemId,
            Value: preference,
        })
    }
    a := make([]model.RecommendedItem, tops.Len())
    for i, item := range(tops.get()) {
        a[i] = model.RecommendedItem(item.(pqRecommendedItem))
    }
    return a, nil
}

func getTopUsers(howMany int, userIds []uint64, e estimator) ([]uint64, error) {
    tops := newPQ(howMany)
    for _, userId := range(userIds) {
        sim, err := e.estimate(userId)
        if err != nil {
            continue
        }
        tops.add(pqRecommendedItem{
            ItemId: userId,
            Value: sim,
        })
    }
    a := make([]uint64, tops.Len())
    for i, item := range(tops.get()) {
        x := item.(pqRecommendedItem)
        a[i] = x.ItemId
    }
    return a, nil
}

