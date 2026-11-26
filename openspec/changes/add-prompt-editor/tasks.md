# Tasks

Every behaviour task below comes in a pair: first the failing test, then the code. Run each test after writing it and confirm it fails for the expected reason before writing the code. Tick a test-writing task once its tests fail for that reason. Tick a code task once its tests pass and `make test` and `make lint` pass.

Every scenario has a test named after it, as `docs/testing.md` describes, and that test asserts the scenario's THEN clauses directly. Unit tests in `internal/editor` are named after what they check. A golden file is a record of the screen, never the thing that makes a test pass.

## 1. Groundwork

- [x] 1.1 On Cami's Mac, in `cmux`, run a throwaway key logger built on Bubble Tea v2, kept in the scratchpad and never committed. Record in `design.md`, under "Keys", which of the bindings listed there reach Chorale and as what, including Shift-Enter, Option-Enter, Cmd-Z, Cmd-Up and Cmd-Down. Also record whether an OSC 52 copy reaches the macOS clipboard, whether `cmux` reports mode 2027, and how the app learns which width method Bubble Tea's renderer is using. Verify by reading the recorded list back to Cami, because she runs the logger. If any binding in the key table never arrives, stop and bring the change to the key table to Cami as an update to this change, before any test is written against it
- [x] 1.2 Write unit tests for `scripts/check-scenarios` covering each rule in `design.md` under "The scenario check and replaced requirements": a removed requirement is skipped, a modified one is skipped and the change's copy is required, a renamed and modified one is skipped, a bare rename is not skipped, a delta in one capability doesn't skip a same-named requirement in another, a name that differs only in case is not matched, a delta naming a missing requirement is an error, and every skip is printed. Run them and confirm they fail because the check has none of these rules
- [x] 1.3 Teach the check those rules, and verify the tests from 1.2 and the existing check tests pass
- [x] 1.4 Write a test for a new `termtest.Screen.Cell` method, feeding known escape sequences for inverse video and a palette foreground colour, and asserting the text, the inverse flag and the colour of chosen cells. Run it and confirm it fails because the method doesn't exist
- [x] 1.5 Implement `Screen.Cell` over `x/vt`, and verify the test from 1.4 passes

## 2. Editor text, layout and movement

- [x] 2.1 Write unit tests in `internal/editor` for inserting text at the insertion point: plain letters, inserting in the middle, wide characters and emoji, and a combining accent inserted on its own joining the letter before it. Run them and confirm they fail because the package doesn't exist
- [x] 2.2 Implement the grapheme text and insertion, with re-segmenting around the insertion point, and verify the tests from 2.1 pass
- [x] 2.3 Write unit tests for `Layout` at a text width of 57: wrapping at a space with the space kept at the end of its row, breaking a word wider than the row, wide characters moving to the next row whole, a row filled exactly to the width, and the heart with its emoji variation selector measured as two columns by graphemes and one by code points. Run them and confirm they fail
- [x] 2.4 Implement `Layout` with both width methods, as `design.md` describes, and verify the tests from 2.3 pass
- [x] 2.5 Write unit tests for movement: by character over an emoji and across a soft wrap, by word with the word rules in `design.md`, by row keeping the remembered column across hard and soft wraps, to the start and end of a row with the end-of-row flag set by End, and to the start and end of the prompt. Run them and confirm they fail
- [x] 2.6 Implement the movements, the remembered column and the end-of-row flag, and verify the tests from 2.5 pass

## 3. Selection, deletion, history and paste

- [x] 3.1 Write unit tests for selection: extending with Shift movement, replacing the selection by typing, selecting everything, clearing it by moving, and counting the selected characters. Run them and confirm they fail
- [x] 3.2 Implement the selection anchor, and verify the tests from 3.1 pass
- [x] 3.3 Write unit tests for deletion: Backspace over an emoji, a word before and after, a word after an emoji, to the start and end of a row, a character across a soft wrap, and clearing the whole prompt. Run them and confirm they fail
- [x] 3.4 Implement the deletions and clearing, and verify the tests from 3.3 pass
- [x] 3.5 Write unit tests for the history: typing undone a word at a time, redo, undo restoring the selection, a new edit discarding redo, at most 200 snapshots kept, and the total stored text kept under 16 MB, read from the history's own size report after 300 edits to a 10,000-line prompt. Run them and confirm they fail
- [x] 3.6 Implement the history, with word-sized typing groups and both limits, and verify the tests from 3.5 pass
- [x] 3.7 Write unit tests for paste cleaning, with one case for each step in `design.md` under "Paste": carriage returns, a colour code, an OSC sequence, a device control string, a lone escape byte, a tab, and C1 and delete characters. Add tests that a paste replaces a selection and is one undo group. Run them and confirm they fail
- [x] 3.8 Implement paste cleaning and insertion, and verify the tests from 3.7 pass
- [x] 3.9 Add a benchmark that inserts a 10,000-line paste, lays it out at width 57 and moves to its end. Run it once and record the timings in the commit, as a record to compare against

## 4. Key table and the prompt box on screen

