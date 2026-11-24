# Spec Delta

## RENAMED Requirements

- FROM: `### Requirement: Quit with q`
- TO: `### Requirement: Quitting`

## ADDED Requirements

### Requirement: Empty session list with the prompt box

When Chorale starts and there are no sessions, it SHALL show a full-screen view with the Chorale name on the first line, a message that there are no sessions yet, an empty prompt box labelled "Prompt" with focus, and a last line listing the keys that work. The real terminal cursor SHALL be visible in the prompt box, so macOS Zoom follows it there.

#### Scenario: Empty session list at 60 by 20 with colour

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** the first line starts with "♪ chorale"
- **AND** the screen shows "No sessions yet."
- **AND** the screen shows the prompt box's label, "Prompt", below that message
- **AND** the last line shows "Enter send · Ctrl-C quit"
- **AND** no line is wider than 60 columns

#### Scenario: Empty session list at 120 by 30 with colour

- **GIVEN** a terminal 120 columns wide and 30 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** the first line starts with "♪ chorale"
- **AND** the screen shows "No sessions yet."
- **AND** the screen shows the prompt box's label, "Prompt"
- **AND** "Enter send · Ctrl-C quit" is on the last row
- **AND** no line is wider than 120 columns

#### Scenario: Cursor in the empty prompt box

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions and has finished drawing
- **THEN** the terminal cursor is visible
- **AND** it sits on the first column of the row below the "Prompt" label

## MODIFIED Requirements

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
- **AND** after Ctrl-C is pressed on the empty prompt, the cursor style is back to the terminal's default

#### Scenario: Startup draws the session list directly

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** nothing Chorale writes before the empty session list appears contains "Chorale needs at least"

### Requirement: Quitting

Chorale SHALL quit when Ctrl-C is pressed on the session list with an empty prompt, and when `q` or Ctrl-C is pressed on the window-too-small message. It SHALL leave the terminal as it found it: every terminal mode it turned on is turned off again, and the cursor is visible.

#### Scenario: Quitting restores the terminal

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the empty session list on the alternate screen, with an empty prompt
- **WHEN** Ctrl-C is pressed
- **THEN** Chorale exits with status 0
- **AND** every terminal mode Chorale turned on during the run is off, including the alternate screen, bracketed paste, mouse reporting and keyboard enhancements
- **AND** the cursor is visible

#### Scenario: Quitting from the too-small screen

- **GIVEN** a terminal 40 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the window-too-small message
- **WHEN** `q` is pressed
- **THEN** Chorale exits with status 0
- **AND** every terminal mode Chorale turned on during the run is off
- **AND** the cursor is visible

#### Scenario: q types on the session list

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on, showing the empty session list
- **WHEN** `q` is pressed
- **THEN** Chorale is still running
- **AND** the prompt is "q"

## REMOVED Requirements

### Requirement: Empty session list on first start

**Reason**: The empty session list now has a focused prompt box, so the cursor sits in it rather than on "No sessions yet.", and the key line changes. OpenSpec doesn't let a modified requirement drop a scenario, so the requirement is replaced.

**Migration**: Its replacement is "Empty session list with the prompt box", in this spec.
