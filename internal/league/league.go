package league

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/davidcurrie/maprun-league/internal/config"
	"github.com/davidcurrie/maprun-league/internal/maprun"
)

// Runner holds the data for a single participant.
type Runner struct {
	Name    string
	Results map[int]int // map[eventIndex]points
}

// NewRunner creates a new runner.
func NewRunner(name string) *Runner {
	return &Runner{Name: name, Results: make(map[int]int)}
}

// Score calculates the runner's total score from their best runs.
func (r *Runner) Score(maxRuns int) int {
	var scores []int
	for _, s := range r.Results {
		scores = append(scores, s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(scores)))

	total := 0
	for i := 0; i < len(scores) && i < maxRuns; i++ {
		total += scores[i]
	}
	return total
}

// Runs returns the number of events the runner has participated in.
func (r *Runner) Runs() int {
	return len(r.Results)
}

// ProcessLeague calculates the league standings and returns the results as an HTML string.
func ProcessLeague(cfg *config.Config) (string, error) {
	runners := make(map[string]*Runner)

	for i, event := range cfg.Events {
		results, err := maprun.GetResults(event.Name)
		if err != nil {
			return "", fmt.Errorf("failed to get results for event %s: %w", event.Name, err)
		}

		// Filter results by closing date
		var validResults []maprun.Result
		for _, r := range results {
			if event.ClosingDate.After(r.Date) {
				validResults = append(validResults, r)
			}
		}

		// Get unique latest run for each person
		sort.Slice(validResults, func(i, j int) bool {
			return validResults[i].Date.Before(validResults[j].Date)
		})
		uniqueResults := make(map[string]maprun.Result)
		for _, r := range validResults {
			uniqueResults[r.Name] = r
		}

		var sortedResults []maprun.Result
		for _, r := range uniqueResults {
			sortedResults = append(sortedResults, r)
		}

		// Sort by score, then time
		sort.Slice(sortedResults, func(i, j int) bool {
			if sortedResults[i].Score == sortedResults[j].Score {
				return sortedResults[i].TimeSecs < sortedResults[j].TimeSecs
			}
			return sortedResults[i].Score > sortedResults[j].Score
		})

		// Assign points
		for pos, res := range sortedResults {
			runner, ok := runners[res.Name]
			if !ok {
				runner = NewRunner(res.Name)
				runners[res.Name] = runner
			}
			points := cfg.Scoring.MaxPoints - pos
			if points < 1 {
				points = 1
			}
			runner.Results[i] = points
		}
	}

	// Sort runners for the final league table
	var sortedRunners []*Runner
	for _, r := range runners {
		sortedRunners = append(sortedRunners, r)
	}
	sort.Slice(sortedRunners, func(i, j int) bool {
		r1 := sortedRunners[i]
		r2 := sortedRunners[j]
		minRuns := min(r2.Runs(), r1.Runs())

		if minRuns < cfg.Scoring.MaxEventsToCount && r1.Runs() != r2.Runs() {
			return r1.Runs() > r2.Runs()
		}
		return r1.Score(cfg.Scoring.MaxEventsToCount) > r2.Score(cfg.Scoring.MaxEventsToCount)
	})

	return generateHTML(cfg, sortedRunners)
}

func generateHTML(cfg *config.Config, sortedRunners []*Runner) (string, error) {
	const tmpl = `
<p>MapRun League results are based on an individual's best {{.MaxEventsToCount}} positions. Individuals with only {{sub .MaxEventsToCount 1}} runs are ranked below those with {{.MaxEventsToCount}} and so on.</p>
<table>
  <tr>
    <th></th>
    <th>Name</th>
    <th>Total</th>
    {{- range $i, $event := .Events}}
    <th><a href="https://results.maprun.net/#/event_results?eventName={{$event.Name}}">{{add $i 1}}</a></th>
    {{- end}}
  </tr>
  {{- range $i, $runner := .Runners}}
  {{- if and (gt $i 0) (eq ($runner.Score $.MaxEventsToCount) ((index $.Runners (sub $i 1)).Score $.MaxEventsToCount))}}
  <tr>
    <td>=</td>
  {{- else}}
  <tr>
    <td>{{add $i 1}}</td>
  {{- end}}
    <td>{{$runner.Name}}</td>
    <td>{{$runner.Score $.MaxEventsToCount}}</td>
    {{- range $j, $event := $.Events }}
    <td>{{emoji (index $runner.Results $j)}}</td>
    {{- end}}
  </tr>
  {{- end}}
</table>
<p>Last updated {{.Timestamp}}</p>`

	t := template.Must(template.New("league").Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"emoji": func(points int) string {
			switch points {
			case cfg.Scoring.MaxPoints:
				return fmt.Sprintf("%d🥇", cfg.Scoring.MaxPoints)
			case cfg.Scoring.MaxPoints - 1:
				return fmt.Sprintf("%d🥈", cfg.Scoring.MaxPoints-1)
			case cfg.Scoring.MaxPoints - 2:
				return fmt.Sprintf("%d🥉", cfg.Scoring.MaxPoints-2)
			}
			if points > 0 {
				return fmt.Sprintf("%d", points)
			}
			return "-"
		},
	}).Parse(tmpl))

	data := struct {
		Runners          []*Runner
		Events           []config.Event
		MaxEventsToCount int
		Timestamp        string
	}{
		Runners:          sortedRunners,
		Events:           cfg.Events,
		MaxEventsToCount: cfg.Scoring.MaxEventsToCount,
		Timestamp:        time.Now().Format("Monday, January 2 2006, 3:04 pm"),
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