- [x] 4.1 Write unit tests for the key table: every binding in `design.md` maps to its action, including the Cmd forms, Tab maps to nothing, and no key maps to two actions. Run them and confirm they fail
- [x] 4.2 Implement the key table, and verify the tests from 4.1 pass
- [x] 4.3 Write interface snapshot tests for the `app-shell` scenarios "Empty session list at 60 by 20 with colour", "Empty session list at 120 by 30 with colour", "Cursor in the empty prompt box", "Empty session list with colour off" and "Window widened to 60 columns", and the `prompt-editor` scenario "Cursor in an empty prompt box". In the same step, rename or rewrite the existing tests for the removed scenario "Cursor on the empty message" and for the "q quit" key line, so that none still asserts the old screen. Run them and confirm they fail because there is no prompt box
- [x] 4.4 Add the prompt box to the empty session list as `design.md` lays it out, with the "▸ Prompt" label and the cursor at the insertion point, and verify the tests from 4.3 pass. Update the existing golden files for the empty session list with `make golden`, and say in the commit that the prompt box and the new key line appear
- [x] 4.5 Write model tests through `termtest.StartModel` for the `app-shell` scenarios "q types on the session list", "The first Ctrl-C asks before quitting", "Another key cancels quitting" and "Clearing then Ctrl-C does not quit". In the same step, change the existing quitting tests, in `internal/app` and at the terminal level, to press Ctrl-C twice for "Quitting restores the terminal", "Quitting from the too-small screen" and "Cursor does not blink", and add a terminal-level test for "q does nothing on the too-small screen". Run them and confirm they fail, the Ctrl-C ones because Ctrl-C doesn't quit yet
- [x] 4.6 Route keys and paste messages from the app to the editor through the key table, and add the pending-quit flag on the session list and the window-too-small message, as `design.md` describes under "Quitting asks first". Verify the tests from 4.5 pass, along with "Idle screen stays still" and "Startup draws the session list directly"
- [x] 4.7 Write model tests for every `prompt-editor` scenario about typing, the cursor, movement, selection, deletion, undo, wrapping, clearing, the key line, sending and the Cmd keys, each named after its scenario and reading the cursor's real position where the scenario names one. Run them and confirm each fails for the expected reason
- [x] 4.8 Add the status row and its messages, the key line that follows the prompt, and the send hook with the "no agent" status. Fix whatever else the tests from 4.7 show, and verify they all pass
- [x] 4.9 Write interface snapshot tests for "A long prompt scrolls at 60 by 20", "Scrolling follows the insertion point" and "Resize keeps the insertion point", with golden files at 60 by 20 for an empty box, a wrapped prompt, a selection and a scrolled prompt. In the same step, change the 4.7 tests that check a text row through the editor's rows, because the box showed only one row, so they read that row from the screen instead. Run them and confirm they fail
- [x] 4.10 Implement growing, scrolling and the counts on the label row, and verify the tests from 4.9 pass
- [ ] 4.11 Write tests for the `app-shell` window-too-small scenarios as modified, which now say "Ctrl-C quit", and for "The prompt is kept while the window is too small". Run them and confirm they fail
- [ ] 4.12 Update the window-too-small message, as `design.md` describes, and verify the tests from 4.11 pass

## 5. Terminal-level behaviour

- [ ] 5.1 Write terminal-level tests for the scenarios "Selection with colour on", "Selection with colour off" and "Long prompt with colour off", reading cells with `Screen.Cell` and checking the raw bytes for colour codes. Run them and confirm they fail
- [ ] 5.2 Draw the selection in inverse video and the accent colour, with only inverse video when `NO_COLOR` is set. Show the count of selected characters on the status row, as `design.md` describes. Verify the tests from 5.1 pass
- [ ] 5.3 Write terminal-level tests for "Copying a selection" and "Cutting a selection", decoding the OSC 52 payload from the raw bytes, and for "Pasting several lines", "Escape sequences in a paste are removed" and "Pasting over a selection", sending real bracketed paste bytes. Run them and confirm they fail
- [ ] 5.4 Wire copy and cut to `tea.SetClipboard`, and fix whatever the paste tests show. Verify the tests from 5.3 pass
- [ ] 5.5 Write a terminal-level test for "Keyboard enhancements requested and released". Bubble Tea already sends and releases the request, so this test records existing behaviour and passes at once. To show it can fail, first write a unit test for its byte-matching helper, fed bytes with no request in them, and confirm the helper reports the request missing
- [ ] 5.6 Write a terminal-level test for "Shift-Enter in the Kitty keyboard form". In the same task, write one that sends Delete's standard escape sequence, `ESC [ 3 ~`, and checks that the character after the insertion point is deleted, so the byte form of a key the logger only saw once is covered. Run it and confirm it fails, then fix whatever it shows and verify it passes
- [ ] 5.7 Write the key help in `docs/keys.md`: every binding, that Ctrl-J is the new line key that works everywhere, that Shift-Enter needs the Kitty keyboard protocol and Option-Enter needs Option sent as Alt, which keys need a terminal that passes Cmd through, that `cmux` keeps Cmd-Up, Cmd-Down and Cmd-A for itself and what to press instead, the select-all gap in Terminal.app, and that copy needs OSC 52 allowed. Link it from `README.md`. Verify every key it names is in the key table test from 4.1

## 6. Finish

- [ ] 6.1 Run `make check` and `openspec validate --strict --all`, and verify both pass with no uncovered scenarios, and that the scenario check lists only the skips this change expects
- [ ] 6.2 Update `STATE.md`: the editor is built, the next change is the service, and the third is the session screen. Verify it names both later changes

## Workflow follow-up

- Archive the change once every task is done and `make check` passes. As part of the archive commit, update the Purpose of `openspec/specs/app-shell/spec.md`, which still says the empty session list's last line is "q quit". It should say "Enter send · Ctrl-C quit", and mention the "▸ Prompt" box.
- Cami runs the milestone 1 "done when" test only after the third change, when a real session exists. For this change, she types and edits a prompt in `cmux` at her usual font size and says whether it feels like a normal text box.
