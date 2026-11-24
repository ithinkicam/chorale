# Spec Delta

## RENAMED Requirements

- FROM: `### Requirement: Quit with q`
- TO: `### Requirement: Quitting`

## ADDED Requirements

### Requirement: Empty session list with the prompt box

When Chorale starts and there are no sessions, it SHALL show a full-screen view with the Chorale name on the first line, a message that there are no sessions yet, an empty prompt box with focus, and a last line listing the keys that work. The prompt box's label SHALL be "▸ Prompt", where the triangle marks it as focused. The real terminal cursor SHALL be visible in the prompt box, so macOS Zoom follows it there.

#### Scenario: Empty session list at 60 by 20 with colour

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** the first line starts with "♪ chorale"
- **AND** the screen shows "No sessions yet."
- **AND** the screen shows the prompt box's label, "▸ Prompt", below that message
- **AND** the last line shows "Enter send · Ctrl-C quit"
- **AND** no line is wider than 60 columns

#### Scenario: Empty session list at 120 by 30 with colour

- **GIVEN** a terminal 120 columns wide and 30 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** the first line starts with "♪ chorale"
- **AND** the screen shows "No sessions yet."
- **AND** the screen shows the prompt box's label, "▸ Prompt"
- **AND** "Enter send · Ctrl-C quit" is on the last row
- **AND** no line is wider than 120 columns

#### Scenario: Cursor in the empty prompt box

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions and has finished drawing
- **THEN** the terminal cursor is visible
- **AND** it sits on the row below the "▸ Prompt" label, on the third column of the screen

### Requirement: Ctrl-C asks before quitting

Ctrl-C with an empty prompt and nothing selected, on the session list, SHALL NOT quit at once. It SHALL show "Press Ctrl-C again to quit." A second Ctrl-C straight after SHALL quit. Any other key, paste or resize in between SHALL cancel it, and the message SHALL go. On the window-too-small message, Ctrl-C SHALL always work this way, whatever the prompt holds, with the message in place of "Ctrl-C quit".

#### Scenario: The first Ctrl-C asks before quitting

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list with an empty prompt
- **WHEN** Ctrl-C is pressed once
- **THEN** Chorale is still running
- **AND** the screen shows "Press Ctrl-C again to quit."

#### Scenario: Another key cancels quitting

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list with an empty prompt
- **WHEN** Ctrl-C, Left and Ctrl-C are pressed in that order
- **THEN** Chorale is still running
- **AND** the screen shows "Press Ctrl-C again to quit."

#### Scenario: Clearing then Ctrl-C does not quit

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "half a thought"
- **WHEN** Ctrl-C is pressed twice
- **THEN** Chorale is still running
- **AND** the prompt is empty
- **AND** the screen shows "Press Ctrl-C again to quit."
- **AND** after Ctrl-Z, the prompt is "half a thought"

### Requirement: The prompt is kept while the window is too small

While the window is too small and the prompt has anything in it, the window-too-small message SHALL say "Your prompt is kept." after the size it needs, and before "Ctrl-C quit" when there is room for both. The prompt SHALL be unchanged when the normal screen comes back.

#### Scenario: The prompt is kept while the window is too small

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, and the prompt is "half a thought"
- **WHEN** the window is resized to 40 columns by 20 rows, then back to 60 columns by 20 rows
- **THEN** while it is 40 columns wide, the screen shows "Your prompt is kept."
- **AND** afterwards the prompt is "half a thought"

## MODIFIED Requirements

### Requirement: Readable without colour

Chorale SHALL respect the `NO_COLOR` environment variable by sending no colour codes, and every piece of information on screen SHALL still be shown in words or symbols. With colour allowed, the default theme SHALL use only the terminal's own palette.

#### Scenario: Empty session list with colour off

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with `NO_COLOR=1`
- **WHEN** Chorale starts with no sessions
- **THEN** the bytes Chorale writes to the terminal contain no colour codes
- **AND** the screen shows "♪ chorale", "No sessions yet.", "▸ Prompt" and "Enter send · Ctrl-C quit"

#### Scenario: Colour on uses the terminal palette

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** the bytes Chorale writes contain at least one colour code
- **AND** every colour code uses one of the terminal's 16 palette colours, with no 256-colour or 24-bit forms

### Requirement: No motion

Chorale SHALL NOT write anything to the terminal unless something has happened, such as input, a resize or new data. It SHALL draw the correct screen for the window's size from the start, without first showing a different one. The cursor SHALL NOT blink, and Chorale SHALL put the cursor style back to the terminal's default when it exits.

#### Scenario: Idle screen stays still

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the empty session list and has written nothing for half a second
- **WHEN** two seconds pass with no input and no new data
- **THEN** Chorale writes zero bytes to the terminal

