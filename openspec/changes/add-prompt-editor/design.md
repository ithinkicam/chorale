# Design: add the keyboard prompt editor

## Context

Before this change, Chorale's only screens are the empty session list and the window-too-small message, both drawn by `internal/app` from small render functions in `internal/screens`. The app quits on `q`. The test harness can drive a model in a test, through `termtest.StartModel`, and the real binary in a pseudo-terminal, through `termtest.Start`. Both read the screen through the `x/vt` emulator.

The plan rules out the Bubbles text area, because it can't do mouse selection, and milestone 2 adds the mouse to this same editor. So the editor is Chorale's own code from the start.

Cami's terminal is `cmux`, which is built on Ghostty. Ghostty supports the Kitty keyboard protocol, so Shift-Enter and Ctrl-Shift keys can be told apart there. By default Ghostty also turns Option-Left and Option-Right into the word-movement codes, and Cmd-Left, Cmd-Right and Cmd-Backspace into Ctrl-A, Ctrl-E and Ctrl-U. Ghostty keeps Cmd-C, Cmd-V and Cmd-A for itself unless they are rebound. These defaults are from Ghostty's documentation, and `cmux` may change them, so task 1.1 checks them on her Mac.

## Goals / Non-Goals

**Goals:**

1. An editor whose logic has no screen in it, so nearly every behaviour is a fast unit test.
2. One table of key bindings, so remapping later is a settings change, not a rewrite.
3. Wrapping and cursor placement that are correct for wide characters, emoji and combining accents, because milestone 2 maps mouse clicks through the same layout.
4. Every editor action reachable with a Ctrl or Option key in any terminal, and with the Cmd key where the terminal passes it through.

**Non-Goals:**

1. Mouse editing, which is milestone 2.
2. Sending a prompt anywhere. Enter reaches a send hook that, in this change, only shows the "no agent" message.
3. A settings file or remapping.
4. Very large prompts beyond what a person types or pastes. The editor stays responsive up to 10,000 lines, which task 3.9 checks, but it isn't built for megabytes.

## Decisions

### The editor is a plain Go package with no screen

`internal/editor` holds the prompt as a slice of graphemes, the insertion point and the selection anchor as grapheme indexes, and the edit history. Its operations are ordinary methods, such as insert, delete a word back, move by word and undo, and each returns the new state. A separate `Layout` function takes the text and a width and returns the rows on screen, with the grapheme index where each row starts and the column of every grapheme. Vertical movement, Home and End, Ctrl-U and Ctrl-K, and the cursor's screen position all use the layout. Milestone 2's click mapping will use it too.

The alternative was to edit a string and recompute positions from byte offsets. Grapheme indexes make "never split a character" true by construction, and keep the code that moves by character simple.

Graphemes and widths come from `github.com/charmbracelet/x/ansi`, which is already a direct dependency, and which uses `rivo/uniseg` underneath. A wide character counts as two columns, and an emoji with a variation selector is measured the way `x/ansi` measures it, so the layout and the renderer always agree.

### Words

A word is a run of letters, digits and underscores. Moving or deleting by word skips any spaces and punctuation first, then the word. That is how macOS text fields behave with Option-arrow, and it makes "run the tests" move one word per press.

### Wrapping

The layout breaks a logical line at the last space that fits. A word wider than the box is broken at the last grapheme that fits. The space at a wrap point stays at the end of the row it ends, so the text is unchanged, and the cursor can sit after it. A row is never wider than the box, and a wide character that would straddle the edge moves to the next row whole.

### Undo history

The history is a list of snapshots of text, insertion point and anchor, with a pointer into it. A prompt is small, so whole snapshots are simpler and safer than storing reversible operations, and the history keeps at most 200 snapshots. Typing extends the current snapshot group until a word boundary: a space or new line typed after a non-space ends the group. Any movement, paste, cut, delete-word or clear starts a new group, so each is one undo. A new edit after an undo truncates the redo part.

### Paste

