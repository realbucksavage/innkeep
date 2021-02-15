package innkeep

import (
	"errors"
	"strconv"
)

var (
	ErrNoSuchKey = errors.New("no such key")
)

type MetadataMap map[string]string

func (m MetadataMap) GetInt(key string) (int, error) {
	val, ok := m[key]
	if !ok {
		return 0, ErrNoSuchKey
	}

	if res, err := strconv.Atoi(val); err != nil {
		return 0, err
	} else {
		return res, nil
	}
}
