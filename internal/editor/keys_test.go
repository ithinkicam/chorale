package editor

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func press(code rune, mod tea.KeyMod) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code, Mod: mod}
}

func move(m Movement, extend bool) Action {
	return Action{Kind: ActionMove, Move: m, Extend: extend}
}

func plain(kind ActionKind) Action { return Action{Kind: kind} }

const (
	shift = tea.ModShift
	alt   = tea.ModAlt
	ctrl  = tea.ModCtrl
	meta  = tea.ModMeta
	super = tea.ModSuper
)

func assertBound(t *testing.T, msg tea.KeyPressMsg, want Action) {
	t.Helper()
	got, ok := Lookup(msg.String())
	if !ok {
		t.Errorf("%q is bound to nothing, want %+v", msg.String(), want)
		return
	}
	if got != want {
		t.Errorf("%q is bound to %+v, want %+v", msg.String(), got, want)
	}
}

func TestEveryKeyInTheDesignIsBoundToItsAction(t *testing.T) {
	tests := []struct {
		msg  tea.KeyPressMsg
		want Action
	}{
		{press(tea.KeyLeft, 0), move(MoveLeft, false)},
		{press('b', ctrl), move(MoveLeft, false)},
		{press(tea.KeyRight, 0), move(MoveRight, false)},
		{press('f', ctrl), move(MoveRight, false)},
		{press('b', alt), move(MoveWordLeft, false)},
		{press(tea.KeyLeft, alt), move(MoveWordLeft, false)},
		{press('f', alt), move(MoveWordRight, false)},
		{press(tea.KeyRight, alt), move(MoveWordRight, false)},
		{press(tea.KeyUp, 0), move(MoveUp, false)},
		{press(tea.KeyDown, 0), move(MoveDown, false)},
		{press(tea.KeyHome, 0), move(MoveRowStart, false)},
		{press('a', ctrl), move(MoveRowStart, false)},
		{press(tea.KeyEnd, 0), move(MoveRowEnd, false)},
		{press('e', ctrl), move(MoveRowEnd, false)},
		{press(tea.KeyUp, alt), move(MoveStart, false)},
		{press(tea.KeyHome, ctrl), move(MoveStart, false)},
		{press(tea.KeyDown, alt), move(MoveEnd, false)},
		{press(tea.KeyEnd, ctrl), move(MoveEnd, false)},
		{press(tea.KeyBackspace, 0), plain(ActionBackspace)},
		{press('h', ctrl), plain(ActionBackspace)},
		{press(tea.KeyDelete, 0), plain(ActionDelete)},
		{press('d', ctrl), plain(ActionDelete)},
		{press(tea.KeyBackspace, alt), plain(ActionDeleteWordBack)},
		{press('w', ctrl), plain(ActionDeleteWordBack)},
		{press(tea.KeyDelete, alt), plain(ActionDeleteWordForward)},
		{press('d', alt), plain(ActionDeleteWordForward)},
		{press('u', ctrl), plain(ActionDeleteToRowStart)},
		{press('k', ctrl), plain(ActionDeleteToRowEnd)},
		{press('z', ctrl), plain(ActionUndo)},
		{press('y', ctrl), plain(ActionRedo)},
		{press('z', ctrl|shift), plain(ActionRedo)},
		{press('c', ctrl), plain(ActionInterrupt)},
		{press('x', ctrl), plain(ActionCut)},
		{press('a', ctrl|shift), plain(ActionSelectAll)},
		{press(tea.KeyEnter, 0), plain(ActionSend)},
		{press(tea.KeyEnter, shift), plain(ActionNewLine)},
		{press(tea.KeyEnter, alt), plain(ActionNewLine)},
		{press('j', ctrl), plain(ActionNewLine)},
	}
	for _, tt := range tests {
		assertBound(t, tt.msg, tt.want)
	}
}

func TestShiftWithAMovementKeyExtendsTheSelection(t *testing.T) {
	tests := []struct {
		msg  tea.KeyPressMsg
		want Action
	}{
		{press(tea.KeyLeft, shift), move(MoveLeft, true)},
		{press('b', ctrl|shift), move(MoveLeft, true)},
		{press(tea.KeyRight, shift), move(MoveRight, true)},
		{press('f', ctrl|shift), move(MoveRight, true)},
		{press(tea.KeyLeft, alt|shift), move(MoveWordLeft, true)},
		{press('b', alt|shift), move(MoveWordLeft, true)},
		{press(tea.KeyRight, alt|shift), move(MoveWordRight, true)},
		{press('f', alt|shift), move(MoveWordRight, true)},
		{press(tea.KeyUp, shift), move(MoveUp, true)},
		{press(tea.KeyDown, shift), move(MoveDown, true)},
		{press(tea.KeyHome, shift), move(MoveRowStart, true)},
		{press(tea.KeyEnd, shift), move(MoveRowEnd, true)},
		{press('e', ctrl|shift), move(MoveRowEnd, true)},
		{press(tea.KeyUp, alt|shift), move(MoveStart, true)},
		{press(tea.KeyHome, ctrl|shift), move(MoveStart, true)},
		{press(tea.KeyDown, alt|shift), move(MoveEnd, true)},
		{press(tea.KeyEnd, ctrl|shift), move(MoveEnd, true)},
	}
	for _, tt := range tests {
		assertBound(t, tt.msg, tt.want)
	}
}

