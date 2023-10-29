package suckutils

import (
	"fmt"
)

type WithId[T comparable] interface {
	GetId() T
}

func ToMap[T WithId[IdT], IdT comparable](a []T) (map[IdT]T, error) {
	result := make(map[IdT]T, len(a))
	for _, item := range a {
		id := item.GetId()
		if _, ok := result[id]; ok {
			return nil, fmt.Errorf("Dublicate Id=%v", id)
		}
		result[id] = item
	}
	return result, nil
}
