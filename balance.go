package balance

import (
	"errors"
	"sync"
)

var (
	// ErrDuplicateID error is thrown when attempt to add an ID
	// which is already added to the balancer.
	ErrDuplicateID = errors.New("entry already added")
	// ErrIDNotFound is thrown when removing a non-existent ID.
	ErrIDNotFound = errors.New("id not found")
)

// Balance represents a smooth weighted round-robin load balancer.
type Balance struct {
	sync.RWMutex

	// items is the list of items to balance
	items []*Item
	// next is the index of the next item to use.
	next *Item
}

// NewBalance creates a new load balancer.
func NewBalance() *Balance {
	return &Balance{
		items: make([]*Item, 0),
	}
}

// Item represents the item in the list.
type Item struct {
	// id is the id of the item.
	id string
	// weight is the weight of the item that is given by the user.
	weight int
	// current is the current weight of the item.
	current int
}

func NewItem(id string, weight int) *Item {
	return &Item{
		id:      id,
		weight:  weight,
		current: 0,
	}
}

func (b *Balance) Add(id string, weight int) error {
	b.Lock()
	defer b.Unlock()
	for _, v := range b.items {
		if v.id == id {
			return ErrDuplicateID
		}
	}

	b.items = append(b.items, NewItem(id, weight))

	return nil
}

func (b *Balance) Get() string {
	b.Lock()
	defer b.Unlock()

	if len(b.items) == 0 {
		return ""
	}

	// Total weight of all items.
	var total int

	// Loop through the list of items and add the item's weight to the current weight.
	// Also increment the total weight counter.
	var max *Item
	for _, item := range b.items {
		item.current += item.weight
		total += item.weight

		// Select the item with max weight.
		if max == nil || item.current > max.current {
			max = item
		}
	}

	// Select the item with the max weight.
	b.next = max
	// Reduce the current weight of the selected item by the total weight.
	max.current -= total

	return max.id
}

// Remove deletes an item by ID from the balancer.
func (b *Balance) Remove(id string) error {
	b.Lock()
	defer b.Unlock()

	for i, item := range b.items {
		if item.id == id {
			b.items = append(b.items[:i], b.items[i+1:]...)
			return nil
		}
	}

	return ErrIDNotFound
}

// ItemIDs returns a list of all item IDs in the balancer.
func (b *Balance) ItemIDs() []string {
	b.RLock()
	defer b.RUnlock()

	ids := make([]string, len(b.items))
	for i, item := range b.items {
		ids[i] = item.id
	}
	return ids
}
