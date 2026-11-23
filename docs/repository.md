# Chorale: how the repository is set up

This is every setting the Chorale repo relies on, on GitHub and in `.github/`. The repo is public and owned by Cami's personal GitHub account. If a setting here and the live repo disagree, one of them is wrong, and the difference is worth finding.

## The two rulesets on main

`main` has two rulesets. Both are kept in `.github/rulesets/` as the payload GitHub accepts when a ruleset is created.

"Protect main" blocks deleting `main` and force-pushing to it. It requires a pull request with one approval and every review thread resolved. A new push dismisses earlier approvals, and an extra approval is required for unattributed changes. Code owner review is not required. The admin role can bypass it, for pull requests only.

"Required checks" requires three checks to pass: `check (macos)`, `check (linux)` and `dependency-review`. Each is tied to the GitHub Actions app, so no other app or status can satisfy it. It has no bypass actors, and it doesn't require a branch to be up to date with `main`.

There are two because bypassing a ruleset skips every rule in it. Agents push as Cami's account, and GitHub won't let her approve her own pull requests, so she merges them with her bypass of "Protect main". The checks sit in a ruleset nobody can bypass, so even her bypass can't merge a red pull request.

To compare a file with the live ruleset, keep the same six fields from the API's answer:

```
gh api repos/ithinkicam/chorale/rulesets --jq '.[] | "\(.id) \(.name)"'
gh api repos/ithinkicam/chorale/rulesets/<id> --jq '{name,target,enforcement,conditions,rules,bypass_actors}'
```

## When a required check is stuck

Sometimes a required check fails for every pull request, through no fault of the pull request. `vuln.go.dev` can be down, so `make vulncheck` can't download the database. Or a new advisory can name code Chorale calls, with no fixed version yet. Nobody can bypass "Required checks", so nothing can merge.

The way out is for Cami to turn the ruleset off, merge the pull request that works around the problem, and turn it straight back on. `govulncheck` can't ignore a single finding, so for an advisory the workaround is a code change, such as no longer calling the affected function.

```
id=$(gh api repos/ithinkicam/chorale/rulesets --jq '.[] | select(.name == "Required checks") | .id')
gh api --method PUT "repos/ithinkicam/chorale/rulesets/$id" -f enforcement=disabled
gh api --method PUT "repos/ithinkicam/chorale/rulesets/$id" -f enforcement=active
```

The first line finds the ruleset, the second turns it off, and the third turns it back on. Only Cami runs these.

## Security features

These are turned on in the repo's settings:

1. Secret scanning, with push protection, which refuses a push that contains a known kind of secret.
2. Dependabot alerts, and Dependabot security updates, which open pull requests that fix known vulnerabilities.
3. Private vulnerability reporting, which `SECURITY.md` points to.

These run as workflows in `.github/workflows/`:

1. `codeql.yml` runs CodeQL on the Go code and on the workflow files, on every pull request, on pushes to `main`, and every Monday. Its findings appear on the pull request and on the repo's security page. They don't block merging.
2. `dependency-review.yml` fails a pull request that adds a package with a known vulnerability of any severity. It is a required check.
3. `scorecard.yml` runs OpenSSF Scorecard on pushes to `main` and every Monday. It publishes the results to the OpenSSF API and to the security page. It doesn't block merging.
4. `make vulncheck`, inside `make check`, runs `govulncheck` on every pull request through `check.yml`.

When CodeQL or Scorecard reports something, an agent reads the finding and brings Cami a short summary with a recommendation. A real problem becomes an OpenSpec change like any other work. A false alarm is dismissed on the security page, with the reason written in the dismissal.

## Merging

Merge commits, squash merges and rebase merges are all allowed. Auto-merge is allowed, and a branch is deleted once its pull request merges. GitHub offers no merge queue for repos owned by a personal account.

`.github/CODEOWNERS` names Cami as the owner of everything, so she is asked to review every contributor's pull request. Her approval is required for those through "Protect main", not through code owner review.

