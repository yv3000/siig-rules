// Command validate checks every rule file in db/ against the shape SIIG
// expects, so a rule that SIIG would refuse to apply can never be merged.
//
// The limits mirror SIIG's rules.Validate and the fetcher's decode settings:
// SIIG parses remote rules with unknown fields rejected, caps the entry count
// and length, and treats an entry containing a newline as hostile.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	maxNameLength     = 80
	maxEntriesPerKind = 512
	maxEntryLength    = 256
	maxVSCodeKeys     = 256
)

var stackID = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type rule struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Detect *struct {
		Files []string `json:"files"`
		Dirs  []string `json:"dirs"`
	} `json:"detect"`
	Rules struct {
		Gitignore    []string `json:"gitignore"`
		Vercelignore []string `json:"vercelignore"`
		Dockerignore []string `json:"dockerignore"`
		Npmignore    []string `json:"npmignore"`
		VSCode       *struct {
			FilesExclude map[string]bool `json:"files.exclude"`
		} `json:"vscode"`
	} `json:"rules"`
}

func main() {
	if err := run("db"); err != nil {
		fmt.Fprintln(os.Stderr, "validate:", err)
		os.Exit(1)
	}
}

func run(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", dir, err)
	}

	var failures []string
	seen := map[string]string{}
	checked := 0

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		checked++

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", path, err))
			continue
		}

		var r rule
		dec := json.NewDecoder(strings.NewReader(string(data)))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&r); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", path, err))
			continue
		}

		for _, problem := range r.check(entry.Name(), seen) {
			failures = append(failures, fmt.Sprintf("%s: %s", path, problem))
		}
	}

	if len(failures) > 0 {
		sort.Strings(failures)
		return fmt.Errorf("%d problem(s) in %d rule file(s):\n  %s", len(failures), checked, strings.Join(failures, "\n  "))
	}

	fmt.Printf("✓ %d rule files valid\n", checked)
	return nil
}

// check returns every reason this rule would be rejected by SIIG or is
// inconsistent with its own filename.
func (r *rule) check(filename string, seen map[string]string) []string {
	var problems []string

	if want := strings.TrimSuffix(filename, ".json"); r.ID != want {
		problems = append(problems, fmt.Sprintf("id %q does not match filename %q", r.ID, want))
	}
	if !stackID.MatchString(r.ID) {
		problems = append(problems, fmt.Sprintf("id %q is not lowercase kebab-case", r.ID))
	}
	if other, dup := seen[r.ID]; dup {
		problems = append(problems, fmt.Sprintf("id %q already used by %s", r.ID, other))
	}
	seen[r.ID] = filename

	if strings.TrimSpace(r.Name) == "" {
		problems = append(problems, "name is required")
	}
	if len(r.Name) > maxNameLength {
		problems = append(problems, fmt.Sprintf("name is longer than %d characters", maxNameLength))
	}

	switch {
	case r.ID == "global" && r.Detect != nil:
		problems = append(problems, "the global rule must use detect: null")
	case r.ID == "global":
		// validated above; global has no detection config by design
	case r.Detect == nil:
		problems = append(problems, "detect is required (null is reserved for global)")
	case len(r.Detect.Files) == 0 && len(r.Detect.Dirs) == 0:
		problems = append(problems, "detect must list at least one file or directory")
	default:
		problems = append(problems, checkList("detect.files", r.Detect.Files)...)
		problems = append(problems, checkList("detect.dirs", r.Detect.Dirs)...)
	}

	problems = append(problems, checkList("rules.gitignore", r.Rules.Gitignore)...)
	problems = append(problems, checkList("rules.vercelignore", r.Rules.Vercelignore)...)
	problems = append(problems, checkList("rules.dockerignore", r.Rules.Dockerignore)...)
	problems = append(problems, checkList("rules.npmignore", r.Rules.Npmignore)...)

	if r.Rules.VSCode == nil {
		return problems
	}
	if len(r.Rules.VSCode.FilesExclude) > maxVSCodeKeys {
		problems = append(problems, fmt.Sprintf("vscode has more than %d exclude keys", maxVSCodeKeys))
	}
	for key := range r.Rules.VSCode.FilesExclude {
		switch {
		case strings.TrimSpace(key) == "":
			problems = append(problems, "vscode exclude key is empty")
		case strings.ContainsAny(key, "\r\n"):
			problems = append(problems, fmt.Sprintf("vscode exclude key %q contains a newline", key))
		case strings.HasSuffix(key, "/"):
			problems = append(problems, fmt.Sprintf("vscode exclude key %q must not end with a slash", key))
		}
	}

	return problems
}

func checkList(kind string, entries []string) []string {
	var problems []string

	if len(entries) > maxEntriesPerKind {
		problems = append(problems, fmt.Sprintf("%s has more than %d entries", kind, maxEntriesPerKind))
	}
	for i, entry := range entries {
		switch {
		case strings.TrimSpace(entry) == "":
			problems = append(problems, fmt.Sprintf("%s entry %d is empty", kind, i))
		case len(entry) > maxEntryLength:
			problems = append(problems, fmt.Sprintf("%s entry %d is longer than %d characters", kind, i, maxEntryLength))
		case strings.ContainsAny(entry, "\r\n"):
			problems = append(problems, fmt.Sprintf("%s entry %d contains a newline", kind, i))
		}
	}

	return problems
}
