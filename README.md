# abfallkalender-api

Go library and CLI for the insert-it.de Abfallkalender API used by several
German municipalities to publish their waste-pickup schedules: **Mannheim,
Hattingen, Herne, Kassel, Krefeld, Lübeck, Offenbach**. Zero dependencies.

## cli

Prints a human-readable list by default, or machine-readable JSON with
`-json` for use from other languages.

**Install:**

- [Download](https://github.com/ByteSizedMarius/abfallkalender-api/releases/latest) a release
- Or install via Go: `go install github.com/ByteSizedMarius/abfallkalender-api/cmd@latest`
- Or clone and build: `go build -o abfallkalender ./cmd`

Verify: `abfallkalender --help`

```
Usage: abfallkalender [-city NAME] [-json] <command> [flags]

Commands:
  streets        list streets; "-filter PREFIX" filters by name prefix
  housenumbers   house-number ranges for a street; needs "-street"
  calendar       full pickup calendar for an address; needs "-street -number"; "-size" filters by bin liters
  next           next pickup per waste type; needs "-street -number"; "-size" filters by bin liters
  pointtypes     service-point categories
  points         all service points (glass containers, recycling, ...)
```

```
$ abfallkalender -city Kassel streets -filter Wilhelms
Wilhelmshöher Allee
Wilhelmshöher Weg
Wilhelmsstraße
Wilhelmsthaler Straße

$ abfallkalender -city Kassel calendar -street "Wilhelmshöher Allee" -number 1
29.05.2026  Altpapier (240 Liter, wöchentlich)
29.05.2026  DSD (240 Liter, 14-täglich)
30.05.2026  Bioabfall (240 Liter, 14-täglich)
30.05.2026  Restabfall (80 Liter / 120 Liter / 240 Liter, 14-täglich)
...

$ abfallkalender -city Kassel -json pointtypes
[
  { "ID": 1, "AppDisplayName": "Recyclinghof" },
  { "ID": 2, "AppDisplayName": "Altglas" }
]
```

## library

```
go get github.com/ByteSizedMarius/abfallkalender-api
```

```go
package main

import (
	"fmt"
	"log"

	"github.com/ByteSizedMarius/abfallkalender-api"
)

func main() {
	abfallkalender.Region = abfallkalender.Regions["Kassel"]

	streets, err := abfallkalender.GetStreetFilter("Wilhelms")
	if err != nil {
		log.Fatal(err)
	}

	houseNumbers, _ := abfallkalender.GetHouseNumbers(streets[0])
	calendar, _ := abfallkalender.GetCalendar(streets[0], houseNumbers[0])
	for _, pickup := range calendar {
		fmt.Println(pickup)
	}
}
```

See [`example/main.go`](https://github.com/ByteSizedMarius/abfallkalender-api-dev/blob/master/example/main.go) for a more complete demo that exercises every endpoint.

### public API

- `GetStreets() ([]string, error)` - every street in the selected region
- `GetStreetFilter(prefix string) ([]string, error)` - case-insensitive prefix filter
- `GetHouseNumbers(street string) ([]HouseNumber, error)` - house-number ranges
- `GetCalendar(street string, hn HouseNumber) ([]TrashDate, error)` - full year of pickups
- `GetNextEmptyings(street string, hn HouseNumber) ([]NextTrashDate, error)` - next pickup per waste type
- `GetServicePointTypes() ([]PointObjectType, error)` - container / Recyclinghof categories
- `GetServicePoints() ([]PointObject, error)` - all service points with coordinates

Switch municipality at any time by reassigning `Region`:

```go
abfallkalender.Region = abfallkalender.Regions["Mannheim"]
```

The [`Regions`](https://github.com/ByteSizedMarius/abfallkalender-api-dev/blob/master/const.go#L18)
map keys match the city names verbatim, except Lübeck uses the ASCII spelling
`Luebeck`.

## known issues

- **Year-end outage:** the upstream API breaks around the year change (last
  week of December through the first week of January). Affects all consumers
  including the official municipal apps. Service resumes once the operator
  publishes the new year.
- **Empty calendar for a valid-looking address:** make sure the `HouseNumber`
  comes from `GetHouseNumbers` for the matching street. The API silently
  returns an empty array for addresses it cannot resolve.

## disclaimer

This project is not affiliated with [insert Infotech GmbH](https://insert-infotech.de/)
or any of the listed municipalities. It is an unofficial client built from publicly observable
network traffic against undocumented endpoints, which may change without
notice.
