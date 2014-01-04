package cf

import (
    "errors"
    "fmt"
)

type NoSuchUserError uint64
type NoSuchItemError uint64

func (e NoSuchUserError) Error() string {
    return fmt.Sprintf("No such user id: %u", e)
}

func (e NoSuchItemError) Error() string {
    return fmt.Sprintf("No such item id: %u", e)
}

var (
    NaNError = errors.New("Not a number")
)
