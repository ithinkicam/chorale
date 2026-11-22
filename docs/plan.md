# Chorale: plan

As of 2026-11-22. Draft for Cami's review, revised after a first research pass and the first interface sketches. Nothing has been built yet. The findings and sources behind this plan are in `docs/research/2026-11-22.md`.

## Outcome

Chorale is a fast native terminal app on Cami's Mac for running Claude Code, Codex and Antigravity CLI sessions side by side. It is built for low vision and for speaking rather than typing. A small background service holds the sessions, so an iOS companion can later show them, answer permission prompts and take dictated prompts from the phone.

The name is a choir piece, and it sounds like "corral", which is what you do with a herd of agents.

## Design principles

These are requirements, not preferences. A feature that breaks one of them isn't finished.

### Low vision first

1. **Works at large font sizes.** The whole interface must be usable at 60 columns and 20 rows. Text wraps to the real window width, never to a fixed column. The layout is a single column that stacks panes, and side panels appear only when the window is wide enough. Claude Code's open low-vision bug reports, about text running off screen and layouts that ignore the window size, are the cautionary tale.
2. **Your terminal's colours by default.** The default theme uses the terminal's own 16-colour palette, so Chorale inherits whatever contrast Cami has already set up. Themes with fixed colour values are optional, and a test checks each of them meets WCAG AAA contrast (Web Content Accessibility Guidelines, a ratio of at least 7 to 1).
3. **Never colour alone.** Every state, such as running, waiting for approval, failed or done, has a word or a symbol as well as a colour.
4. **An obvious focus that macOS Zoom can follow.** The focused row has a marker character and inverse video, and the real terminal cursor sits on it, so the system magnifier tracks it.
5. **No motion.** No spinners, shimmer or flicker. A steady status line shows progress in words, such as "Claude is working, 40 seconds". Only the parts of the screen that changed are redrawn.
6. **Sound and notification cues.** A distinct sound, and optionally a macOS notification, when an agent finishes or needs approval, so she doesn't have to watch the screen.
7. **Every key binding can be remapped,** with sensible defaults.

### Voice first

1. **Dictation is built in, not bolted on.** Press a key, speak, and the text lands in the prompt box for editing before it is sent. Nothing is sent without her say.
2. **Transcription is local.** Audio never leaves the Mac.
3. **Custom vocabulary.** Chorale knows the project names, file names and agent names in play, and corrects the transcript towards them. "Chorale" is the first entry, because a speech model will hear it as "corral".
4. **Dictation should cope with natural speech:** pauses, restarts and self-corrections.
5. **Say it again.** With text selected in the prompt box, the talk key replaces just the selection with new dictation, so one misheard phrase can be fixed without retyping the prompt.

### Speed

The terminal app starts in well under a second and never blocks on an agent. The point of building it natively is that it feels instant.

## Interface

The first sketches are on a private design canvas: https://claude.ai/artifact/CYQdsBMrHpnmQ3xfrJQRRQ. They show one session, a permission prompt, dictation, the session list, a wide two-pane window and the prompt editor, each at 60 columns and in three themes.

### A real editor in the prompt box

The prompt box works like a normal text box, not like terminal line editing where the cursor only moves with the arrow keys. This is a requirement.

1. Click to place the cursor. Drag to select. Double-click selects a word, triple-click a line, and Shift-click extends the selection.
2. Option with the arrow keys moves by word. Undo and redo work across the whole prompt.
3. Long prompts wrap softly, and clicks map to the right place across the wraps.
4. Copy goes to the Mac clipboard through OSC 52, the terminal escape code for setting the clipboard. Paste arrives as one block through bracketed paste.
5. Selected text uses the theme's accent colour, and the real terminal cursor stays at the insertion point, so macOS Zoom follows it.

Chorale turns on the terminal's mouse reporting, in its SGR mode, so it receives every click, drag and release with its position. Bubble Tea v2 delivers these as messages. The Bubbles text area doesn't handle mouse selection, so Chorale needs its own editor component. While mouse reporting is on, Option-drag still gives the terminal's own selection, which is the escape hatch for copying from the transcript.

### Session names

Each session is named for its work, such as "Bearer-check tests", not for its agent. The agent, model and project sit underneath as a subheading, such as "claude · opus 5.5 · commonplace-code".

The service asks a small, inexpensive model for a name from the first prompt, and may rename the session once if the work clearly changes. Cami has the last word: the `r` key renames a session, and a name she sets is never replaced.

### Choristers

Each session has a small chorister sprite in choir robes, with a white ruff collar and a hymnal. The agents differ in shape as well as robe colour, so they can be told apart without relying on colour: Claude has plain hair and a pink robe, Codex a dark choir cap with a gold tassel and a blue robe, Antigravity a gold halo and a green robe.

The pose follows the session's state and never animates. Singing with the hymnal open means working. A raised hand means it needs her. Eyes closed with the hymnal put away means done. The state is always written in words beside the sprite as well.

Sprites are 12 by 12 pixels. In terminals that support the Kitty graphics protocol, such as Ghostty and Kitty, they are drawn as real images scaled to fit 3 rows beside the session header. Other terminals get an 8 by 8 version drawn with half-block characters, two pixels per character cell. Both versions come from the same sprite source in the repo.

