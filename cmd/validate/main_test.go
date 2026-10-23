package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func goodRule() *rule {
	r := &rule{ID: "gleam", Name: "Gleam"}
	r.Detect = &struct {
		Files []string `json:"files"`
		Dirs  []string `json:"dirs"`
	}{Files: []string{"gleam.toml"}}
	r.Rules.Gitignore = []string{"build/"}
	return r
}

func TestCheck_acceptsWellFormedRule(t *testing.T) {
	if problems := goodRule().check("gleam.json", map[string]string{}); len(problems) != 0 {
		t.Fatalf("unexpected problems: %v", problems)
	}
}

func TestCheck_rejects(t *testing.T) {
	withDetect := func(files, dirs []string) *rule {
		r := goodRule()
		r.Detect = &struct {
			Files []string `json:"files"`
			Dirs  []string `json:"dirs"`
		}{Files: files, Dirs: dirs}
		return r
	}

	tests := []struct {
		name     string
		rule     *rule
		filename string
		seen     map[string]string
		want     string
	}{
		{name: "id and filename disagree", rule: goodRule(), filename: "other.json", seen: map[string]string{}, want: "does not match filename"},
		{name: "uppercase id", rule: withUppercaseID(), filename: "gleam.json", seen: map[string]string{}, want: "not lowercase kebab-case"},
		{name: "duplicate id", rule: goodRule(), filename: "gleam.json", seen: map[string]string{"gleam": "first.json"}, want: "already used by first.json"},
		{name: "missing name", rule: func() *rule { r := goodRule(); r.Name = ""; return r }(), filename: "gleam.json", seen: map[string]string{}, want: "name is required"},
		{name: "long name", rule: func() *rule { r := goodRule(); r.Name = strings.Repeat("n", maxNameLength+1); return r }(), filename: "gleam.json", seen: map[string]string{}, want: "name is longer"},
		{name: "global must not detect", rule: globalWithDetect(), filename: "global.json", seen: map[string]string{}, want: "must use detect: null"},
		{name: "stack needs detect", rule: func() *rule { r := goodRule(); r.Detect = nil; return r }(), filename: "gleam.json", seen: map[string]string{}, want: "detect is required"},
		{name: "empty detect", rule: withDetect(nil, nil), filename: "gleam.json", seen: map[string]string{}, want: "at least one file or directory"},
		{name: "newline entry", rule: func() *rule { r := goodRule(); r.Rules.Gitignore = []string{"build/\n*"}; return r }(), filename: "gleam.json", seen: map[string]string{}, want: "contains a newline"},
		{name: "blank entry", rule: func() *rule { r := goodRule(); r.Rules.Dockerignore = []string{"  "}; return r }(), filename: "gleam.json", seen: map[string]string{}, want: "entry 0 is empty"},
		{name: "oversized entry", rule: func() *rule {
			r := goodRule()
			r.Rules.Npmignore = []string{strings.Repeat("a", maxEntryLength+1)}
			return r
		}(), filename: "gleam.json", seen: map[string]string{}, want: "longer than"},
		{name: "too many entries", rule: func() *rule {
			r := goodRule()
			r.Rules.Vercelignore = make([]string, maxEntriesPerKind+1)
			for i := range r.Rules.Vercelignore {
				r.Rules.Vercelignore[i] = "x"
			}
			return r
		}(), filename: "gleam.json", seen: map[string]string{}, want: "more than 512"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problems := tt.rule.check(tt.filename, tt.seen)
			if !containsSubstring(problems, tt.want) {
				t.Fatalf("problems %v do not mention %q", problems, tt.want)
			}
		})
	}
}

func TestCheck_vscodeKeys(t *testing.T) {
	r := goodRule()
	r.Rules.VSCode = &struct {
		FilesExclude map[string]bool `json:"files.exclude"`
	}{FilesExclude: map[string]bool{"build/": true, "": true, "docs": true}}

	problems := r.check("gleam.json", map[string]string{})
	if !containsSubstring(problems, "must not end with a slash") {
		t.Fatalf("trailing slash key not reported: %v", problems)
	}
	if !containsSubstring(problems, "vscode exclude key is empty") {
		t.Fatalf("blank key not reported: %v", problems)
	}
}

func TestRun(t *testing.T) {
	dir := t.TempDir()
	writeRule(t, dir, "gleam.json", `{"id":"gleam","name":"Gleam","detect":{"files":["gleam.toml"],"dirs":[]},"rules":{"gitignore":["build/"]}}`)
	writeRule(t, dir, "nim.json", `{"id":"nim","name":"Nim","detect":{"files":["*.nimble"],"dirs":[]},"rules":{"gitignore":["nimcache/"]}}`)

	if err := run(dir); err != nil {
		t.Fatalf("valid rules should pass: %v", err)
	}

	writeRule(t, dir, "broken.json", `{"id":"nim","name":"Nim","detect":{"files":["x"],"dirs":[]}}`)
	err := run(dir)
	if err == nil {
		t.Fatal("duplicate id should fail the run")
	}
	if !strings.Contains(err.Error(), "broken.json") {
		t.Fatalf("error should name the offending file: %v", err)
	}
}

func TestRun_missingDirectory(t *testing.T) {
	if err := run(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Fatal("expected error for a missing db directory")
	}
}

func withUppercaseID() *rule {
	r := goodRule()
	r.ID = "Gleam"
	return r
}

func globalWithDetect() *rule {
	r := goodRule()
	r.ID = "global"
	return r
}

func writeRule(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func containsSubstring(values []string, want string) bool {
	for _, v := range values {
		if strings.Contains(v, want) {
			return true
		}
	}
	return false
}
