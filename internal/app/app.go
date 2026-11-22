// Package app holds the Bubble Tea model for the whole of Chorale.
package app

import tea "charm.land/bubbletea/v2"

// Model is Chorale's top-level Bubble Tea model.
type Model struct{}

// New returns the model Chorale starts with.
func New() Model { return Model{} }

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

// View implements tea.Model.
func (m Model) View() tea.View { return tea.NewView("") }
