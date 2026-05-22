// Command abfallkalender is a CLI for the abfallkalender library. It queries
// the insert-it.de waste-collection API (Müllabfuhr / Abfallkalender) for one
// of seven German municipalities.
//
// Every command prints a human-readable list by default, or machine-readable
// JSON with -json, which makes the tool easy to call from other languages by
// shelling out and parsing stdout.
//
// Exit status is 0 on success and non-zero on error.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/ByteSizedMarius/abfallkalender-api"
)

func main() {
	city := flag.String("city", "Mannheim", "municipality to query")
	jsonOut := flag.Bool("json", false, "output JSON instead of text")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 || flag.Arg(0) == "help" {
		usage()
		os.Exit(0)
	}

	region, ok := abfallkalender.Regions[*city]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown city %q; supported: %s\n", *city, supportedCities())
		os.Exit(1)
	}
	abfallkalender.Region = region

	data, err := dispatch(flag.Arg(0), flag.Args()[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if *jsonOut {
		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
		return
	}
	printText(data)
}

// dispatch runs one subcommand and returns its result for formatting.
func dispatch(command string, args []string) (any, error) {
	switch command {
	case "streets":
		fs := flag.NewFlagSet("streets", flag.ExitOnError)
		filter := fs.String("filter", "", "case-insensitive name prefix")
		_ = fs.Parse(args)
		if *filter == "" {
			return abfallkalender.GetStreets()
		}
		return abfallkalender.GetStreetFilter(*filter)

	case "housenumbers":
		fs := flag.NewFlagSet("housenumbers", flag.ExitOnError)
		street := fs.String("street", "", "street name (required)")
		_ = fs.Parse(args)
		if *street == "" {
			return nil, fmt.Errorf("housenumbers: -street is required")
		}
		return abfallkalender.GetHouseNumbers(*street)

	case "calendar":
		street, hn, err := addressFlags("calendar", args)
		if err != nil {
			return nil, err
		}
		return abfallkalender.GetCalendar(street, hn)

	case "next":
		street, hn, err := addressFlags("next", args)
		if err != nil {
			return nil, err
		}
		return abfallkalender.GetNextEmptyings(street, hn)

	case "pointtypes":
		return abfallkalender.GetServicePointTypes()

	case "points":
		return abfallkalender.GetServicePoints()

	default:
		usage()
		return nil, fmt.Errorf("unknown command %q", command)
	}
}

// addressFlags parses the -street/-number pair shared by calendar and next.
func addressFlags(name string, args []string) (string, abfallkalender.HouseNumber, error) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	street := fs.String("street", "", "street name (required)")
	number := fs.Int("number", 0, "house number (required)")
	_ = fs.Parse(args)
	if *street == "" {
		return "", abfallkalender.HouseNumber{}, fmt.Errorf("%s: -street is required", name)
	}
	if *number <= 0 {
		return "", abfallkalender.HouseNumber{}, fmt.Errorf("%s: -number is required and must be positive", name)
	}
	return *street, abfallkalender.HouseNumber{HouseNumberStart: *number}, nil
}

// printText prints a slice one element per line, or any other value directly.
func printText(data any) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			fmt.Println("(no results)")
			return
		}
		for i := 0; i < v.Len(); i++ {
			fmt.Println(v.Index(i).Interface())
		}
		return
	}
	fmt.Println(data)
}

func supportedCities() string {
	names := make([]string, 0, len(abfallkalender.Regions))
	for name := range abfallkalender.Regions {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func usage() {
	bin := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
	fmt.Printf(`%s - waste-collection schedules for German municipalities (insert-it.de)

Usage: %s [-city NAME] [-json] <command> [flags]

Global flags:
  -city NAME   municipality to query (default: Mannheim)
  -json        emit JSON instead of text (for calling from other languages)

Commands:
  streets        list streets; -filter PREFIX filters by name prefix
  housenumbers   house-number ranges for a street; needs -street
  calendar       full pickup calendar for an address; needs -street -number
  next           next pickup per waste type; needs -street -number
  pointtypes     service-point categories
  points         all service points (glass containers, recycling, ...)

Supported cities: %s

Examples:
  %s -city Kassel streets -filter Wilhelms
  %s -city Kassel housenumbers -street "Wilhelmshöher Allee"
  %s -city Mannheim -json calendar -street "Aachener Straße" -number 1
`, bin, bin, supportedCities(), bin, bin, bin)
}
