package test

import (
	"github.com/rdnply/url-shortener/internal/counter"
	"github.com/rdnply/url-shortener/internal/link"
)

var _ counter.Storage = &MockCounterStorage{}

type MockCounterStorage struct {
	Value uint
}

func (m *MockCounterStorage) Init() error {
	m.Value = 0

	return nil
}

func (m *MockCounterStorage) Increment() (uint, error) {
	m.Value++
	return m.Value, nil
}

var _ link.Storage = &MockLinkStorage{}

type MockLinkStorage struct {
	Items []*link.Link
}

func (m *MockLinkStorage) AddLink(link *link.Link) (uint, error) {
	m.Items = append(m.Items, link)

	return uint(len(m.Items)), nil
}

func (m *MockLinkStorage) GetLinkByShortID(shortID string) (*link.Link, error) {
	for _, l := range m.Items {
		if l.ShortID == shortID {
			return l, nil
		}
	}

	return nil, nil
}

func (m *MockLinkStorage) IncrementLinkCounter(link *link.Link) (uint, error) {
	for _, l := range m.Items {
		if l.ShortID == link.ShortID {
			l.Clicks++

			return l.Clicks, nil
		}
	}

	return 0, nil
}
