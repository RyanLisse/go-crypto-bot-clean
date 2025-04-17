package service

import (
    "sort"

    "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// ActionType defines the type of queue action
// ActionCreateCoin queues creation of a new coin
// ActionUpdateCoin queues a status update for an existing coin
type ActionType int

const (
    ActionCreateCoin ActionType = iota + 1
    ActionUpdateCoin
)

// EventQueueItem holds the information for an event to be processed
type EventQueueItem struct {
    Event    *model.NewCoinEvent
    Coin     *model.NewCoin
    Priority int
    Action   ActionType
}

// EventPriorityQueue implements a simple priority queue for EventQueueItem
// Items are sorted by Priority ascending (lower value = higher priority)
type EventPriorityQueue struct {
    items   []*EventQueueItem
    maxSize int
}

// NewEventPriorityQueue creates a new EventPriorityQueue with the specified maximum size
// If maxSize is 0 or negative, the queue has unlimited capacity
func NewEventPriorityQueue(maxSize int) *EventPriorityQueue {
    return &EventPriorityQueue{
        items:   make([]*EventQueueItem, 0),
        maxSize: maxSize,
    }
}

// Push adds a new item to the priority queue. If the queue has reached its maxSize, the item is dropped
func (q *EventPriorityQueue) Push(item *EventQueueItem) {
    if q.maxSize > 0 && len(q.items) >= q.maxSize {
        return
    }
    q.items = append(q.items, item)
    sort.SliceStable(q.items, func(i, j int) bool {
        return q.items[i].Priority < q.items[j].Priority
    })
}

// Pop removes and returns the highest-priority item (lowest Priority value). Returns nil if the queue is empty.
func (q *EventPriorityQueue) Pop() *EventQueueItem {
    if len(q.items) == 0 {
        return nil
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item
}

// IsEmpty returns true if the queue contains no items
func (q *EventPriorityQueue) IsEmpty() bool {
    return len(q.items) == 0
}

// Size returns the current number of items in the queue
func (q *EventPriorityQueue) Size() int {
    return len(q.items)
}
