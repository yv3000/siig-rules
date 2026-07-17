# Contributing to siig-rules

Thanks for helping SIIG recognize more project stacks and generate safer ignore files.

## What this repo is

`siig-rules` is the community-maintained rules database for [SIIG](https://github.com/yv3000/SIIG). Every JSON file here gets embedded into the SIIG binary on each release.

Each file in `db/` describes:

- How SIIG detects a stack.
- Which generated, secret, cache, or build files belong in supported ignore files.
- Which noisy paths VS Code may hide from the Explorer.

## How to add a new stack

1. Fork this repository.
2. Create `db/<stack-name>.json`. Use a lowercase filename matching the stack's `id`, with hyphens instead of spaces.
3. Follow [`schema.json`](./schema.json) and the structure below.

   ```json
   {
     "id": "lowercase-stack-id",
     "name": "Human-readable Stack Name",
     "detect": {
       "files": ["specific-file.ext"],
       "dirs": ["specific-directory"]
     },
     "rules": {
       "gitignore": ["build/"],
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

   Field reference:

   | Field | Required | Meaning |
   | --- | --- | --- |
   | `id` | Yes | Lowercase identifier. Letters, numbers, and separating hyphens are allowed; spaces are not. It must match the JSON filename. |
   | `name` | Yes | Human-readable stack name shown by SIIG. |
   | `detect` | Yes | Object containing the evidence SIIG uses to detect the stack. The reserved `global.json` file is the only exception and uses `null`. |
   | `detect.files` | Yes | Array of specific filenames or file patterns. Use `[]` when no file signal is needed. |
   | `detect.dirs` | Yes | Array of specific directory names. Use `[]` when no directory signal is needed. |
   | `rules` | Yes | Object containing ignore rules. All rule groups inside it are optional. |
   | `rules.gitignore` | No | Entries SIIG may add to `.gitignore`. |
   | `rules.vercelignore` | No | Entries SIIG may add to `.vercelignore`. |
   | `rules.dockerignore` | No | Entries SIIG may add to `.dockerignore`. |
   | `rules.npmignore` | No | Entries SIIG may add to `.npmignore`. |
   | `rules.vscode` | No | VS Code settings object containing a `files.exclude` map of paths to `true` or `false`. |

4. Use this complete Gleam example as a model:

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
       "dockerignore": ["build/"],
       "vscode": { "files.exclude": { "build": true } }
     }
   }
   ```

5. Open a pull request titled `feat: add <stack-name> rules`.
6. In the pull request description, explain:
   - What the stack is.
   - What files or directories it generates that should be ignored.

## Rules for good contributions

- Make `detect.files` and `detect.dirs` specific enough to avoid false positives in unrelated projects.
- Only add entries that are genuinely dangerous to track, such as secrets and keys, or noisy to track, such as build artifacts and caches.
- `.env`, `*.key`, `*.pem`, and `*.p12` are already in `global.json`. Do not repeat them in a stack file.
- End `gitignore` directory entries with `/`.
- Never add source files to ignore lists. Only add build outputs, caches, secrets, and generated files.
- Test locally: put your JSON in `SIIG/rules/db/`, rebuild with `go build -o siig .`, then run `siig --dry-run --verbose` in a project that uses your stack.

## Currently supported stacks (110)

- `android`
- `angular`
- `ansible`
- `apache`
- `astro`
- `aws-cdk`
- `bazel`
- `bun`
- `c`
- `capacitor`
- `clojure`
- `conda`
- `cpp`
- `csharp`
- `dart`
- `deno`
- `django`
- `docker`
- `dotnet`
- `drizzle`
- `electron`
- `elixir`
- `erlang`
- `expo`
- `express`
- `fastapi`
- `firebase`
- `flask`
- `flutter`
- `flyio`
- `fsharp`
- `gatsby`
- `github-actions`
- `global`
- `go`
- `godot`
- `gradle`
- `graphql`
- `groovy`
- `haskell`
- `ios`
- `java`
- `julia`
- `jupyter`
- `kotlin`
- `kotlin-jvm`
- `kubernetes`
- `laravel`
- `lerna`
- `lua`
- `matlab`
- `maven`
- `mongodb`
- `mysql`
- `nestjs`
- `netlify`
- `nextjs`
- `nginx`
- `nix`
- `node`
- `nuxt`
- `nx`
- `objectivec`
- `packer`
- `perl`
- `php`
- `pnpm`
- `postgres`
- `prisma`
- `pulumi`
- `python`
- `pytorch`
- `qwik`
- `r`
- `railway`
- `react`
- `react-native`
- `redis`
- `remix`
- `ruby`
- `rust`
- `scala`
- `shopify`
- `solid`
- `spring`
- `sqlite`
- `storybook`
- `stripe`
- `supabase`
- `svelte`
- `swift`
- `symfony`
- `t3`
- `tamagui`
- `tauri`
- `terraform`
- `testing`
- `trpc`
- `turborepo`
- `unity`
- `unreal`
- `vagrant`
- `vercel`
- `vite`
- `vscode-extension`
- `vue`
- `wails`
- `wordpress`
- `yarn`
- `zig`

## What not to add

- Duplicate stacks already listed above.
- Overly broad detection or ignore rules that match too many projects.
- Entries that ignore important source files.
