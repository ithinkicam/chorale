package editor

// ActionKind is what a key does to the prompt.
type ActionKind int

// The things a key can do.
const (
	ActionMove ActionKind = iota + 1
	ActionSelectAll
	ActionBackspace
	ActionDelete
	ActionDeleteWordBack
	ActionDeleteWordForward
	ActionDeleteToRowStart
	ActionDeleteToRowEnd
	ActionUndo
	ActionRedo
	ActionCopy
	ActionCut
	ActionInterrupt
	ActionSend
	ActionNewLine
)

// Action is a key's effect. Move and Extend apply only to ActionMove.
type Action struct {
	Kind   ActionKind
	Move   Movement
	Extend bool
}

// Binding gives the keys that do an action, named as Bubble Tea names them.
type Binding struct {
	Action Action
	Keys   []string
}

// Bindings returns every key binding.
func Bindings() []Binding { return nil }

// Lookup returns the action bound to key, which is a name from tea.KeyPressMsg's String.
func Lookup(_ string) (Action, bool) { return Action{}, false }
