# siig-rules

Community-maintained rules database for [SIIG](https://github.com/yv3000/SIIG) — the zero-config CLI that auto-populates ignore files based on your project's tech stack.

Every JSON file in `db/` defines one stack's detection logic and ignore rules. These files are automatically embedded into the SIIG binary on each release.

## Structure

```text
db/
├── node.json       ← Node.js rules
├── python.json     ← Python rules
├── rust.json       ← Rust rules
└── ...110 stacks total
```

## Want to add a stack?

See [CONTRIBUTING.md](./CONTRIBUTING.md) for the full guide.

## How it works

When a PR is merged into this repo, the SIIG main repo automatically syncs the updated rules and cuts a new release with the updated binary.

Made by [@yv3000](https://github.com/yv3000)
