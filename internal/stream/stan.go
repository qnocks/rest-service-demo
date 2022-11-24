package stream

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"team-task/internal/dto"
	"team-task/internal/storage"
	"team-task/pkg/logger"
)

type STANClient struct {
	Conn    stan.Conn
	Subject string
}

func NewSTANClient(clusterID, clientID, natsURL string, subject string) (*STANClient, error) {
	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}

	return &STANClient{Conn: conn, Subject: subject}, nil
}

func (c *STANClient) Publish(subject string, bytes []byte) error {
	return c.Conn.Publish(subject, bytes)
}

func (c *STANClient) Listen(store *storage.Storage) error {
	if _, err := c.Conn.Subscribe(c.Subject, func(msg *stan.Msg) {
		c.handleSubscribe(msg, store)
	}, stan.StartWithLastReceived()); err != nil {
		return err
	}

	return nil
}

func (c *STANClient) handleSubscribe(msg *stan.Msg, s *storage.Storage) {
	logger.Infof("receiving from stan: %s\n", string(msg.Data))

	var userGrade dto.UserGrade
	if err := json.Unmarshal(msg.Data, &userGrade); err != nil {
		logger.Warnf("error converting streamed bytes[] to UserGrade: %s\n", err.Error())
		return
	}

	_, err := s.Get(userGrade.UserId)
	if err != nil {
		s.Set(userGrade)
	}
}

func (c *STANClient) Shutdown() error {
	if err := c.Conn.Close(); err != nil {
		return err
	}

	return nil
}
