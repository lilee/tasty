package model

import (
    "fmt"
)

type GenericPreference struct {
    userId uint64
    itemId uint64
    value float64
}

func NewGenericPreference(userId, itemId uint64, value float64) *GenericPreference {
    return &GenericPreference{
        userId: userId,
        itemId: itemId,
        value: value,
    }
}

func (this GenericPreference) UserId() uint64 {
    return this.userId
}

func (this GenericPreference) ItemId() uint64 {
    return this.itemId
}

func (this GenericPreference) Value() float64 {
    return this.value
}

type basePreferenceArray struct {
    data []Preference
    id uint64
    ids []uint64
    values []float64
}

func (this basePreferenceArray) Size() int {
    return len(this.data)
}

func (this *basePreferenceArray) Get(i int) Preference {
    return this.data[i]
}

func (this *basePreferenceArray) Id() uint64 {
    return this.id
}

func (this *basePreferenceArray) Ids() []uint64 {
    return this.ids
}

func (this *basePreferenceArray) Values() []float64 {
    return this.values
}

func (this *basePreferenceArray) Raw() []Preference {
    return this.data
}

func (this *basePreferenceArray) add(p Preference) {
    this.data = append(this.data, p)
}

// impletation for sort.Interface
func (p basePreferenceArray) Len() int { return len(p.data) }
func (p basePreferenceArray) Swap(i, j int) {
    p.data[i], p.data[j] = p.data[j], p.data[i]
}

type userPreferenceArray struct {
    basePreferenceArray
}

func newUserPreferenceArray() *userPreferenceArray {
    return &userPreferenceArray{
        basePreferenceArray{
            data: []Preference{},
        },
    }
}

func (this *userPreferenceArray) buildCache() {
    this.id = this.data[0].UserId()
    n := this.Size()
    this.ids = make([]uint64, n)
    this.values = make([]float64, n)
    for i, p := range(this.data) {
        this.ids[i] = p.ItemId()
        this.values[i] = p.Value()
    }
}

func (this *userPreferenceArray) Less(i, j int) bool {
    return this.data[i].ItemId() < this.data[j].ItemId()
}

type itemPreferenceArray struct {
    basePreferenceArray
}

func newItemPreferenceArray() *itemPreferenceArray {
    return &itemPreferenceArray{
        basePreferenceArray{
            data: []Preference{},
        },
    }
}

func (this *itemPreferenceArray) buildCache() {
    this.id = this.data[0].ItemId()
    n := this.Size()
    this.ids = make([]uint64, n)
    this.values = make([]float64, n)
    for i, p := range(this.data) {
        this.ids[i] = p.UserId()
        this.values[i] = p.Value()
    }
}

func (this *itemPreferenceArray) Less(i, j int) bool {
    return this.data[i].UserId() < this.data[j].UserId()
}

type userPreferenceArrayMap map[uint64]*userPreferenceArray

func NewUserPreferenceArrayMap() *userPreferenceArrayMap {
    return &userPreferenceArrayMap{}
}

func (this *userPreferenceArrayMap) Set(p Preference) {
    id := p.UserId()
    prefs, ok := (*this)[id]
    if !ok {
        prefs = newUserPreferenceArray()
        (*this)[id] = prefs
    }
    prefs.add(p)
}

func (this *userPreferenceArrayMap) Raw() map[uint64]PreferenceArray {
    a := make(map[uint64]PreferenceArray)
    for k, v := range(*this) {
        a[k] = v
    }
    return a
}

type itemPreferenceArrayMap map[uint64]*itemPreferenceArray

func NewItemPreferenceArrayMap() *itemPreferenceArrayMap {
    return &itemPreferenceArrayMap{}
}

func (this *itemPreferenceArrayMap) Set(p Preference) {
    id := p.ItemId()
    prefs, ok := (*this)[id]
    if !ok {
        prefs = newItemPreferenceArray()
        (*this)[id] = prefs
    }
    prefs.add(p)
}

func (this *itemPreferenceArrayMap) Raw() map[uint64]PreferenceArray {
    a := make(map[uint64]PreferenceArray)
    for k, v := range(*this) {
        a[k] = v
    }
    return a
}

func (this RecommendedItem) String() string {
    return fmt.Sprintf("(%d, %f)", this.ItemId, this.Value)
}


