package abfallkalender

import (
	"encoding/json"
	"fmt"
	"time"
)

// TrashDate is a single waste pickup, returned by GetCalendar and
// GetNextEmptyings.
type TrashDate struct {
	BmsWasteTypeId     int
	BmsWasteTypeName   string
	Deadline           time.Time
	AppCalendarIconUrl string
	EmptyingCountCycle int
	SizeName           string
	Size               int
	Cycle              int
	EmptyingCount      int
	CycleAsText        string
}

func (td *TrashDate) UnmarshalJSON(data []byte) error {
	type Alias TrashDate
	aux := &struct {
		Deadline           string
		AppCalendarIconUrl string
		*Alias
	}{
		Alias: (*Alias)(td),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02T15:04:05", aux.Deadline)
	if err != nil {
		return err
	}
	td.Deadline = t
	td.AppCalendarIconUrl = imgURL() + aux.AppCalendarIconUrl

	return nil
}

// String renders a TrashDate as a single line. The receiver is a value so fmt
// formats slice elements with it; UnmarshalJSON above needs a pointer receiver
// to populate the struct, hence the deliberate receiver-type mismatch.
//
//goland:noinspection GoMixedReceiverTypes
func (td TrashDate) String() string {
	return fmt.Sprintf("BmsWasteTypeId: %d, BmsWasteTypeName: %s, Deadline: %s, EmptyingCountCycle: %d, SizeName: %s, Size: %d, Cycle: %d, EmptyingCount: %d, CycleAsText: %s",
		td.BmsWasteTypeId, td.BmsWasteTypeName, td.Deadline.Format("02.01.2006"), td.EmptyingCountCycle, td.SizeName, td.Size, td.Cycle, td.EmptyingCount, td.CycleAsText)
}

// HouseNumber is a house-number range as returned by GetHouseNumbers. The API
// uses the same struct for a single number and a range; the End fields are
// zero/empty for a single number.
type HouseNumber struct {
	HouseNumberStart      int
	HouseNumberStartExtra string
	HouseNumberEnd        int
	HouseNumberEndExtra   string
}

func (hn HouseNumber) String() string {
	str := fmt.Sprintf("%d%s", hn.HouseNumberStart, hn.HouseNumberStartExtra)
	if hn.HouseNumberEnd != 0 || hn.HouseNumberEndExtra != "" {
		str += fmt.Sprintf("-%d%s", hn.HouseNumberEnd, hn.HouseNumberEndExtra)
	}
	return str
}
