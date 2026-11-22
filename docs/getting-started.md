# Chorale: getting started on your Mac

As of 2026-11-22. These steps set up a clone of Chorale on Cami's M4 Mac mini, ready for the first OpenSpec change. Nothing is built yet, so there is no app to run.

## 1. Clone the repo

The repo is public. Clone it with the GitHub CLI, which also sets up your login for pushing branches and opening pull requests:

```
gh auth status
gh repo clone ithinkicam/chorale ~/Dev/chorale
cd ~/Dev/chorale
```

If `gh auth status` says you aren't logged in, run `gh auth login` first.

## 2. Install the tools

The `Brewfile` lists Go, golangci-lint, Node and the GitHub CLI. Install them with:

```
brew bundle
go version
```

`go version` must report Go 1.27 or later, which Bubble Tea v2 needs.

## 3. Install OpenSpec

Chorale pins OpenSpec 2.1.1. Install it and turn off its anonymous usage reports:

```
npm install -g @fission-ai/openspec@2.1.1
openspec config set telemetry.enabled false
openspec --version
```

## 4. Check the specs

These should list one change and pass validation:

```
openspec list
openspec validate --strict --all
```

## 5. Review the first change

The first change is `add-project-skeleton`, which is milestone 0. It sets up the Go project and every test layer, proved on the empty session list. Read it before any code is written. It's short:

```
openspec show add-project-skeleton
```

The files are in `openspec/changes/add-project-skeleton/`. The proposal and the spec are the parts to review. The design and tasks say how it will be built.

## 6. Start building

When you're happy with the change, open an agent in the repo and run the apply skill. In Claude Code:

```
/openspec-apply-change add-project-skeleton
```

In Codex, use `$openspec-apply-change add-project-skeleton`. In Antigravity CLI, use `/openspec-apply-change add-project-skeleton`.

The agent works on a branch, through `tasks.md` test first: it writes each test, shows it failing, then writes the code. It opens a pull request when it's done, because `main` only accepts pull requests. The rules it follows are in `AGENTS.md`. Claude Code reads them through `CLAUDE.md`.

## Good to know

1. The OpenSpec skills live in `.agents/skills/`. The folders in `.claude/skills/` are links to them, so edit them in one place only.
2. `make check` doesn't exist until milestone 0 creates it. Until then, `openspec validate --strict --all` is the only gate.
3. The interface sketches are on the design canvas linked from `docs/plan.md`.
4. If a cloud Claude session is also working on the repo, run `git pull` before starting, so you aren't working on an old copy.