#### Scenario: Cursor does not blink

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions and has finished drawing
- **THEN** the cursor style Chorale has set is a steady one, not a blinking one
- **AND** after Ctrl-C is pressed twice on the empty prompt, the cursor style is back to the terminal's default

#### Scenario: Startup draws the session list directly

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** nothing Chorale writes before the empty session list appears contains "Chorale needs at least"

### Requirement: Quitting

Chorale SHALL quit on a second Ctrl-C, as "Ctrl-C asks before quitting" describes. The `q` key SHALL NOT quit anywhere. When Chorale quits, it SHALL leave the terminal as it found it: every terminal mode it turned on is turned off again, and the cursor is visible.

#### Scenario: Quitting restores the terminal

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the empty session list on the alternate screen, with an empty prompt
- **WHEN** Ctrl-C is pressed twice
- **THEN** Chorale exits with status 0
- **AND** every terminal mode Chorale turned on during the run is off, including the alternate screen, bracketed paste, mouse reporting and keyboard enhancements
- **AND** the cursor is visible

#### Scenario: Quitting from the too-small screen

- **GIVEN** a terminal 40 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the window-too-small message
- **WHEN** Ctrl-C is pressed twice
- **THEN** Chorale exits with status 0
- **AND** every terminal mode Chorale turned on during the run is off
- **AND** the cursor is visible

#### Scenario: q types on the session list

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** `q` is pressed
- **THEN** Chorale is still running
- **AND** the prompt is "q"

#### Scenario: q does nothing on the too-small screen

- **GIVEN** a terminal 40 columns wide and 20 rows tall, in the default theme, with colour on, showing the window-too-small message
- **WHEN** `q` is pressed
- **THEN** Chorale is still running
- **AND** the screen still shows "Chorale needs at least 60 by 20."

### Requirement: Window too small

When the window is narrower than 60 columns or shorter than 20 rows, Chorale SHALL replace the screen with a short message saying the size it needs, wrapped to the window width, followed by "Ctrl-C quit" when there is room for it. The message comes first: when the window is too short for both, "Ctrl-C quit" is left out. Chorale SHALL switch between this message and the normal screen whenever a resize crosses that size.

#### Scenario: Started at 40 columns

- **GIVEN** a terminal 40 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts
- **THEN** the screen shows "Chorale needs at least 60 by 20."
- **AND** still shows "Ctrl-C quit"
- **AND** no line is wider than 40 columns

#### Scenario: Started at 15 rows

- **GIVEN** a terminal 80 columns wide and 15 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts
- **THEN** the screen shows "Chorale needs at least 60 by 20."
- **AND** still shows "Ctrl-C quit"
- **AND** no line is wider than 80 columns

#### Scenario: Message wraps in a tiny window

- **GIVEN** a terminal 20 columns wide and 10 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts
- **THEN** no line is wider than 20 columns
- **AND** the screen's lines, joined with single spaces, contain "Chorale needs at least 60 by 20."
- **AND** the screen shows "Ctrl-C quit"

#### Scenario: Message comes first in a very short window

- **GIVEN** a terminal 20 columns wide and 2 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts
- **THEN** the screen's lines, joined with single spaces, contain "Chorale needs at least 60 by 20."
- **AND** no line is wider than 20 columns

#### Scenario: Too-small message with colour off

- **GIVEN** a terminal 40 columns wide and 20 rows tall, in the default theme, with `NO_COLOR=1`
- **WHEN** Chorale starts
- **THEN** the bytes Chorale writes contain no colour codes
- **AND** the screen shows "Chorale needs at least 60 by 20." and "Ctrl-C quit"

#### Scenario: Window widened to 60 columns

- **GIVEN** a terminal 40 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the window-too-small message
- **WHEN** the window is resized to 60 columns by 20 rows
- **THEN** the screen shows "♪ chorale", "No sessions yet.", "▸ Prompt" and "Enter send · Ctrl-C quit"
- **AND** no longer shows "Chorale needs at least"

#### Scenario: Window narrowed while running

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the empty session list
- **WHEN** the window is resized to 40 columns by 20 rows
- **THEN** the screen shows "Chorale needs at least 60 by 20." and "Ctrl-C quit"
- **AND** no line is wider than 40 columns

## REMOVED Requirements

### Requirement: Empty session list on first start

**Reason**: The empty session list now has a focused prompt box, so the cursor sits in it rather than on "No sessions yet.", and the key line changes. OpenSpec doesn't let a modified requirement drop a scenario, so the requirement is replaced.

**Migration**: Its replacement is "Empty session list with the prompt box", in this spec.