func TestCtrlShiftAIsSelectAllNotShiftCtrlA(t *testing.T) {
	assertBound(t, press('a', ctrl|shift), plain(ActionSelectAll))
	assertBound(t, press('a', ctrl), move(MoveRowStart, false))
}

func TestCmdKeysAsTheTerminalSendsThem(t *testing.T) {
	tests := []struct {
		msg  tea.KeyPressMsg
		want Action
	}{
		{press('a', ctrl), move(MoveRowStart, false)},
		{press('e', ctrl), move(MoveRowEnd, false)},
		{press(tea.KeyLeft, shift|meta), move(MoveRowStart, true)},
		{press(tea.KeyRight, shift|meta), move(MoveRowEnd, true)},
		{press('z', super), plain(ActionUndo)},
		{press('z', shift|super), plain(ActionRedo)},
		{press('c', super), plain(ActionCopy)},
		{press('x', super), plain(ActionCut)},
	}
	for _, tt := range tests {
		assertBound(t, tt.msg, tt.want)
	}
}

func TestCmdKeysAsASuperKey(t *testing.T) {
	tests := []struct {
		msg  tea.KeyPressMsg
		want Action
	}{
		{press(tea.KeyLeft, super), move(MoveRowStart, false)},
		{press(tea.KeyRight, super), move(MoveRowEnd, false)},
		{press(tea.KeyUp, super), move(MoveStart, false)},
		{press(tea.KeyDown, super), move(MoveEnd, false)},
		{press('a', super), plain(ActionSelectAll)},
		{press(tea.KeyLeft, shift|super), move(MoveRowStart, true)},
		{press(tea.KeyRight, shift|super), move(MoveRowEnd, true)},
		{press(tea.KeyUp, shift|super), move(MoveStart, true)},
		{press(tea.KeyDown, shift|super), move(MoveEnd, true)},
		{press('a', shift|super), plain(ActionSelectAll)},
	}
	for _, tt := range tests {
		assertBound(t, tt.msg, tt.want)
	}
}

func requireBindings(t *testing.T) []Binding {
	t.Helper()
	all := Bindings()
	if len(all) == 0 {
		t.Fatal("the key table has no bindings")
	}
	return all
}

func TestTabAndCtrlVAreBoundToNothing(t *testing.T) {
	requireBindings(t)
	for _, msg := range []tea.KeyPressMsg{press(tea.KeyTab, 0), press('v', ctrl)} {
		if got, ok := Lookup(msg.String()); ok {
			t.Errorf("%q is bound to %+v, want nothing", msg.String(), got)
		}
	}
}

func TestNoKeyIsBoundToTwoActions(t *testing.T) {
	seen := map[string]Action{}
	for _, b := range requireBindings(t) {
		for _, key := range b.Keys {
			if prev, dup := seen[key]; dup {
				t.Errorf("%q is bound to both %+v and %+v", key, prev, b.Action)
			}
			seen[key] = b.Action
		}
	}
}

func TestLookupAgreesWithTheBindings(t *testing.T) {
	for _, b := range requireBindings(t) {
		if len(b.Keys) == 0 {
			t.Errorf("%+v has no keys", b.Action)
		}
		for _, key := range b.Keys {
			if got, ok := Lookup(key); !ok || got != b.Action {
				t.Errorf("Lookup(%q) is %+v, %v; want %+v", key, got, ok, b.Action)
			}
		}
	}
}

func TestEveryActionHasAKey(t *testing.T) {
	bound := map[Action]bool{}
	for _, b := range requireBindings(t) {
		bound[b.Action] = true
	}
	for kind := ActionSelectAll; kind <= ActionNewLine; kind++ {
		if !bound[plain(kind)] {
			t.Errorf("action %d has no key", kind)
		}
	}
	for m := MoveLeft; m <= MoveEnd; m++ {
		for _, extend := range []bool{false, true} {
			if !bound[move(m, extend)] {
				t.Errorf("movement %d, extending %v, has no key", m, extend)
			}
		}
	}
}
