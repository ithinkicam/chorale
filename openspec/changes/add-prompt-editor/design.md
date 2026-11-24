# Design: add the keyboard prompt editor

## Context

Before this change, Chorale's only screens are the empty session list and the window-too-small message, both drawn by `internal/app` from small render functions in `internal/screens`. The app quits on `q`. The test harness can drive a model in a test, through `termtest.StartModel`, and the real binary in a pseudo-terminal, through `termtest.Start`. Both read the screen, including the cursor, through the `x/vt` emulator.

The plan rules out the Bubbles text area, because it can't do mouse selection, and milestone 2 adds the mouse to this same editor. So the editor is Chorale's own code from the start.

Cami's terminal is `cmux`, which is built on Ghostty. Ghostty supports the Kitty keyboard protocol, so Shift-Enter and Ctrl-Shift keys can be told apart there. By default Ghostty also turns Option-Left and Option-Right into the word-movement codes, and Cmd-Left, Cmd-Right and Cmd-Backspace into Ctrl-A, Ctrl-E and Ctrl-U. Ghostty keeps Cmd-C, Cmd-V and Cmd-A for itself unless they are rebound, and may keep Cmd-Z, Cmd-Up and Cmd-Down too. These defaults are from Ghostty's documentation, and `cmux` may change them, so task 1.1 checks them on her Mac.

## Goals / Non-Goals

**Goals:**

1. An editor whose logic has no screen in it, so nearly every behaviour has a fast unit test.
2. One table of key bindings, so remapping later is a settings change, not a rewrite.
3. Wrapping and cursor placement that are correct for wide characters, emoji and combining accents, because milestone 2 maps mouse clicks through the same layout.
4. Every editor action reachable with a Ctrl or Option key in any terminal, and with the Cmd key where the terminal passes it through.
5. No prompt lost to a slip of the finger: clearing is undoable, quitting asks first, and every state the box is in is said in words.

**Non-Goals:**

1. Mouse editing, which is milestone 2.
2. Sending a prompt anywhere. Enter reaches a send hook that, in this change, only shows the "no agent" message.
3. A settings file or remapping.
4. Very large prompts beyond what a person types or pastes. The editor stays responsive up to 10,000 lines, which a benchmark tracks, but it isn't built for megabytes.

## Decisions

### The editor is a plain Go package with no screen

`internal/editor` holds the prompt as a slice of graphemes, the insertion point and the selection anchor as grapheme indexes, and the edit history. Its operations are ordinary methods, such as insert, delete a word back, move by word and undo. A separate `Layout` function takes the text, a width and a width method, and returns the rows on screen, with the grapheme index where each row starts and the column of every grapheme. Vertical movement, Home and End, Ctrl-U and Ctrl-K, and the cursor's screen position all use the layout. Milestone 2's click mapping will use it too.

The alternative was to edit a string and recompute positions from byte offsets. Grapheme indexes make "never split a character" true by construction, and keep the code that moves by character simple.

Inserting text re-segments the graphemes around the insertion point: the grapheme before it, the new text and the grapheme after it are joined and split again. So a combining accent typed on its own joins the letter before it, as the spec requires, and an emoji built from several keystrokes ends up as one grapheme.

### Character widths match the renderer

Graphemes and widths come from `github.com/charmbracelet/x/ansi`, which is already a direct dependency. Bubble Tea's renderer measures text one of two ways. When the terminal reports support for grapheme clustering, mode 2027, it measures whole graphemes. Otherwise it measures code point by code point, the way `wcwidth` does. The two disagree on some emoji, such as a heart followed by the variation selector that asks for emoji style, which is two columns as a grapheme and one column code point by code point. Ghostty reports mode 2027, and Terminal.app doesn't.

So `Layout` takes the width method as an argument, and the app passes whichever one the renderer is using. Unit tests cover both methods on the heart with its variation selector. Task 1.1 confirms how the app learns which method is in use, and whether `cmux` reports mode 2027.

### Words

A word is a run of letters, digits and underscores, and an emoji is a word of its own. Every other character, including symbols such as `+`, `=`, `<` and `|`, counts as punctuation. Moving or deleting by word skips any spaces and punctuation first, then one word. That is how macOS text fields behave with Option-arrow and Option-Backspace: "run the tests" moves one word per press, and Option-Backspace after "go 🎵" removes only the emoji.

