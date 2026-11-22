package termtest

// quote shows escape sequences readably in test messages, with ⎋ for ESC.
func quote(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r == 0x1b {
			out = append(out, '⎋')
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
