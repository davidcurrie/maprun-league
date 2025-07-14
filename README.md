# MapRun League

MapRun League is a tool for calculating and displaying league standings for MapRun events. It fetches results, applies custom scoring rules, and generates an HTML leaderboard.

## Features
- Fetches event results from MapRun
- Customizable scoring system (max points, runs to count)
- Handles ties and displays medals for top scores
- Generates a compact HTML table for publishing

## Requirements
- Go 1.18+
- Internet access to fetch MapRun results

## Setup
1. Clone the repository:
   ```sh
   git clone https://github.com/davidcurrie/maprun-league.git
   cd maprun-league
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Configure your league in `config.yaml` (see example in repo).

## Usage
Run the league processor:
```sh
go run ./cmd/league
```
This will output the HTML results table to stdout or a file, depending on your configuration.

## Customization
- Edit `config.yaml` to set events, scoring, and other options.
- Adjust templates in `internal/league/league.go` for custom HTML output.

## License
MIT
