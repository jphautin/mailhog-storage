package storage

import (
	"github.com/mailhog/data"
	"sort"
	"strings"
)

// Storage represents a storage backend
type Storage interface {
	Store(m *data.Message) (string, error)
	List(start, limit int, field, order string) (*data.Messages, error)
	Search(kind, query string, start, limit int, field, order string) (*data.Messages, int, error)
	Count() int
	DeleteOne(id string) error
	DeleteAll() error
	Load(id string) (*data.Message, error)
}

func sortMessages(messages []data.Message, field, order string) {
	if field == "time" {
		sort.SliceStable(messages, func(i, j int) bool {
			return messages[i].Created.After(messages[j].Created)
		})
	}
	if field == "size" {
		sort.SliceStable(messages, func(i, j int) bool {
			return len(messages[i].Raw.Data) > len(messages[j].Raw.Data)
		})
	}
	if field == "to" {
		sort.SliceStable(messages, func(i, j int) bool {
			return strings.Compare(recipients(messages[i]), recipients(messages[j])) >= 0
		})
	}
	if field == "from" {
		sort.SliceStable(messages, func(i, j int) bool {
			return strings.Compare(sender(messages[i]), sender(messages[j])) >= 0
		})
	}
}

func email(path *data.Path) string {
	return path.Mailbox + "@" + path.Domain
}

func sender(message data.Message) string {
	return email(message.From)
}

func recipients(message data.Message) string {
	var emails = make([]string, 0)
	for _, to := range message.To {
		emails = append(emails, email(to))
	}
	return strings.Join(emails, ",")
}
