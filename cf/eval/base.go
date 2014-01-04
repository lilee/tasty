package eval

import (
    "github.com/lilee/tasty/cf/model"
    "github.com/lilee/tasty/cf/recommender"
)

type baseRecommenderEvaluator struct {
}

func (this *baseRecommenderEvaluator) Evaluate(recommenderBuilder RecommenderBuilder,
                                               dataModelBuilder DataModelBuilder,
                                               dataModel model.DataModel,
                                               trainingPercentage, evaluationPercentage float32) (float64, error) {
    trainingPrefs, testPrefs := SplitTrainingAndTest(dataModel, trainingPercentage, evaluationPercentage)
    trainingModel, err := dataModelBuilder.Build(trainingPrefs)
    if err != nil {
        return 0, err
    }
    recommender, err := recommenderBuilder.Build(trainingModel)
    if err != nil {
        return 0, err
    }
    return this.getEvaluation(testPrefs, recommender)
}

func (this *baseRecommenderEvaluator) getEvaluation(tests model.PreferenceArrayMap, r recommender.Recommender) (float64, error) {
    for uid, prefs := range(tests.Raw()) {
        go func() {
            for _, p := range(prefs.Raw()) {
                estimatedPreference, err := r.EstimatePreference(uid, p.ItemId())
                if err != nil {
                } else {
                    //TODO: capper
                    this.processOneEstimate(estimatedPreference, p)
                }
            }
        }()
    }
    return 0, nil
}

func (this *baseRecommenderEvaluator) processOneEstimate(estimatedPreference float64, pref model.Preference) {
}
