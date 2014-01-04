package recommender

import (
    "github.com/lilee/tasty/cf"
    "github.com/lilee/tasty/cf/model"
    "github.com/lilee/tasty/cf/similarity"
)

type NearestNUserNeighborhood struct {
    n int
    minSim float64
    samplingRate float64
    dataModel model.DataModel
    similarity similarity.UserSimilarity
}

func NewNearestNUserNeighborhood(n int, minSim float64, m model.DataModel, s similarity.UserSimilarity) *NearestNUserNeighborhood {
    return &NearestNUserNeighborhood{
        n: n,
        minSim: minSim,
        dataModel: m,
        similarity: s,
    }
}

func (this *NearestNUserNeighborhood) Neighborhoods(userId uint64) ([]uint64, error) {
    e := estimatorFunc(func(uid uint64) (float64, error) {
        if uid == userId {
            return 0.0, cf.NaNError
        }
        sim, err := this.similarity.UserSimilarity(userId, uid)
        if err != nil {
            return 0.0, err
        }
        if sim < this.minSim {
            return 0.0, cf.NaNError
        }
        return sim, nil
    })
    return getTopUsers(this.n, this.dataModel.UserIds(), e)
}

