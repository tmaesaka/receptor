// Package subscription implements subscription management.
package subscription

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/toru/dexter/feed"
)

// Subscription represents a subscription to a data feed.
type Subscription struct {
	ID      [sha256.Size224]byte // Unique ID
	FeedURL url.URL              // URL of the data endpooint

	unreachable  bool // Consider using a enum
	checksum     [sha256.Size224]byte
	createdAt    time.Time
	lastSyncedAt time.Time
}

// New returns a new Subscription.
func New(feedURL string) (*Subscription, error) {
	u, err := url.Parse(feedURL)
	if err != nil {
		return nil, err
	}

	s := &Subscription{
		ID:      sha256.Sum224([]byte(feedURL)),
		FeedURL: *u,
	}
	return s, nil
}

// IsOffline returns a boolean indicating the data feed reachability.
func (s *Subscription) IsOffline() bool {
	// TODO(toru): Somehow allow to retry. Maybe exponential backoff.
	return s.unreachable
}

// Sync downloads the data feed and parses it.
func (s *Subscription) Sync() error {
	if len(s.FeedURL.String()) == 0 {
		return fmt.Errorf("subscription has no FeedURL")
	}
	if s.unreachable {
		return fmt.Errorf("%s is unreachable", s.FeedURL.String())
	}

	// TODO(toru): This is only for dev-purpose. Craft a proper HTTP
	// client with defensive settings like network timeout.
	resp, err := http.Get(s.FeedURL.String())
	s.lastSyncedAt = time.Now().UTC()
	if err != nil {
		s.unreachable = true
		return err
	}
	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		s.unreachable = true
		return fmt.Errorf("sync failure (%d): %s", resp.StatusCode, s.FeedURL.String())
	}

	checksum := sha256.Sum224(payload)
	if bytes.Equal(s.checksum[:], checksum[:]) {
		// There's no new content to process.
		return nil
	}
	s.checksum = checksum

	if feed.IsAtomFeed(payload) {
		af, err := feed.ParseAtomFeed(payload)
		if err != nil {
			return err
		}
		af.SubscriptionID = s.ID
		// TODO(toru): Store the delta to persistent storage
	} else {
		return fmt.Errorf("unknown syndication format")
	}

	return nil
}
