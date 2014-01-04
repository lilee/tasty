package model

// Implementations represent a repository of information about users and their associated Preferences for items.
type DataModel interface {
    // Add a particular preference for a user.
    AddPreference(p Preference)

    // Get user's preferences, ordered by item id
    GetUserPreferences(userId uint64) (PreferenceArray, error)

    // Get all existing preferences expressed for that item, ordered by user id, as an array
    GetItemPreferences(itemId uint64) (PreferenceArray, error)

    PreferenceValue(userId, itemId uint64) (float64, error)

    // The maximum preference value that is possible in the current problem domain being evaluated.
    // For example, if the domain is movie ratings on a scale of 1 to 5, this should be 5. While a
    // Recommender may estimate a preference value above 5.0, it isn't "fair" to consider that
    // the system is actually suggesting an impossible rating of, say, 5.4 stars.
    // In practice the application would cap this estimate to 5.0. Since evaluators evaluate
    // the difference between estimated and actual value, this at least prevents this effect from unfairly
    // penalizing a Recommender
    MaxPreferenceValue() float64

    // See MaxPreferenceValue
    MinPreferenceValue() float64

    // All item ids array in the model
    ItemIds() []uint64

    // Get total number of items known to the model. This is generally the union of all items preferred by
    NumItems() int

    // All user ids array in the model
    UserIds() []uint64

    // Get total number of users known to the model.
    NumUsers() int
}

type Preference interface {
    // Id of user who prefers the item
    UserId() uint64

    // Item id that is preferred
    ItemId() uint64

    // Strength of the preference for that item. Zero should indicate "no preference either way";
    // positive values indicate preference and negative values indicate dislike
    Value() float64
}

// An alternate representation of an array of Preference.
// Implementations, in theory, can produce a more memory-efficient representation.
type PreferenceArray interface {
    // Get a materialized Preference representation of the preference at i
    Get(i int) Preference

    // Size of preference arary
    Size() int

    // Get user or item id
    Id() uint64

    // Get all user or item ids
    Ids() []uint64

    // Get preferences values
    Values() []float64

    Raw() []Preference
}

type PreferenceArrayMap interface {
    Set(p Preference)
    Raw() map[uint64]PreferenceArray
}

type RecommendedItem struct {
    ItemId uint64
    Value float64
}

const (
    MinPreferenceValue = -1000.0
)