### Wrapping and the insertion point at a wrap

The text width is the window width minus three: two columns of indent and one column kept free at the right edge. The free column means that a cursor after the last character of a full row is always drawn on screen, never past the edge, where terminals behave inconsistently.

The layout breaks a logical line after the last space that fits. A word wider than the box is broken at the last grapheme that fits. The space at a wrap point stays at the end of the row it ends when it fits there, so the text is unchanged. A row is never wider than the text width, so when a word exactly fills the rest of a row and the space after it doesn't fit, the word moves down to the next row. When a row has no earlier space, the space that doesn't fit starts the next row. A wide character that would straddle the edge moves to the next row whole.

At a soft wrap, the end of one row and the start of the next are the same grapheme index. The insertion point carries a flag saying which of the two it is drawn at. End, Ctrl-E and Cmd-Right set it to the end of the row. Everything else sets it to the start of the next row. Up and Down keep a remembered screen column until some other key is pressed.

The row a key moves from is the row the cursor is drawn on, worked out before the flag changes. So Home after End goes to the start of the same row, and a second End stays put.

Up on the first row moves to the start of the prompt, and Down on the last row moves to its end, as in macOS text fields. Left and Right with a selection, and without Shift, move to the selection's start or end and clear it, rather than moving one character.

### Undo history

The history is a list of snapshots, each holding the text as one string, the insertion point and the anchor, with a pointer into it. Storing the text as a string, rather than the grapheme slice, keeps a snapshot close to the size of the text itself. Whole snapshots are simpler and safer than storing reversible operations.

The history keeps at most 200 snapshots, and drops the oldest when the total size of their text passes 16 MB. A 10,000-line prompt of about 500 KB still keeps around 30 undo steps under that cap. The history reports its own size, and a unit test checks both limits directly.

Typing extends the current snapshot group until a word boundary: a space or new line typed after a non-space ends the group. A run of Backspace presses is one group too, and so is a run of Delete presses, as in macOS text fields, so holding Backspace can't push earlier text out of the history. Any movement, paste, cut, delete-word or clear starts a new group, so each is one undo. A new edit after an undo truncates the redo part.

### Paste

Bubble Tea v2 already turns on bracketed paste and delivers the pasted text as one `tea.PasteMsg`. The editor cleans it before inserting, in this order:

1. Carriage-return and new-line pairs, and lone carriage returns, become new lines.
2. Every escape sequence is removed whole, using the `x/ansi` parser. That covers control sequences such as colour codes, operating system commands such as clipboard and title sequences, device control strings, and a lone escape byte.
3. Tabs become four spaces.
4. Every remaining control character is removed, including the C1 range and delete.

Removing whole sequences matters: removing only the escape byte would leave text such as `[31m` behind. But a paste must never lose ordinary text to a broken sequence. An escape byte followed by a control character, such as a new line, is a lone escape, and the text after it is kept. A control string that never ends, such as an operating system command with no terminator, loses only its introducer, and the rest of the paste is kept. C1 introducers written as Unicode characters, such as U+009B, are treated the same as their escape forms.

Undo restores a snapshot by splitting its text into graphemes again, so it keeps the insertion point and the anchor within the restored text. The paste is one undo group, and it replaces any selection.

### Copy and cut

Copying uses `tea.SetClipboard`, which writes an OSC 52 sequence for the system clipboard. Nothing reads the clipboard: paste always comes from the terminal's own paste command. Ghostty accepts OSC 52 writes by default, and task 1.1 checks that `cmux` does too.

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
11. Ctrl-C and Cmd-C copy a selection. With no selection, Ctrl-C clears the prompt or, if it is empty, asks before quitting.
12. Ctrl-X and Cmd-X cut.
13. Enter sends. Shift-Enter, Option-Enter and Ctrl-J insert a new line.
14. Tab does nothing. It is kept for moving between areas once the session screen has more than one.

