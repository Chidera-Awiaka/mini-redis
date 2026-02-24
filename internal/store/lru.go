package store

type node struct {
	key  string
	prev *node
	next *node
}

type LRU struct {
	capacity int
	nodes    map[string]*node
	head     *node
	tail     *node
}

func NewLRU(cap int) *LRU {
	return &LRU{
		capacity: cap,
		nodes:    make(map[string]*node),
	}
}

// Touch moves an existing key to the front, or inserts it at the front if new.
func (l *LRU) Touch(key string) {
	if n, ok := l.nodes[key]; ok {
		l.remove(n)
		l.addToFront(n)
		return
	}

	n := &node{key: key}
	l.nodes[key] = n
	l.addToFront(n)
}

func (l *LRU) Remove(key string) {
	n, ok := l.nodes[key]
	if !ok {
		return
	}
	l.remove(n)
	delete(l.nodes, key)
}

func (l *LRU) RemoveOldest() string {
	if l.tail == nil {
		return ""
	}
	oldest := l.tail
	l.remove(oldest)
	delete(l.nodes, oldest.key)
	return oldest.key
}

func (l *LRU) Size() int {
	return len(l.nodes)
}

func (l *LRU) addToFront(n *node) {
	n.prev = nil
	n.next = l.head

	if l.head != nil {
		l.head.prev = n
	}
	l.head = n

	if l.tail == nil {
		l.tail = n
	}
}

func (l *LRU) remove(n *node) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		l.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		l.tail = n.prev
	}

	n.prev = nil
	n.next = nil
}
