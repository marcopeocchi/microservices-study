package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fuu/v/cmd/internal"
	"net/http"

	"go.uber.org/zap"
)

const consumerName = "knight"

type Server struct {
	srv    *http.Server
	logger *zap.SugaredLogger
	rmq    *internal.RabbitMQ
	done   chan struct{}
}

func (s *Server) ListenAndServe() error {
	queue, err := s.rmq.Channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = s.rmq.Channel.QueueBind(
		queue.Name,
		"gallery.event.*",
		"images",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := s.rmq.Channel.Consume(
		queue.Name,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			s.logger.Infow("consumer", "message", msg.RoutingKey)

			nack := false
			switch msg.RoutingKey {
			case "gallery.event.convert":
				var res string

				err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&res)
				if err != nil {
					s.logger.Errorw("decoding error", "error", err)
					return
				}
				go convert(res, "webp", s.logger)
			default:
				nack = true
			}

			if nack {
				s.logger.Warnw("consumer", "nack", nack)
				msg.Nack(false, nack)
			} else {
				s.logger.Infow("consumer", "ack", !nack)
				msg.Ack(false)
			}
		}

		s.done <- struct{}{}
	}()

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.rmq.Channel.Cancel(consumerName, false)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.done:
			return nil
		}
	}
}
