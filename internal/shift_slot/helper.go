package shiftslot

import "time"

func GenerateSlots(start, end time.Time) []ShiftSlot {
	var slots []ShiftSlot

	// Overnight shift
	if end.Before(start) {
		end = end.Add(24 * time.Hour)
	}

	current := start
	index := 1

	for current.Before(end) {
		next := current.Add(time.Hour)

		if next.After(end) {
			next = end
		}

		slots = append(slots, ShiftSlot{
			Start: current,
			End:   next,
			Index: index,
		})

		current = next
		index++
	}

	return slots
}
