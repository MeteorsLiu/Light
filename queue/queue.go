package queue

import (
	"context"
	"errors"
	"time"
)

//This is a simple queue implements the GO's FIFO Buffer Channel

type Queue struct {
	Queue chan interface{}
}

const (
	MAXSIZE    = 200
	RETRY_TIME = 5
	WAIT_TIME  = 100 * time.Microsecond
)

func New(ctx context.Context) *Queue {
	return &Queue{
		Queue: make(chan interface{}, MAXSIZE),
	}
}

func (q *Queue) Push(elem interface{}) error {
	count := 0
	for {
		select {
		case q.Queue <- elem:
			return nil
		default:
			//Full wait until the channel is available
			if count > RETRY_TIME {
				return errors.New("Queue is full")
			}
			count++
			time.Sleep(WAIT_TIME)
		}
	}
}

func (q *Queue) Reset() {
	q.Queue = nil
}
