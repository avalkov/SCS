package datastructures

type OrderedMap struct {
	data map[string]*Node
	list *DoublyLinkedList
	size int
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		data: make(map[string]*Node),
		list: NewDoublyLinkedList(),
		size: 0,
	}
}

func (om *OrderedMap) Add(key string, value interface{}) {
	if node, exists := om.data[key]; exists {
		node.value = value
	} else {
		node := om.list.Append(key, value)
		om.data[key] = node
		om.size++
	}
}

func (om *OrderedMap) Remove(key string) {
	if node, exists := om.data[key]; exists {
		om.list.Remove(node)
		delete(om.data, key)
		om.size--
	}
}

func (om *OrderedMap) Get(key string) (interface{}, bool) {
	if node, exists := om.data[key]; exists {
		return node.value, true
	}
	return nil, false
}

func (om *OrderedMap) GetAll() []KeyValue {
	items := make([]KeyValue, 0, om.size)
	for node := om.list.head; node != nil; node = node.next {
		items = append(items, KeyValue{node.key, node.value})
	}
	return items
}

type KeyValue struct {
	Key   string
	Value interface{}
}
