package service

import "testing"

func TestEventPriorityQueue_PushPop(t *testing.T) {
    q := NewEventPriorityQueue(0)
    items := []*EventQueueItem{
        {Priority: 3},
        {Priority: 1},
        {Priority: 2},
    }
    for _, item := range items {
        q.Push(item)
    }
    wantOrder := []int{1, 2, 3}
    for i, want := range wantOrder {
        item := q.Pop()
        if item == nil {
            t.Fatalf("expected item with priority %d, got nil", want)
        }
        if item.Priority != want {
            t.Errorf("at index %d, got priority %d, want %d", i, item.Priority, want)
        }
    }
    if !q.IsEmpty() {
        t.Error("expected queue to be empty after popping all items")
    }
}

func TestEventPriorityQueue_MaxSize(t *testing.T) {
    q := NewEventPriorityQueue(2)
    q.Push(&EventQueueItem{Priority: 1})
    q.Push(&EventQueueItem{Priority: 2})
    q.Push(&EventQueueItem{Priority: 3}) // should be dropped
    if got := q.Size(); got != 2 {
        t.Errorf("expected size 2, got %d", got)
    }
}
