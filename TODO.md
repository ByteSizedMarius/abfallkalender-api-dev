# TODO

Open items the user opted to do later. The in-repo rename
(`insert_it` -> `abfallkalender-api`, Go package `abfallkalender`) is done;
these are the GitHub-side and filesystem-side tasks that can't be done from a
running Claude Code session.

## GitHub-side

### Rename the repos
```bash
# Dev repo (origin now points here):
gh repo rename abfallkalender-api-dev --repo ByteSizedMarius/insert_it-dev

# Public mirror (whatever name it currently has on GitHub):
gh repo rename abfallkalender-api --repo ByteSizedMarius/insert_it
```

(After the dev repo rename, the local `origin` URL is already correct - it was
pre-pointed at `https://github.com/ByteSizedMarius/abfallkalender-api-dev.git`.
Pushes will start working once the GitHub-side rename happens.)

### Set repo description + topics (both repos)
```bash
gh repo edit ByteSizedMarius/abfallkalender-api-dev \
  --description "Go API client + CLI for German municipal waste collection (Abfallkalender / Müllabfuhr). Mannheim, Kassel, Krefeld, Lübeck, Herne, Hattingen, Offenbach. Zero dependencies." \
  --add-topic abfallkalender \
  --add-topic muellabfuhr \
  --add-topic muellkalender \
  --add-topic waste-collection \
  --add-topic germany \
  --add-topic kassel \
  --add-topic mannheim \
  --add-topic krefeld \
  --add-topic luebeck \
  --add-topic herne \
  --add-topic hattingen \
  --add-topic offenbach \
  --add-topic golang \
  --add-topic api-client \
  --add-topic cli

# Repeat for the public repo once it's renamed:
gh repo edit ByteSizedMarius/abfallkalender-api \
  --description "..." --add-topic ...
```

The description and the README H1 carry most of the SEO weight; topics help
GitHub's own search. "Müll API" / "Abfallkalender" / "Kassel Müll" type
queries should hit either the README or the repo card.

## Filesystem-side

### Rename the local working directory
```
C:\Users\byte\repos\5-others\insert_it
  -> C:\Users\byte\repos\5-others\abfallkalender-api
```
Has to happen outside Claude Code (the CWD can't move under a running
session). Close any editor or shell open in the directory first.

### After the local rename, update CLAUDE.md
The project-root path in the subagent-guardrails block at the bottom of
`CLAUDE.md` still references the old path:
```
Project root is `C:\Users\byte\repos\5-others\insert_it`.
```
Update it to `\abfallkalender-api` once the directory is renamed, otherwise
subagents will get the wrong project-root constraint.

## Verification

### Run govulncheck
The earlier verification couldn't fetch the vuln DB from this sandbox. Run it
once with network access:
```bash
govulncheck ./...
```
Module has zero third-party dependencies, so the surface is just whatever Go
stdlib symbols we import.

### Trigger pkg.go.dev indexing
After the GitHub rename, the new module path is `github.com/ByteSizedMarius/abfallkalender-api`.
It gets indexed automatically on first `go get`, but you can prod the proxy:
```bash
GOPROXY=https://proxy.golang.org go get github.com/ByteSizedMarius/abfallkalender-api@latest
```
Then check `https://pkg.go.dev/github.com/ByteSizedMarius/abfallkalender-api`.
The package doc comment in `const.go` is what shows up at the top of that page.
