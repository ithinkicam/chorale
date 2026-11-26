# Spec Delta

## Purpose

The prompt editor is the box where Cami writes what she asks an agent. It behaves like a normal text box, driven from the keyboard, so a prompt can be written, fixed and sent without fighting terminal line editing.

Every scenario that types into the prompt box starts from the empty session list, where the prompt box has focus. Key names follow macOS: Option is the Alt key, and Cmd is the Command key. "The prompt is" gives the prompt's whole text, with a new line written as `\n`.

The prompt box has a label row, then its text rows, then a status row. The label row reads "▸ Prompt", followed by how many rows are out of sight when there are any. The label calls them lines. The text rows are indented by two columns, and the last column of the window is kept free, so at 60 columns each text row holds 57 columns of text. Screen columns are counted from 1 at the left edge, so the first text column is the third column of the screen.

## ADDED Requirements

### Requirement: Typing inserts at the cursor

Characters typed into the prompt box SHALL be inserted at the insertion point, and the insertion point SHALL move past them. Any printable character, including accented letters, wide characters and emoji, SHALL be accepted. A combining accent typed on its own SHALL join the character before it. The Tab key SHALL do nothing, and SHALL NOT insert anything.

#### Scenario: Typing a word

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "hello" is typed
- **THEN** the prompt box shows "hello"
- **AND** the terminal cursor is visible on the eighth column of the screen, the one after the "o"

#### Scenario: Typing in the middle

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "helo" with the insertion point between "l" and "o"
- **WHEN** "l" is typed
- **THEN** the prompt is "hello"
- **AND** the insertion point is before the "o"

#### Scenario: Wide characters and emoji

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "合唱 🎵 é" is typed
- **THEN** the prompt is "合唱 🎵 é"
- **AND** the terminal cursor sits on the twelfth column of the screen, counting each of "合", "唱" and "🎵" as two columns wide

#### Scenario: An accent typed on its own joins its letter

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "e" with the insertion point at the end
- **WHEN** the combining acute accent, U+0301, is typed on its own
- **THEN** the prompt is "e" followed by U+0301, shown as "é"
- **AND** after Left is pressed once, the insertion point is at the start of the prompt

#### Scenario: Tab does nothing

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run" with the insertion point at the end
- **WHEN** Tab is pressed
- **THEN** the prompt is "run"
- **AND** the insertion point is at the end

### Requirement: Cursor at the insertion point

While the prompt box has focus, the real terminal cursor SHALL be visible, steady, and on the insertion point, including while text is selected. With a selection, the insertion point is the end the cursor last moved to.

#### Scenario: Cursor in an empty prompt box

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions and has finished drawing
- **THEN** the terminal cursor is visible and steady
- **AND** it sits on the third column of the screen, on the prompt box's first text row

### Requirement: The cursor at a soft wrap

At a soft wrap, the end of one row and the start of the next are the same place in the text. The cursor SHALL be drawn at the start of the next row, except right after End, Ctrl-E or Cmd-Right, when it SHALL be drawn after the last character of the row it was on.

#### Scenario: Cursor after a soft wrap

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** twelve copies of "word", separated by single spaces, followed by " end", are typed
- **THEN** the prompt box's first text row shows eleven copies of "word"
- **AND** its second text row shows "word end"
- **AND** the terminal cursor is on the second text row, on the eleventh column of the screen

#### Scenario: Cursor at the end of a full row

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is 57 copies of "a" followed by "bbb", with the insertion point at the start
- **WHEN** End is pressed
- **THEN** the terminal cursor is on the prompt box's first text row, on the sixtieth column of the screen
- **AND** after Right is pressed, the terminal cursor is on the second text row, on the fourth column of the screen

### Requirement: Moving by character

Left and Right, and Ctrl-B and Ctrl-F, SHALL move the insertion point by one character, treating a character as a whole grapheme, so an emoji or an accented letter is never split. Moving across a soft wrap SHALL take one press, like moving past any other character.

#### Scenario: Moving over an emoji

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "a🎵b" with the insertion point at the end
- **WHEN** Left is pressed twice
- **THEN** the insertion point is between "a" and "🎵"

#### Scenario: Moving left across a soft wrap

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is twelve copies of "word", separated by single spaces, followed by " end", with the insertion point before the twelfth "word"
- **WHEN** Left is pressed
- **THEN** the insertion point is between the eleventh "word" and the space after it
- **AND** the terminal cursor is on the prompt box's first text row, on the fifty-seventh column of the screen

