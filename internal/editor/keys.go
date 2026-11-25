package editor

import "slices"

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

var bindings = []Binding{
	{bindMove(MoveLeft, false), []string{"left", "ctrl+b"}},
	{bindMove(MoveLeft, true), []string{"shift+left", "ctrl+shift+b"}},
	{bindMove(MoveRight, false), []string{"right", "ctrl+f"}},
	{bindMove(MoveRight, true), []string{"shift+right", "ctrl+shift+f"}},
	{bindMove(MoveWordLeft, false), []string{"alt+left", "alt+b"}},
	{bindMove(MoveWordLeft, true), []string{"alt+shift+left", "alt+shift+b"}},
	{bindMove(MoveWordRight, false), []string{"alt+right", "alt+f"}},
	{bindMove(MoveWordRight, true), []string{"alt+shift+right", "alt+shift+f"}},
	{bindMove(MoveUp, false), []string{"up"}},
	{bindMove(MoveUp, true), []string{"shift+up"}},
	{bindMove(MoveDown, false), []string{"down"}},
	{bindMove(MoveDown, true), []string{"shift+down"}},
	{bindMove(MoveRowStart, false), []string{"home", "ctrl+a", "super+left"}},
	{bindMove(MoveRowStart, true), []string{"shift+home", "shift+meta+left", "shift+super+left"}},
	{bindMove(MoveRowEnd, false), []string{"end", "ctrl+e", "super+right"}},
	{bindMove(MoveRowEnd, true), []string{"shift+end", "ctrl+shift+e", "shift+meta+right", "shift+super+right"}},
	{bindMove(MoveStart, false), []string{"alt+up", "ctrl+home", "super+up"}},
	{bindMove(MoveStart, true), []string{"alt+shift+up", "ctrl+shift+home", "shift+super+up"}},
	{bindMove(MoveEnd, false), []string{"alt+down", "ctrl+end", "super+down"}},
	{bindMove(MoveEnd, true), []string{"alt+shift+down", "ctrl+shift+end", "shift+super+down"}},
	{Action{Kind: ActionSelectAll}, []string{"ctrl+shift+a", "super+a", "shift+super+a"}},
	{Action{Kind: ActionBackspace}, []string{"backspace", "ctrl+h"}},
	{Action{Kind: ActionDelete}, []string{"delete", "ctrl+d"}},
	{Action{Kind: ActionDeleteWordBack}, []string{"alt+backspace", "ctrl+w"}},
	{Action{Kind: ActionDeleteWordForward}, []string{"alt+delete", "alt+d"}},
	{Action{Kind: ActionDeleteToRowStart}, []string{"ctrl+u"}},
	{Action{Kind: ActionDeleteToRowEnd}, []string{"ctrl+k"}},
	{Action{Kind: ActionUndo}, []string{"ctrl+z", "super+z"}},
	{Action{Kind: ActionRedo}, []string{"ctrl+y", "ctrl+shift+z", "shift+super+z"}},
	{Action{Kind: ActionCopy}, []string{"super+c"}},
	{Action{Kind: ActionCut}, []string{"ctrl+x", "super+x"}},
	{Action{Kind: ActionInterrupt}, []string{"ctrl+c"}},
	{Action{Kind: ActionSend}, []string{"enter"}},
	{Action{Kind: ActionNewLine}, []string{"shift+enter", "alt+enter", "ctrl+j"}},
}

var byKey = func() map[string]Action {
	table := map[string]Action{}
	for _, b := range bindings {
		for _, key := range b.Keys {
			table[key] = b.Action
		}
	}
	return table
}()

func bindMove(m Movement, extend bool) Action {
	return Action{Kind: ActionMove, Move: m, Extend: extend}
}

// Bindings returns every key binding.
func Bindings() []Binding { return slices.Clone(bindings) }

// Lookup returns the action bound to key, which is a name from tea.KeyPressMsg's String.
func Lookup(key string) (Action, bool) {
	a, ok := byKey[key]
	return a, ok
}
