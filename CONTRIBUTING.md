# Contributing to Chorale

Thank you for your interest. Chorale is a personal tool, built for one low-vision user, so the bar for changes is that they fit her needs and the plan in `docs/plan.md`.

All work starts as an OpenSpec change in `openspec/changes/`, with a proposal, specs, a design where needed, and tasks. Behaviour is agreed in the specs before code is written. The working rules, including the review steps, are in `AGENTS.md`. They apply to people as well as agents.

Development is test first, including the interface. Write the test, watch it fail for the reason you expect, then write the code. The test layers are described in `docs/testing.md`, and the code style is in `docs/style.md`.

`make check` is the one gate, locally and in continuous integration. A pull request needs it to pass on macOS and Linux, and needs the maintainer's approval. How the repo is set up on GitHub is in `docs/repository.md`.

To report a security problem, follow `SECURITY.md` instead of opening an issue.
