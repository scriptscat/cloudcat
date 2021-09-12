package utils

import "strconv"

func Errs(err ...error) error {
	for _, v := range err {
		if v != nil {
			return v
		}
	}
	return nil
}

func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}
