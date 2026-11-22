package termtest

import (
	"io"
	"sync"
)

// replyQueue hands the emulator's replies to their writer on its own
// goroutine, so a program that never reads its input can't block the screen.
type replyQueue struct {
	mu      sync.Mutex
	ready   *sync.Cond
	pending [][]byte
	writer  io.Writer
	closed  bool
	done    chan struct{}
}

func newReplyQueue() *replyQueue {
	q := &replyQueue{writer: io.Discard, done: make(chan struct{})}
	q.ready = sync.NewCond(&q.mu)
	go q.run()
	return q
}

func (q *replyQueue) setWriter(w io.Writer) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.writer = w
}

func (q *replyQueue) push(reply []byte) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.pending = append(q.pending, append([]byte(nil), reply...))
	q.ready.Signal()
}

func (q *replyQueue) close() {
	q.mu.Lock()
	q.closed = true
	q.ready.Signal()
	q.mu.Unlock()
	<-q.done
}

func (q *replyQueue) run() {
	defer close(q.done)
	for {
		q.mu.Lock()
		for len(q.pending) == 0 && !q.closed {
			q.ready.Wait()
		}
		if q.closed {
			q.mu.Unlock()
			return
		}
		reply, w := q.pending[0], q.writer
		q.pending = q.pending[1:]
		q.mu.Unlock()
		_, _ = w.Write(reply)
	}
}
