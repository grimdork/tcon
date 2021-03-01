package tcon_test

import (
	"testing"

	"github.com/grimdork/tcon"
)

var (
	text         = `WordWrap returns one or more lines, based on the desired width. This will be split at spaces or punctuation as far as possible, but if using small widths and long words they will be split and a hyphen added.`
	annoyingtext = `ybOQDvjm4aY7ASzUDckEEsVliwZBNtBBXcsEaXHuHtOUMsA4YPuLQbjxLHumMRakvTDJTBUY4MKluMNUGbTMbky5G4SjMJjKpqtQjt-c6d5R9ZKg7yJUbAjXTPxd8cIXfDSAHlTkJhjfSv9NFBV5LFUlEr9ynqNaoXzvtfeedCbB_73bQDi38VS_oULXCVMx7UXKUtE98LdJVr7UJnDuLISyRR2tky8Fe0o8AbtSY2xO_pGYiFlNmOdubVuaaojijmFLx_EJRSMsIMBqPuel0KfNESMAxfijONytqLrNfKdr`
)

func TestWordWrap(t *testing.T) {
	t.Logf("Splitting text string of %d characters into 80-character lines.", len(text))
	lines := tcon.WordWrap(text, 80)
	for _, l := range lines {
		t.Logf("%s", l)
	}

	t.Logf("Splitting text string of %d characters into 10-character lines.", len(text))
	lines = tcon.WordWrap(text, 10)
	for _, l := range lines {
		t.Logf("%s", l)
	}

	t.Logf("Splitting spacing- and punctuation-less text string of %d characters into 80-character lines.", len(annoyingtext))
	lines = tcon.WordWrap(annoyingtext, 80)
	for _, l := range lines {
		t.Logf("%s", l)
	}
}
