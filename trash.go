package abfallkalender

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

// GetNextEmptyings returns the next pickup date for each waste type at the
// given address.
func GetNextEmptyings(street string, hn HouseNumber) ([]TrashDate, error) {
	resp, err := httpGet(emptyingsURL(getNextEmptyings, street, hn))
	if err != nil {
		return nil, err
	}
	return emptyingsFromBody(resp.Body)
}

// GetCalendar returns every waste pickup for the current year at the given
// address.
func GetCalendar(street string, hn HouseNumber) ([]TrashDate, error) {
	resp, err := httpGet(emptyingsURL(getEmptyings, street, hn))
	if err != nil {
		return nil, err
	}
	return emptyingsFromBody(resp.Body)
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

func emptyingsFromBody(body io.ReadCloser) ([]TrashDate, error) {
	defer body.Close()
	var trashDates []TrashDate
	if err := json.NewDecoder(body).Decode(&trashDates); err != nil {
		return nil, err
	}
	return trashDates, nil
}
