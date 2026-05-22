package abfallkalender

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestNextTrashDateMatchesNextEndpointJSON decodes a representative sample
// from the GetNextEmptyingsFromLocation response. Earlier code pointed
// GetNextEmptyings at TrashDate, whose field names (BmsWasteTypeName,
// Deadline, Size...) do not match the next endpoint's keys (name,
// executionDate, binSize). The result was every decoded entry coming back
// all-zero - this test would have caught it.
func TestNextTrashDateMatchesNextEndpointJSON(t *testing.T) {
	Region = Regions["Mannheim"]

	payload := `[
		{"id": 1, "name": "Restmüll", "binSize": 240, "executionDate": "2026-05-22T00:00:00", "appCalendarIconUrl": "ic_rest.png"}
	]`
	var got []NextTrashDate
	if err := json.Unmarshal([]byte(payload), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d entries, want 1", len(got))
	}
	if got[0].Name != "Restmüll" {
		t.Errorf("Name: got %q, want %q", got[0].Name, "Restmüll")
	}
	if got[0].BinSize != 240 {
		t.Errorf("BinSize: got %d, want 240", got[0].BinSize)
	}
	want := time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC)
	if !got[0].ExecutionDate.Equal(want) {
		t.Errorf("ExecutionDate: got %v, want %v", got[0].ExecutionDate, want)
	}
	if !strings.HasPrefix(got[0].AppCalendarIconUrl, "https://www.insert-it.de/") {
		t.Errorf("AppCalendarIconUrl: %q (want region img URL prefix)", got[0].AppCalendarIconUrl)
	}
}

// TestTrashDateMatchesCalendarEndpointJSON decodes a representative sample
// from the GetEmptyingsByStreetNameAndNumber response. The calendar shape is
// distinct from the next-emptyings shape and must not be conflated.
func TestTrashDateMatchesCalendarEndpointJSON(t *testing.T) {
	Region = Regions["Mannheim"]

	payload := `[
		{"bmsWasteTypeId": 1, "bmsWasteTypeName": "Altpapier", "deadline": "2026-05-29T00:00:00", "sizeName": "240 Liter", "size": 240, "cycle": 8, "cycleAsText": "wöchentlich", "emptyingCount": 8, "emptyingCountCycle": 1, "appCalendarIconUrl": "ic_alt.png"}
	]`
	var got []TrashDate
	if err := json.Unmarshal([]byte(payload), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d entries, want 1", len(got))
	}
	if got[0].BmsWasteTypeName != "Altpapier" {
		t.Errorf("BmsWasteTypeName: got %q, want %q", got[0].BmsWasteTypeName, "Altpapier")
	}
	if got[0].SizeName != "240 Liter" {
		t.Errorf("SizeName: got %q, want %q", got[0].SizeName, "240 Liter")
	}
	want := time.Date(2026, 5, 29, 0, 0, 0, 0, time.UTC)
	if !got[0].Deadline.Equal(want) {
		t.Errorf("Deadline: got %v, want %v", got[0].Deadline, want)
	}
}

// TestNextTrashDateEmptyEntryDoesNotError covers the all-zero entries the
// Mannheim Katharinenstr. probe surfaced: the API may return entries with an
// empty executionDate, and the decoder must accept them with a zero time
// instead of failing the whole decode.
func TestNextTrashDateEmptyEntryDoesNotError(t *testing.T) {
	Region = Regions["Mannheim"]

	payload := `[
		{"id": 1, "name": "Restmüll", "binSize": 240, "executionDate": "2026-05-22T00:00:00"},
		{"id": 0, "name": "", "binSize": 0, "executionDate": ""}
	]`
	var got []NextTrashDate
	if err := json.Unmarshal([]byte(payload), &got); err != nil {
		t.Fatalf("unmarshal failed on empty entry: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d entries, want 2", len(got))
	}
	if !got[1].ExecutionDate.IsZero() {
		t.Errorf("ExecutionDate for empty entry should be zero, got %v", got[1].ExecutionDate)
	}
}

// TestTrashDateEmptyDeadlineDoesNotError mirrors the above for the calendar
// type.
func TestTrashDateEmptyDeadlineDoesNotError(t *testing.T) {
	Region = Regions["Mannheim"]

	payload := `[
		{"bmsWasteTypeId": 1, "bmsWasteTypeName": "Restmüll", "deadline": "2026-05-29T00:00:00"},
		{"bmsWasteTypeId": 0, "deadline": ""}
	]`
	var got []TrashDate
	if err := json.Unmarshal([]byte(payload), &got); err != nil {
		t.Fatalf("unmarshal failed on empty entry: %v", err)
	}
	if !got[1].Deadline.IsZero() {
		t.Errorf("Deadline for empty entry should be zero, got %v", got[1].Deadline)
	}
}

// TestHouseNumberEndIsInt ensures houseNumberEnd decodes from a JSON number.
// The original library declared HouseNumberEnd as string, which broke decoding
// for every range-style address.
func TestHouseNumberEndIsInt(t *testing.T) {
	payload := `[
		{"houseNumberStart": 1, "houseNumberEnd": 9},
		{"houseNumberStart": 11}
	]`
	var got []HouseNumber
	if err := json.Unmarshal([]byte(payload), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d, want 2", len(got))
	}
	if got[0].HouseNumberStart != 1 || got[0].HouseNumberEnd != 9 {
		t.Errorf("range entry: got start=%d end=%d, want 1/9", got[0].HouseNumberStart, got[0].HouseNumberEnd)
	}
	if got[1].HouseNumberEnd != 0 {
		t.Errorf("single-number entry should have HouseNumberEnd=0, got %d", got[1].HouseNumberEnd)
	}
}

// TestGetNextEmptyingsSignature is a compile-time regression: GetNextEmptyings
// must return []NextTrashDate. Reverting to []TrashDate fails to build.
func TestGetNextEmptyingsSignature(_ *testing.T) {
	var _ func(string, HouseNumber) ([]NextTrashDate, error) = GetNextEmptyings
}

// TestGetCalendarSignature is the analogous check for GetCalendar.
func TestGetCalendarSignature(_ *testing.T) {
	var _ func(string, HouseNumber) ([]TrashDate, error) = GetCalendar
}
