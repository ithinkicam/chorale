package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"unicode"
)

// Scenario is one "#### Scenario:" in a spec.
type Scenario struct {
	Name   string
	File   string
	Reason string
}

// Report sorts every scenario into covered, uncovered and checked by hand.
type Report struct {
	Covered, Uncovered, ByHand []Scenario
	Duplicates                 []Scenario
}

// Failed reports whether any scenario is uncovered.
func (r Report) Failed() bool { return len(r.Uncovered) > 0 || len(r.Duplicates) > 0 }

// Print writes the report for a person to read.
func (r Report) Print(w io.Writer) {
	for _, s := range r.Duplicates {
		_, _ = fmt.Fprintf(w, "Duplicate scenario name: %q in %s\n", s.Name, s.File)
	}
	for _, s := range r.ByHand {
		_, _ = fmt.Fprintf(w, "Checked by hand: %q in %s, because %s\n", s.Name, s.File, s.Reason)
	}
	for _, s := range r.Uncovered {
		_, _ = fmt.Fprintf(w, "No passing test: %q in %s\n", s.Name, s.File)
	}
	_, _ = fmt.Fprintf(w, "%d scenarios covered, %d checked by hand, %d with no passing test.\n",
		len(r.Covered), len(r.ByHand), len(r.Uncovered))
}

const (
	scenarioPrefix = "#### Scenario:"
	byHandPrefix   = "- **Checked by hand:**"
)

// Check reads the scenarios in openspecDir/specs and in each open change
// under openspecDir/changes, and reports which the test results cover.
func Check(openspecDir string, results ...io.Reader) (Report, error) {
	scenarios, err := readScenarios(openspecDir)
	if err != nil {
		return Report{}, err
	}
	outcomes, err := readOutcomes(results)
	if err != nil {
		return Report{}, err
	}

	var report Report
	report.Duplicates = duplicates(scenarios)
	for _, s := range scenarios {
		switch {
		case s.Reason != "":
			report.ByHand = append(report.ByHand, s)
		case outcomes[normalise(s.Name)] == passed:
			report.Covered = append(report.Covered, s)
		default:
			report.Uncovered = append(report.Uncovered, s)
		}
	}
	return report, nil
}

func readScenarios(openspecDir string) ([]Scenario, error) {
	archive := filepath.Join(openspecDir, "changes", "archive")
	var scenarios []Scenario
	for _, root := range []string{filepath.Join(openspecDir, "specs"), filepath.Join(openspecDir, "changes")} {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			if err != nil {
				return err
			}
			if d.IsDir() && path == archive {
				return filepath.SkipDir
			}
			if d.IsDir() || filepath.Ext(path) != ".md" || !slices.Contains(strings.Split(filepath.ToSlash(path), "/"), "specs") {
				return nil
			}
			found, err := scenariosIn(path)
			scenarios = append(scenarios, found...)
			return err
		})
		if err != nil {
			return nil, err
		}
	}
	sort.SliceStable(scenarios, func(i, j int) bool { return scenarios[i].File < scenarios[j].File })
	return scenarios, nil
}

func scenariosIn(path string) ([]Scenario, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}
	var scenarios []Scenario
	current := -1
	inFence := false
	for _, line := range strings.Split(data, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		switch {
		case strings.HasPrefix(trimmed, scenarioPrefix):
			name := strings.TrimSpace(strings.TrimPrefix(trimmed, scenarioPrefix))
			scenarios = append(scenarios, Scenario{Name: name, File: path})
			current = len(scenarios) - 1
		case strings.HasPrefix(trimmed, "#"):
			current = -1
		case current >= 0 && strings.HasPrefix(trimmed, byHandPrefix):
			scenarios[current].Reason = strings.TrimSpace(strings.TrimPrefix(trimmed, byHandPrefix))
		}
	}
	return scenarios, nil
}

type outcome int

const (
	unknown outcome = iota
	passed
	notPassed
)

type event struct {
	Action string
	Test   string
}

// readOutcomes maps each normalised test and subtest name to whether it
// passed. Any failure or skip under a name outweighs a pass.
func readOutcomes(results []io.Reader) (map[string]outcome, error) {
	outcomes := map[string]outcome{}
	for _, r := range results {
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 || line[0] != '{' {
				continue
			}
			var e event
			if err := json.Unmarshal(line, &e); err != nil {
				return nil, fmt.Errorf("reading test results: %w", err)
			}
			if e.Test == "" {
				continue
			}
			var result outcome
			switch e.Action {
			case "pass":
				result = passed
			case "fail", "skip":
				result = notPassed
			default:
				continue
			}
			for _, segment := range strings.Split(e.Test, "/") {
				name := normalise(segment)
				if outcomes[name] != notPassed {
					outcomes[name] = result
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading test results: %w", err)
		}
	}
	return outcomes, nil
}

// normalise turns a scenario or test name into the form names are matched
// in: a leading Go "Test" removed, only letters and digits kept, all in lower case.
func normalise(name string) string {
	if rest, ok := strings.CutPrefix(name, "Test"); ok {
		if rest == "" || unicode.IsUpper(rune(rest[0])) || unicode.IsDigit(rune(rest[0])) || rest[0] == '_' {
			name = rest
		}
	}
	var b strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

func duplicates(scenarios []Scenario) []Scenario {
	count := map[string]int{}
	for _, s := range scenarios {
		count[normalise(s.Name)]++
	}
	var out []Scenario
	for _, s := range scenarios {
		if count[normalise(s.Name)] > 1 {
			out = append(out, s)
		}
	}
	return out
}
