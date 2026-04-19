package shifttiming

import (
	"fmt"
	"time"
)

// func convertToMinutes(start, end string) (int, int, error) {

// 	s, err := convertTime(start)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	e, err := convertTime(end)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	return s, e, nil
// }

// func convertTime(t string) (int, error) {
// 	parts := strings.Split(t, ":")
// 	h, _ := strconv.Atoi(parts[0])
// 	m, _ := strconv.Atoi(parts[1])
// 	return h*60 + m, nil
// }

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

func convertToMinutes(start, end string) (int, int, error) {
	s, err := time.Parse("15:04", start)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start time: %s", start)
	}

	e, err := time.Parse("15:04", end)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end time: %s", end)
	}

	return s.Hour()*60 + s.Minute(), e.Hour()*60 + e.Minute(), nil
}

func isOverlapping(existing [][2]int, newStart, newEnd int) bool {
	for _, e := range existing {
		if newStart < e[1] && newEnd > e[0] {
			return true
		}
	}
	return false
}

func (req BulkCreateShiftRequest) ValidateSingleTenant() error {
	if len(req) == 0 {
		return fmt.Errorf("empty request")
	}

	code := req[0].TenantCode

	for _, r := range req {
		if r.TenantCode != code {
			return fmt.Errorf("multiple tenant codes not allowed")
		}
	}
	return nil
}
