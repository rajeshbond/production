package shifttiming

import (
	"strconv"
	"strings"
)

func convertToMinutes(start, end string) (int, int, error) {

	s, err := convertTime(start)
	if err != nil {
		return 0, 0, err
	}

	e, err := convertTime(end)
	if err != nil {
		return 0, 0, err
	}

	return s, e, nil
}

func convertTime(t string) (int, error) {
	parts := strings.Split(t, ":")
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	return h*60 + m, nil
}

// func convertToMinutes(start, end string) (int, int, error) {
// 	parse := func(t string) (int, error) {
// 		parts := strings.Split(t, ":")
// 		if len(parts) != 2 {
// 			return 0, fmt.Errorf("invalid time format")
// 		}
// 		h, _ := strconv.Atoi(parts[0])
// 		m, _ := strconv.Atoi(parts[1])
// 		return h*60 + m, nil
// 	}

// 	s, err := parse(start)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	e, err := parse(end)

// 	if err != nil {
// 		return 0, 0, err

// 	}

// 	return s, e, nil

// }
