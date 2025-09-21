package blocker

import (
	"fmt"
	"sort"
	"strings"

	"fogos/internal/hosts"
)

type Blocker struct {
	hostsManager *hosts.Manager
}

func New() *Blocker {
	return &Blocker{
		hostsManager: hosts.New(),
	}
}

func (b *Blocker) Block(website string) error {
	if blocked, err := b.IsBlocked(website); err != nil {
		return err
	} else if blocked {
		fmt.Printf("Website %s is already blocked\n", website)
		return nil
	}

	entries := []hosts.Entry{
		{
			IP:      "127.0.0.1",
			Domain:  website,
			Comment: hosts.FogosComment,
		},
		{
			IP:      "127.0.0.1",
			Domain:  "www." + website,
			Comment: hosts.FogosComment,
		},
	}

	return b.hostsManager.AddEntries(entries)
}

func (b *Blocker) Unblock(website string) error {
	if blocked, err := b.IsBlocked(website); err != nil {
		return err
	} else if !blocked {
		fmt.Printf("Website %s was not blocked\n", website)
		return nil
	}

	return b.hostsManager.RemoveEntriesByDomain(website)
}

func (b *Blocker) IsBlocked(website string) (bool, error) {
	entries, err := b.hostsManager.GetFogosEntries()
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.Domain == website || entry.Domain == "www."+website {
			return true, nil
		}
	}

	return false, nil
}

func (b *Blocker) ListBlocked() ([]string, error) {
	entries, err := b.hostsManager.GetFogosEntries()
	if err != nil {
		return nil, err
	}

	websiteSet := make(map[string]bool)
	for _, entry := range entries {
		domain := entry.Domain
		if strings.HasPrefix(domain, "www.") {
			domain = strings.TrimPrefix(domain, "www.")
		}
		websiteSet[domain] = true
	}

	websites := make([]string, 0, len(websiteSet))
	for website := range websiteSet {
		websites = append(websites, website)
	}

	sort.Strings(websites)
	return websites, nil
}