### Requirement: Moving by word

Option-Left and Option-Right, and Alt-B and Alt-F, SHALL move to the start of the previous word and the end of the next word, skipping any spaces and punctuation in between. A word is a run of letters, digits and underscores, and an emoji counts as a word of its own. Other symbols, such as `+`, `=` and `|`, are skipped like punctuation.

#### Scenario: Moving by word

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests now" with the insertion point at the end
- **WHEN** Option-Left is pressed twice
- **THEN** the insertion point is before "tests"
- **AND** after Option-Right is pressed once, it is after "tests"

### Requirement: Moving by row and to the ends

Up and Down SHALL move between the rows as they appear on screen, across soft wraps, keeping the screen column the insertion point had when vertical movement began. Home and End, and Ctrl-A and Ctrl-E, SHALL move to the start and end of the row on screen. Option-Up and Option-Down SHALL move to the start and end of the whole prompt. Up on the first row SHALL move to the start of the prompt, and Down on the last row to its end.

#### Scenario: Down keeps the column

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "abcdef\nab\nabcdef" with the insertion point after "abcde" on the first line
- **WHEN** Down is pressed twice
- **THEN** the insertion point is after "abcde" on the third line

#### Scenario: Up moves across a soft wrap

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is twelve copies of "word", separated by single spaces, followed by " end", with the insertion point at the end
- **WHEN** Up is pressed
- **THEN** the terminal cursor is on the prompt box's first text row, on the eleventh column of the screen
- **AND** the insertion point is between "wor" and "d" in the second "word"

#### Scenario: Up on the first row

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "one\ntwo" with the insertion point after "on"
- **WHEN** Up is pressed
- **THEN** the insertion point is before "one"
- **AND** after Down is pressed twice, it is after "two"

#### Scenario: Start and end of the prompt

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "one\ntwo\nthree" with the insertion point after "tw"
- **WHEN** Option-Up is pressed
- **THEN** the insertion point is before "one"
- **AND** after Option-Down is pressed, it is after "three"

### Requirement: Selecting with the keyboard

Holding Shift with any movement key SHALL extend a selection from where the insertion point was. Ctrl-Shift-A SHALL select the whole prompt. Moving without Shift SHALL clear the selection, and Left or Right without Shift SHALL move to the selection's start or end. Typing, pasting or deleting while text is selected SHALL replace the selection.

#### Scenario: Selecting a word

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with the insertion point at the end
- **WHEN** Shift-Option-Left is pressed
- **THEN** "tests" is selected
- **AND** the terminal cursor sits on the "t" at the start of "tests"

#### Scenario: Typing replaces the selection

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" selected
- **WHEN** "build" is typed
- **THEN** the prompt is "run the build"
- **AND** nothing is selected

#### Scenario: Right moves to the end of a selection

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests now" with "tests" selected and the insertion point before it
- **WHEN** Right is pressed
- **THEN** nothing is selected
- **AND** the insertion point is after "tests"

#### Scenario: Selecting everything

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "one\ntwo"
- **WHEN** Ctrl-Shift-A is pressed
- **THEN** "one\ntwo" is selected

### Requirement: The selection is shown without relying on colour

Selected text SHALL be drawn in inverse video. With colour on, it SHALL also use the theme's accent colour. With colour off, the inverse video alone SHALL mark it. While text is selected, the status row SHALL say in words how many characters are selected, counting a new line as one character, such as "5 characters selected" or "1 character selected".

#### Scenario: Selection with colour on

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" selected
- **WHEN** the screen is drawn
- **THEN** the cells showing "tests" are in inverse video and in the accent colour
- **AND** the cells showing "run the " are not in inverse video
- **AND** the screen shows "5 characters selected"

#### Scenario: Selection with colour off

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with `NO_COLOR=1`, and the prompt is "run the tests" with "tests" selected
- **WHEN** the screen is drawn
- **THEN** the cells showing "tests" are in inverse video
- **AND** the screen shows "5 characters selected"
- **AND** the bytes Chorale writes contain no colour codes

### Requirement: Deleting

Backspace and Delete SHALL remove one character before or after the insertion point, or the selection if there is one. Option-Backspace and Ctrl-W SHALL remove the word before it, and Option-Delete and Alt-D the word after it, using the same words as Option-Left and Option-Right. Ctrl-U SHALL remove to the start of the row on screen, and Ctrl-K to its end. Deleting across a soft wrap SHALL work like deleting any other character.

