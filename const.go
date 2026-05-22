// Package abfallkalender is a client for the insert-it.de waste-collection API
// (Müllabfuhr / Abfallkalender) used by several German municipalities:
// Mannheim, Hattingen, Herne, Kassel, Krefeld, Lübeck and Offenbach.
//
// Select a municipality by assigning Region before calling the package
// functions:
//
//	abfallkalender.Region = abfallkalender.Regions["Kassel"]
//	streets, err := abfallkalender.GetStreets()
package abfallkalender

// Region selects the municipality every API call targets. Assign it any key of
// Regions before calling the package functions. URLs are rebuilt per request,
// so switching Region at runtime takes effect immediately.
var Region = Regions["Mannheim"]

// Regions maps a municipality name to its insert-it.de webservice path segment.
var Regions = map[string]string{
	"Hattingen": "BmsAbfallkalenderHattingen/",
	"Herne":     "BmsAbfallkalenderHerne/",
	"Kassel":    "BmsAbfallkalenderKassel/",
	"Krefeld":   "BmsAbfallkalenderKrefeld/",
	"Luebeck":   "BmsAbfallkalenderLuebeck/",
	"Mannheim":  "BmsAbfallkalenderMannheim/",
	"Offenbach": "BmsAbfallkalenderOffenbach/",
}

const (
	baseurl = "https://www.insert-it.de/"
	service = "Webservice/"

	getAllStreets            = "GetAllStreets"
	getHouseNumbersForStreet = "GetHouseNumbersByStreetName?streetName=%s"
	getNextEmptyings         = "GetNextEmptyingsFromLocation?"
	getEmptyings             = "GetEmptyingsByStreetNameAndNumber?"

	servicePointTypes = "GetAllPointObjectTypes"
	servicePoints     = "GetAllPointObjects"
)

// svcURL returns the webservice base URL for the currently selected Region.
func svcURL() string { return baseurl + Region + service }

// imgURL returns the image base URL for the currently selected Region.
func imgURL() string { return baseurl + Region + "img/" }
