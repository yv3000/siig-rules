# Contributing to siig-rules

Thanks for helping SIIG recognize more stacks and generate safer ignore files. Every JSON file you add here gets embedded into the SIIG binary on the next release — so your contribution reaches every SIIG user automatically.

---

## Before you start

Check the [supported stacks list](./README.md#supported-stacks-110) — if your stack is already there, no need to add it again.

---

## How to add a new stack

**1. Fork this repo and create your file**

```
db/<stack-name>.json
```

Use a lowercase filename that matches the stack's `id`. Hyphens are allowed, spaces are not. Examples: `gleam.json`, `react-native.json`.

**2. Follow this schema**

```json
{
  "id": "gleam",
  "name": "Gleam",
  "detect": {
    "files": ["gleam.toml", "*.gleam"],
    "dirs": ["build"]
  },
  "rules": {
    "gitignore": ["build/", "*.beam"],
    "vercelignore": ["build/"],
    "dockerignore": ["build/"],
    "npmignore": ["build/"],
    "vscode": {
      "files.exclude": {
        "build": true
      }
    }
  }
}
```

**Field reference:**

| Field | Required | Notes |
|---|---|---|
| `id` | ✅ | Lowercase, hyphens allowed, no spaces. Must match the filename. |
| `name` | ✅ | Human-readable name shown in SIIG output. |
| `detect` | ✅ | How SIIG identifies the stack. Use `null` only for `global.json`. |
| `detect.files` | ✅ | Filenames or glob patterns (`*.toml`). Use `[]` if none needed. |
| `detect.dirs` | ✅ | Directory names. Use `[]` if none needed. |
| `rules` | ✅ | All rule groups inside are optional — only include what's relevant. |
| `rules.gitignore` | No | Entries added to `.gitignore`. |
| `rules.vercelignore` | No | Entries added to `.vercelignore`. |
| `rules.dockerignore` | No | Entries added to `.dockerignore`. |
| `rules.npmignore` | No | Entries added to `.npmignore`. |
| `rules.vscode` | No | `files.exclude` map — hides noisy dirs in VS Code Explorer. |

**3. Test it locally**

```bash
# Copy your file into SIIG
cp db/gleam.json path/to/SIIG/rules/db/

# Rebuild SIIG
cd path/to/SIIG
go build -o siig .

# Test in a real Gleam project
cd your-gleam-project
siig --dry-run --verbose
```

Your stack should appear in the "Detected Stacks" list.

**4. Open a pull request**

- Title: `feat: add <stack-name> rules`
- Description: briefly explain what the stack is and what files it generates that should be ignored.

---

## Rules for good contributions

✅ **Detection must be specific** — `detect.files` and `detect.dirs` should only match projects that actually use this stack. A file like `main.py` alone is too broad; `manage.py` (Django-specific) is fine.

✅ **Only ignore generated files** — build outputs, caches, secrets, and IDE artifacts. Never ignore source files.

✅ **Directory entries end with `/`** — `build/` not `build`.

✅ **Don't repeat global rules** — `.env`, `*.key`, `*.pem`, `*.p12`, `.DS_Store` are already in `global.json`. Don't add them again.

✅ **vscode entries without trailing slash** — `"build": true` not `"build/": true`.

❌ **Don't add overly broad rules** — if an entry would affect unrelated projects, leave it out.

❌ **Don't add source files** — only outputs, caches, and generated files.

---

## What happens after your PR is merged

```
PR merged → Monday sync → new SIIG release → users get your stack automatically
```

No further action needed from you. SIIG's CI picks up the new file, rebuilds the binary with it embedded, and publishes a new release. Every SIIG user gets your stack on their next update.

---

## What not to add

- Stacks already in the [supported list](./README.md#supported-stacks-110)
- Rules that are too broad and would match unrelated projects
- Source files in any ignore list
- Entries already covered by `global.json`