The setting that lets GitHub Actions create and approve pull requests is off, and stays off. With it on, any pull request from a branch in this repo could add a workflow step that approves itself.

## Dependabot

`.github/dependabot.yml` asks for version updates every week, for Go modules, GitHub Actions and npm. Within each, patch and minor updates come grouped in one pull request. Major updates come one per pull request.

Each ecosystem has a cooldown of seven days, so Dependabot doesn't propose a version until it has been public for a week. Most malicious releases are found and pulled within days. Security updates ignore the cooldown, because they move to the patched version a reviewed advisory names.

Dependabot can't update the golangci-lint version in the `Makefile`, the Node version in `.nvmrc`, or the `toolchain` line in `go.mod`. Those are updated by hand, through an OpenSpec change.

## Auto-merge for Dependabot

`.github/workflows/dependabot-auto-merge.yml` approves Dependabot's pull requests and turns on auto-merge, so they merge on their own once every required check passes. It covers version updates and security updates alike.

Its `decide` job reads the update with `dependabot/fetch-metadata` and allows a merge only when all of these hold:

1. The update is a patch or a minor update. A major update, or an update whose size can't be told, waits.
2. It names at least one dependency.
3. None of its dependencies is held back.

Three dependencies are held back, and each comes in a pull request of its own. The first two are Charm's experimental packages, `github.com/charmbracelet/x/exp/teatest/v2` and `github.com/charmbracelet/x/exp/golden`. They are pinned to pseudo-versions, and `fetch-metadata` calls any change between two pseudo-versions a patch. The third is `dependabot/fetch-metadata` itself, so that a new version of it never judges its own update.

When `decide` allows it, the `merge` job approves the pull request and turns on auto-merge. It runs no actions, only the `gh` command, so no third-party code ever shares a runner with the token it uses.

## The auto-merge token

The `merge` job uses a fine-grained personal access token on Cami's account, stored as the Dependabot secret `DEPENDABOT_AUTOMERGE_TOKEN`. Dependabot secrets reach only workflow runs that Dependabot starts. So approvals on Dependabot's pull requests show Cami's name, though the workflow gave them.

The token has these settings:

1. Resource owner: Cami's account. Repository access: only `ithinkicam/chorale`.
2. Permissions: Contents, read and write. Pull requests, read and write. Metadata, read-only, which GitHub always adds. Nothing else.
3. Expiry: one year. A reminder on Cami's calendar comes a week before.

A personal token is used, rather than the workflow's own, for three reasons. The workflow's token could only approve with the setting described under "Merging" turned on. An approval from `github-actions` might not count toward "Protect main", while Cami's does, because she isn't the pull request's author. And a merge enabled by the workflow's token doesn't start the workflows that run on pushes to `main`.

To renew it, Cami creates a new token with the same settings at https://github.com/settings/personal-access-tokens/new and stores it:

```
gh secret set DEPENDABOT_AUTOMERGE_TOKEN --app dependabot --repo ithinkicam/chorale
```

`gh` asks for the token, so it never lands in the shell history. From a phone, she can store it on the website instead, at https://github.com/ithinkicam/chorale/settings/secrets/dependabot, by updating `DEPENDABOT_AUTOMERGE_TOKEN`. Then she deletes the old token on GitHub and moves the calendar reminder on a year.

If the token lapses, the `merge` job fails, and Dependabot's pull requests wait until it is renewed. Nothing else breaks. If the token leaks, Cami revokes it on GitHub straight away.

## Updates that wait for Cami

Major updates and updates to the held-back dependencies wait for Cami. She isn't asked to review the code. An agent reads the release notes, runs `make check` on the pull request's branch, and gives her a short summary: what changes, what breaks, and what it recommends. She makes the decision. Agents never merge one of these on their own.

If a grouped update fails `make check`, an agent finds the package that broke it, and tells Dependabot to skip that version with a `@dependabot ignore` comment on the pull request.
