package convert

import "strconv"

func StringToUInt64(str string) (uint64, error) {
	uInt64, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return uInt64, nil
}
