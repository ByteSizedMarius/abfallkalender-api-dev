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
	"time"

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
		street, hn, size, err := addressFlags("calendar", args)
		if err != nil {
			return nil, err
		}
		data, err := abfallkalender.GetCalendar(street, hn)
		if err != nil {
			return nil, err
		}
		if size > 0 {
			data = filterCalendarBySize(data, size)
		}
		return data, nil

	case "next":
		street, hn, size, err := addressFlags("next", args)
		if err != nil {
			return nil, err
		}
		data, err := abfallkalender.GetNextEmptyings(street, hn)
		if err != nil {
			return nil, err
		}
		if size > 0 {
			data = filterNextBySize(data, size)
		}
		return data, nil

	case "pointtypes":
		return abfallkalender.GetServicePointTypes()

	case "points":
		return abfallkalender.GetServicePoints()

	default:
		usage()
		return nil, fmt.Errorf("unknown command %q", command)
	}
}

// addressFlags parses -street/-number/-size shared by calendar and next.
func addressFlags(name string, args []string) (string, abfallkalender.HouseNumber, int, error) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	street := fs.String("street", "", "street name (required)")
	number := fs.Int("number", 0, "house number (required)")
	size := fs.Int("size", 0, "filter by bin size in liters (optional)")
	_ = fs.Parse(args)
	if *street == "" {
		return "", abfallkalender.HouseNumber{}, 0, fmt.Errorf("%s: -street is required", name)
	}
	if *number <= 0 {
		return "", abfallkalender.HouseNumber{}, 0, fmt.Errorf("%s: -number is required and must be positive", name)
	}
	return *street, abfallkalender.HouseNumber{HouseNumberStart: *number}, *size, nil
}

func filterCalendarBySize(in []abfallkalender.TrashDate, size int) []abfallkalender.TrashDate {
	out := make([]abfallkalender.TrashDate, 0, len(in))
	for _, td := range in {
		if td.Size == size {
			out = append(out, td)
		}
	}
	return out
}

func filterNextBySize(in []abfallkalender.NextTrashDate, size int) []abfallkalender.NextTrashDate {
	out := make([]abfallkalender.NextTrashDate, 0, len(in))
	for _, ntd := range in {
		if ntd.BinSize == size {
			out = append(out, ntd)
		}
	}
	return out
}

// printText prints a slice line-by-line. Calendar/next entries that share the
// same date, waste type and cycle but only differ by bin size collapse into
// one line with sizes joined by "/". Other slice types use generic printing.
func printText(data any) {
	switch d := data.(type) {
	case []abfallkalender.TrashDate:
		printTrashDates(d)
	case []abfallkalender.NextTrashDate:
		printNextTrashDates(d)
	default:
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
}

func printTrashDates(in []abfallkalender.TrashDate) {
	if len(in) == 0 {
		fmt.Println("(no results)")
		return
	}
	type key struct {
		date  time.Time
		name  string
		cycle string
	}
	type sizeEntry struct {
		size    int
		display string
	}
	groups := map[key][]sizeEntry{}
	order := []key{}
	for _, td := range in {
		k := key{td.Deadline, td.BmsWasteTypeName, td.CycleAsText}
		if _, ok := groups[k]; !ok {
			order = append(order, k)
		}
		groups[k] = append(groups[k], sizeEntry{td.Size, td.SizeName})
	}
	for _, k := range order {
		entries := groups[k]
		sort.Slice(entries, func(i, j int) bool { return entries[i].size < entries[j].size })
		displays := make([]string, len(entries))
		for i, e := range entries {
			displays[i] = e.display
		}
		deadline := "-"
		if !k.date.IsZero() {
			deadline = k.date.Format("02.01.2006")
		}
		fmt.Printf("%s  %s (%s, %s)\n", deadline, k.name, strings.Join(displays, " / "), k.cycle)
	}
}

func printNextTrashDates(in []abfallkalender.NextTrashDate) {
	if len(in) == 0 {
		fmt.Println("(no results)")
		return
	}
	type key struct {
		date time.Time
		name string
	}
	groups := map[key][]int{}
	order := []key{}
	for _, ntd := range in {
		k := key{ntd.ExecutionDate, ntd.Name}
		if _, ok := groups[k]; !ok {
			order = append(order, k)
		}
		groups[k] = append(groups[k], ntd.BinSize)
	}
	for _, k := range order {
		sizes := groups[k]
		sort.Ints(sizes)
		displays := make([]string, len(sizes))
		for i, s := range sizes {
			displays[i] = fmt.Sprintf("%d L", s)
		}
		deadline := "-"
		if !k.date.IsZero() {
			deadline = k.date.Format("02.01.2006")
		}
		fmt.Printf("%s  %s (%s)\n", deadline, k.name, strings.Join(displays, " / "))
	}
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
  calendar       full pickup calendar for an address; needs -street -number; -size filters by bin liters
  next           next pickup per waste type; needs -street -number; -size filters by bin liters
  pointtypes     service-point categories
  points         all service points (glass containers, recycling, ...)

Text output collapses calendar/next entries that differ only in bin size into
one line: "26.05.2026  Rest (80 L / 240 L, 1x 14-täglich)". JSON output keeps
the raw entries intact.

Supported cities: %s

Examples:
  %s -city Kassel streets -filter Wilhelms
  %s -city Mannheim calendar -street "Katharinenstr." -number 47
  %s -city Mannheim calendar -street "Katharinenstr." -number 47 -size 240
  %s -city Mannheim -json next -street "Katharinenstr." -number 47
`, bin, bin, supportedCities(), bin, bin, bin, bin)
}
