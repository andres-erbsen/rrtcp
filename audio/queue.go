package audio

import (
	"sync"
)

const messageDoesntExist = "NONE"

type Queue struct {
	sync.Mutex
  messages    map[int]interface{}
  lastPulled  int
}

func NewQueue() *Queue {
	return &Queue{lastPulled: -1, messages: make(map[int]interface{})}
}

// someone is asking for a message with an index, return it if
// we have it and update lastPulled
// Pull should be called in increasing index
func (q *Queue) Pull(index int) interface{} {
	q.Lock()
	defer q.Unlock()
  q.lastPulled = index
  message, exists := q.messages[index]
  if !exists {
    return messageDoesntExist
  } else {
    delete(q.messages, index)
    return message
  }
}

// put a message with an index into the queue
// if the index is < the last thing pulled, ignore it
func (q *Queue) Push(index int, message interface{}) {
  q.Lock()
	defer q.Unlock()
  // make sure it isn't a stale message
  if q.lastPulled < index {
    q.messages[index] = message
  }
}
