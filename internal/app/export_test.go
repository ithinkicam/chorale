package app

import tea "charm.land/bubbletea/v2"

// PendingQuit reports whether m, which must be a Model, is waiting for a second Ctrl-C.
func PendingQuit(m tea.Model) bool {
	return m.(Model).pendingQuit
}
