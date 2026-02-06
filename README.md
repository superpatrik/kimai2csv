![License: CC BY-NC 4.0](https://img.shields.io/badge/license-CC--BY--NC--4.0-lightgrey)

# kimai2csv
> Happy-path CLI — non-commercial license

A small Go CLI tool that exports **personal Kimai timesheets to CSV**.

Designed for a single-user workflow, this tool focuses on speed:
generate structured CSV data or get a clean console overview in seconds.

No UI clicking. No manual exports. Just run the CLI.

---

## Motivation

This tool exists because exporting timesheets from Kimai for client reporting is annoying.

As a freelancer I repeatedly had to:
- click through Kimai export menus
- clean up formatting
- restructure data
- paste everything into client-specific timesheet templates

That friction adds up.

This CLI turns the workflow into:

> track time in Kimai → export CSV via CLI → transfer into client timesheet → invoice done

No UI gymnastics. No manual cleanup.  
Just structured data, ready to move into the format your client expects.

It’s optimized for developer workflows where speed and repeatability matter.

---

## Design Philosophy

This is a **happy-path tool**.

It intentionally optimizes for the common case and assumes a working Kimai setup.
Error handling is minimal and the code favors simplicity over abstraction.

The goal is not to be a production-grade SDK — it’s a pragmatic personal utility.

> This is a "weekend tool that stuck around", not a framework.

If something breaks, it will probably fail loudly — and that’s acceptable
for a tool built around a predictable personal workflow.

---

## Features

- Single-user export workflow
- Export Kimai → CSV optimized for client timesheets
- Grouping by date → project → (optional) activity
- Default export range: **current month**
- Automatic **last-month** export (`-lastMonth`)
- Extended output mode with activities (`-extended`)
- Automatic totals printed to console:
    - total hours
    - net + gross amount (gross = net * 1.19)
- Fast console overview for quick inspection

---

## Requirements

- Go (latest version recommended)
- Access to a Kimai instance + API token
- Internal package `timesheet/kimai` (Kimai client)

---

## Installation

### Build locally

~~~bash
go build -o kimai2csv .
~~~

### Or install via Go

~~~bash
go install
~~~

Then run from anywhere:

~~~bash
kimai2csv -apiToken "YOUR_TOKEN" -url "https://kimai.example.com"
~~~

---

## Makefile (recommended)

The repository includes a Makefile to simplify building and releasing binaries.

### Targets

- `make clean` → removes the `build/` directory
- `make tidy` → runs `go mod tidy`
- `make compile` → builds stripped static binary
- `make release` → cross-compiles + packages tarballs + checksums
- `make install` → installs binary to `/usr/local/bin`

### Common commands

~~~bash
make compile
make release
sudo make install
~~~

---

## Usage

### Required parameters

- `-apiToken` — Kimai API token
- `-url` — base URL of your Kimai instance

### Optional parameters

- `-projects` — comma-separated project IDs
- `-user` — user to export (typically yourself)
- `-begin` — start time in format `2006-01-02T15:04:05`
- `-end` — end time in format `2006-01-02T15:04:05`
- `-extended` — show activities in output
- `-csv` — output CSV file path
- `-lastMonth` — automatically use last month

---

## CSV Output

Semicolon-separated (`;`) CSV tailored for client timesheets.

### Extended mode

~~~
Datum;Beginn;Ende;Dauer;Projekt;Tätigkeit;Beschreibung;Preis
~~~

### Grouped mode

~~~
Datum;Beginn;Ende;Dauer;Projekt;Beschreibung;Preis
~~~

Grouped mode merges entries per day/project and aggregates
duration, descriptions and total price.

---

## Notes

- If `-begin` and `-end` are empty → current month
- `-lastMonth` overrides default month selection
- Date format: `YYYY-MM-DDTHH:MM:SS`
- Designed for personal workflows
- Optimized for speed, not enterprise robustness

---

## License

This project is licensed under  
**Creative Commons Attribution–NonCommercial 4.0 (CC BY-NC 4.0)**.

You may use, modify and share this software for personal and non-commercial purposes.

Commercial use is **not allowed** without explicit permission.

## Commercial Licensing

If you want to use this tool in a commercial product, service, or company environment,
please contact the author to obtain a commercial license.
