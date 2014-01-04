tasty是一个帮助你构建推荐系统的工具包，目前特性，测试都不完善，不能用于生产环境，可以用于推荐系统入门学习

tasty采用了对商业友好的[Apache License v2](/LICENSE)发布

tasty现实了协同过滤推荐算法，目前包括：

1. 基于用户的协同过滤
2. 基于项目的协同过滤

相似度算法目前仅实现了皮尔逊相关系数，后续会不断完善

# 使用

tasty参考了Mahout的接口设计，使用方法和Mahout比较相似。

下面是一个基于项目协同过滤的例子

```go
package main

import (
    "github.com/lilee/tasty/cf/model"
    "github.com/lilee/tasty/cf/recomender"
    "github.com/lilee/tasty/cf/similarity"
)

func main() {
    m := model.NewGenericDataModel()
    add := func(userId, itemId uint64, value float64) {
        p := model.NewGenericPreference(userId, itemId, value)
        m.AddPreference(p)
    }
    // A用户喜欢物品A,C
    // B用户喜欢物品A,B,C
    // C用户喜欢物品A
    // 将C物品推荐给用户C
    add(1, 1, 5.0)
    add(1, 3, 4.0)
    add(2, 1, 4.0)
    add(2, 2, 5.0)
    add(2, 3, 5.0)
    add(3, 1, 5.0)
    // 使用皮尔逊相关系数相似度算法
    s := similarity.NewPearsonCorrelationSimilarity(m)
    // 基于项目的协同过滤
    r := recommender.NewGenericItemBasedRecommender(m, s)
    items, err := r.Recommend(3, 1)
    if err != nil {
        log.Fatal(err)
    }
    for _, item := range(items) {
        log.Println("item:", item)
    }
}
```

# TODO

未来一个月内的计划

1. 完善文档
2. 完善单元测试
3. 实现基本的离线评测方法
3. 实现基于隐语义模型的协同过滤算法
4. 完善向量相似度算法
