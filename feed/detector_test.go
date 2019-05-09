package feed

import (
	"testing"
)

var rssDoc = []byte(`<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"></rss>`)
var atomDoc = []byte(`<?xml version="1.0" encoding="UTF-8"?><feed xmlns="http://www.w3.org/2005/Atom"></feed>`)

func TestIsAtomFeed(t *testing.T) {
	t.Run("with an atom document", func(t *testing.T) {
		if !IsAtomFeed(atomDoc) {
			t.Errorf("Got: false, Want: true")
		}
	})

	t.Run("with an rss document", func(t *testing.T) {
		if IsAtomFeed(rssDoc) {
			t.Errorf("Got: true, Want: false")
		}
	})
}