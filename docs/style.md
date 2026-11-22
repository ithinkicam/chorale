# Chorale: code style

As of 2026-11-22. Every agent follows this guide, and `AGENTS.md` points here. It builds on Google's Go Style Guide, https://google.github.io/styleguide/go/, which applies wherever this guide says nothing.

## The principle

Code explains itself. The reasons live elsewhere: behaviour in the OpenSpec specs, decisions in each change's `design.md`, working rules in `AGENTS.md`, and the story of a change in its commit message. A comment repeats none of that.

## Two kinds of comment

Doc comments sit directly above a package clause or a top-level declaration. They are documentation: `go doc` and pkgsite turn them into reference pages, and a future Chorale doc site would publish them. Write them for someone reading the docs, not for someone reviewing the diff.

Every other comment is a note for whoever maintains the code. It never appears in the docs, so it should be rare.

## Doc comments

1. Every package has a package comment, and every exported name has a doc comment, as Google's guide asks.
2. A doc comment is one sentence by default. It starts with the name and says what the thing is or does.
3. Add a second sentence only for what a caller needs and the signature doesn't show, such as cleanup the caller must do, an error to compare against, or a surprising zero value.
4. A doc comment never explains how something is implemented, why it was designed that way, or which task or scenario it came from. Those belong in the spec, the design or the commit message.
5. Unexported names get a doc comment only when their behaviour isn't obvious.

## Other comments

1. Make the code clearer before writing a comment. Rename, extract a function, or add a named constant or type first.
2. Inside a function, write no comments by default. Write one only for a reason the code can't show, in one line where possible. Examples: why an ignored error is safe, why a value is what it is, a workaround for a library bug, or when a goroutine exits.
3. Never write a comment that says what the code does, such as "loop over the sessions" or "return the result".
4. Never name a task, a scenario or a ticket in a comment. A test's name already ties it to its scenario.
5. Don't restate standard behaviour, such as context cancellation, safety for concurrent use, or the list of parameters.
6. Leftover work goes in the OpenSpec tasks, never in a TODO comment.
7. Don't add, reword or remove comments in code your task doesn't otherwise change.

## Go

1. Format with `gofmt`.
2. Handle every error. When one really can be ignored, assign it to `_` and say why in one line.
3. A name's length grows with its scope. Receiver names are short and the same across a type's methods. A name doesn't repeat its package, its type or its receiver, and getters don't start with `Get`.
4. Keep changes small, and in the style of the code around them.

## Tests

1. A test that covers a scenario is named after it, as `docs/testing.md` describes.
2. Test helpers call `t.Helper()`.
3. A failure message says what was found and what was wanted, in that order.

## What the tools check

`make lint` checks every rule a tool can check: formatting, that every exported name and package has a doc comment that starts with its name, comment spacing, commented-out code, placeholder doc comments, TODO comments, repeated words, and a reason on every `//nolint`. No tool can tell whether a comment only repeats the code, or whether a doc comment is chatty, so those rules are checked in review.
