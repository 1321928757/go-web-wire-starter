package compo

import (
	"github.com/sony/sonyflake"
	"time"
)

// NewSonyFlake .
func NewSonyFlake() *sonyflake.Sonyflake {
	return sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime:      time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC),
		MachineID:      nil,
		CheckMachineID: nil,
	})
}
