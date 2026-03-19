package utils

import (
	"strconv"
)

func ParseHymnWithVariant(raw string) (int, string, error) {
	variant := ""
	num, err := strconv.Atoi(raw)
	if err != nil {
		if len(raw) == 0 {
			return 0, "", err
		}
		num, err = strconv.Atoi(raw[0 : len(raw)-1])
		if err != nil {
			return 0, "", err

		}
		variant = string(raw[len(raw)-1])
	}

	return num, variant, nil

}