Bubble Tea v2 already turns on bracketed paste and delivers the pasted text as one `tea.PasteMsg`. The editor cleans it before inserting: carriage-return and new-line pairs and lone carriage returns become new lines, tabs become four spaces, and every other control character is removed, including the escape byte. Removing the escape byte is what keeps pasted escape sequences from ever reaching the screen. The paste is one undo group.

### Copy and cut

Copying uses `tea.SetClipboard`, which writes an OSC 52 sequence for the system clipboard. Nothing reads the clipboard: paste always comes from the terminal's own paste command. Ghostty accepts OSC 52 writes by default.

### Keys

All bindings live in one table in `internal/editor`, mapping a key, as Bubble Tea names it, to an action. The defaults:

1. Left and Right, or Ctrl-B and Ctrl-F, move by character.
2. Option-Left and Option-Right, or Alt-B and Alt-F, move by word.
3. Up and Down move by row on screen.
4. Home and End, Ctrl-A and Ctrl-E, and Cmd-Left and Cmd-Right, move to the start and end of the row.
5. Option-Up and Option-Down, Ctrl-Home and Ctrl-End, and Cmd-Up and Cmd-Down, move to the start and end of the prompt.
6. Shift with any of these extends the selection. Ctrl-Shift-A and Cmd-A select everything.
7. Backspace and Delete, or Ctrl-H and Ctrl-D, delete a character.
8. Option-Backspace and Ctrl-W delete the word before. Option-Delete and Alt-D delete the word after.
9. Ctrl-U and Ctrl-K delete to the start and end of the row.
10. Ctrl-Z and Cmd-Z undo. Ctrl-Y, Ctrl-Shift-Z and Cmd-Shift-Z redo.
11. Ctrl-C and Cmd-C copy a selection. With no selection, Ctrl-C clears the prompt or, if it is empty, quits.
12. Ctrl-X and Cmd-X cut.
13. Enter sends. Shift-Enter, Option-Enter and Ctrl-J insert a new line.

Ctrl-Shift-A and Ctrl-Shift-Z only reach Chorale in terminals that speak the Kitty keyboard protocol, which Ghostty does. Elsewhere, Ctrl-Y still redoes. Select-all has no form that works in every terminal, because Ctrl-A is the long-standing key for the start of a line, and Ghostty sends it for Cmd-Left. In Terminal.app, select-all is Option-Up then Shift-Option-Down. That gap is accepted for now and written in the key help in `docs/`.

Ctrl-C does three things depending on the state. That follows the editors Cami already uses, where Ctrl-C copies when there is a selection, and Claude Code, where it clears a prompt before it quits. A prompt is never lost to one key press, because clearing is undoable.

### Keyboard enhancements

The view asks for key disambiguation by setting the Kitty keyboard flags in `tea.View`. Bubble Tea turns them off on exit, and the existing clean-exit test already checks that keyboard enhancements are off afterwards. The app doesn't need to know whether the terminal agreed. Without them, Shift-Enter arrives as Enter, so Option-Enter and Ctrl-J are always there.

### The prompt box on screen

The empty session list stacks, from the top: the "♪ chorale" line, a blank line, "No sessions yet.", the flexible space, the prompt box, and the key line "Enter send · Ctrl-C quit". The prompt box is a label row, "Prompt", followed by its text rows, then a status row when there is something to say. The label row shows "3 lines above" when rows are hidden above, and the status row shows "2 lines below" when rows are hidden below. The status row also shows "No agent is connected yet. Your prompt is kept." after a send, until the next edit. The status row is empty when there is nothing to say, so the box's height stays steady while typing.

The box shows one text row when empty and grows up to 8 text rows at 20 rows tall, or 12 at 30 rows and above. In between, the limit is the window height minus 12. When the text is taller, the box scrolls by the least amount that keeps the insertion point's row in view. The box's rows are indented by two columns, so the text width is the window width minus two.

The selection is drawn with inverse video, plus the accent colour of the default theme, which is the terminal palette's blue. With `NO_COLOR` set, only the inverse video is written. Inverse video is not a colour code, so the colour-off tests still pass.

### Where the send hook lives

