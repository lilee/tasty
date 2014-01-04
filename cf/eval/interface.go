package eval

import (
    "github.com/lilee/tasty/cf/model"
    "github.com/lilee/tasty/cf/recommender"
)

type RecommenderEvaluator interface {
    Evaluate(recommenderBuilder RecommenderBuilder, dataModelBuilder DataModelBuilder,
             dataModel model.DataModel, trainingPercentage, evaluationPercentage float32) (float64, error)
}

type RecommenderBuilder interface {
    Build(dataModel model.DataModel) (recommender.Recommender, error)
}

type RecommenderBuilderFunc func(m model.DataModel) (recommender.Recommender, error)

func (f RecommenderBuilderFunc) Build(m model.DataModel) (recommender.Recommender, error) {
    return f(m)
}

type DataModelBuilder interface {
    Build(m model.PreferenceArrayMap) (model.DataModel, error)
}

type DataModelBuilderFunc func(m model.PreferenceArrayMap) (model.DataModel, error)

func (f DataModelBuilderFunc) Build(m model.PreferenceArrayMap) (model.DataModel, error) {
    return f(m)
}
