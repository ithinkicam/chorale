# Spec Delta

## Purpose

The app shell is how Chorale starts, quits, reports its version, and what it shows before there is any work in it. It also covers how Chorale behaves when the window is too narrow to use.

## ADDED Requirements

### Requirement: Empty session list on first start

When Chorale starts and there are no sessions, it SHALL show a full-screen view with the Chorale name on the first line, a message that there are no sessions yet, and a line listing the keys that work.

#### Scenario: Empty session list at 60 by 20 with colour

- **GIVEN** a terminal 60 columns wide and 20 rows tall, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** the first line starts with "♪ chorale"
- **AND** the screen shows the words "No sessions yet."
- **AND** the last line shows "q quit"
- **AND** no line is wider than 60 columns

#### Scenario: Empty session list at 120 by 30 with colour

- **GIVEN** a terminal 120 columns wide and 30 rows tall, with colour on
- **WHEN** Chorale starts with no sessions
- **THEN** the screen shows the same words as at 60 by 20
- **AND** the key line is on the last row

### Requirement: Readable without colour

Chorale SHALL respect the `NO_COLOR` environment variable by sending no colour escape codes, and every piece of information on screen SHALL still be shown in words or symbols.

#### Scenario: Empty session list with colour off

- **GIVEN** a terminal 60 columns wide and 20 rows tall, and `NO_COLOR=1`
- **WHEN** Chorale starts with no sessions
- **THEN** the output contains no colour escape codes
- **AND** the screen shows the same words as with colour on

### Requirement: No motion

Chorale SHALL NOT change what is on screen unless something has happened, such as input, a resize or new data.

#### Scenario: Idle screen stays still

- **GIVEN** Chorale is showing the empty session list
- **WHEN** one second passes with no input and no new data
- **THEN** the screen is identical to what it was before

### Requirement: Quit with q

Chorale SHALL quit when `q` is pressed on the session list, and SHALL leave the terminal as it found it.

#### Scenario: Quitting restores the terminal

- **GIVEN** Chorale is showing the empty session list
- **WHEN** `q` is pressed
- **THEN** Chorale exits with status 0
- **AND** the terminal has left the alternate screen, the cursor is visible, and mouse reporting is turned off

### Requirement: Report the version

`chorale --version` SHALL print the version and exit without opening the interface.

#### Scenario: Version flag

- **WHEN** `chorale --version` is run
- **THEN** it prints one line starting with "chorale " followed by the version
- **AND** exits with status 0
- **AND** does not switch to the alternate screen

### Requirement: Window too narrow

When the window is narrower than 60 columns or shorter than 20 rows, Chorale SHALL replace the screen with a short message saying the size it needs, wrapped to the window width, and SHALL redraw normally once the window is big enough.

#### Scenario: Started at 40 columns

- **GIVEN** a terminal 40 columns wide and 20 rows tall
- **WHEN** Chorale starts
- **THEN** the screen shows "Chorale needs at least 60 by 20." with no line wider than 40 columns
- **AND** still shows "q quit"

#### Scenario: Window widened to 60 columns

- **GIVEN** Chorale is showing the window-too-narrow message
- **WHEN** the window is resized to 60 columns by 20 rows
- **THEN** the empty session list is shown

#### Scenario: Started at 15 rows

- **GIVEN** a terminal 80 columns wide and 15 rows tall
- **WHEN** Chorale starts
- **THEN** the screen shows "Chorale needs at least 60 by 20." with no line wider than 80 columns
- **AND** still shows "q quit"
