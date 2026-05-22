// Command example demonstrates the abfallkalender library: it lists the
// supported municipalities and queries one of them for streets, house numbers,
// the waste calendar and service points.
package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/ByteSizedMarius/abfallkalender-api"
)

func main() {
	// The API covers seven municipalities. Region selects which one every
	// call targets; reassigning it at runtime takes effect immediately.
	names := make([]string, 0, len(abfallkalender.Regions))
	for name := range abfallkalender.Regions {
		names = append(names, name)
	}
	sort.Strings(names)
	fmt.Println("Supported municipalities:", names)

	abfallkalender.Region = abfallkalender.Regions["Kassel"]

	// GetStreetFilter fetches every street in the region and filters locally
	// by a case-insensitive prefix; GetStreets returns the unfiltered list.
	streets, err := abfallkalender.GetStreetFilter("Wilhelms")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nKassel streets starting with \"Wilhelms\": %v\n", streets)
	if len(streets) == 0 {
		return
	}
	street := streets[0]

	// A street maps to one or more house-number ranges.
	houseNumbers, err := abfallkalender.GetHouseNumbers(street)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("House numbers for %s: %v\n", street, houseNumbers)
	hn := houseNumbers[0]

	// The full pickup calendar for the year at one address.
	calendar, err := abfallkalender.GetCalendar(street, hn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nCalendar for %s %s: %d pickups\n", street, hn, len(calendar))

	// The next pickup per waste type.
	next, err := abfallkalender.GetNextEmptyings(street, hn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nNext pickup per waste type:")
	for _, td := range next {
		fmt.Println(" ", td)
	}

	// Service points: glass containers, recycling centres, ...
	types, err := abfallkalender.GetServicePointTypes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nService point types:")
	for _, t := range types {
		fmt.Println(" ", t)
	}

	points, err := abfallkalender.GetServicePoints()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nTotal service points in Kassel: %d\n", len(points))
}
