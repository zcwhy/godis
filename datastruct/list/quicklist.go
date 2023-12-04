package list

import "container/list"

const pageSize = 1024

// QuickList page链表，数据存放在page中
type QuickList struct {
	data *list.List
	size int
}

// 标识一个val的具体信息（page，off）的二元组
type posInfo struct {
	page   *list.Element
	offset int
}

func NewQuickList() *QuickList {
	return &QuickList{
		data: list.New(),
	}
}

func (ql *QuickList) Add(val interface{}) {
	ql.size++
	if ql.data.Len() == 0 {
		// empty list新分配一个page大小的slice
		page := make([]interface{}, 0, pageSize)
		page = append(page, val)
		ql.data.PushBack(page)

		return
	}

	lastNode := ql.data.Back()
	lastPage := lastNode.Value.([]interface{})
	if len(lastPage) == cap(lastPage) {
		// full
		page := make([]interface{}, 0, pageSize)
		page = append(page, val)
		ql.data.PushBack(page)

		return
	}

	lastPage = append(lastPage, val)
	lastNode.Value = lastPage
}

func (ql *QuickList) find(index int) *posInfo {
	if ql == nil {
		panic("list is nil")
	}
	if index < 0 || index >= ql.size {
		panic("index out of bound")
	}

	pageIdx := index / pageSize
	pageOff := index % pageSize

	target := ql.data.Front()
	for i := 0; i < pageIdx; i++ {
		target = target.Next()
	}

	return &posInfo{
		page:   target,
		offset: pageOff,
	}
}

func (p *posInfo) getPage() []interface{} {
	return p.page.Value.([]interface{})
}

func (ql *QuickList) Get(index int) (val interface{}) {
	n := ql.find(index)
	return n.getPage()[n.offset]
}

func (ql *QuickList) Set(index int, val interface{}) {
	n := ql.find(index)
	n.getPage()[n.offset] = val
}

func (ql *QuickList) Insert(index int, val interface{}) {
	if index == ql.size { // insert at
		ql.Add(val)
		return
	}
	pos := ql.find(index)
	page := pos.getPage()

	if len(page) < pageSize {
		// not full
		page = append(page[:pos.offset+1], page[pos.offset:]...)
		page[pos.offset] = val
		pos.page.Value = page
		ql.size++
		return
	}
}

func (ql *QuickList) Remove(index int) (val interface{}) {
	return nil
}

func (ql *QuickList) RemoveLast() (val interface{}) {
	return nil
}
