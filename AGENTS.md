## Goal
- Fork Gitea to CreeperCoding with a built-in AI bot (`@ccopilot`) that responds to issue/PR mentions, can be prompted via a repo Agent tab, and has file read/write/branch/PR tools.

## Constraints & Preferences
- Full fork rebrand: Go module, import paths, CLI binary name, theme files, Dockerfiles, docs, metadata all renamed from `gitea` to `creepercoding`.
- `LICENSE` attribution: existing holders preserved, added `Copyright (c) creepernet.qzz.io`.
- Theme: dark is the default (`creepercoding-dark`). All 4 GitHub-style themes provided from `/home/creeper/Downloads/theme/dist/`.
- AI agent account is `ccopilot`, non-admin, built-in system bot (`UserTypeBot`), now also a DB user created at startup.
- Only owner / edit-access / admin users can trigger code edits (branch creation + PR).
- Branch naming: `ccopilot-{date}-{name}`.
- OpenAI-compatible API endpoint + key configured in admin panel.
- ccopilot is instance-wide, configurable in admin panel (base URL, model name, API key, max tokens).
- Per-repo disable toggle in repo advanced settings (enabled by default).
- Config file (`app.ini`) editable from admin panel; **restart from UI removed** (was dangerous — killed desktop session).
- Migration sources include Gitea, GitHub, GitLab, Gogs, OneDev, GitBucket, Codebase, CodeCommit, **CreeperCoding**, **Forgejo**.
- Tagline: "Open Source, where you need it most".

## Progress
### Done
- Fork rebranded: `go.mod` → `creepercoding.dev`, all Go import paths updated, external deps (`gitea.dev/sdk`, `gitea.dev/actions-proto-go`) preserved as-is.
- CLI name and binary: `creepercoding`.
- Brand strings in `cmd/main.go`, Makefile, Dockerfiles, workflow files, locale, `custom/conf/app.example.ini` all updated.
- Copyright comments in `.go` files: `The CreeperCoding Authors`.
- `modules/setting/`: `DefaultTheme` set to `creepercoding-dark`; custom emojis `creepercoding`.
- Docs: `CREECODE.md`, `AGENTS.md`, `CLAUDE.md`, `README.md`, `README.zh-*.md`, `CONTRIBUTING.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`, `docs/*.md` updated.
- 4 GitHub-style CSS themes → `web_src/css/themes/theme-creepercoding-*.css`; old 9 theme files removed.
- `LICENSE`: `Copyright (c) creepernet.qzz.io` appended.
- `models/user/user_system.go`: `CcopilotUser` constants (ID -3, name `ccopilot`, `UserTypeBot`, `IsAdmin: false`).
- `modules/setting/ccopilot.go`: config struct with `Enabled`, `APIKey`, `Endpoint`, `ModelName`, `MaxTokens` options.
- `modules/setting/config.go`: `Ccopilot` field in `ConfigStruct`.
- `templates/admin/config_settings/ccopilot.tmpl`: admin settings form (password & number field types now supported by JS).
- `services/ccopilot/notifier.go`: notifier hook for `CreateIssueComment` — detects `@ccopilot`, enqueues task; queue init moved from `init()` → `Init()` (fixes panic at startup).
- `services/ccopilot/ai.go`: OpenAI-compatible chat API caller; added `queryAIMessages` for full conversation history.
- `services/ccopilot/agent.go`: `AgentChat` — per-repo session memory, sends full conversation to AI.
- `services/ccopilot/processor.go`: `processMention`, `buildContext`, `postReply`, `createBranchAndPR`, `loadPRDiff`, `suggestCommitMessage`, `reviewPullRequest`.
- `processor.go` fixes: import alias `issue_service`, `GetTree` → `commit.Tree.ListEntries`, `CcopilotUserName` qualifier, `FindComments` instead of `GetCommentsByID`.
- `models/repo/repo.go`: `CcopilotDisabled bool` field.
- `models/migrations/v1_27/v337.go`: migration adding `CcopilotDisabled` column.
- `models/migrations/migrations.go`: registered migration 337.
- `services/forms/repo_form.go`: `DisableCcopilot bool` form field.
- `routers/web/repo/setting/setting.go`: handler saves `CcopilotDisabled`.
- `templates/repo/settings/options.tmpl`: per-repo ccopilot disable checkbox.
- Tagline changed across 9 files.
- Pre-existing import-name bugs fixed: `ParseGiteaSiteURL`, `GiteaHooks{New,Edit}Post`, `GiteaAvatarLink`.
- `modules/structs/repo.go`: `CreeperCodingService` (10) and `ForgejoService` (11) in `GitServiceType` enum with support methods.
- `services/migrations/creepercoding.go` + `forgejo.go`: downloader factories delegating to `NewGiteaDownloader`.
- Migration form templates: `creepercoding.tmpl` + `forgejo.tmpl`.
- SVG icons: `gitea-creepercoding.svg`, `gitea-forgejo.svg`.
- **`/SKILL.md` page** — served dynamically at root, documents how AI agents interact with the CreeperCoding API (Gitea-compatible).
- **Config Editor** (`/-/admin/config/editor`): textarea with INI validation + `.bak` backup. **Restart button removed** (was calling `DoGracefulRestart()` which killed terminal/desktop sessions; replaced with manual-restart info box).
- **Locale-wide Gitea → CreeperCoding replacement** in all locale value strings (keys preserved, URLs/system prefixes kept intact).
- **`ccopilot` service wired in** — `routers/init.go` imports `services/ccopilot`, calls `Init()` late in startup (after graceful manager, before actions). Registers notifier + creates DB user + starts mention queue.
- **Agent tab** on repo pages (between Wiki and Activity) with chat interface: `GET /{repo}/agent` + `POST /{repo}/agent/chat` (JSON, session-memory, AI replies).
- `web_src/js/features/admin/config.ts`: added `password` and `number` to supported input types (was crashing on ccopilot settings form).