#### Scenario: Backspace over an emoji

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "go 🎵" with the insertion point at the end
- **WHEN** Backspace is pressed once
- **THEN** the prompt is "go "

#### Scenario: Deleting a word

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with the insertion point at the end
- **WHEN** Option-Backspace is pressed
- **THEN** the prompt is "run the "

#### Scenario: Deleting a word after an emoji

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "go 🎵" with the insertion point at the end
- **WHEN** Option-Backspace is pressed
- **THEN** the prompt is "go "

#### Scenario: Deleting to the end of the row

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "keep this\nand this" with the insertion point after "keep"
- **WHEN** Ctrl-K is pressed
- **THEN** the prompt is "keep\nand this"

#### Scenario: Backspace at the start of a wrapped row

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is twelve copies of "word", separated by single spaces, followed by " end", with the insertion point before the twelfth "word"
- **WHEN** Backspace is pressed
- **THEN** the prompt is ten copies of "word", separated by single spaces, followed by " wordword end"

### Requirement: Undo and redo

Ctrl-Z SHALL undo the last edit and Ctrl-Y or Ctrl-Shift-Z SHALL redo it, across the whole prompt, restoring the text, the insertion point and the selection. Typing is undone a word at a time: a run of typing ends at a space or new line after a word, at any movement, and at any other kind of edit. A run of Backspace presses, or of Delete presses, SHALL be undone in one step. A new edit after an undo SHALL discard what could be redone.

#### Scenario: Undo typing a word at a time

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "run the tests" is typed and Ctrl-Z is pressed
- **THEN** the prompt is "run the "
- **AND** after Ctrl-Z is pressed twice more, the prompt is empty

#### Scenario: Undoing a run of Backspace

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with the insertion point at the end
- **WHEN** Backspace is pressed five times and then Ctrl-Z once
- **THEN** the prompt is "run the tests"

#### Scenario: Redo after undo

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and "run the tests" was typed, then Option-Backspace was pressed, leaving "run the "
- **WHEN** Ctrl-Z is pressed and then Ctrl-Y
- **THEN** the prompt is "run the "
- **AND** the insertion point is at the end

#### Scenario: Undo restores the selection

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" selected
- **WHEN** "build" is typed and Ctrl-Z is pressed
- **THEN** the prompt is "run the tests"
- **AND** "tests" is selected again

#### Scenario: A new edit discards redo

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and "run the tests" was typed and then undone once
- **WHEN** "x" is typed and Ctrl-Y is pressed
- **THEN** the prompt is "run the x"

### Requirement: Soft wrapping

The prompt box SHALL wrap text to its real width at spaces, breaking a word only when it is wider than the box. A space at a wrap stays at the end of the row it ends when it fits there, and otherwise starts the next row. It SHALL never split a character, and SHALL count wide characters as two columns, measuring every character the same way the screen draws it. Wrapping SHALL NOT change the prompt's text. A resize SHALL rewrap it and keep the insertion point on the same character.

#### Scenario: Wrapping at a space at 60 columns

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** fifteen copies of "word", separated by single spaces, are typed
- **THEN** no row of the prompt box is wider than 60 columns
- **AND** no "word" is split across two rows
- **AND** the prompt has no new lines in it

#### Scenario: A long word is broken

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** 100 copies of "x", with no spaces, are typed
- **THEN** the prompt box's first text row shows 57 copies of "x"
- **AND** its second text row shows 43 copies of "x"

#### Scenario: Wide characters wrap whole

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** 40 copies of "合" are typed with no spaces
- **THEN** the prompt box's first text row shows 28 copies of "合"
- **AND** its second text row shows 12 copies of "合"
- **AND** no row of the prompt box is wider than 60 columns

#### Scenario: Resize keeps the insertion point

- **GIVEN** a terminal 120 columns wide and 30 rows tall, in the default theme, with colour on, and the prompt is eighteen copies of "word", separated by single spaces, with the insertion point after the fourteenth "word"
- **WHEN** the window is resized to 60 columns by 20 rows
- **THEN** the prompt is unchanged
- **AND** the terminal cursor is on the prompt box's second text row, on the seventeenth column of the screen

### Requirement: Pasting

Text pasted through the terminal's bracketed paste SHALL be inserted at the insertion point as one block, replacing any selection, and one undo SHALL remove it. New lines in it SHALL stay as new lines and SHALL NOT send the prompt. Carriage returns SHALL become new lines, tabs SHALL become four spaces, and every escape sequence SHALL be removed whole, along with any other control character.

