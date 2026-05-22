package abfallkalender

import (
	"encoding/json"
	"fmt"
)

// PointObjectType is a category of service point, e.g. "Altglascontainer".
type PointObjectType struct {
	ID             int
	AppDisplayName string
}

func (pt PointObjectType) String() string {
	return fmt.Sprintf("ID: %d, %s", pt.ID, pt.AppDisplayName)
}

// PointObject is a single service point (glass container, recycling centre,
// etc.) with its coordinates.
type PointObject struct {
	ID                   int
	Lat                  float64
	Lon                  float64
	Remark               *string
	BmsPointObjectTypeId int
}

func (p PointObject) String() string {
	remark := ""
	if p.Remark != nil {
		remark = ", Remark: " + *p.Remark
	}
	return fmt.Sprintf("ID: %d, TypeId: %d, Lat: %g, Lon: %g%s", p.ID, p.BmsPointObjectTypeId, p.Lat, p.Lon, remark)
}

// GetServicePointTypes returns the service point categories available in the
// current Region.
func GetServicePointTypes() ([]PointObjectType, error) {
	resp, err := httpGet(svcURL() + servicePointTypes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var types []PointObjectType
	if err = json.NewDecoder(resp.Body).Decode(&types); err != nil {
		return nil, err
	}
	return types, nil
}

// GetServicePoints returns all service points in the current Region.
func GetServicePoints() ([]PointObject, error) {
	resp, err := httpGet(svcURL() + servicePoints)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var points []PointObject
	if err = json.NewDecoder(resp.Body).Decode(&points); err != nil {
		return nil, err
	}
	return points, nil
}
