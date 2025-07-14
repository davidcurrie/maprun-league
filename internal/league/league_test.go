package league

import "testing"

func TestRunner_Score(t *testing.T) {
	runner := NewRunner("Test Runner")
	runner.Results = map[int]int{
		0: 50, // event 0, 50 points
		1: 45, // event 1, 45 points
		2: 48, // event 2, 48 points
		3: 30, // event 3, 30 points
	}

	// Test with maxRuns = 3
	// Should sum the top 3 scores: 50 + 48 + 45 = 143
	score := runner.Score(3)
	if score != 143 {
		t.Errorf("expected score to be 143 with 3 max runs, but got %d", score)
	}

	// Test with maxRuns = 2
	// Should sum the top 2 scores: 50 + 48 = 98
	score = runner.Score(2)
	if score != 98 {
		t.Errorf("expected score to be 98 with 2 max runs, but got %d", score)
	}

	// Test with maxRuns > number of runs
	// Should sum all scores: 50 + 48 + 45 + 30 = 173
	score = runner.Score(5)
	if score != 173 {
		t.Errorf("expected score to be 173 with 5 max runs, but got %d", score)
	}
}