## Architecture

Chorale has three parts on the Mac, and later a fourth on the phone.

### The service, `chorale serve`

A long-running background process, started at login by `launchd`. It starts and supervises agent sessions, stores transcripts, and sends events to whichever clients are connected. Agents keep running when the terminal app is closed.

It talks to every agent through the Agent Client Protocol (ACP): JSON-RPC 2.0 over each agent's standard input and output. The stable protocol version is 1. A version 2 draft exists, and Chorale targets version 1 first.

All three agents speak ACP:

1. **Claude Code** through `claude-agent-acp`, from the `agentclientprotocol` GitHub organisation. It is built on the Claude Agent SDK and uses Cami's existing `claude` login. Chorale points it at her installed CLI with `CLAUDE_CODE_EXECUTABLE`.
2. **Codex** through `codex-acp`, from the same organisation. It drives the Codex App Server. Chorale points it at her installed CLI with `CODEX_PATH`.
3. **Antigravity CLI** natively, with `agy acp`, from version 1.5.0 onwards. Her Mac had 1.2.16 as of 2026-10-03, so it needs upgrading first.

Two rules follow from the research. First, Chorale never wraps or screen-scrapes an agent's terminal interface. Omnara tried that with Claude Code and called it "unfeasible to maintain". Second, the service removes `ANTHROPIC_API_KEY` and `OPENAI_API_KEY` from the environment it gives the adapters, because either one silently switches billing from her subscription to paid API use.

### The terminal app, `chorale`

A client of the service, written in Go with Bubble Tea v2, Lip Gloss v2 and Bubbles v2. Bubble Tea v2 has been stable since February 2026. It reports key releases in terminals that support them, can set the cursor position for the magnifier, and redraws only what changed.

### The listening helper, `chorale-listen`

A small Swift program that records from the microphone and transcribes locally on Apple Silicon. The terminal app starts it and reads the text back. Writing it in Swift means the iOS companion can share the speech code.

The engine sits behind one interface so it can be swapped. The order of trial is:

1. **FluidAudio running NVIDIA's Parakeet models.** Apache-2.0 code with permissively licensed models, runs on the Neural Engine, works on macOS and iOS, and is what Hex and many other dictation apps use. It has rare-word biasing for custom vocabulary in batch mode.
2. **Apple's SpeechTranscriber,** built into macOS 26 and iOS 26, so there is nothing to bundle.
3. **WhisperKit** as the fallback. Its own vocabulary feature is paid, but plain prompt conditioning works.

Whatever the engine, a final correction step matches the transcript against the vocabulary list: project names, file names in the current repo, agent names. This works with every engine.

### How clients talk to the service

The service does not pass raw ACP through to clients. It exposes its own smaller JSON-RPC API over WebSocket: list sessions, send a prompt, stream updates, list pending approvals, answer an approval, and replay what was missed after a reconnect. The terminal app connects over a local Unix socket. The iOS app connects over Tailscale, through Tailscale Serve, which gives the service an HTTPS address with a real certificate, so iOS will accept it.

### Code to study from Petal

Petal is MIT licensed and uses the same stack: Bubble Tea v2 and Coder's Go ACP library, `github.com/coder/acp-go-sdk` at v0.20.0. Parts worth studying or adapting, with the licence notice kept:

1. `internal/processgroup`, for starting adapters in their own process group and stopping the whole group cleanly.
2. `internal/acphost`, which confines an agent's file reads and writes to the project folder with `os.OpenRoot`, and caps terminal output.
3. The permission flow, which turns an agent's permission request into a message carrying a reply channel.
4. `internal/acpv2`, a transcript model that applies session updates in order.

The Go ACP library is pre-1.0 and its last release broke its API, so Chorale pins an exact version and keeps all ACP code in one package.

## Push-to-talk in a terminal

Most terminals don't tell a program when a key is released, only when it's pressed. There are three ways round it:

1. **Toggle mode.** Press once to start listening and once to stop. Works in every terminal, including Terminal.app. This is the default.
2. **True hold-to-talk** in terminals that report key releases through the Kitty keyboard protocol. Kitty and Ghostty do. iTerm2 probably does. WezTerm needs a setting turned on. Terminal.app does not. Bubble Tea v2 tells Chorale at startup whether the terminal supports it. The talk key must not be Enter, Tab or Backspace, which only report releases in a stricter mode.
3. **Double-tap to lock,** borrowed from Hex. Hold to talk briefly, or double-tap to keep listening hands-free until the next tap.

A global hotkey that works when the terminal isn't focused is possible later, but needs the macOS Input Monitoring permission.

## Microphone permission

This is the biggest risk in the voice design, so it gets its own section.

macOS charges microphone access to the app that launched the program, which for Chorale is Cami's terminal app. If that terminal doesn't declare that it uses the microphone, macOS denies access silently, with no prompt. Claude Code has open bug reports about exactly this.

