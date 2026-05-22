package abfallkalender

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

// GetNextEmptyings returns the next pickup per waste type at the given
// address. The endpoint's JSON shape differs from GetCalendar's - results
// come back as NextTrashDate, not TrashDate.
func GetNextEmptyings(street string, hn HouseNumber) ([]NextTrashDate, error) {
	resp, err := httpGet(emptyingsURL(getNextEmptyings, street, hn))
	if err != nil {
		return nil, err
	}
	return decodeArray[NextTrashDate](resp.Body)
}

// GetCalendar returns every waste pickup for the current year at the given
// address.
func GetCalendar(street string, hn HouseNumber) ([]TrashDate, error) {
	resp, err := httpGet(emptyingsURL(getEmptyings, street, hn))
	if err != nil {
		return nil, err
	}
	return decodeArray[TrashDate](resp.Body)
}

// emptyingsURL builds an emptyings endpoint URL with URL-encoded query
// parameters. HouseNumberEnd is sent empty (not "0") for single addresses,
// because the API treats 0 as a literal range bound and returns nothing.
func emptyingsURL(endpoint, street string, hn HouseNumber) string {
	endStr := ""
	if hn.HouseNumberEnd != 0 {
		endStr = strconv.Itoa(hn.HouseNumberEnd)
	}
	q := url.Values{
		"streetName":            {street},
		"houseNumberStart":      {strconv.Itoa(hn.HouseNumberStart)},
		"houseNumberStartExtra": {hn.HouseNumberStartExtra},
		"houseNumberEnd":        {endStr},
		"houseNumberEndExtra":   {hn.HouseNumberEndExtra},
	}
	return svcURL() + endpoint + q.Encode()
}

// decodeArray reads a JSON array from body into a slice of T, closing the
// body when done. Used by both emptyings endpoints, which return different
// element types.
func decodeArray[T any](body io.ReadCloser) ([]T, error) {
	defer body.Close()
	var result []T
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
