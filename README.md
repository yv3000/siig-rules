<p align="center">
  <img src="https://img.shields.io/badge/JSON-000000?style=for-the-badge&logo=json&logoColor=white" alt="JSON" />
  <img src="https://img.shields.io/badge/license-MIT-green?style=for-the-badge" alt="MIT License" />
  <img src="https://img.shields.io/github/v/release/yv3000/SIIG?style=for-the-badge&color=blueviolet&label=siig" alt="SIIG Release" />
</p>

<h1 align="center">siig-rules</h1>
<p align="center"><b>Community-maintained rules database for <a href="https://github.com/yv3000/SIIG">SIIG</a>.</b></p>
<p align="center">Every JSON file here defines one stack's detection logic and ignore rules — and reaches SIIG users in two ways: embedded into the next release, or fetched on demand the moment it is merged.</p>

---

## What is this?

[SIIG](https://github.com/yv3000/SIIG) is a CLI that auto-populates `.gitignore`, `.vercelignore`, `.dockerignore`, and `.vscode/settings.json` based on your project's tech stack. **siig-rules** is the database it reads from.

Each file in `db/` tells SIIG two things:

- **How to detect the stack** — which files or directories indicate this stack is present.
- **What to ignore** — which build artifacts, caches, and secrets belong in ignore files.

```
siig-rules/db/node.json
       ↓
SIIG detects package.json → applies Node.js ignore rules
```

---

## Structure

```
siig-rules/
├── db/
│   ├── global.json      ← always applied, regardless of stack
│   ├── node.json
│   ├── python.json
│   ├── rust.json
│   └── ...122 stacks total
├── cmd/validate/        ← checks every file in db/ against schema.json
├── schema.json          ← JSON Schema for validating rule files
├── LICENSE              ← MIT
└── CONTRIBUTING.md      ← full guide for adding new stacks
```

---

## Two ways a rule reaches a user

**1. Fetched on demand — immediately after merge.** When SIIG finds a file or directory that none of its embedded rules recognise, it takes the name as a stack id and requests `https://raw.githubusercontent.com/yv3000/siig-rules/main/db/<id>.json`. A hit is applied and cached in `~/.siig/cache/` for 7 days; a 404 is skipped silently. This is how a stack added here today works for users on an older SIIG binary tonight.

**2. Embedded — after the sync PR lands.** The sync workflow copies `db/*.json` into [SIIG](https://github.com/yv3000/SIIG)`/rules/db/`, where it is compiled into the binary with `go:embed`, so the stack then works fully offline.

```
rule merged into siig-rules main
        ↓
        ├── SIIG fetches db/<id>.json on demand   (works right away, needs network)
        ↓
weekly sync opens a PR against SIIG  →  CI + maintainer merge  →  tagged release
        ↓
rule is embedded — works offline in every future binary
```

The sync opens a **pull request**; it never pushes to `main` or cuts a release by itself. What ships is decided by SIIG's CI and a human.

---

## Validating a rule

`go run ./cmd/validate` checks every file in `db/` against the shape SIIG expects: filename must equal `id`, `id` must be lowercase kebab-case, `detect` is required (except for `global.json`, which must be `null`), unknown JSON fields are rejected, and entries must be non-empty, ≤256 characters and free of newlines. CI runs it on every push and pull request, because SIIG rejects at runtime anything that fails these — a rule that passes review here but breaks there would simply never apply.

---

## Want to add a stack?

See [CONTRIBUTING.md](./CONTRIBUTING.md) for the full guide. The short version:

1. Fork this repo
2. Create `db/<your-stack>.json`
3. Open a PR titled `feat: add <stack-name> rules`

---

## Supported stacks (122)

`android` · `angular` · `ansible` · `apache` · `astro` · `aws-cdk` · `bazel` · `bun` · `c` · `capacitor` · `clojure` · `conda` · `coq` · `cpp` · `crystal` · `csharp` · `dart` · `deno` · `django` · `dlang` · `docker` · `dotnet` · `drizzle` · `electron` · `elixir` · `erlang` · `expo` · `express` · `fastapi` · `firebase` · `flask` · `flutter` · `flyio` · `fortran` · `fsharp` · `gatsby` · `github-actions` · `gleam` · `global` · `go` · `godot` · `gradle` · `graphql` · `groovy` · `haskell` · `ios` · `java` · `julia` · `jupyter` · `kotlin` · `kotlin-jvm` · `kubernetes` · `laravel` · `latex` · `lean` · `lerna` · `lua` · `matlab` · `maven` · `mongodb` · `mysql` · `nestjs` · `netlify` · `nextjs` · `nginx` · `nim` · `nix` · `node` · `nuxt` · `nx` · `objectivec` · `ocaml` · `packer` · `perl` · `php` · `platformio` · `pnpm` · `postgres` · `prisma` · `pulumi` · `purescript` · `python` · `pytorch` · `qwik` · `r` · `railway` · `react` · `react-native` · `redis` · `remix` · `ruby` · `rust` · `scala` · `shopify` · `solid` · `solidity` · `spring` · `sqlite` · `storybook` · `stripe` · `supabase` · `svelte` · `swift` · `symfony` · `t3` · `tamagui` · `tauri` · `terraform` · `testing` · `trpc` · `turborepo` · `unity` · `unreal` · `vagrant` · `vercel` · `vite` · `vscode-extension` · `vue` · `wails` · `wordpress` · `yarn` · `zig`

---

## License

MIT — [@yv3000](https://github.com/yv3000)