package model

import (
    "bufio"
    "errors"
    "io"
    "sort"
    "strconv"
    "strings"

    "github.com/lilee/tasty/cf"
    "github.com/lilee/tasty/cf/utils"
)

type genericDataModelImpl struct {
    userPreferenceMap userPreferenceArrayMap
    itemPreferenceMap itemPreferenceArrayMap
    userSorted utils.IdSet
    itemSorted utils.IdSet
    maxPreferenceValue float64
    minPreferenceValue float64
}

// 一个简单的DataModel实现，用户可以通过AddPreference方法向其中添加用户偏好，
// 该实现主要用于数据量很小的例子代码中
func NewGenericDataModel() DataModel {
    return &genericDataModelImpl {
        userPreferenceMap: *NewUserPreferenceArrayMap(),
        itemPreferenceMap: *NewItemPreferenceArrayMap(),
        userSorted: utils.NewIdSet(),
        itemSorted: utils.NewIdSet(),
    }
}

func (this *genericDataModelImpl) AddPreference(p Preference) {
    // insert preference into user and item preferences map
    this.userPreferenceMap.Set(p)
    this.itemPreferenceMap.Set(p)

    value := p.Value()
    // set max & min preference value
    if this.NumItems() == 0 {
        this.maxPreferenceValue = value
        this.minPreferenceValue = value
    }
    if value > this.maxPreferenceValue {
        this.maxPreferenceValue = value
    }
    if value < this.minPreferenceValue {
        this.minPreferenceValue = value
    }
}

func (this *genericDataModelImpl) GetUserPreferences(userId uint64) (PreferenceArray, error) {
    if prefs, ok := this.userPreferenceMap[userId]; ok {
        if (!this.userSorted.Contains(userId)) {
            sort.Sort(prefs)
            prefs.buildCache()
            this.userPreferenceMap[userId] = prefs
            this.userSorted.Add(userId)
        }
        return prefs, nil
    }
    return nil, cf.NoSuchUserError(userId)
}

func (this *genericDataModelImpl) GetItemPreferences(itemId uint64) (PreferenceArray, error) {
    if prefs, ok := this.itemPreferenceMap[itemId]; ok {
        if (!this.itemSorted.Contains(itemId)) {
            sort.Sort(prefs)
            prefs.buildCache()
            this.itemPreferenceMap[itemId] = prefs
            this.itemSorted.Add(itemId)
        }
        return prefs, nil
    }
    return nil, cf.NoSuchItemError(itemId)
}

func (this *genericDataModelImpl) PreferenceValue(userId, itemId uint64) (float64, error) {
    prefs, err := this.GetUserPreferences(userId)
    if err != nil {
        return 0, err
    }
    values := prefs.Values()
    for i, id := range(prefs.Ids()) {
        if id == itemId {
            return values[i], nil
        }
    }
    return 0.0, cf.NoSuchItemError(itemId)
}

func (this genericDataModelImpl) MaxPreferenceValue() float64 {
    return this.maxPreferenceValue
}

func (this genericDataModelImpl) MinPreferenceValue() float64 {
    return this.minPreferenceValue
}

func (this *genericDataModelImpl) UserIds() []uint64 {
    // TODO: build sorted ids in cache
    ids := make(utils.IdSlice, this.NumUsers())
    i := 0
    for id := range(this.userPreferenceMap) {
        ids[i] = id
        i++
    }
    ids.Sort()
    return ids
}

func (this genericDataModelImpl) NumUsers() int {
    return len(this.userPreferenceMap)
}

func (this *genericDataModelImpl) ItemIds() []uint64 {
    // TODO: build sorted ids in cache
    ids := make(utils.IdSlice, this.NumItems())
    i := 0
    for id := range(this.itemPreferenceMap) {
        ids[i] = id
        i++
    }
    ids.Sort()
    return ids
}

func (this genericDataModelImpl) NumItems() int {
    return len(this.itemPreferenceMap)
}

type fileDataModelImpl struct {
    DataModel
    file io.Reader
    sep string
}

/*
读取文件内容构建DataModel的实现，他要求的文件格式为：每行一条用户偏好信息，
每条偏好信息为(UserId,ItemId,PreferenceValue)的三元组，使用逗号分隔，如：

    1,24,0.3
    1,17,0.4
    2,6,0.5

*/
func NewFileDataModel(file io.Reader, sep string) (DataModel, error) {
    m := &fileDataModelImpl{
        DataModel: NewGenericDataModel(),
        file: file,
        sep: sep,
    }
    if err := m.buildModel(); err != nil {
        return nil, err
    }
    return m, nil
}

func (this *fileDataModelImpl) buildModel() error {
    s := bufio.NewScanner(this.file)
    var err error
    for s.Scan() {
        // TODO: filter comment line and empty line
        parts := strings.Split(s.Text(), this.sep)
        if len(parts) < 3 {
            return errors.New("invalid file")
        }
        var userId, itemId uint64
        var value float64
        if userId, err = strconv.ParseUint(parts[0], 10, 0); err != nil {
            continue
        }
        if itemId, err = strconv.ParseUint(parts[1], 10, 0); err != nil {
            continue
        }
        if value, err = strconv.ParseFloat(parts[2], 0); err != nil {
            continue
        }
        p := NewGenericPreference(userId, itemId, value)
        this.AddPreference(p)
    }
    return nil
}

