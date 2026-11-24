# Proposal: add the keyboard prompt editor

## Why

Everything Cami asks an agent goes through the prompt box, whether she types it or, from milestone 2, dictates it. `docs/plan.md` requires that box to behave like a normal text box rather than terminal line editing. Building it first, before any agent is connected, means it can be tested thoroughly on its own, and the next two changes can rely on it.

This is the first of three changes that make up milestone 1, "one Claude session", in `docs/plan.md`. The second adds the background service and its connection to Claude Code. The third adds the session screen, which joins the two.

## What Changes

- A prompt box appears at the bottom of the empty session list, focused, with the real terminal cursor at the insertion point.
- The prompt box is a keyboard editor, written for Chorale:
  - It moves the cursor by character, by word, by line, and to the start or end of the prompt.
  - It selects with Shift and the movement keys, and selects everything with one key.
  - It deletes by character, by word, and to the start or end of a line.
  - It has undo and redo across the whole prompt.
  - It wraps soft at word boundaries to the real window width, counting wide characters correctly.
  - It accepts bracketed paste as one block, which one undo removes.
  - It copies and cuts the selection to the Mac clipboard through OSC 52.
  - It shows the selection in the accent colour and in inverse video, so the selection is never colour alone.
  - It grows with the prompt up to a limit, then scrolls to keep the cursor in view, saying in words how many lines are out of sight.
- Every editor action has a Ctrl or Option key, which works in every terminal. In Ghostty, which `cmux` is built on, the usual Cmd keys work too where the terminal passes them through.
- Enter sends the prompt, and Shift-Enter, Option-Enter or Ctrl-J starts a new line. With no agent connected yet, sending keeps the prompt in the box and says that no agent is connected. The third change replaces this with a real send.
- **BREAKING:** `q` no longer quits from the session list, because the focused prompt box types it. Ctrl-C quits instead. If the prompt has text, the first Ctrl-C clears it, and undo brings it back. `q` still quits from the window-too-small message, where there is no prompt box.

## Capabilities

### New Capabilities

- `prompt-editor`: the prompt box and how it edits text. It covers movement, selection, deletion, undo and redo, soft wrapping, paste, copy and cut, sending, scrolling, and where the cursor sits.

### Modified Capabilities

- `app-shell`: the empty session list now shows the prompt box, with the cursor in it, and its key line changes. Quitting moves from `q` to Ctrl-C on the session list, and stays on `q` for the window-too-small message.

## Out of scope

1. Mouse editing: clicks, drags, double and triple clicks, and Shift-click. These come in milestone 2, as `docs/plan.md` says.
2. Sending a prompt to an agent, and anything about sessions. These come in the second and third changes of milestone 1.
3. Remapping keys. Every binding sits in one table, so a later change can make them configurable, but this change has no settings file.
4. The high-contrast fixed-colour theme, which arrives with the session screen.
5. Dictation and "say it again", which come in milestone 2.
6. Reading the clipboard. Paste arrives from the terminal's own paste command, as bracketed paste.

## Impact

1. New package `internal/editor`, holding the text buffer, the edit history, the wrapping and the key table, with unit tests.
2. A new screen part for the prompt box in `internal/screens`, and changes to `internal/app` to route keys, paste and quitting.
3. New interface, interaction and terminal-level tests, and new golden files. The golden files for the empty session list change, because the prompt box and the new key line appear on it.
4. No new dependencies. Bubble Tea v2.1.2 already provides bracketed paste, OSC 52 copy, the Cmd modifier and keyboard enhancements, and `github.com/charmbracelet/x/ansi`, already in `go.mod`, measures text width.
