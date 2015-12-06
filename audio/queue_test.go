package audio

import (
	"testing"
)

func TestInOrder(t *testing.T) {
	queue := NewQueue()
  queue.Push(2, "a")
  queue.Push(3, "b")
  queue.Push(4, "c")
  if queue.Pull(2) != "a" {
    t.Fail()
  }
  if queue.Pull(3) != "b" {
    t.Fail()
  }
  if queue.Pull(4) != "c" {
    t.Fail()
  }
}

func TestOutOfOrder(t *testing.T) {
	queue := NewQueue()
  queue.Push(7, "f")
  queue.Push(6, "e")
  queue.Push(5, "d")
  if queue.Pull(5) != "d" {
    t.Fail()
  }
  if queue.Pull(6) != "e" {
    t.Fail()
  }
  if queue.Pull(7) != "f" {
    t.Fail()
  }
}

func TestNotInQueue(t *testing.T) {
	queue := NewQueue()
  if queue.Pull(0) != "NONE" {
    t.Fail()
  }
  if queue.Pull(1) != "NONE" {
    t.Fail()
  }
}
