package recommender

import (
    "github.com/lilee/tasty/cf/model"
)

type Recommender interface {
    Recommend(userId uint64, howMany int) ([]model.RecommendedItem, error)
    EstimatePreference(userId, itemId uint64) (float64, error)
}

/*
实现这个接口，可以计算一个指定用户的邻居用户，邻居用户信息帮助实现基于用户的协同过滤推荐
*/
type UserNeighborhood interface {
    // 获取指定用户的邻居用户ID
    Neighborhoods(userId uint64) ([]uint64, error)
}

