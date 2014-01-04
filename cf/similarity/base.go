package similarity

import (
    "github.com/lilee/tasty/cf"
    "github.com/lilee/tasty/cf/model"
)

// Several subclasses in this package implement this method to actually compute the similarity from figures
// computed over users or items. Note that the computations in this class "center" the data, such that X and
// Y's mean are 0.
//
// Note that the sum of all X and Y values must then be 0. This value isn't passed down into the standard
// similarity computations as a result.
type resultComputer interface {
    computeResult(count int, sumXY, sumX2, sumY2, sumXYdiff2 float64) (float64, error)
}

type resultComputerFunc func(count int, sumXY, sumX2, sumY2, sumXYdiff2 float64) (float64, error)

func (f resultComputerFunc) computeResult(count int, sumXY, sumX2, sumY2, sumXYdiff2 float64) (float64, error) {
    return f(count, sumXY, sumX2, sumY2, sumXYdiff2)
}

// Base class encapsulating functionality that is common to most implementations in this package
type baseSimilarity struct {
    dataModel model.DataModel
    centerData bool
    weighted bool
    resultComputer resultComputer
}

func (this *baseSimilarity) ItemSimilarity(itemId1, itemId2 uint64) (float64, error) {
    xPrefsArray, err := this.dataModel.GetItemPreferences(itemId1)
    if err != nil {
        return 0.0, err
    }
    yPrefsArray, err := this.dataModel.GetItemPreferences(itemId2)
    if err != nil {
        return 0.0, err
    }
    xUserIds, yUserIds := xPrefsArray.Ids(), yPrefsArray.Ids()
    xValues, yValues := xPrefsArray.Values(), yPrefsArray.Values()
    xLength, yLength := xPrefsArray.Size(), yPrefsArray.Size()

    if xLength == 0 || yLength  == 0 {
        return 0.0, cf.NaNError
    }

    xUserId := xUserIds[0]
    yUserId := yUserIds[0]

    var xIndex, yIndex, count int
    var sumX, sumX2, sumY, sumY2, sumXY, sumXYdiff2 float64

    for {
        compare := int64(xUserId) - int64(yUserId)
        if compare == 0 {
            // Both users expressed a preference for the item
            x, y := xValues[xIndex], yValues[yIndex]
            sumXY += x * y
            sumX += x
            sumX2 += x * x
            sumY += y
            sumY2 += y * y
            diff := x -y
            sumXYdiff2 += diff * diff
            count++
        }
        if compare <= 0 {
            xIndex++
            if xIndex == xLength {
                break
            }
            xUserId = xUserIds[xIndex]
        }
        if compare >= 0 {
            yIndex++
            if yIndex == yLength {
                break
            }
            yUserId = yUserIds[yIndex]
        }
    }
    var result float64
    if this.centerData {
        n := float64(count)
        meanX := sumX / n
        meanY := sumY / n
        centeredSumXY := sumXY - meanY * sumX
        centeredSumX2 := sumX2 - meanX * sumX
        centeredSumY2 := sumY2 - meanY * sumY
        result, err = this.resultComputer.computeResult(count, centeredSumXY, centeredSumX2, centeredSumY2, sumXYdiff2)
    } else {
        result, err  = this.resultComputer.computeResult(count, sumXY, sumX2, sumY2, sumXYdiff2);
    }
    if err != nil {
        return 0.0, err
    }
    result = normalizeWeightResult(result, this.weighted, count, this.dataModel.NumItems());
    return result, nil
}

func (this *baseSimilarity) ItemSimilarities(itemId1 uint64, itemId2 []uint64) ([]float64, error) {
    results := make([]float64, len(itemId2))
    for i, id := range(itemId2) {
        results[i], _ = this.ItemSimilarity(itemId1, id)
    }
    return results, nil
}

func (this *baseSimilarity) UserSimilarity(userId1, userId2 uint64) (float64, error) {
    xPrefsArray, err := this.dataModel.GetUserPreferences(userId1)
    if err != nil {
        return 0.0, err
    }
    yPrefsArray, err := this.dataModel.GetUserPreferences(userId2)
    if err != nil {
        return 0.0, err
    }
    xItemIds, yItemIds := xPrefsArray.Ids(), yPrefsArray.Ids()
    xValues, yValues := xPrefsArray.Values(), yPrefsArray.Values()
    xLength, yLength := xPrefsArray.Size(), yPrefsArray.Size()

    if xLength == 0 || yLength  == 0 {
        return 0.0, cf.NaNError
    }

    xItemId := xItemIds[0]
    yItemId := yItemIds[0]

    var xIndex, yIndex, count int
    var sumX, sumX2, sumY, sumY2, sumXY, sumXYdiff2 float64

    for {
        compare := int64(xItemId) - int64(yItemId)
        if (compare == 0) {
            x, y := xValues[xIndex], yValues[yIndex]
            sumXY += x * y
            sumX += x
            sumX2 += x * x
            sumY += y
            sumY2 += y * y
            diff := x -y
            sumXYdiff2 += diff * diff
            count++
        }
        if compare <= 0 {
            xIndex++
            if xIndex == xLength {
                break
            }
            xItemId= xItemIds[xIndex]
        }
        if compare >= 0 {
            yIndex++
            if yIndex == yLength {
                break
            }
            yItemId = yItemIds[yIndex]
        }
    }
    var result float64
    if this.centerData {
        n := float64(count)
        meanX := sumX / n
        meanY := sumY / n
        centeredSumXY := sumXY - meanY * sumX
        centeredSumX2 := sumX2 - meanX * sumX
        centeredSumY2 := sumY2 - meanY * sumY
        result, err = this.resultComputer.computeResult(count, centeredSumXY, centeredSumX2, centeredSumY2, sumXYdiff2)
    } else {
        result, err  = this.resultComputer.computeResult(count, sumXY, sumX2, sumY2, sumXYdiff2);
    }
    if err != nil {
        return 0.0, err
    }
    result = normalizeWeightResult(result, this.weighted, count, this.dataModel.NumUsers());
    return result, nil
}

func normalizeWeightResult(result float64, weighted bool, count, num int) float64 {
    normalizedResult := result
    if weighted {
        scaleFactor := 1.0 - float64(count) / float64(num + 1)
        if normalizedResult < 0.0 {
            normalizedResult = -1.0 + scaleFactor * (1.0 + normalizedResult);
        } else {
            normalizedResult = 1.0 - scaleFactor * (1.0 - normalizedResult);
        }
    }
    // Make sure the result is not accidentally a little outside [-1.0, 1.0] due to rounding:
    if normalizedResult < -1.0 {
        normalizedResult = -1.0;
    } else if normalizedResult > 1.0 {
        normalizedResult = 1.0;
    }
    return normalizedResult;
}

