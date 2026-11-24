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

func TestDuplicateScenarioNamesFail(t *testing.T) {
	report, err := Check("testdata/dupes", strings.NewReader(`{"Action":"pass","Test":"TestSameName"}`))
	if err != nil {
		t.Fatal(err)
	}
	assertNames(t, "duplicate", report.Duplicates, "Same name", "Same name")
	if !report.Failed() {
		t.Error("report.Failed() is false with duplicate scenario names")
	}
}

func TestScenarioHeadingsInCodeBlocksAreIgnored(t *testing.T) {
	report, err := Check("testdata/dupes", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range slices.Concat(report.Covered, report.Uncovered, report.ByHand) {
		if s.Name == "Inside a code block" {
			t.Errorf("the heading inside a code block was read as a scenario")
		}
	}
}

func TestARelativeOpenSpecPathWorks(t *testing.T) {
	t.Chdir("testdata/openspec")
	report, err := Check(".", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if got := len(report.Uncovered) + len(report.ByHand); got != 8 {
		t.Errorf("found %d scenarios from the current folder, want 8", got)
	}
}

func TestAChangeMayRepeatScenariosFromTheMainSpecs(t *testing.T) {
	report, err := Check("testdata/modified", strings.NewReader(`{"Action":"pass","Test":"TestSharedScenario"}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Duplicates) != 0 {
		t.Errorf("Duplicates = %+v, want none, because a modified requirement repeats its scenarios", report.Duplicates)
	}
	if report.Failed() {
		t.Error("report.Failed() is true for a change that modifies a requirement")
	}
}

func TestAChangeReplacesMainSpecRequirements(t *testing.T) {
	report, err := Check("testdata/replaced", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	assertNames(t, "uncovered", report.Uncovered,
		"Same requirement name in another capability",
		"Scenario after the change",
		"Scenario kept through a rename",
		"Scenario under the new name",
		"Untouched scenario")
}

func TestSkippedRequirementsAreListed(t *testing.T) {
	report, err := Check("testdata/replaced", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	report.Print(&out)
	for _, want := range []string{
		`Skipped "Removed one" in testdata/replaced/specs/shell/spec.md, because change rework removes it`,
		`Skipped "Renamed and modified" in testdata/replaced/specs/shell/spec.md, because change rework renames and modifies it`,
		`Skipped "Modified one" in testdata/replaced/specs/shell/spec.md, because change rework modifies it`,
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("the report doesn't say %q; it says:\n%s", want, out.String())
		}
	}
	if got := strings.Count(out.String(), "Skipped "); got != 3 {
		t.Errorf("the report lists %d skips, want 3:\n%s", got, out.String())
	}
}

func TestAChangeNamingAMissingRequirementIsAnError(t *testing.T) {
	_, err := Check("testdata/stale", strings.NewReader(""))
	if err == nil || !strings.Contains(err.Error(), `"Long gone"`) || !strings.Contains(err.Error(), "old") {
		t.Fatalf("err = %v, want an error naming the requirement \"Long gone\" and the change old", err)
	}
}

func TestRequirementNamesAreMatchedExactly(t *testing.T) {
	_, err := Check("testdata/case", strings.NewReader(""))
	if err == nil || !strings.Contains(err.Error(), `"quit with Q"`) {
		t.Fatalf("err = %v, want an error, because \"quit with Q\" is not \"Quit with q\"", err)
	}
}
