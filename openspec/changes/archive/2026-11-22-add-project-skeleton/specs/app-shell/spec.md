# Spec Delta

## Purpose

The app shell is how Chorale starts, quits, reports its version, and what it shows before there is any work in it. It also covers how Chorale behaves when the window is too small to use.

Every scenario about the screen uses the default theme, which draws with the terminal's own palette. It is the only theme in this change.

The exact words on the empty session list are "♪ chorale" on the first line, "No sessions yet." in the body, and "q quit" on the last line. A scenario that says "the empty session list words" means all three.

A colour code is a Select Graphic Rendition escape sequence that sets a foreground or background colour: parameters 30 to 39, 40 to 49, 90 to 97 or 100 to 107, including the extended forms that start with 38 or 48.

## ADDED Requirements

### Requirement: Empty session list on first start

When Chorale starts and there are no sessions, it SHALL show a full-screen view with the Chorale name on the first line, a message that there are no sessions yet, and a line listing the keys that work. The real terminal cursor SHALL be visible at the start of the "No sessions yet." line, so macOS Zoom follows it there.

#### Scenario: Empty session list at 60 by 20 with colour

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** the first line starts with "♪ chorale"
- **AND** the screen shows "No sessions yet."
- **AND** the last line shows "q quit"
- **AND** no line is wider than 60 columns

#### Scenario: Empty session list at 120 by 30 with colour

- **GIVEN** a terminal 120 columns wide and 30 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** the first line starts with "♪ chorale"
- **AND** the screen shows "No sessions yet."
- **AND** "q quit" is on the last row
- **AND** no line is wider than 120 columns

#### Scenario: Cursor on the empty message

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions and has finished drawing
- **THEN** the terminal cursor is visible
- **AND** it sits on the first character of "No sessions yet."

### Requirement: Readable without colour

Chorale SHALL respect the `NO_COLOR` environment variable by sending no colour codes, and every piece of information on screen SHALL still be shown in words or symbols. With colour allowed, the default theme SHALL use only the terminal's own palette.

#### Scenario: Empty session list with colour off

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with `NO_COLOR=1`
- **WHEN** Chorale starts with no sessions
- **THEN** the bytes Chorale writes to the terminal contain no colour codes
- **AND** the screen shows the empty session list words

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
- **AND** after `q` is pressed, the cursor style is back to the terminal's default

#### Scenario: Startup draws the session list directly

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** nothing Chorale writes before the empty session list appears contains "Chorale needs at least"

### Requirement: Quit with q

Chorale SHALL quit when `q` is pressed on the session list or on the window-too-small message, and SHALL leave the terminal as it found it: every terminal mode it turned on is turned off again, and the cursor is visible.

#### Scenario: Quitting restores the terminal

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the empty session list on the alternate screen
- **WHEN** `q` is pressed
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

### Requirement: Report the version

`chorale --version` SHALL print the version and exit without opening the interface. A build with no version set SHALL report the version `dev`.

#### Scenario: Version flag

- **GIVEN** a `chorale` binary built with the version set to `0.0.0-test`
- **WHEN** `chorale --version` is run in a terminal 60 columns wide and 20 rows tall
- **THEN** it prints exactly one line, "chorale 0.0.0-test"
- **AND** exits with status 0
- **AND** never switches to the alternate screen

#### Scenario: Development build version

- **GIVEN** a `chorale` binary built with no version set
- **WHEN** `chorale --version` is run
- **THEN** it prints exactly one line, "chorale dev"
- **AND** exits with status 0

### Requirement: Window too small

When the window is narrower than 60 columns or shorter than 20 rows, Chorale SHALL replace the screen with a short message saying the size it needs, wrapped to the window width, followed by "q quit" when there is room for it. The message comes first: when the window is too short for both, "q quit" is left out. Chorale SHALL switch between this message and the normal screen whenever a resize crosses that size.

#### Scenario: Started at 40 columns

- **GIVEN** a terminal 40 columns wide and 20 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts
- **THEN** the screen shows "Chorale needs at least 60 by 20."
- **AND** still shows "q quit"
- **AND** no line is wider than 40 columns

#### Scenario: Started at 15 rows

- **GIVEN** a terminal 80 columns wide and 15 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts
- **THEN** the screen shows "Chorale needs at least 60 by 20."
- **AND** still shows "q quit"
- **AND** no line is wider than 80 columns

#### Scenario: Message wraps in a tiny window

- **GIVEN** a terminal 20 columns wide and 10 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts
- **THEN** no line is wider than 20 columns
- **AND** the screen's lines, joined with single spaces, contain "Chorale needs at least 60 by 20."
- **AND** the screen shows "q quit"

#### Scenario: Message comes first in a very short window

- **GIVEN** a terminal 20 columns wide and 2 rows tall, in the default theme, with colour on
- **WHEN** Chorale starts
- **THEN** the screen's lines, joined with single spaces, contain "Chorale needs at least 60 by 20."
- **AND** no line is wider than 20 columns

#### Scenario: Too-small message with colour off

- **GIVEN** a terminal 40 columns wide and 20 rows tall, in the default theme, with `NO_COLOR=1`
- **WHEN** Chorale starts
- **THEN** the bytes Chorale writes contain no colour codes
- **AND** the screen shows "Chorale needs at least 60 by 20." and "q quit"

#### Scenario: Window widened to 60 columns

- **GIVEN** a terminal 40 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the window-too-small message
- **WHEN** the window is resized to 60 columns by 20 rows
- **THEN** the screen shows the empty session list words
- **AND** no longer shows "Chorale needs at least"

#### Scenario: Window narrowed while running

- **GIVEN** a terminal 60 columns wide and 20 rows tall, in the default theme, with colour on
- **AND** Chorale is showing the empty session list
- **WHEN** the window is resized to 40 columns by 20 rows
- **THEN** the screen shows "Chorale needs at least 60 by 20." and "q quit"
- **AND** no line is wider than 40 columns
