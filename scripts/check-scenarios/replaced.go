package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Skip is a main-spec requirement whose scenarios aren't checked, because
// an open change removes or replaces it.
type Skip struct {
	Requirement string
	File        string
	Change      string
	Verb        string
}

const requirementPrefix = "### Requirement:"

type delta struct {
	removed, modified []string
	renamed           map[string]string
}

// replaced reads every open change's delta specs and returns the main-spec
// requirements they remove or replace. A requirement is named exactly, and
// only within its own capability. A bare rename replaces nothing, because
// its scenarios keep their names. A delta naming a requirement the main
// spec doesn't have is an error.
func replaced(openspecDir string) ([]Skip, error) {
	changes, err := os.ReadDir(filepath.Join(openspecDir, "changes"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var skips []Skip
	for _, change := range changes {
		if !change.IsDir() || change.Name() == "archive" {
			continue
		}
		specsDir := filepath.Join(openspecDir, "changes", change.Name(), "specs")
		err := filepath.WalkDir(specsDir, func(path string, d fs.DirEntry, err error) error {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			if err != nil || d.IsDir() || d.Name() != "spec.md" {
				return err
			}
			capability, err := filepath.Rel(specsDir, filepath.Dir(path))
			if err != nil {
				return err
			}
			found, err := skipsFor(openspecDir, change.Name(), capability, path)
			skips = append(skips, found...)
			return err
		})
		if err != nil {
			return nil, err
		}
	}
	sort.SliceStable(skips, func(i, j int) bool {
		if skips[i].File != skips[j].File {
			return skips[i].File < skips[j].File
		}
		return skips[i].Requirement < skips[j].Requirement
	})
	return skips, nil
}

func skipsFor(openspecDir, change, capability, deltaPath string) ([]Skip, error) {
	text, err := readFile(deltaPath)
	if err != nil {
		return nil, err
	}
	d := parseDelta(text)
	if len(d.removed)+len(d.modified)+len(d.renamed) == 0 {
		return nil, nil
	}

	mainPath := filepath.Join(openspecDir, "specs", capability, "spec.md")
	mainText, err := readFile(mainPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	existing := map[string]bool{}
	for _, name := range requirementsIn(mainText) {
		existing[name] = true
	}
	missing := func(name string) error {
		return fmt.Errorf("change %s names the requirement %q, which %s doesn't have", change, name, mainPath)
	}

	var skips []Skip
	skip := func(name, verb string) {
		skips = append(skips, Skip{Requirement: name, File: mainPath, Change: change, Verb: verb})
	}
	for _, name := range d.removed {
		if !existing[name] {
			return nil, missing(name)
		}
		skip(name, "removes")
	}
	renamedTo := map[string]string{}
	for from, to := range d.renamed {
		if !existing[from] {
			return nil, missing(from)
		}
		renamedTo[to] = from
	}
	for _, name := range d.modified {
		if from, ok := renamedTo[name]; ok {
			skip(from, "renames and modifies")
			continue
		}
		if !existing[name] {
			return nil, missing(name)
		}
		skip(name, "modifies")
	}
	return skips, nil
}

func parseDelta(text string) delta {
	d := delta{renamed: map[string]string{}}
	section := ""
	from := ""
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			section = strings.ToUpper(strings.Fields(trimmed)[1])
			continue
		}
		switch section {
		case "REMOVED", "MODIFIED":
			name, ok := requirementName(trimmed)
			if !ok {
				continue
			}
			if section == "REMOVED" {
				d.removed = append(d.removed, name)
			} else {
				d.modified = append(d.modified, name)
			}
		case "RENAMED":
			if rest, ok := strings.CutPrefix(trimmed, "- FROM:"); ok {
				from, _ = requirementName(strings.Trim(strings.TrimSpace(rest), "`"))
			} else if rest, ok := strings.CutPrefix(trimmed, "- TO:"); ok && from != "" {
				to, _ := requirementName(strings.Trim(strings.TrimSpace(rest), "`"))
				d.renamed[from] = to
				from = ""
			}
		}
	}
	return d
}

func requirementsIn(text string) []string {
	var names []string
	for _, line := range strings.Split(text, "\n") {
		if name, ok := requirementName(strings.TrimSpace(line)); ok {
			names = append(names, name)
		}
	}
	return names
}

func requirementName(line string) (string, bool) {
	rest, ok := strings.CutPrefix(line, requirementPrefix)
	return strings.TrimSpace(rest), ok
}
