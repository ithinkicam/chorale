# Spec Delta

## Purpose

The prompt editor is the box where Cami writes what she asks an agent. It behaves like a normal text box, driven from the keyboard, so a prompt can be written, fixed and sent without fighting terminal line editing.

Every scenario that types into the prompt box starts from the empty session list, where the prompt box has focus. Key names follow macOS: Option is the Alt key, and Cmd is the Command key. "The prompt is" gives the prompt's whole text, with a new line written as `\n`.

## ADDED Requirements

### Requirement: Typing inserts at the cursor

Characters typed into the prompt box SHALL be inserted at the insertion point, and the insertion point SHALL move past them. Any printable character, including accented letters, wide characters and emoji, SHALL be accepted.

#### Scenario: Typing a word

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "hello" is typed
- **THEN** the prompt box shows "hello"
- **AND** the terminal cursor is visible on the column after the "o"

#### Scenario: Typing in the middle

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "helo" with the insertion point between "l" and "o"
- **WHEN** "l" is typed
- **THEN** the prompt is "hello"
- **AND** the insertion point is before the "o"

#### Scenario: Wide characters and emoji

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "合唱 🎵 é" is typed
- **THEN** the prompt is "合唱 🎵 é"
- **AND** the terminal cursor sits on the column after the "é", counting each of "合", "唱" and "🎵" as two columns wide

### Requirement: Cursor at the insertion point

While the prompt box has focus, the real terminal cursor SHALL be visible, steady, and on the insertion point, including after a soft wrap and while text is selected. With a selection, the insertion point is the end the cursor last moved to.

#### Scenario: Cursor in an empty prompt box

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions and has finished drawing
- **THEN** the terminal cursor is visible and steady
- **AND** it sits on the first column of the prompt box's first text row

#### Scenario: Cursor after a soft wrap

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** enough words are typed to wrap onto a second row, ending with "end"
- **THEN** the terminal cursor is on the prompt box's second text row
- **AND** it sits on the column after "end"

### Requirement: Moving by character and by word

Left and Right SHALL move the insertion point by one character, treating a character as a whole grapheme, so an emoji or an accented letter is never split. Option-Left and Option-Right SHALL move to the start of the previous word and the end of the next word. Ctrl-B and Ctrl-F SHALL work as Left and Right.

#### Scenario: Moving over an emoji

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "a🎵b" with the insertion point at the end
- **WHEN** Left is pressed twice
- **THEN** the insertion point is between "a" and "🎵"

#### Scenario: Moving by word

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests now" with the insertion point at the end
- **WHEN** Option-Left is pressed twice
- **THEN** the insertion point is before "tests"
- **AND** after Option-Right is pressed once, it is after "tests"

### Requirement: Moving by row and to the ends

Up and Down SHALL move between the rows as they appear on screen, across soft wraps, keeping the column the insertion point had when vertical movement began. Home and End, and Ctrl-A and Ctrl-E, SHALL move to the start and end of the row on screen. Option-Up and Option-Down SHALL move to the start and end of the whole prompt.

#### Scenario: Down keeps the column

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "abcdef\nab\nabcdef" with the insertion point after "abcde" on the first line
- **WHEN** Down is pressed twice
- **THEN** the insertion point is after "abcde" on the third line

#### Scenario: Up moves across a soft wrap

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and a prompt with no new lines that wraps onto two rows, with the insertion point at the end
- **WHEN** Up is pressed
- **THEN** the terminal cursor moves to the prompt box's first text row

#### Scenario: Start and end of the prompt

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "one\ntwo\nthree" with the insertion point after "tw"
- **WHEN** Option-Up is pressed
- **THEN** the insertion point is before "one"
- **AND** after Option-Down is pressed, it is after "three"

### Requirement: Selecting with the keyboard

Holding Shift with any movement key SHALL extend a selection from where the insertion point was. Ctrl-Shift-A SHALL select the whole prompt. Moving without Shift SHALL clear the selection. Typing, pasting or deleting while text is selected SHALL replace the selection.

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

