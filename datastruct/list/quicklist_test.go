package list

import "testing"

func TestQuickList_Add(t *testing.T) {
	list := NewQuickList()
	for i := 0; i < pageSize*10; i++ {
		list.Add(i)
	}

	for i := 0; i < pageSize*10; i++ {
		v := list.Get(i).(int)
		if v != i {
			t.Errorf("wrong value at: %d", i)
		}
	}
}

func TestQuickList_Set(t *testing.T) {
	list := NewQuickList()
	for i := 0; i < pageSize*10; i++ {
		list.Add(i)
	}
	for i := 0; i < pageSize*10; i++ {
		list.Set(i, 2*i)
	}
	for i := 0; i < pageSize*10; i++ {
		v := list.Get(i).(int)
		if v != 2*i {
			t.Errorf("wrong value at: %d", i)
		}
	}
}
