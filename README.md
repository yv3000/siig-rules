<p align="center">
  <img src="https://img.shields.io/badge/JSON-000000?style=for-the-badge&logo=json&logoColor=white" alt="JSON" />
  <img src="https://img.shields.io/badge/license-MIT-green?style=for-the-badge" alt="MIT License" />
  <img src="https://img.shields.io/github/v/release/yv3000/SIIG?style=for-the-badge&color=blueviolet&label=siig" alt="SIIG Release" />
</p>

<h1 align="center">siig-rules</h1>
<p align="center"><b>Community-maintained rules database for <a href="https://github.com/yv3000/SIIG">SIIG</a>.</b></p>
<p align="center">Every JSON file here defines one stack's detection logic and ignore rules — and gets embedded into the SIIG binary on each release.</p>

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
│   └── ...110 stacks total
├── schema.json          ← JSON Schema for validating rule files
└── CONTRIBUTING.md      ← full guide for adding new stacks
```

---

## How it syncs with SIIG

When a PR is merged into this repo, the SIIG main repo automatically picks up the new rules:

```
PR merged into siig-rules
        ↓
Every Monday at 2am UTC — sync workflow runs
        ↓
New JSON files copied into SIIG/rules/db/
        ↓
New SIIG release built and published automatically
        ↓
Users get the new stack on next update — no action needed
```

---

## Want to add a stack?

See [CONTRIBUTING.md](./CONTRIBUTING.md) for the full guide. The short version:

1. Fork this repo
2. Create `db/<your-stack>.json`
3. Open a PR titled `feat: add <stack-name> rules`

---

## Supported stacks (110)

`android` · `angular` · `ansible` · `apache` · `astro` · `aws-cdk` · `bazel` · `bun` · `c` · `capacitor` · `clojure` · `conda` · `cpp` · `csharp` · `dart` · `deno` · `django` · `docker` · `dotnet` · `drizzle` · `electron` · `elixir` · `erlang` · `expo` · `express` · `fastapi` · `firebase` · `flask` · `flutter` · `flyio` · `fsharp` · `gatsby` · `github-actions` · `global` · `go` · `godot` · `gradle` · `graphql` · `groovy` · `haskell` · `ios` · `java` · `julia` · `jupyter` · `kotlin` · `kotlin-jvm` · `kubernetes` · `laravel` · `lerna` · `lua` · `matlab` · `maven` · `mongodb` · `mysql` · `nestjs` · `netlify` · `nextjs` · `nginx` · `nix` · `node` · `nuxt` · `nx` · `objectivec` · `packer` · `perl` · `php` · `pnpm` · `postgres` · `prisma` · `pulumi` · `python` · `pytorch` · `qwik` · `r` · `railway` · `react` · `react-native` · `redis` · `remix` · `ruby` · `rust` · `scala` · `shopify` · `solid` · `spring` · `sqlite` · `storybook` · `stripe` · `supabase` · `svelte` · `swift` · `symfony` · `t3` · `tamagui` · `tauri` · `terraform` · `testing` · `trpc` · `turborepo` · `unity` · `unreal` · `vagrant` · `vercel` · `vite` · `vscode-extension` · `vue` · `wails` · `wordpress` · `yarn` · `zig`

---

## License

MIT — [@yv3000](https://github.com/yv3000)