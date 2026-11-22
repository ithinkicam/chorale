package main

import (
	"os"
	"slices"
	"strings"
	"testing"
)

func TestCheckAgainstFixtures(t *testing.T) {
	results, err := os.Open("testdata/results.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = results.Close() }()

	report, err := Check("testdata/openspec", results)
	if err != nil {
		t.Fatal(err)
	}

	assertNames(t, "covered", report.Covered,
		"Matching subtest", "Passing test", "Started at 40 columns")
	assertNames(t, "uncovered", report.Uncovered,
		"Failed test", "Golden file only", "No test at all", "Skipped test")
	assertNames(t, "checked by hand", report.ByHand, "Checked by hand only")

	for _, s := range slices.Concat(report.Covered, report.Uncovered, report.ByHand) {
		if strings.Contains(s.File, "archive") {
			t.Errorf("scenario %q from an archived change was checked", s.Name)
		}
	}
	if !report.Failed() {
		t.Error("report.Failed() is false with uncovered scenarios")
	}
}

func TestByHandKeepsTheReason(t *testing.T) {
	report, err := Check("testdata/openspec", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if len(report.ByHand) != 1 || report.ByHand[0].Reason != "needs a real microphone" {
		t.Fatalf("ByHand = %+v, want one scenario with the reason %q", report.ByHand, "needs a real microphone")
	}
}

func TestNormalise(t *testing.T) {
	cases := map[string]string{
		"Empty session list at 60 by 20 with colour": "emptysessionlistat60by20withcolour",
		"TestEmptySessionListAt60By20WithColour":     "emptysessionlistat60by20withcolour",
		"Too-small message with colour off":          "toosmallmessagewithcolouroff",
		"Matching_subtest":                           "matchingsubtest",
		"Testing the tester":                         "testingthetester",
	}
	for in, want := range cases {
		if got := normalise(in); got != want {
			t.Errorf("normalise(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAReportWithEverythingCoveredPasses(t *testing.T) {
	report := Report{Covered: []Scenario{{Name: "x"}}, ByHand: []Scenario{{Name: "y"}}}
	if report.Failed() {
		t.Error("report.Failed() is true with nothing uncovered")
	}
}

func assertNames(t *testing.T, what string, got []Scenario, want ...string) {
	t.Helper()
	var names []string
	for _, s := range got {
		names = append(names, s.Name)
	}
	slices.Sort(names)
	if !slices.Equal(names, want) {
		t.Errorf("%s scenarios = %q, want %q", what, names, want)
	}
}
