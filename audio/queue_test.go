package audio

import (
	"testing"
)

func Test(t *testing.T) {
  ch := make(chan interface{}, 1)
	queue := NewQueue(ch)
  queue.Put(0, "abc")
  queue.Pull(0)
}
