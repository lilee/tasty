package main

import (
    "bufio"
    "errors"
    "flag"
    "log"
    "os"
    "strings"
    "strconv"
    "github.com/lilee/tasty/cf/model"
    "github.com/lilee/tasty/cf/recommender"
    "github.com/lilee/tasty/cf/similarity"
)

var (
    userId = flag.Uint64("u", 25, "user id")
    howMany = flag.Int("n", 10, "how many recommended items")
    what = flag.String("t", "user", "")
)

func main() {
    flag.Parse()

    itemInfos, err := loadMovieData()
    if err != nil {
        log.Fatal(err)
    }

    file, err := os.Open("ratings.dat")
    if err != nil {
        log.Fatal(err)
    }
    m, err := model.NewFileDataModel(file, "::")
    if err != nil {
        log.Fatal(err)
    }
    s := similarity.NewPearsonCorrelationSimilarity(m)

    log.Println("User", *userId, "preferences: ")
    userPrefs, _ := m.GetUserPreferences(*userId)
    for i, itemId := range(userPrefs.Ids()) {
        itemName, ok := itemInfos[itemId]
        if !ok {
            continue
        }
        log.Println(itemId, itemName, userPrefs.Values()[i])
    }

    if *what == "user" {
        log.Println("<UserCF>")
        n := recommender.NewNearestNUserNeighborhood(100, -100.0, m, s)
        ur := recommender.NewGenericUserBasedRecommender(m, n, s)
        rItems, err := ur.Recommend(*userId, *howMany)
        if err != nil {
            log.Fatal(err)
        }
        for _, item := range(rItems) {
            itemName, ok := itemInfos[item.ItemId]
            if !ok {
                continue
            }
            log.Println("item:", item, itemName)
        }
    }

    if *what == "item" {
        log.Println("<ItemCF>")
        r := recommender.NewGenericItemBasedRecommender(m, s)
        rItems, err := r.Recommend(*userId, *howMany)
        if err != nil {
            log.Fatal(err)
        }
        for _, item := range(rItems) {
            itemName, ok := itemInfos[item.ItemId]
            if !ok {
                continue
            }
            log.Println("item:", item, itemName)
        }
    }
}

func simpleItemCF() {
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
    s := similarity.NewPearsonCorrelationSimilarity(m)
    r := recommender.NewGenericItemBasedRecommender(m, s)
    rItems, err := r.Recommend(3, 1)
    if err != nil {
        log.Fatal(err)
    }
    for _, item := range(rItems) {
        log.Println("item:", item)
    }
}

func loadMovieData() (map[uint64]string, error) {
    m := map[uint64]string{}

    file, err := os.Open("movies.dat")
    if err != nil {
        return nil, err
    }

    s := bufio.NewScanner(file)
    for s.Scan() {
        parts := strings.Split(s.Text(), "::")
        if len(parts) < 2 {
            return nil, errors.New("invalid file")
        }
        itemId, err := strconv.ParseUint(parts[0], 10, 0)
        if err != nil {
            continue
        }
        m[itemId] = parts[1] + ":" + parts[2]
    }
    return m, nil
}
