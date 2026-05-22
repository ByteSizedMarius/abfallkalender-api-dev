package abfallkalender

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// GetStreets returns the names of all streets in the currently selected Region.
func GetStreets() ([]string, error) {
	resp, err := httpGet(svcURL() + getAllStreets)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []struct{ Name string }
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	streets := make([]string, len(data))
	for i, item := range data {
		streets[i] = item.Name
	}

	return streets, nil
}

// GetStreetFilter returns the streets in the current Region whose name starts
// with prefix, matched case-insensitively.
func GetStreetFilter(prefix string) ([]string, error) {
	streets, err := GetStreets()
	if err != nil {
		return nil, err
	}

	prefix = strings.ToLower(prefix)
	var result []string
	for _, street := range streets {
		if strings.HasPrefix(strings.ToLower(street), prefix) {
			result = append(result, street)
		}
	}

	return result, nil
}

// GetHouseNumbers returns the house-number ranges registered for the given
// street in the current Region.
func GetHouseNumbers(street string) ([]HouseNumber, error) {
	resp, err := httpGet(fmt.Sprintf(svcURL()+getHouseNumbersForStreet, url.QueryEscape(street)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var houseNumbers []HouseNumber
	if err = json.NewDecoder(resp.Body).Decode(&houseNumbers); err != nil {
		return nil, err
	}

	if len(houseNumbers) == 0 {
		return nil, fmt.Errorf("no house numbers found for street %s", street)
	}

	return houseNumbers, nil
}
