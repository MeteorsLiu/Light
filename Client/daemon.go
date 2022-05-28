package main

import (
	"context"

	"log"

	"github.com/MeteorsLiu/Light/queue"
)

func NewPollEventDaemon(ctx context.Context, q *queue.Queue, handler func(*interfaces.UploadPayload) error) {
	for {
		select {
		case elem := <-q.Queue:
			if err := handler(elem.(*interfaces.UploadPayload)); err != nil {
				log.Println(err.Error())
			}
		case <-ctx.Done():
			return
		}

	}
}
