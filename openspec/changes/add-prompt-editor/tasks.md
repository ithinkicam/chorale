# Tasks

Every behaviour task below comes in a pair: first the failing test, then the code. Run each test after writing it and confirm it fails for the expected reason before writing the code. Every test asserts its scenario's THEN clauses directly, and is named after its scenario. A golden file is a record of the screen, never the thing that makes a test pass. Run `make test` and `make lint` before ticking a task.

## 1. Groundwork

- [ ] 1.1 On Cami's Mac, in `cmux`, run a throwaway key logger built on Bubble Tea v2 with key disambiguation on, kept in the scratchpad and never committed. Record in `design.md`, under "Keys", which of the bindings listed there reach Chorale and as what. Change the binding table in `design.md` first if a key never arrives. Verify by reading the recorded list back to Cami, because she runs the logger
- [ ] 1.2 Write a unit test for `scripts/check-scenarios` with a main spec and an open change that removes one requirement, renames another and modifies a third. Assert that the main spec's scenarios under all three are not required, and that the change's own scenarios are. Run it and confirm it fails because the check requires them
- [ ] 1.3 Teach the check to skip main-spec requirements that an open change removes, renames or modifies, as `design.md` describes. Verify the test from 1.2 and the existing check tests pass
- [ ] 1.4 Write a test for a new `termtest.Screen.Cell` method, feeding known escape sequences for inverse video and a palette foreground colour, and asserting the text, the inverse flag and the colour of chosen cells. Run it and confirm it fails because the method doesn't exist
- [ ] 1.5 Implement `Screen.Cell` over `x/vt`, and verify the test from 1.4 passes

## 2. Editor text and movement

- [ ] 2.1 Write unit tests in `internal/editor` for the scenarios "Typing a word", "Typing in the middle" and "Wide characters and emoji", asserting the text and the insertion point. Run them and confirm they fail because the package doesn't exist
- [ ] 2.2 Implement the editor's grapheme text and insertion, and verify the tests from 2.1 pass
- [ ] 2.3 Write unit tests for the layout: the scenarios "Wrapping at a space at 60 columns", "A long word is broken" and "Wide characters wrap whole", at a text width of 58, asserting each row's width and content. Run them and confirm they fail
- [ ] 2.4 Implement `Layout` as `design.md` describes, and verify the tests from 2.3 pass
- [ ] 2.5 Write unit tests for the scenarios "Moving over an emoji", "Moving by word", "Down keeps the column", "Up moves across a soft wrap" and "Start and end of the prompt". Run them and confirm they fail
- [ ] 2.6 Implement character, word, row and whole-prompt movement, with the remembered column for Up and Down, and verify the tests from 2.5 pass

## 3. Selection, deletion, history and paste

- [ ] 3.1 Write unit tests for the scenarios "Selecting a word", "Typing replaces the selection" and "Selecting everything". Run them and confirm they fail
- [ ] 3.2 Implement the selection anchor and Shift movement, and verify the tests from 3.1 pass
- [ ] 3.3 Write unit tests for the scenarios "Backspace over an emoji", "Deleting a word" and "Deleting to the end of the row". Run them and confirm they fail
- [ ] 3.4 Implement the deletions, and verify the tests from 3.3 pass
- [ ] 3.5 Write unit tests for the scenarios "Undo typing a word at a time", "Redo after undo", "Undo restores the selection" and "A new edit discards redo". Run them and confirm they fail
- [ ] 3.6 Implement the history with word-sized typing groups, and verify the tests from 3.5 pass
- [ ] 3.7 Write unit tests for the paste cleaning, covering the scenarios "Pasting several lines" and "Escape sequences in a paste are removed" at the editor level. Run them and confirm they fail
- [ ] 3.8 Implement paste cleaning and insertion as one undo group, and verify the tests from 3.7 pass
- [ ] 3.9 Write a unit test that inserting a 10,000-line paste, laying it out at width 58 and moving to its end each take under 200 milliseconds, and that the history stays under 16 MB after 300 such edits. Run it and confirm it fails or passes. If it passes on the first run, record that in the commit. Then fix whatever it shows

## 4. Key table and the prompt box on screen

