package recommender

import (
    "github.com/lilee/tasty/cf"
    "github.com/lilee/tasty/cf/model"
    "github.com/lilee/tasty/cf/similarity"
    "github.com/lilee/tasty/cf/utils"
)

func NewGenericUserBasedRecommender(m model.DataModel, n UserNeighborhood, s similarity.UserSimilarity) Recommender {
    return &GenericUserBasedRecommender{
        dataModel: m,
        neighborhood: n,
        similarity: s,
    }
}

type GenericUserBasedRecommender struct {
    dataModel model.DataModel
    neighborhood UserNeighborhood
    similarity similarity.UserSimilarity
}

func (this *GenericUserBasedRecommender) Recommend(userId uint64, howMany int) ([]model.RecommendedItem, error) {
    neighborhoods, err := this.neighborhood.Neighborhoods(userId)
    if err != nil {
        return nil, err
    }

    itemIds, err := this.getAllOtherItemIds(neighborhoods, userId)
    if err != nil {
        return nil, err
    }

    e := this.newEstimator(userId, neighborhoods)

    return getTopItems(howMany, itemIds, e)
}

func (this *GenericUserBasedRecommender) EstimatePreference(userId, itemId uint64) (float64, error) {
    actualPref, err := this.dataModel.PreferenceValue(userId, itemId)
    if err == nil {
        return actualPref, nil
    }
    neighborhoods, err := this.neighborhood.Neighborhoods(userId)
    if err != nil {
        return 0, err
    }
    e := this.newEstimator(userId, neighborhoods)
    return e.estimate(itemId)
}

func (this *GenericUserBasedRecommender) getAllOtherItemIds(neighborhoods[]uint64, userId uint64) ([]uint64, error) {
    possibleIdSet := utils.IdSet{}
    for _, uid := range(neighborhoods) {
        userPrefs, err := this.dataModel.GetUserPreferences(uid)
        if err != nil {
            continue
        }
        possibleIdSet.AddArray(userPrefs.Ids())
    }
    userPrefs, err := this.dataModel.GetUserPreferences(userId)
    if err != nil {
        return nil, err
    }
    possibleIdSet.RemoveArray(userPrefs.Ids())
    return possibleIdSet.ToArray(), nil
}

func (this *GenericUserBasedRecommender) newEstimator(userId uint64, neighborhoods []uint64) estimator {
    return estimatorFunc(func (itemId uint64) (float64, error) {
        var preference, totalSimilarity float64
        var count int
        for _, uid := range(neighborhoods) {
            if userId == uid {
                continue
            }
            pref, err := this.dataModel.PreferenceValue(uid, itemId)
            if err != nil {
                continue
            }
            sim ,err := this.similarity.UserSimilarity(uid, userId)
            preference += sim * pref;
            totalSimilarity += sim;
            count++;
        }
        if count <= 1 {
            return 0, cf.NaNError
        }
        estimate := preference / totalSimilarity;
        estimate = capper(estimate, this.dataModel)
        return estimate, nil
    })
}