Only Ctrl-J inserts a new line in every terminal. Shift-Enter needs the Kitty keyboard protocol, which Ghostty speaks. Without it, Shift-Enter arrives as Enter and sends. Option-Enter needs the terminal to send Option as Alt, which is a setting in most macOS terminals. The key help says so, and names Ctrl-J as the key that always works.

Ctrl-Shift-A and Ctrl-Shift-Z only reach Chorale in terminals that speak the Kitty keyboard protocol. Elsewhere, Ctrl-Y still redoes. Select-all has no form that works in every terminal, because Ctrl-A is the long-standing key for the start of a line, and Ghostty sends it for Cmd-Left. In Terminal.app, select-all is Option-Up then Shift-Option-Down. That gap is accepted for now and written in the key help in `docs/keys.md`.

Ctrl-C does three things depending on the state. That follows the editors Cami already uses, where Ctrl-C copies when there is a selection, and Claude Code, where it clears a prompt before it quits. The key line always names what Ctrl-C will do next, and the status row confirms what it did.

### Quitting asks first

Ctrl-C on an empty prompt, with nothing selected, shows "Press Ctrl-C again to quit." and sets a pending-quit flag. A second Ctrl-C quits. Any other key, paste or resize clears the flag and the message. There is no time limit, because a timer would change the screen without any input, which the "No motion" requirement rules out. Clearing and then quitting takes three presses, so a prompt and its undo history are never lost to a quick double tap.

The window-too-small message uses the same flag. There, Ctrl-C always asks first, whatever the prompt holds, because the prompt can't be seen. `q` no longer quits anywhere.

### Keyboard enhancements

Bubble Tea v2.1.2 always asks for key disambiguation, the first Kitty keyboard flag, and turns it off on exit. Chorale doesn't need to set anything. A terminal-level test records that the request is sent and released, so a Bubble Tea upgrade that changes it is caught. Where the terminal supports the protocol, Shift-Enter arrives as its own key.

### The prompt box on screen

The empty session list stacks, from the top: the "♪ chorale" line, a blank line, "No sessions yet.", the flexible space, the prompt box, and the key line.

The prompt box is a label row, its text rows, and a status row:

1. The label row reads "▸ Prompt". The triangle is the focus marker, there for when the session screen adds a second area. When rows are out of sight, the label row adds the counts, such as "▸ Prompt · 3 lines above · 2 lines below". Both counts sit on the label row so that they never compete with a message for the status row's width.
2. The text rows are indented by two columns, so the text lines up under the "P" of "Prompt".
3. The status row shows one message at a time. While text is selected and no other message is showing, it says how many characters are selected. The other messages are "Copied.", "Cut.", "Prompt cleared. Ctrl-Z brings it back.", "No agent is connected yet. Your prompt is kept." and "Press Ctrl-C again to quit." Each one stays until the next key or paste. The longest is 47 characters, which fits at 60 columns. The status row is empty when there is nothing to say, so the box's height stays steady while typing.

The key line reads "Enter send · Ctrl-C copy" while text is selected, "Enter send · Ctrl-C clear" while the prompt has anything in it, and "Enter send · Ctrl-C quit" while it is empty.

The box shows one text row when empty and grows to at most the window height minus 12 text rows, and never more than 12. That is 8 rows at 20 rows tall, and 12 from 24 rows up. When the text is taller, the box scrolls by the least amount that keeps the insertion point's row in view.

The selection is drawn with inverse video, plus the accent colour of the default theme, which is the terminal palette's blue. With `NO_COLOR` set, only the inverse video is written. Inverse video is not a colour code, so the colour-off tests still pass.

The "▸" triangle and the "♪" in the title are both characters whose width Unicode leaves to the terminal. Ghostty and Terminal.app both draw them one column wide, and `x/ansi` measures them that way.

### The window-too-small message

The message keeps its order: the size Chorale needs, then "Your prompt is kept." when the prompt has anything in it, then "Ctrl-C quit" or "Press Ctrl-C again to quit.". When the window is too short for all of them, the later ones are left out. The editor's state is untouched while the message shows, so the prompt, the insertion point and the history are all there when the window grows back.

### Where the send hook lives

Enter calls a `Send` function that the app passes to the prompt box. In this change, the app's function sets the "no agent" status and keeps the prompt. The session screen change replaces it with a real send, and that change will decide what happens to the prompt then.

