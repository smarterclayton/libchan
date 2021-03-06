package utils

import (
	"github.com/docker/libchan"
)

type Queue struct {
	*libchan.PipeSender
	dst libchan.Sender
	ch  chan *libchan.Message
}

func NewQueue(dst libchan.Sender, size int) *Queue {
	r, w := libchan.Pipe()
	q := &Queue{
		PipeSender: w,
		dst:        dst,
		ch:         make(chan *libchan.Message, size),
	}
	go func() {
		defer close(q.ch)
		for {
			msg, err := r.Receive(libchan.Ret)
			if err != nil {
				r.Close()
				return
			}
			q.ch <- msg
		}
	}()
	go func() {
		for msg := range q.ch {
			_, err := dst.Send(msg)
			if err != nil {
				r.Close()
				return
			}
		}
	}()
	return q
}