#### Scenario: Pasting several lines

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "first\r\nsecond\nthird" arrives as one bracketed paste
- **THEN** the prompt is "first\nsecond\nthird"
- **AND** the prompt has not been sent
- **AND** after Ctrl-Z, the prompt is empty

#### Scenario: Escape sequences in a paste are removed

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "safe", then the escape sequence that sets the foreground red, `ESC [ 31 m`, then "text", arrives as one bracketed paste
- **THEN** the prompt is "safetext"
- **AND** the bytes Chorale writes to show the prompt contain no part of the pasted escape sequence

#### Scenario: Pasting over a selection

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" selected
- **WHEN** "build" arrives as one bracketed paste
- **THEN** the prompt is "run the build"
- **AND** after Ctrl-Z, the prompt is "run the tests" with "tests" selected

### Requirement: Copying and cutting through OSC 52

With text selected, Ctrl-C SHALL copy it to the system clipboard by writing an OSC 52 sequence, the terminal escape code that sets the clipboard, SHALL keep the selection, and SHALL show "Copied." in the status row. Ctrl-X SHALL copy the selection the same way, remove it from the prompt, and show "Cut." in the status row. With nothing selected, Ctrl-X SHALL do nothing. A status message stays until the next key or paste.

#### Scenario: Copying a selection

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" selected
- **WHEN** Ctrl-C is pressed
- **THEN** the bytes Chorale writes include an OSC 52 sequence setting the clipboard to "tests"
- **AND** the prompt is still "run the tests" with "tests" selected
- **AND** the screen shows "Copied."
- **AND** Chorale is still running

#### Scenario: Cutting a selection

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" selected
- **WHEN** Ctrl-X is pressed
- **THEN** the bytes Chorale writes include an OSC 52 sequence setting the clipboard to "tests"
- **AND** the prompt is "run the "
- **AND** the screen shows "Cut."
- **AND** after Ctrl-Z, the prompt is "run the tests"

### Requirement: Clearing the prompt with Ctrl-C

With anything in the prompt, even only spaces or new lines, and nothing selected, Ctrl-C SHALL clear the prompt, as one edit that undo can reverse, and SHALL show "Prompt cleared. Ctrl-Z brings it back." in the status row. With an empty prompt, Ctrl-C starts to quit Chorale, as the app shell describes.

#### Scenario: Ctrl-C clears the prompt

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "half a thought" with nothing selected
- **WHEN** Ctrl-C is pressed
- **THEN** the prompt is empty
- **AND** the screen shows "Prompt cleared. Ctrl-Z brings it back."
- **AND** Chorale is still running
- **AND** after Ctrl-Z, the prompt is "half a thought"

#### Scenario: Ctrl-C clears a prompt of only spaces

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "  \n "
- **WHEN** Ctrl-C is pressed
- **THEN** the prompt is empty
- **AND** the screen shows "Prompt cleared. Ctrl-Z brings it back."

### Requirement: The key line follows the prompt

The last line of the screen SHALL name what Ctrl-C does at that moment: "Enter send · Ctrl-C copy" while text is selected, "Enter send · Ctrl-C clear" while the prompt has anything in it and nothing is selected, and "Enter send · Ctrl-C quit" while the prompt is empty.

#### Scenario: Key line with text

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "run" is typed
- **THEN** the last line shows "Enter send · Ctrl-C clear"

#### Scenario: Key line with a selection

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with the insertion point at the end
- **WHEN** Shift-Option-Left is pressed
- **THEN** the last line shows "Enter send · Ctrl-C copy"

### Requirement: New lines and sending

Enter SHALL send the prompt. Shift-Enter, Option-Enter and Ctrl-J SHALL insert a new line instead. Enter on a prompt that is empty or only spaces and new lines SHALL do nothing. While no agent is connected, sending SHALL keep the prompt as it is and show "No agent is connected yet. Your prompt is kept." in the status row.

#### Scenario: New line keys

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "a", Shift-Enter, "b", Option-Enter, "c", Ctrl-J and "d" are typed in that order
- **THEN** the prompt is "a\nb\nc\nd"

#### Scenario: Sending with no agent connected

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests"
- **WHEN** Enter is pressed
- **THEN** the screen shows "No agent is connected yet. Your prompt is kept."
- **AND** the prompt is still "run the tests"