- [ ] 4.1 Write unit tests for the key table: every binding in `design.md` maps to its action, including the Cmd forms, and no key maps to two actions. Run them and confirm they fail
- [ ] 4.2 Implement the key table, and verify the tests from 4.1 pass
- [ ] 4.3 Write interface snapshot tests for the `app-shell` scenarios "Empty session list at 60 by 20 with colour", "Empty session list at 120 by 30 with colour" and "Cursor in the empty prompt box", and the `prompt-editor` scenario "Cursor in an empty prompt box". Run them and confirm they fail because there is no prompt box
- [ ] 4.4 Add the prompt box to the empty session list as `design.md` lays it out, with the cursor at the insertion point, and verify the tests from 4.3 pass. Update the existing golden files for the empty session list with `make golden`, and say in the commit that the prompt box and the new key line appear
- [ ] 4.5 Write interaction tests, sending key messages to the app, for the scenarios "Cursor after a soft wrap", "Resize keeps the insertion point", "Ctrl-C clears the prompt", "New line keys", "Sending with no agent connected", "Enter on an empty prompt", "Cmd-Z undoes", "Cmd-A selects everything" and "q types on the session list", and one interaction test per movement, selection, deletion and undo scenario from groups 2 and 3, named after the scenario with the suffix "ThroughKeys". Run them and confirm they fail
- [ ] 4.6 Route keys and paste messages from the app to the editor through the key table, add the send hook with the "no agent" status, and move quitting to Ctrl-C on an empty prompt. Verify the tests from 4.5 pass
- [ ] 4.7 Write interface snapshot tests for the scenarios "A long prompt scrolls at 60 by 20" and "Scrolling follows the insertion point", with golden files at 60 by 20 for an empty box, a wrapped prompt, a selection and a scrolled prompt. Run them and confirm they fail
- [ ] 4.8 Implement growing, scrolling and the "lines above" and "lines below" words, and verify the tests from 4.7 pass

## 5. Terminal-level behaviour

- [ ] 5.1 Write terminal-level tests for the scenarios "Selection with colour on", "Selection with colour off" and "Long prompt with colour off", reading cells with `Screen.Cell` and checking the raw bytes for colour codes. Run them and confirm they fail
- [ ] 5.2 Draw the selection in inverse video and the accent colour, with only inverse video when `NO_COLOR` is set, and verify the tests from 5.1 pass
- [ ] 5.3 Write terminal-level tests for the scenarios "Copying a selection" and "Cutting a selection", decoding the OSC 52 payload from the raw bytes, and a terminal-level version of "Pasting several lines" that sends real bracketed paste bytes. Run them and confirm they fail
- [ ] 5.4 Wire copy and cut to `tea.SetClipboard`, and verify the tests from 5.3 pass
- [ ] 5.5 Write a terminal-level test for the scenario "Keyboard enhancements requested and released", and one that a Shift-Enter sent in the Kitty keyboard protocol's form inserts a new line. Run them and confirm they fail
- [ ] 5.6 Set the keyboard enhancement flags in the view, and verify the tests from 5.5 pass
- [ ] 5.7 Update the existing terminal-level tests for "Quitting restores the terminal", "Quitting from the too-small screen", "Cursor does not blink", "Idle screen stays still" and "Startup draws the session list directly" to use Ctrl-C and the prompt box. Run them against the code before 4.6 and confirm the Ctrl-C ones fail, then verify they all pass now
- [ ] 5.8 Write the key help in `docs/keys.md`: every binding, which ones need the Kitty keyboard protocol or a terminal that passes Cmd through, the select-all gap in Terminal.app, and that copy needs OSC 52 allowed. Link it from `README.md`. Verify every key it names is in the key table test from 4.1

## 6. Finish

- [ ] 6.1 Run `make check` and `openspec validate --strict --all`, and verify both pass with no uncovered scenarios
- [ ] 6.2 Update `STATE.md`: the editor is built, the next change is the service, and the third is the session screen. Verify it names both later changes

## Workflow follow-up

- Archive the change once every task is done and `make check` passes. As part of the archive commit, update the Purpose of `openspec/specs/app-shell/spec.md`, which still says the empty session list's last line is "q quit", to say "Enter send · Ctrl-C quit" and to mention the prompt box.
- Cami runs the milestone 1 "done when" test only after the third change, when a real session exists. For this change, she types and edits a prompt in `cmux` at her usual font size and says whether it feels like a normal text box.
