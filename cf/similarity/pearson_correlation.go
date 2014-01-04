package similarity

import (
    "math"
    "github.com/lilee/tasty/cf"
    "github.com/lilee/tasty/cf/model"
)

/*
An implementation of the Pearson correlation. For users X and Y, the following values are calculated:

      sumX2: sum of the square of all X's preference values
      sumY2: sum of the square of all Y's preference values
      sumXY: sum of the product of X and Y's preference value for all items for which both X and Y express a preference

The correlation is then:

      sumXY / sqrt(sumX2 * sumY2)

Note that this correlation "centers" its data, shifts the user's preference values so that each of their
means is 0. This is necessary to achieve expected behavior on all data sets.

This correlation implementation is equivalent to the cosine similarity since the data it receives
is assumed to be centered -- mean is 0. The correlation may be interpreted as the cosine of the angle
between the two vectors defined by the users' preference values.
*/
type PearsonCorrelationSimilarity struct {
    baseSimilarity
}

func NewPearsonCorrelationSimilarity(m model.DataModel) Similarity {
    s := &PearsonCorrelationSimilarity{
        baseSimilarity{
            dataModel: m,
            centerData: true,
            weighted: false,
            resultComputer: resultComputerFunc(computePearsonCorrelationResult),
        },
    }
    return s
}

func computePearsonCorrelationResult(n int, sumXY, sumX2, sumY2, sumXYdiff2 float64) (float64, error) {
    if n == 0 {
        return 0.0, cf.NaNError
    }
    denominator := math.Sqrt(sumX2) * math.Sqrt(sumY2);
    if denominator == 0.0 {
        return 0.0, cf.NaNError
    }
    return sumXY / denominator, nil
}

