package list

type Expected func(a interface{}) bool

type List interface {
	Add(val interface{})
	Get(index int) (val interface{})
	Set(index int, val interface{})
	Insert(index int, val interface{})
	Remove(index int) (val interface{})
	RemoveLast() (val interface{})
	RemoveAllByVal(expected Expected) int
	RemoveByVal(expected Expected, count int) int
	ReverseRemoveByVal(expected Expected, count int) int
	Len() int
	//ForEach(consumer Consumer)
	Contains(expected Expected) bool
	Range(start int, stop int) []interface{}
}