Enter calls a `Send` function that the app passes to the prompt box. In this change, the app's function sets the "no agent" status and keeps the prompt. The session screen change replaces it with a real send, and that change will decide what happens to the prompt then.

### The scenario check and replaced requirements

`scripts/check-scenarios` reads every scenario in `openspec/specs/` and in each open change, and matches them to passing tests by name. This change removes one `app-shell` requirement and renames another. Their old scenarios would stay in the main spec until archive, with no test to match. So the check learns to skip a main-spec requirement when an open change removes, renames or modifies it. That change's own copy of the requirement is checked instead. It reads only the `## REMOVED`, `## RENAMED` and `## MODIFIED` headings in the change's delta specs, and matches requirement names the way OpenSpec does, ignoring case and spacing.

### The harness reads cell attributes

The selection scenarios check inverse video and colour on particular cells. `termtest.Screen` gains a `Cell(x, y)` method that returns the cell's text, whether it is inverse, and its foreground and background colour, read from `x/vt`. It is tested against a known escape sequence before any editor test uses it.

## Risks / Trade-offs

**Ghostty and `cmux` key defaults may differ from what's written here.** Task 1.1 records what reaches Chorale on Cami's Mac, using a small key logging test program that the change doesn't ship. If a default key never arrives there, the binding table changes before any test is written against it.

**Ctrl-C has three meanings.** That is efficient, but a copy might be mistaken for a quit. Quitting needs an empty prompt with nothing selected, so the worst case is a second press.

**Snapshot undo uses more memory than operation undo.** At 200 snapshots of a 10,000-line prompt, that is about 100 MB in the worst case. The history drops the oldest snapshots when their total size passes 16 MB, which bounds it. Task 3.9 covers this.

**OSC 52 can be turned off in a terminal.** If it is, copy silently does nothing. Chorale can't detect this. The key help in `docs/` mentions it.

**The "no agent" message is temporary.** It ships to `main` and stays until the session screen change. That is honest about what Chorale can do in between, and costs one requirement that the third change modifies.

## Test strategy

Each requirement and the layer from `docs/testing.md` that covers it:

1. Typing inserts at the cursor: unit tests on the editor, and interface snapshot tests for what shows on screen and where the cursor is.
2. Cursor at the insertion point: interface snapshot tests, and a terminal-level test for the empty prompt box in the real binary.
3. Moving by character and by word: unit tests, plus one interaction test per scenario through key messages.
4. Moving by row and to the ends: unit tests on the layout and movement, plus interaction tests.
5. Selecting with the keyboard: unit tests and interaction tests.
6. The selection is shown without relying on colour: terminal-level tests reading cell attributes and raw bytes, with colour on and with `NO_COLOR=1`.
7. Deleting: unit tests and interaction tests.
8. Undo and redo: unit tests and interaction tests.
9. Soft wrapping: unit tests on the layout, with golden snapshots at 60 by 20 and 120 by 30.
10. Pasting: unit tests on the cleaning, an interaction test sending a paste message, and a terminal-level test sending a real bracketed paste.
11. Copying and cutting through OSC 52: terminal-level tests reading the raw bytes.
12. Clearing the prompt with Ctrl-C: interaction test.
13. New lines and sending: interaction tests, and a terminal-level test that Shift-Enter arrives as a new line when the Kitty keyboard protocol is on.
14. Growing and scrolling: interface snapshot tests at 60 by 20, with colour on and off.
15. Cmd keys where the terminal passes them: interaction tests sending key messages with the Cmd modifier.
16. Telling Shift-Enter apart from Enter: terminal-level test reading the raw bytes and the modes at exit.
17. Empty session list with the prompt box: interface snapshot tests at 60 by 20 and 120 by 30, and the terminal-level cursor test.
18. No motion and Quitting: the existing terminal-level tests, updated for Ctrl-C, plus a test that `q` types on the session list.

Golden snapshots of the empty session list with the prompt box, a wrapped prompt, a selection and a scrolled prompt are taken at 60 by 20 in the default theme, and with colour off. The default theme is the only one shipped in this change.
