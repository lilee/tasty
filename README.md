tasty是一个帮助你构建推荐系统的工具包，目前特性，测试都不完善，不能用于生产环境。但可以用于推荐系统入门学习

tasty采用了对商业友好的[Apache License v2](/LICENSE)发布

tasty目前实现了最基本的协同过滤推荐算法，目前包括：

1. 基于用户的协同过滤
2. 基于项目的协同过滤

相似度算法目前仅实现了皮尔逊相关系数，后续会不断完善，如：

1. 欧几里德距离
2. 余弦相似度
3. Tanimoto系数

# 使用

tasty参考了[Mahout](http://mahout.apache.org)的接口设计，使用方法和[Mahout](http://mahout.apache.org)比较相似。

下面是一个基于项目协同过滤的例子

```go
package main

import (
    "fmt"
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
    // A用户喜欢物品a,c
    // B用户喜欢物品a,b,c
    // C用户喜欢物品a
    // 物品a,c同时被用户喜欢，说明物品a,c的相似度很高
    // C用户喜欢a物品，则可以将c物品推荐给用户C
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
        panic(err)
    }
    for _, item := range(items) {
        fmt.Println("item:", item)
    }

    //Output:
    // item: (3, 5.000000)
}
```

- DataModel提供了存储及访问用户偏好的接口
- Similarity实现了计算两个用户相似度的方法
- Recommender使用上面两个模块向指定用户提供TopN推荐

程序运行输出是：

    item: (3, 5.000000)

表示将推荐物品c，并且系统预测用户C对物品c的可能评分是5分

# TODO

未来一个月内的计划

1. 完善文档
2. 完善单元测试
3. 实现基本的离线评测方法
3. 实现基于隐语义模型的协同过滤算法
4. 完善向量相似度算法