### In Progress
- (none)

### Fixed Recently
- **`update_mirrors` startup crash** — Root cause: duplicate key `"repo.agent"` in `locale_en-US.json`. Silent in `encoding/json` v1, but Go 1.26's `jsonv2` (`goexperiment.jsonv2`) strictly rejects duplicate object member names, causing ALL translations to fail loading.

## Key Decisions
- Module name `creepercoding.dev`; external deps (`gitea.dev/sdk`, `gitea.dev/actions-proto-go`) left unchanged.
- `ccopilot` user is `UserTypeBot` (4), same pattern as `gitea-actions`. Now persisted in DB at startup.
- Code changes use internal Go repo API (`repo_service.CreateNewBranch`, `pull.NewPullRequest`).
- Queue-based async processing (`queue.WorkerPoolQueue[*mentionTask]`) avoids blocking comment creation.
- Theme files replaced entirely; dark set as default.
- CreeperCoding and Forgejo migration sources reuse `NewGiteaDownloader`.
- Config editor writes raw ini with `.bak` backup; INI validated with `gopkg.in/ini.v1`.
- **Restart button removed entirely** — `DoGracefulRestart()` kills parent process, which cascaded to kill the desktop session when run from a terminal. Replaced with manual-restart instructions.
- `ccopilot` queue init moved from `init()` → `Init()` called late in startup to avoid panic (queue settings not loaded at package-init time).
- Agent chat uses per-repo in-memory session with full conversation history; no persistence yet.
- `type="password"` and `type="number"` fields on admin config form required JS changes — both now supported in the config form handler.

## Next Steps
- Implement file-content reading/writing tools for the agent (currently only reads root tree entries).
- Handle case where ccopilot is not a repo collaborator (read-only response fallback).
- Debug the pre-existing `update_mirrors` translation error — verify `options/locale/` is accessible from the running binary's working directory.
- Add session branch creation + PR creation flow from agent chat.
- User will specify additional features later.

## Critical Context
- `UserTypeBot` (4) is shared by `gitea-actions`; ccopilot uses the same type.
- Import paths use `creepercoding.dev/...`; `gitea.dev/` deps are in `go.mod` require block.
- `notify_service.RegisterNotifier()` in `init()` registers the ccopilot notifier.
- `modules/setting/config.Option[T]` stores dynamic (DB-backed) settings; `setting.Config().Ccopilot.*` exposes them.
- The ccopilot DB user is created by `ccopilot.Init()` at startup (called from `routers/init.go` after graceful manager is ready).
- Queue initialization (`queue.CreateUniqueQueue`) must happen in `Init()`, not `init()`, because setting configs aren't ready during package-init time.
- PR review and commit message features share `loadPRDiff` helper (git diff truncated to 8KB).
- Notifier gates on: not ccopilot's own comment, contains `@ccopilot`, `Enabled` true, endpoint+API key+model name non-empty, repository's `CcopilotDisabled` false.
- Agent chat is per-repo in-memory only (lost on server restart).
- `make build` includes frontend (Vite) + backend (Go); JS changes require frontend rebuild.
- `dev/` directory is used for local dev deployment; needs `public/`, `templates/`, `options/` alongside the binary.
- `routers/init.go:170`: `ccopilot.Init` is called before `actions_service.Init`.

## Relevant Files
- `routers/web/misc/misc.go`: `SkillMD` handler.
- `routers/web/admin/config.go`: `ConfigRestart` changed to safe message (no `DoGracefulRestart`).
- `templates/admin/config_editor.tmpl`: restart button removed, replaced with manual restart info.
- `web_src/js/features/admin/config.ts`: added `password` and `number` support in form element handlers.
- `services/ccopilot/notifier.go`: `Init()` creates DB user + initializes queue; notifier hooks `CreateIssueComment`.
- `services/ccopilot/ai.go`: `queryAIMessages` added for full conversation history.
- `services/ccopilot/agent.go`: `AgentChat` — in-memory session, sends full history to AI.
- `routers/web/repo/agent.go`: `Agent` handler for agent page.
- `routers/web/repo/agent_chat.go`: `AgentChat` handler for POST chat endpoint.
- `templates/repo/agent.tmpl`: chat UI template with JS for send/receive.
- `templates/repo/header.tmpl`: Agent tab added after Wiki, before Activity.
- `routers/web/web.go`: agent routes registered; `ccopilot` import added.
- `routers/init.go`: imports `services/ccopilot`, calls `ccopilot.Init` late in startup.
- `models/user/user_system.go`: `CcopilotUser` constants and `NewCcopilotUser()`, `GetSystemUserByName`.
- `modules/setting/ccopilot.go`: `CcopilotStruct` with config options.
- `modules/setting/config.go`: `ConfigStruct.Ccopilot` registration.
- `models/repo/repo.go`: `CcopilotDisabled` field on Repository.
- `models/migrations/v1_27/v337.go`: migration for `CcopilotDisabled` column.
- `models/migrations/migrations.go`: migration 337 registration.
- `templates/repo/settings/options.tmpl`: per-repo ccopilot checkbox.
- `templates/admin/config_settings/ccopilot.tmpl`: admin ccopilot settings form.
- `templates/admin/navbar.tmpl`: Config Editor link.
- `options/locale/locale_en-US.json`: all locale keys, Gitea values replaced with CreeperCoding.
