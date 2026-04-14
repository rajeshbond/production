package shiftslot

import "time"

type ShiftSlot struct {
	Start time.Time
	End   time.Time
	Index int
}
