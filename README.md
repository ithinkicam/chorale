# Chorale

A voice-first terminal workbench for coding agents, built for low vision. Chorale runs Claude Code, Codex and Antigravity CLI sessions side by side on a Mac. You dictate prompts with local transcription. An iOS companion is planned.

Milestone 0, the project skeleton, is built: the test harness, `make check`, continuous integration and the first screen. To set up a clone on your Mac, follow `docs/getting-started.md`. Then start here:

- `STATE.md` says where things stand and what's next.
- `docs/plan.md` is the roadmap.
- `docs/testing.md` is the test standard. All development is test first, including the interface.
- `AGENTS.md` holds the working rules for every agent.
- `docs/repository.md` describes how the repo is set up on GitHub: the required checks, security scanning and Dependabot.
- `openspec/` holds the specs and changes. Chorale is built with OpenSpec, and the first change, now archived, is in `openspec/changes/archive/2026-11-22-add-project-skeleton/`.
