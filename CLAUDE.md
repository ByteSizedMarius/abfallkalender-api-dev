# abfallkalender-api

Go library and CLI for the insert-it.de waste-collection API
(Müllabfuhr / Abfallkalender). Public/unauthenticated endpoints only.

## Project structure

```
const.go            Region selection, endpoint constants, svcURL/imgURL functions
http.go             Shared http.Client with timeout + httpGet helper (status check)
models.go           TrashDate, HouseNumber (custom UnmarshalJSON for date parsing)
search.go           GetStreets, GetStreetFilter, GetHouseNumbers
trash.go            GetCalendar, GetNextEmptyings, emptyingsURL helper
service_points.go   GetServicePointTypes, GetServicePoints
example/main.go     library demo
cmd/main.go         CLI (stdlib flag, six subcommands, -json output)
```

## Supported municipalities

Seven, each addressed via a `BmsAbfallkalender<Name>/` URL segment:
Hattingen, Herne, Kassel, Krefeld, Luebeck (Lübeck), Mannheim, Offenbach.

The `Regions` map in `const.go` is the source of truth; `Region` selects which
one all package functions hit.

## Dependencies

Zero. Standard library only. Don't add deps without a strong reason.

## API host

`https://www.insert-it.de/BmsAbfallkalender<City>/Webservice/<endpoint>`

## Design decisions

### Region selection via a global var, URLs built per call
`Region` is a package-level `string`. `svcURL()` and `imgURL()` recompute the
URL on every call. An earlier version cached `svcUrl` at package init, which
silently froze every request to Mannheim regardless of subsequent reassignment.
Not concurrency-safe by design - set `Region` once at startup, or guard
yourself.

### HouseNumberEnd is `int`, not `string`
The API returns `houseNumberEnd` as a JSON number. Confirmed empirically by
running the CLI; an earlier assumption that it was a string caused unmarshal
failures on every range-style address.

### emptyingsURL sends empty (not "0") for absent end
The API treats `houseNumberEnd=0` as a literal range bound `1->0` and returns
an empty array. For single-house queries (`HouseNumberEnd == 0`) the URL sends
`houseNumberEnd=` (empty). Built via `url.Values` so encoding is handled
uniformly.

### URL encoding everywhere
All user-supplied strings going into URLs route through `url.QueryEscape` (or
`url.Values`). Street names like `Aachener Straße` contain spaces and umlauts
that produce malformed requests if interpolated raw.

### Date format
`String()` methods format dates as DD.MM.YYYY (`02.01.2006`), matching German
conventions. JSON output uses `time.Time`'s default RFC3339 marshaller for
machine-readability.

### Mixed receivers on TrashDate
`UnmarshalJSON` has a pointer receiver because it mutates; `String()` has a
value receiver so `fmt.Println` formats slice elements with it. The
`//goland:noinspection GoMixedReceiverTypes` comment suppresses the IDE
warning - the mismatch is intentional, not a smell.

### No HTTP retries
`httpGet` makes one attempt with a 30-second timeout and returns. Year-end
outages and transient failures are surfaced to the caller; this library does
not loop or back off.

## Year-end outage (recurring upstream issue)

The upstream insert-it.de API routinely breaks around the year change (roughly
the last week of December through the first week of January). Affects every
consumer including the official municipal apps. The library will return HTTP
or decode errors during that window; this is not a bug here. Service resumes
once the operator publishes the new year's schedule.

## Subagent guardrails

These apply to any Explore / Plan / general-purpose agent operating on this
project. The orchestrator should still paste them into the agent prompt, but
mirroring them here means they propagate even when the orchestrator forgets.

- Project root is `C:\Users\byte\repos\5-others\insert_it`. Stay inside it.
- Never `ls`, `find`, `dir`, `tree`, or `grep` against any parent directory:
  no `C:\Users\byte\repos\`, no `C:\Users\byte\`, no `C:\`.
- Never use `..`, `../../`, or any relative parent path.
- Use the Glob and Grep tools, NOT Bash with `ls` / `find` / `grep`.
- No `cd <path> && <cmd>` chaining; it bypasses permission detection.
- If you don't know where something lives within this project, start at the
  project root and use Glob/Grep. Do not look outside.
