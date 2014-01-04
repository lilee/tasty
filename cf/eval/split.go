package eval

import (
    "math/rand"
    "github.com/lilee/tasty/cf/model"
)

func SplitTrainingAndTest(m model.DataModel, trainingPercentage, evaluationPercentage float32) (model.PreferenceArrayMap, model.PreferenceArrayMap) {
    trainings := model.NewUserPreferenceArrayMap()
    tests := model.NewUserPreferenceArrayMap()
    for _, uid := range(m.UserIds()) {
        if rand.Float32() < evaluationPercentage {
            prefs, err := m.GetUserPreferences(uid)
            if err != nil {
                continue
            }
            for _, p := range(prefs.Raw()) {
                if rand.Float32() < trainingPercentage {
                    trainings.Set(p)
                } else {
                    tests.Set(p)
                }
            }
        }
    }
    return trainings, tests
}

