package broker

import (
	"context"
	"fmt"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/database"
)

type Module struct {
	isClosed    bool
	subscribers map[string][]chan broker.Message
	db          database.Database
}

func NewModule() broker.Broker {
	return &Module{
		subscribers: make(map[string][]chan broker.Message),
		db:          database.NewInMemory(),
	}
}

func (m *Module) Close() error {
	for _, v := range m.subscribers {
		for _, ch := range v {
			close(ch)
		}
	}
	fmt.Println(m.subscribers)
	return nil
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	if m.isClosed {
		return 0, broker.ErrUnavailable
	}
	for _, listener := range m.subscribers[subject] {
		go func(listener chan broker.Message) {
			listener <- msg
		}(listener)
	}
	msg.ID = m.db.SaveMessage(msg, subject)
	return msg.ID, nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	panic("implement me")
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	panic("implement me")
}
