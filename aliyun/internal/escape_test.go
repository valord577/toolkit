package internal

import "testing"

func TestEscape(t *testing.T) {
	t.Logf("%s", Escape(`" ~*`))
	t.Logf("%s", Escape(`a_bwcASdsw`))
}
