package shift

import "fmt"

func toMinutes(t string) (int, error) {
	var h, m int

	_, err := fmt.Sscanf(t, "%d:%d", &h, &m)
	if err != nil {
		return 0, fmt.Errorf("invalid time format: %s", t)
	}

	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("invalid time value: %s", t)
	}

	return h*60 + m, nil
}
