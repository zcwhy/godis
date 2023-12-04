package set

import "godis/datastruct/dict"

// Set val为nil的hash table
type Set struct {
	dict dict.Dict
}

func Make(members ...string) *Set {
	set := &Set{
		dict: dict.MakeSimple(),
	}
	for _, member := range members {
		set.Add(member)
	}
	return set
}

func (set *Set) Add(val string) int {
	return set.dict.Put(val, nil)
}

func (set *Set) Remove(val string) int {
	_, ret := set.dict.Remove(val)
	return ret
}

func (set *Set) Has(val string) bool {
	if set == nil || set.dict == nil {
		return false
	}
	_, exists := set.dict.Get(val)
	return exists
}

func (set *Set) Len() int {
	if set == nil || set.dict == nil {
		return 0
	}
	return set.dict.Len()
}

func (set *Set) ToSlice() []string {
	slice := make([]string, set.Len())
	i := 0
	set.dict.ForEach(func(key string, val interface{}) bool {
		if i < len(slice) {
			slice[i] = key
		} else {
			// 在进行for-each遍历的过程中，set添加新的元素，append在slice后
			slice = append(slice, key)
		}
		i++
		return true
	})
	return slice
}

func (set *Set) ForEach(consumer func(member string) bool) {
	if set == nil || set.dict == nil {
		return
	}
	set.dict.ForEach(func(key string, val interface{}) bool {
		return consumer(key)
	})
}

func (set *Set) ShallowCopy() *Set {
	result := Make()
	set.ForEach(func(member string) bool {
		result.Add(member)
		return true
	})
	return result
}

func (set *Set) RandomMembers(limit int) []string {
	if set == nil || set.dict == nil {
		return nil
	}
	return set.dict.RandomKeys(limit)
}

func (set *Set) RandomDistinctMembers(limit int) []string {
	return set.dict.RandomDistinctKeys(limit)
}
