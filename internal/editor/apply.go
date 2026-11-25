package editor

// Apply does a key's action, if the editor can do it by itself.
func (e *Editor) Apply(a Action) {
	switch a.Kind {
	case ActionMove:
		e.Move(a.Move, a.Extend)
	case ActionSelectAll:
		e.SelectAll()
	case ActionBackspace:
		e.Backspace()
	case ActionDelete:
		e.Delete()
	case ActionDeleteWordBack:
		e.DeleteWordBack()
	case ActionDeleteWordForward:
		e.DeleteWordForward()
	case ActionDeleteToRowStart:
		e.DeleteToRowStart()
	case ActionDeleteToRowEnd:
		e.DeleteToRowEnd()
	case ActionUndo:
		e.Undo()
	case ActionRedo:
		e.Redo()
	case ActionNewLine:
		e.Insert("\n")
	}
}
