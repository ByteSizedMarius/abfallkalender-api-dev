package abfallkalender

import (
	"encoding/json"
	"fmt"
	"time"
)

// NextTrashDate is one entry returned by GetNextEmptyings. The next-emptyings
// endpoint returns a different JSON shape than the calendar endpoint - note
// the distinct field names (Name vs BmsWasteTypeName, ExecutionDate vs
// Deadline, BinSize vs Size/SizeName).
type NextTrashDate struct {
	ID                 int
	Name               string
	AppCalendarIconUrl string
	BinSize            int
	ExecutionDate      time.Time
}

func (ntd *NextTrashDate) UnmarshalJSON(data []byte) error {
	type Alias NextTrashDate
	aux := &struct {
		ExecutionDate      string
		AppCalendarIconUrl string
		*Alias
	}{
		Alias: (*Alias)(ntd),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Some entries have no scheduled next pickup; their ExecutionDate arrives
	// empty. Keep the entry with a zero time instead of failing the whole
	// decode.
	if aux.ExecutionDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05", aux.ExecutionDate)
		if err != nil {
			return err
		}
		ntd.ExecutionDate = t
	}
	if aux.AppCalendarIconUrl != "" {
		ntd.AppCalendarIconUrl = imgURL() + aux.AppCalendarIconUrl
	}

	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (ntd NextTrashDate) String() string {
	when := "-"
	if !ntd.ExecutionDate.IsZero() {
		when = ntd.ExecutionDate.Format("02.01.2006")
	}
	return fmt.Sprintf("%s  %s (%d L)", when, ntd.Name, ntd.BinSize)
}

// TrashDate is one entry returned by GetCalendar. The calendar endpoint
// returns a fuller record than the next-emptyings endpoint; see
// NextTrashDate for the other shape.
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

	// Defensive: accept empty/missing Deadline as a zero time rather than
	// failing the whole decode. Calendar entries normally have a date but the
	// API can surprise.
	if aux.Deadline != "" {
		t, err := time.Parse("2006-01-02T15:04:05", aux.Deadline)
		if err != nil {
			return err
		}
		td.Deadline = t
	}
	if aux.AppCalendarIconUrl != "" {
		td.AppCalendarIconUrl = imgURL() + aux.AppCalendarIconUrl
	}

	return nil
}

// String renders date + waste type + (size, frequency). The internal
// metadata fields (IDs, redundant Size/Cycle ints) stay in the struct for
// JSON output but are dropped from text output.
//
//goland:noinspection GoMixedReceiverTypes
func (td TrashDate) String() string {
	deadline := "-"
	if !td.Deadline.IsZero() {
		deadline = td.Deadline.Format("02.01.2006")
	}
	return fmt.Sprintf("%s  %s (%s, %s)", deadline, td.BmsWasteTypeName, td.SizeName, td.CycleAsText)
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
