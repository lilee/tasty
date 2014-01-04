package similarity

type Similarity interface {
    ItemSimilarity
    UserSimilarity
}

// Implementations of this interface define a notion of similarity between two items. Implementations should
// return values in the range -1.0 to 1.0, with 1.0 representing perfect similarity.
type ItemSimilarity interface {
    // Returns the degree of similarity, of two items, based on the preferences that users have expressed for
    // the items.
    ItemSimilarity(itemId1, itemId2 uint64) (float64, error)

    // A bulk-get version of ItemSimilarity
    ItemSimilarities(itemId1 uint64, itemId2 []uint64) ([]float64, error)
}

type UserSimilarity interface {
    UserSimilarity(userId1, userId2 uint64) (float64, error)
}
