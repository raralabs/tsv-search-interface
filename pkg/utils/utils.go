package utils

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func MapToJSON(input interface{}) datatypes.JSON {
	v, err := json.Marshal(input)
	if err != nil {
		return nil
	}
	return v
}

func Pagination(opts ...int) (int, int) {
	offset := 0
	limit := 15

	if len(opts) >= 2 {
		if opts[1] > 0 {
			limit = opts[1]
		}
	}

	if len(opts) >= 1 {
		if opts[0] > 1 {
			offset = (opts[0] - 1) * limit
		}
	}
	return offset, limit
}