The first task in milestone 2 is to check whether her terminal shows the microphone prompt. If it doesn't, `chorale-listen` becomes a tiny background app bundle with its own permission, which the terminal app talks to over a local socket.

## Milestones

Each milestone ends with something Cami can use, and nothing moves on until its checks pass.

### Milestone 0: skeleton

The repo layout, a Go module, a `Makefile`, continuous integration that runs the tests, the linter and the theme contrast check, and an `AGENTS.md` for agents working in the repo.

Done when `make test` passes in continuous integration.

### Milestone 1: one Claude session

The service starts a single Claude Code session through `claude-agent-acp`. The terminal app shows it in one column, sends prompts and shows replies. Permission requests appear as a clear numbered choice with a sound. The default theme uses the terminal palette, and a high-contrast fixed-colour theme ships alongside. The prompt box is the full mouse-driven editor from the start, because everything else is typed or dictated into it. The session header shows Claude's chorister in its half-block version.

Done when Cami can do a real small task in Chorale at her usual font size without reaching for Claude Code directly.

### Milestone 2: dictation

First the microphone permission check, then `chorale-listen` with FluidAudio and Parakeet, toggle-mode dictation into the prompt box, and the vocabulary correction step. Hold-to-talk and double-tap-to-lock follow where her terminal supports them. Apple's SpeechTranscriber is tested against her own speech as a comparison.

Done when Cami prefers dictating into Chorale over Wispr Flow for prompts.

### Milestone 3: three agents

Several sessions at once, a session list that works at 60 columns, Codex and Antigravity as the second and third agents, and switching between sessions by number or by name. Sessions get their work-based names, all three choristers are drawn, and the Kitty graphics version of the sprites is used where the terminal supports it. Transcripts are saved and searchable by keyword.

Done when she can run Claude, Codex and Antigravity sessions on the same project and move between them without losing her place.

### Milestone 4: agents handing off

One agent can hand a task to another, as Petal does, through a small built-in MCP server the service offers to each agent.

Done when Claude can ask Codex for a second opinion on a change and Cami sees both sessions in the list.

### Milestone 5: iOS companion

A SwiftUI app that connects to the service through Tailscale Serve. First version: see sessions, read the latest output at a large Dynamic Type size, approve or decline permission prompts, and dictate a prompt with the shared speech code. Approvals are shown on the session list itself, a pattern borrowed from the CC Pocket app.

The phone can't keep a connection open in the background, so notifications have to come from the service. Pushover is the first route: a one-off $4.99, with no Apple account needed. Apple's own push service comes later, sent directly from the service, once there is a paid developer account. Happy, an open-source mobile companion for Claude Code and Codex, is the main project to study for pairing and the approval screen.

Done when she can approve an agent's request from her phone without walking to the Mac.

## Not in scope for now

Load balancing across subscriptions, semantic search over transcripts, editor integration, and a web front end. Petal does the first three. Any of them can come back later.

## Open decisions

**Where the service runs.** This is the biggest one, and Cami decides. Her projects live in `~/Dev` on her Mac. The phone can only reach an always-on machine, and a sleeping Mac laptop can't be woken reliably over Tailscale. So there are three options:

1. Run the service on her Mac and accept that the phone works only while the Mac is awake.
2. Run it on plex-server, the always-on Mac mini, with project repos cloned there and synced through git.
3. Run it on both, with the phone talking to whichever is up.

The proposal is option 1 for milestones 1 to 4, and the decision is revisited before milestone 5.

**Which terminal Cami uses.** Ghostty, Kitty or iTerm2 make true hold-to-talk possible, and Terminal.app limits it to toggle mode. It also decides the microphone permission question, whether the 12-pixel sprites can be drawn as images, and whether Cmd shortcuts such as Cmd-Z reach Chorale. Many terminals keep Cmd keys for themselves, so every editor shortcut also has a Ctrl or Option form that works everywhere. Not urgent.

**The speech engine.** Parakeet through FluidAudio is the proposal. The final choice comes from testing accuracy on her own speech, and speed on her Mac. There are no published benchmarks for an M1. Decided in milestone 2.

## How the plan will be checked

Each milestone has a "done when" test that Cami runs on her Mac. Continuous integration runs the Go tests, the linter and the theme contrast check on every push. Transcription accuracy is measured against a small set of her own recordings, which stay on her machine and are never committed.

## Outward-facing steps that need separate approval

1. Creating the private GitHub repo `chorale`, done by Cami.
2. Upgrading Antigravity CLI on her Mac to 1.5 or later.
3. A `launchd` service that starts at login, installed only with her go-ahead in milestone 1.
4. Microphone permission for her terminal or for the `chorale-listen` helper, granted by her.
5. Pushover, a one-off $4.99, in milestone 5.
6. The paid Apple Developer Program, $99 a year, needed for TestFlight and Apple push notifications. A free account can install the app on her phone, but it expires after seven days.

## Policy note

Anthropic's terms let an individual sign in to Claude Code with their own subscription. They forbid third-party products from routing other people's requests through Claude subscription logins. Chorale is a personal tool run by its only user on her own login, which fits. If it is ever shared or hosted for others, that needs a careful reread first.