#### Scenario: Selecting everything

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "one\ntwo"
- **WHEN** Ctrl-Shift-A is pressed
- **THEN** "one\ntwo" is selected

### Requirement: The selection is shown without relying on colour

Selected text SHALL be drawn in inverse video. With colour on, it SHALL also use the theme's accent colour. With colour off, the inverse video alone SHALL mark it.

#### Scenario: Selection with colour on

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" selected
- **WHEN** the screen is drawn
- **THEN** the cells showing "tests" are in inverse video and in the accent colour
- **AND** the cells showing "run the " are not in inverse video

#### Scenario: Selection with colour off

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with `NO_COLOR=1`, and the prompt is "run the tests" with "tests" selected
- **WHEN** the screen is drawn
- **THEN** the cells showing "tests" are in inverse video
- **AND** the bytes Chorale writes contain no colour codes

### Requirement: Deleting

Backspace and Delete SHALL remove one character before or after the insertion point, or the selection if there is one. Option-Backspace and Ctrl-W SHALL remove the word before it, and Option-Delete the word after it. Ctrl-U SHALL remove to the start of the row on screen, and Ctrl-K to its end.

#### Scenario: Backspace over an emoji

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "go 🎵" with the insertion point at the end
- **WHEN** Backspace is pressed once
- **THEN** the prompt is "go "

#### Scenario: Deleting a word

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with the insertion point at the end
- **WHEN** Option-Backspace is pressed
- **THEN** the prompt is "run the "

#### Scenario: Deleting to the end of the row

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "keep this\nand this" with the insertion point after "keep"
- **WHEN** Ctrl-K is pressed
- **THEN** the prompt is "keep\nand this"

### Requirement: Undo and redo

Ctrl-Z SHALL undo the last edit and Ctrl-Y or Ctrl-Shift-Z SHALL redo it, across the whole prompt, restoring the text, the insertion point and the selection. Typing is undone a word at a time: a run of typing ends at a space or new line after a word, at any movement, and at any other kind of edit. A new edit after an undo SHALL discard what could be redone.

#### Scenario: Undo typing a word at a time

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "run the tests" is typed and Ctrl-Z is pressed
- **THEN** the prompt is "run the "
- **AND** after Ctrl-Z is pressed twice more, the prompt is empty

#### Scenario: Redo after undo

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" removed by Option-Backspace
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

The prompt box SHALL wrap text to its real width at spaces, breaking a word only when it is wider than the box. It SHALL never split a character, and SHALL count wide characters as two columns. Wrapping SHALL NOT change the prompt's text. A resize SHALL rewrap it and keep the insertion point on the same character.

#### Scenario: Wrapping at a space at 60 columns

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** a sentence of 15 four-letter words, separated by single spaces, is typed
- **THEN** no row of the prompt box is wider than 60 columns
- **AND** no word is split across two rows
- **AND** the prompt has no new lines in it

#### Scenario: A long word is broken

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** a word of 100 letters with no spaces is typed
- **THEN** the word fills the prompt box's first row and continues on the next
- **AND** no row is wider than 60 columns

#### Scenario: Wide characters wrap whole

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** 40 copies of "合" are typed with no spaces
- **THEN** no row of the prompt box is wider than 60 columns
- **AND** every row shows only whole "合" characters

#### Scenario: Resize keeps the insertion point

- **GIVEN** a terminal 120 columns wide and 30 rows tall, in the default theme, with colour on, and a 90-character prompt with the insertion point after its 70th character
- **WHEN** the window is resized to 60 columns by 20 rows
- **THEN** the prompt is unchanged
- **AND** the terminal cursor sits just after the 70th character, on the prompt box's second text row

### Requirement: Pasting

Text pasted through the terminal's bracketed paste SHALL be inserted at the insertion point as one block, replacing any selection, and one undo SHALL remove it. New lines in it SHALL stay as new lines and SHALL NOT send the prompt. Carriage returns SHALL become new lines, tabs SHALL become four spaces, and other control characters, including escape sequences, SHALL be removed.

#### Scenario: Pasting several lines

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "first\r\nsecond\nthird" arrives as one bracketed paste
- **THEN** the prompt is "first\nsecond\nthird"
- **AND** the prompt has not been sent
- **AND** after Ctrl-Z, the prompt is empty