### The scenario check and replaced requirements

`scripts/check-scenarios` reads every scenario in `openspec/specs/` and in each open change, and matches them to passing tests by name. This change removes one `app-shell` requirement, renames another, and modifies several. Their old scenarios stay in the main spec until archive, with no test to match. So the check learns to skip a main-spec requirement that an open change replaces. The change's own copy of the requirement is checked instead.

The rule is narrow, so it can't hide a real gap:

1. It reads only the `## REMOVED`, `## RENAMED` and `## MODIFIED` sections of the change's delta specs, and only for the capability that delta belongs to.
2. It matches requirement names exactly, after trimming spaces at either end, as OpenSpec does.
3. A removed requirement is skipped. A modified requirement is skipped, because the change's copy is checked. A renamed requirement is skipped only when the change also modifies it under its new name. A bare rename keeps its scenarios checked, because their names don't change.
4. A delta naming a requirement that the main spec doesn't have is an error, so a stale change is caught.
5. The check prints every requirement it skipped, and which change made it skip.

### The harness reads cell attributes

The selection scenarios check inverse video and colour on particular cells. `termtest.Screen` gains a `Cell(x, y)` method that returns the cell's text, whether it is inverse, and its foreground and background colour, read from `x/vt`. It is tested against a known escape sequence before any editor test uses it.

## Risks / Trade-offs

**Ghostty and `cmux` key defaults may differ from what's written here.** Task 1.1 records what reaches Chorale on Cami's Mac, using a small key logging program that the change doesn't ship. If a binding never arrives, the change to the key table goes back to Cami as an update to this change before any test is written against it.

**Ctrl-C has three meanings.** That is efficient, but a copy might be mistaken for a clear, or a clear for a quit. The key line says what it will do, the status row says what it did, clearing is undoable, and quitting always asks first.

**Snapshot undo uses more memory than operation undo.** The 16 MB cap bounds it, at the cost of fewer undo steps on a very large prompt.

**OSC 52 can be turned off in a terminal.** If it is, copy does nothing, though the status row still says "Copied.". Chorale can't detect this. The key help mentions it.

**The "no agent" message is temporary.** It ships to `main` and stays until the session screen change. That is honest about what Chorale can do in between, and costs one requirement that the third change modifies.

## Test strategy

Two kinds of test sit behind each behaviour.

Unit tests in `internal/editor` cover the editor's operations, the layout, the width methods, the paste cleaning and the history limits. They are named after what they check, not after scenarios, because a scenario starts from the session list and presses keys, which a unit test doesn't do.

Every scenario has a test named after it, and that test does what the scenario says. Scenarios that press keys drive the app model through `termtest.StartModel`, sending key and paste messages, and read the screen and the real cursor position from `x/vt`. Scenarios about bytes, colour codes, OSC 52, bracketed paste, keyboard enhancements or quitting run the real binary through `termtest.Start`.

By requirement:

1. Typing, cursor, movement, selection, deletion, undo, wrapping, clearing, the key line, sending, scrolling and the Cmd keys: model tests through key messages, reading the screen and cursor, on top of the unit tests.
2. The selection shown without colour: terminal-level tests reading cell attributes and raw bytes, with colour on and with `NO_COLOR=1`.
3. Pasting: a terminal-level test sending real bracketed paste bytes for "Escape sequences in a paste are removed", and model tests for the rest.
4. Copying and cutting: terminal-level tests decoding the OSC 52 payload from the raw bytes.
5. Telling Shift-Enter apart from Enter: terminal-level tests reading the raw bytes, the modes at exit, and a Shift-Enter sent in the Kitty form.
6. The app shell's requirements: interface snapshot tests at 60 by 20 and 120 by 30, and the existing terminal-level tests, updated for the two Ctrl-C presses.

Golden snapshots of the empty session list with the prompt box, a wrapped prompt, a selection and a scrolled prompt are taken at 60 by 20 in the default theme, and with colour off. The default theme is the only one shipped in this change.

A benchmark in `internal/editor` times inserting a 10,000-line paste, laying it out and moving to its end. It is a record to compare against, not a pass or fail test, because timings in CI vary too much to gate on.
