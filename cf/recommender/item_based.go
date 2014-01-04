package recommender

import (
    "github.com/lilee/tasty/cf/model"
    "github.com/lilee/tasty/cf/similarity"
    "github.com/lilee/tasty/cf/utils"
)

func NewGenericItemBasedRecommender(m model.DataModel, s similarity.ItemSimilarity) Recommender {
    return &GenericItemBasedRecommender{
        dataModel: m,
        similarity: s,
    }
}

type GenericItemBasedRecommender struct {
    dataModel model.DataModel
    similarity similarity.ItemSimilarity
}

func (this *GenericItemBasedRecommender) Recommend(userId uint64, howMany int) ([]model.RecommendedItem, error) {
    userPrefs, err := this.dataModel.GetUserPreferences(userId)
    if err != nil {
        return nil, err
    }
    // TODO: use CandidateItemsStrategy
    // 获取和该用户有关联的ItemIds，计算方法是：
    // 1. 获取该user的评过分的items
    // 2. 获取这些items对应的users
    // 3. 获取这些users评过分的items（除去自己评分过的项目）
    possibleItemIds, err := this.getAllOtherItemIds(userPrefs.Ids(), this.dataModel)
    if err != nil {
        return nil, err
    }
    // 构建评分器。评分器可以根据用户历史偏好信息，计算一个用户可能为一个项目评多少分
    // 传入参数为用户偏好信息和相似度计算算法
    e := this.newEstimator(userPrefs, this.similarity)
    // 获取Top条推荐结果
    results, err := getTopItems(howMany, possibleItemIds, e)
    return results, err
}

func (this *GenericItemBasedRecommender) EstimatePreference(userId, itemId uint64) (float64, error) {
    userPrefs, err := this.dataModel.GetUserPreferences(userId)
    if err != nil {
        return 0, nil
    }
    for i, iid := range(userPrefs.Ids()) {
        if iid == itemId {
            return userPrefs.Values()[i], nil
        }
    }
    e := this.newEstimator(userPrefs, this.similarity)
    return e.estimate(itemId)
}

func (this *GenericItemBasedRecommender) getAllOtherItemIds(preferredItemIds []uint64, dataModel model.DataModel) ([]uint64, error) {
    possibleIdSet := utils.IdSet{}
    for _, itemId := range(preferredItemIds) {
        itemPrefs, err := dataModel.GetItemPreferences(itemId)
        if err != nil {
            continue
        }
        for _, uid := range(itemPrefs.Ids()) {
            up, err := dataModel.GetUserPreferences(uid)
            if err != nil {
                continue
            }
            possibleIdSet.AddArray(up.Ids())
        }
    }
    possibleIdSet.RemoveArray(preferredItemIds)
    return possibleIdSet.ToArray(), nil
}

func (this *GenericItemBasedRecommender) newEstimator(userPrefs model.PreferenceArray, similarity similarity.ItemSimilarity) estimator {
    return estimatorFunc(func (itemId uint64) (float64, error) {
        similarities, err := similarity.ItemSimilarities(itemId, userPrefs.Ids())
        if err != nil {
            return 0.0, err
        }
        var preference, totalSimilarity float64
        values := userPrefs.Values()
        for i, s := range(similarities) {
            preference += s * values[i]
            totalSimilarity += s
        }
        estimate := preference / totalSimilarity;
        estimate = capper(estimate, this.dataModel)
        return estimate, nil
    })
}