#### Scenario: Escape sequences in a paste are removed

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** "safe" followed by the bytes of an escape sequence that sets a colour, then "text", arrives as one bracketed paste
- **THEN** the prompt is "safetext"
- **AND** the bytes Chorale writes to show the prompt contain no part of the pasted escape sequence

### Requirement: Copying and cutting through OSC 52

With text selected, Ctrl-C SHALL copy it to the system clipboard by writing an OSC 52 sequence, the terminal escape code that sets the clipboard, and SHALL keep the selection. Ctrl-X SHALL copy the selection the same way and remove it from the prompt. With nothing selected, Ctrl-X SHALL do nothing.

#### Scenario: Copying a selection

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" selected
- **WHEN** Ctrl-C is pressed
- **THEN** the bytes Chorale writes include an OSC 52 sequence setting the clipboard to "tests"
- **AND** the prompt is still "run the tests" with "tests" selected
- **AND** Chorale is still running

#### Scenario: Cutting a selection

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "run the tests" with "tests" selected
- **WHEN** Ctrl-X is pressed
- **THEN** the bytes Chorale writes include an OSC 52 sequence setting the clipboard to "tests"
- **AND** the prompt is "run the "
- **AND** after Ctrl-Z, the prompt is "run the tests"

### Requirement: Clearing the prompt with Ctrl-C

With text in the prompt and nothing selected, Ctrl-C SHALL clear the prompt, as one edit that undo can reverse. With an empty prompt, Ctrl-C quits Chorale, as the app shell describes.

#### Scenario: Ctrl-C clears the prompt

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "half a thought" with nothing selected
- **WHEN** Ctrl-C is pressed
- **THEN** the prompt is empty
- **AND** Chorale is still running
- **AND** after Ctrl-Z, the prompt is "half a thought"

### Requirement: New lines and sending

Enter SHALL send the prompt. Shift-Enter, Option-Enter and Ctrl-J SHALL insert a new line instead. Enter on a prompt that is empty or only spaces and new lines SHALL do nothing. While no agent is connected, sending SHALL keep the prompt as it is and show "No agent is connected yet. Your prompt is kept." on the screen.

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

The prompt box SHALL grow with its text, up to 8 text rows in a window 20 rows tall and up to 12 in a taller one, and then scroll so the insertion point is always in view. When rows are out of sight, the box SHALL say in words how many, such as "3 lines above" and "2 lines below".

#### Scenario: A long prompt scrolls at 60 by 20

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** 20 short lines, "line 1" to "line 20", are typed with Shift-Enter between them
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

### Requirement: Cmd keys where the terminal passes them

Where the terminal passes the Cmd key through to Chorale, Cmd-Z SHALL undo, Cmd-Shift-Z SHALL redo, Cmd-A SHALL select the whole prompt, Cmd-C SHALL copy, Cmd-X SHALL cut, Cmd-Left and Cmd-Right SHALL move to the start and end of the row, and Cmd-Up and Cmd-Down to the start and end of the prompt.

#### Scenario: Cmd-Z undoes

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, that passes Cmd keys through, and "run the tests" was typed
- **WHEN** Cmd-Z is pressed
- **THEN** the prompt is "run the "

#### Scenario: Cmd-A selects everything

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, that passes Cmd keys through, and the prompt is "one\ntwo"
- **WHEN** Cmd-A is pressed
- **THEN** "one\ntwo" is selected

### Requirement: Telling Shift-Enter apart from Enter

Chorale SHALL ask the terminal to report keys unambiguously, through the Kitty keyboard protocol, so Shift-Enter and Ctrl-Shift keys can be told apart where the terminal supports it. It SHALL turn this off again when it exits.

#### Scenario: Keyboard enhancements requested and released

- **GIVEN** `chorale` running in a pseudo-terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** it starts, and then Ctrl-C is pressed on the empty prompt
- **THEN** the bytes it wrote on starting include a request to turn on Kitty keyboard key disambiguation
- **AND** once it has exited, keyboard enhancements are off
