package audio

import (
	"sync"
)

type Queue struct {
	sync.Mutex
  messages    map[int]string
	subscriber  chan<- interface{}
  lastPulled  int
}

func NewQueue(subscriber chan<- interface{}) *Queue {
	return &Queue{subscriber: subscriber, lastPulled: -1, messages: make(map[int]string)}
}

// someone is asking for a message with an index, return it if
// we have it and update lastPulled
func (q *Queue) Pull(index int) {
	q.Lock()
	defer q.Unlock()
  message, exists := q.messages[index]
  if !exists{
    //todo
  } else {
    q.subscriber <- message
    delete(q.messages, index)
  }
  q.lastPulled = index
}

// put a message with an index into the queue
// if the index is >= the last thing pulled, ignore it
func (q *Queue) Put(index int, message string) {
  q.Lock()
	defer q.Unlock()
  // make sure it isn't a stale message
  if q.lastPulled < index {
    q.messages[index] = message
  }
}
