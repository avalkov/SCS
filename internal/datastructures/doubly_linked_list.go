package datastructures

type Node struct {
	key   string
	value interface{}
	prev  *Node
	next  *Node
}

type DoublyLinkedList struct {
	head *Node
	tail *Node
}

func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{}
}

func (dll *DoublyLinkedList) Append(key string, value interface{}) *Node {
	node := &Node{key: key, value: value}
	if dll.tail == nil {
		dll.head = node
		dll.tail = node
	} else {
		node.prev = dll.tail
		dll.tail.next = node
		dll.tail = node
	}
	return node
}

func (dll *DoublyLinkedList) Remove(node *Node) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		dll.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		dll.tail = node.prev
	}
}
