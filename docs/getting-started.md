# Chorale: getting started on your Mac

These steps set up a clone of Chorale on Cami's M4 Mac mini.

## 1. Clone the repo

The repo is public. Clone it with the GitHub CLI, which also sets up your login for pushing branches and opening pull requests:

```
gh auth status
gh repo clone ithinkicam/chorale ~/Dev/chorale
cd ~/Dev/chorale
```

If `gh auth status` says you aren't logged in, run `gh auth login` first.

## 2. Install the tools

The `Brewfile` lists Go, Node and the GitHub CLI. golangci-lint has to be one exact version, which Homebrew can't pin, so it comes from `go install`:

```
brew bundle
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
go version
golangci-lint version --short
```

`go version` must report Go 1.27 or later. Inside the repo, `make` uses the exact release named in `go.mod` and downloads it if needed. `golangci-lint version --short` must report 2.15.0, the version the `Makefile` names. If it doesn't, check that `$(go env GOPATH)/bin` comes before Homebrew in your `PATH`.

## 3. Install the project's own tools

OpenSpec comes from the repo's `package-lock.json`, so everyone gets the same version. gotestsum and actionlint are Go tools pinned in `go.mod`, and need no install step.

```
npm ci --ignore-scripts
npx openspec --version
```

The OpenSpec skills that agents use call `openspec` directly, so install the same version globally too, and turn off its anonymous usage reports:

```
npm install -g @fission-ai/openspec@2.1.1
openspec config set telemetry.enabled false
openspec --version
```

## 4. Check everything works

```
make check
```

This runs the linters, every test layer, the scenario check and OpenSpec validation. It should end with every scenario covered and validation passing.

## 5. Make a change

All work starts as an OpenSpec change. Open an agent in the repo and propose one with the `openspec-propose` skill. Once it has had its reviews and you've approved it, apply it with the `openspec-apply-change` skill. In Claude Code, these are `/openspec-propose` and `/openspec-apply-change`. In Codex, use `$` instead of `/`. In Antigravity CLI, use `/`.

The agent works on a branch, through `tasks.md` test first: it writes each test, shows it failing, then writes the code. It opens a pull request when it's done, because `main` only accepts pull requests. The rules it follows are in `AGENTS.md`, and the code style is in `docs/style.md`. Claude Code reads the rules through `CLAUDE.md`.

## Good to know

1. The OpenSpec skills live in `.agents/skills/`. The folders in `.claude/skills/` are links to them, so edit them in one place only.
2. `npx openspec list` shows the open changes, and `npx openspec validate --strict --all` checks them. Always pass `--all` or a change name.
3. The interface sketches are on the design canvas linked from `docs/plan.md`.
4. If a cloud Claude session is also working on the repo, run `git pull` before starting, so you aren't working on an old copy.
