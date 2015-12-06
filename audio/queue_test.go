package audio

import (
	"testing"
  "fmt"
)

func Test(t *testing.T) {
	queue := NewQueue()
  fmt.Println("Push nothing")
  if queue.Pull(1) != "NONE" {
    t.Fail()
  }

  fmt.Println("Push in order")
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

  fmt.Println("Push out of order")
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