#### Scenario: Enter on an empty prompt

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "  \n "
- **WHEN** Enter is pressed
- **THEN** the screen does not show "No agent is connected yet."
- **AND** the prompt is still "  \n "

### Requirement: Growing and scrolling

The prompt box SHALL show one text row when empty and grow with its text, up to the window's height minus 12 rows, and never more than 12 text rows. So it reaches 8 text rows in a window 20 rows tall and 12 in a window 24 rows tall or taller. Past that, it SHALL scroll by the least amount that keeps the insertion point in view. When rows are out of sight, the label row SHALL say in words how many, such as "▸ Prompt · 3 lines above · 2 lines below".

#### Scenario: A long prompt scrolls at 60 by 20

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** 20 short lines, "line 1" to "line 20", are typed with Ctrl-J between them
- **THEN** the prompt box shows 8 text rows, the last being "line 20"
- **AND** the screen shows "12 lines above"
- **AND** the terminal cursor is after "line 20"
- **AND** no line is wider than 60 columns

#### Scenario: Scrolling follows the insertion point

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the 20-line prompt "line 1" to "line 20" with the insertion point at the end
- **WHEN** Option-Up is pressed
- **THEN** the prompt box's first text row is "line 1"
- **AND** the screen shows "12 lines below"
- **AND** does not show "lines above"

#### Scenario: Long prompt with colour off

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with `NO_COLOR=1`, and the 20-line prompt "line 1" to "line 20"
- **WHEN** the screen is drawn
- **THEN** the screen shows "12 lines above"
- **AND** the bytes Chorale writes contain no colour codes

### Requirement: Counting rows out of sight

The counts on the label row SHALL be of rows on screen, not of lines ended by a line break, so a line that wraps onto three rows counts as three lines.

#### Scenario: A wrapped line out of sight counts by its rows

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** a line of 150 letter "x" characters, then the lines "line 2" to "line 9", each after a line break, are pasted
- **THEN** the prompt box shows 8 text rows, the first being "line 2" and the last being "line 9"
- **AND** the screen shows "3 lines above"
- **AND** does not show "1 line above"

### Requirement: Cmd keys where the terminal passes them

Where the terminal passes the Cmd key through to Chorale, Cmd-Z SHALL undo, Cmd-Shift-Z SHALL redo, Cmd-A SHALL select the whole prompt, Cmd-C SHALL copy, Cmd-X SHALL cut, Cmd-Left and Cmd-Right SHALL move to the start and end of the row, and Cmd-Up and Cmd-Down to the start and end of the prompt. Shift-Cmd-Left and Shift-Cmd-Right SHALL select to the start and end of the row.

#### Scenario: Cmd-Z undoes

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, that passes Cmd keys through, and "run the tests" was typed
- **WHEN** Cmd-Z is pressed
- **THEN** the prompt is "run the "

#### Scenario: Cmd-A selects everything

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, that passes Cmd keys through, and the prompt is "one\ntwo"
- **WHEN** Cmd-A is pressed
- **THEN** "one\ntwo" is selected

#### Scenario: Shift-Cmd-Left selects to the start of the row

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, that passes Cmd keys through, and the prompt is "run the tests" with the insertion point at the end
- **WHEN** Shift-Cmd-Left is pressed
- **THEN** "run the tests" is selected
- **AND** the terminal cursor sits on screen column 3 of the prompt's first row

### Requirement: Telling Shift-Enter apart from Enter

Chorale SHALL ask the terminal to report keys unambiguously, through the Kitty keyboard protocol, so Shift-Enter and Ctrl-Shift keys can be told apart where the terminal supports it. It SHALL turn this off again when it exits. Where the terminal reports a Shift-Enter in the Kitty keyboard protocol's form, it SHALL insert a new line.

#### Scenario: Keyboard enhancements requested and released

- **GIVEN** `chorale` running in a pseudo-terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** it starts, and then Ctrl-C is pressed twice on the empty prompt
- **THEN** the bytes it wrote on starting include a request to turn on Kitty keyboard key disambiguation
- **AND** once it has exited, keyboard enhancements are off

#### Scenario: Shift-Enter in the Kitty keyboard form

- **GIVEN** `chorale` running in a pseudo-terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and "a" typed into the prompt
- **WHEN** the terminal sends Shift-Enter in the Kitty keyboard protocol's form, then "b"
- **THEN** the prompt is "a\nb"
- **AND** the screen does not show "No agent is connected yet."
