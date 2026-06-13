# CreeCode Agent Instructions — CreeperCoding

This file is parsed as a system prompt for CreeCode, the AI coding assistant for
the **CreeperCoding** project (a fork of Gitea).

---

## Project Identity

- **Project:** CreeperCoding (fork of Gitea) — a painless self-hosted Git service.
- **Module:** `creepercoding.dev`
- **Binary:** `creepercoding`
- **CLI Name:** `creepercoding`
- **Copyright:** The CreeperCoding Authors, The Gitea Authors, The Gogs Authors.

---

## Build & Development

### Targets
```
make help            # list all targets
make build           # build everything (frontend + backend)
make backend         # build Go backend → ./creepercoding
make frontend        # build frontend assets (Vite + Rolldown)
make fmt             # format Go code
make lint-go         # lint Go files
make lint-js         # lint TypeScript/JS files
make tidy            # run go mod tidy after go.mod changes
```

### Testing
```
go test -run '^TestName$' ./modulepath/          # single Go test
pnpm exec vitest <path-filter>                   # single JS test
GITEA_TEST_E2E_FLAGS='<filepath>' make test-e2e  # single Playwright e2e test
```

---

## Code Conventions

### General
- Use Conventional Commits: `type(scope): subject`
  - Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`, `perf`, `style`
  - `!` before colon for breaking changes
- Append `Assisted-By: AGENT_NAME:MODEL_VERSION` to every commit
- Never force-push, amend, or squash unless explicitly asked
- Preserve existing code comments; never remove or rewrite relevant comments
- Keep comments short, prefer same-line, explain *why*, never narrate what the code does
- Add current year to copyright headers of new `.go` files
- Ensure no trailing whitespace in edited files
- Never add HTML-style comments to Go or template files
- Use `make fmt` before committing Go changes

### Go
- Prefer unit tests over integration tests when logic is testable in isolation
- Aim for sub-2s local runtime for integration and e2e tests
- Import external `gitea.dev/*` packages as-is (they are upstream modules)

### TypeScript / Frontend
- Use `!` (non-null assertion) over `?.`/`??` when a value is known to always exist
- Prefer `flex-*` helpers over per-child `tw-ml-*` / `tw-mr-*` margins
- Fall back to `tw-*` utilities when specificity requires `!important`

### CSS
- Theme files: `web_src/css/themes/theme-creepercoding-*.css`
- All themes use custom element `<creepercoding-theme-meta-info>` with attributes:
  - `--theme-display-name`: `"Light"`, `"Dark"`, `"Auto"`
  - `--theme-color-scheme`: `"light"`, `"dark"`, `"auto"`
  - `--theme-colorblind-type`: optional `"red-green"` or `"blue-yellow"`

---

## Branding

- App name in UI / CLI: **CreeperCoding**
- CLI binary: `creepercoding`
- Module: `creepercoding.dev`
- CSS prefix: `cc-` (for CreeperCoding utility classes)
- Theme file prefix: `theme-creepercoding-`
- Logo: `assets/logo.svg` (CreeperCoding brand mark)
- Favicon: `assets/favicon.svg`

---

## Architecture

- **Language:** Go 1.26+ (backend), TypeScript/Vue 3 (frontend)
- **Build:** Vite 8 + Rolldown (frontend), Go toolchain (backend)
- **CSS:** Tailwind CSS (`tw-` prefix) + custom CSS modules
- **Database:** SQLite, PostgreSQL, MySQL, MSSQL — configurable via `app.ini`
- **Templates:** Go `text/template` in `templates/` directory
- **Routing:** Chi router in `routers/`
- **ORM:** Custom models in `models/`, migrations in `models/migrations/`
- **Actions:** CI/CD runner system in `services/actions/`, `models/actions/`

---

## Key File Paths

```
main.go                              -- entry point
cmd/main.go                          -- CLI app definition (NewMainApp)
cmd/web.go                           -- web server command
go.mod                               -- module: creepercoding.dev
Makefile                             -- build orchestration
vite.config.ts                       -- frontend build (Vite + Rolldown)
tailwind.config.ts                   -- Tailwind CSS config
assets/logo.svg                      -- primary logo
assets/favicon.svg                   -- favicon
web_src/css/themes/theme-creepercoding-*.css  -- theme files
web_src/css/index.css                -- main CSS entry
web_src/css/base.css                 -- root CSS variables
web_src/js/index.ts                  -- main JS entry
Dockerfile / Dockerfile.rootless     -- container builds
custom/conf/app.example.ini          -- example config
```

---

## Related Projects

- **CreeperCoding** — this project
- **act_runner** — Actions runner (upstream: `gitea.com/gitea/runner`)
- **CreeperNet** — deployment at `creepernet.qzz.io`
