# Contributing to siig-rules

Thanks for helping SIIG recognize more stacks and generate safer ignore files. A merged rule reaches users two ways: SIIG fetches it from this repo the moment it lands (for stacks its current binary does not embed), and it is compiled into the binary after SIIG's sync pull request is merged and a release is cut.

---

## Before you start

Check the [supported stacks list](./README.md#supported-stacks-122) — if your stack is already there, no need to add it again.

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

**3. Validate it**

```bash
go run ./cmd/validate          # checks every file in db/, including yours
```

CI runs the same check on every pull request, plus the unit tests for the
validator itself. SIIG rejects a fetched rule that fails these limits, so a
rule that merges broken would simply never apply for anyone.

Hard limits enforced by both the validator and SIIG:

| Constraint | Why |
|---|---|
| filename = `id`, lowercase kebab-case, ≤64 chars | the id is the fetch URL and the cache filename |
| `name` present, ≤80 characters | shown in SIIG output |
| ≤512 entries per rule group, ≤256 characters per entry | a fetched rule is written straight into the user's file |
| no entry may contain a newline | `"a\n*"` would inject an extra rule line into `.gitignore` |
| `detect` required; `null` only for `global.json` | a rule with no detection can never match |
| no unknown JSON fields | a typo'd field silently does nothing |

To see how SIIG treats your rule end to end, copy it into a local SIIG build:

```bash
cp db/gleam.json path/to/SIIG/rules/db/
cd path/to/SIIG && go build -o siig .
cd your-gleam-project && siig --dry-run --verbose
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

❌ **Don't use generic directory names for detection** — `build`, `dist`, `output`, `target` and `node_modules` appear in unrelated projects. `nimcache`, `_opam`, `.lake`, `.pio` are safe because only their stack creates them. Prefer a specific marker file (`gleam.toml`, `fpm.toml`, `lean-toolchain`) over a directory when you can.

❌ **Don't rely on entry-point filenames** — SIIG only counts names like `index.js`, `app.py` or `server.ts` at the project root, because a nested one proves nothing. Pick a stack-specific file instead.

❌ **Don't detect on files SIIG writes** — `.gitignore`, `.vercelignore`, `.dockerignore` and `.npmignore` are ignored as detection input when they already carry SIIG's marker, otherwise a second run would detect its own output.

⚠️ **For the on-demand fetch to find your rule, one detect file must be named after the id.** SIIG builds the fetch URL from the unrecognised file's name without its extension: `gleam.toml` asks for `db/gleam.json`, `platformio.ini` asks for `db/platformio.json`. A stack whose only markers are `spago.dhall` or `dune-project` is reachable once it is embedded in a release, but not through the on-demand path — so include a `<id>.<ext>` marker file in `detect.files` where the stack has one.

---

## What happens after your PR is merged

```
PR merged into siig-rules main
        ↓
SIIG fetches db/<your-stack>.json on demand → users with an older binary get it immediately (cached 7 days)
        ↓
SIIG's weekly sync opens a pull request with the new files
        ↓
SIIG CI + a maintainer merge it → a tagged release embeds it → works offline from then on
```

The sync never pushes to `main` or publishes a release by itself — shipping is decided by SIIG's CI and a human review.

---

## What not to add

- Stacks already in the [supported list](./README.md#supported-stacks-122)
- Rules that are too broad and would match unrelated projects
- Source files in any ignore list
- Entries already covered by `global.json`