package utils

func Errs(err ...error) error {
	for _, v := range err {
		if v != nil {
			return v
		}
	}
	return nil
}
