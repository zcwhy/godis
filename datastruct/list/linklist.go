package list

type LinkedList struct {
	first *node
	last  *node
	size  int
}

type node struct {
	val  interface{}
	prev *node
	next *node
}

func (ll *LinkedList) Add(val interface{}) {
	if ll == nil {
		panic("list is nil")
	}

	n := &node{
		val: val,
	}
	if ll.last == nil {
		// empty list
		ll.first = n
		ll.last = n
	} else {
		n.prev = ll.last
		ll.last.next = n
		ll.last = n
	}
	ll.size++
}

func (ll *LinkedList) find(index int) (n *node) {
	if index < ll.size/2 {
		// 从前往后找
		n := ll.first
		for i := 0; i < index; i++ {
			n = n.next
		}
	} else {
		n := ll.last
		for i := ll.size - 1; i > index; i-- {
			n = n.prev
		}
	}
	return n
}

func (ll *LinkedList) Get(index int) (val interface{}) {
	if ll == nil {
		panic("list is nil")
	}
	if index < 0 || index >= ll.size {
		panic("index out of bound")
	}
	return ll.find(index)
}

// Set ll[index] = val, idx should between [0, size]
func (ll *LinkedList) Set(index int, val interface{}) {
	if ll == nil {
		panic("list is nil")
	}
	if index < 0 || index > ll.size {
		panic("index out of bound")
	}
	n := ll.find(index)
	n.val = val
}

func (ll *LinkedList) Insert(index int, val interface{}) {
	if ll == nil {
		panic("list is nil")
	}
	if index < 0 || index > ll.size {
		panic("index out of bound")
	}

	if index == ll.size {
		ll.Add(val)
		return
	}

	pivot := ll.find(index)
	n := &node{
		val:  val,
		prev: pivot.prev,
		next: pivot,
	}
	if pivot.prev == nil {
		ll.first = n
	} else {
		pivot.prev.next = n
	}
	pivot.prev = n
	ll.size++
}

func (ll *LinkedList) removeNode(n *node) {
	if n.prev == nil {
		ll.first = n.next
	} else {
		n.prev.next = n.next
	}
	if n.next == nil {
		ll.last = n.prev
	} else {
		n.next.prev = n.prev
	}

	// for gc
	n.prev = nil
	n.next = nil

	ll.size--
}

func (ll *LinkedList) Remove(index int) (val interface{}) {
	if ll == nil {
		panic("list is nil")
	}
	if index < 0 || index >= ll.size {
		panic("index out of bound")
	}

	n := ll.find(index)
	ll.removeNode(n)
	return n.val
}

func (ll *LinkedList) Len() int {
	if ll == nil {
		panic("list is nil")
	}
	return ll.size
}

func (ll *LinkedList) Range(start int, stop int) []interface{} {
	if ll == nil {
		panic("list is nil")
	}
	if start < 0 || start >= ll.size {
		panic("`start` out of range")
	}
	if stop < start || stop > ll.size {
		panic("`stop` out of range")
	}

	sliceSize := stop - start
	slice := make([]interface{}, sliceSize)
	n := ll.first
	i := 0
	sliceIndex := 0
	for n != nil {
		if i >= start && i < stop {
			slice[sliceIndex] = n.val
			sliceIndex++
		} else if i >= stop {
			break
		}
		i++
		n = n.next
	}
	return slice
}

func Make(vals ...interface{}) *LinkedList {
	list := LinkedList{}
	for _, v := range vals {
		list.Add(v)
	}
	return &list
}
