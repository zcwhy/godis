package sortedset

import (
	"math/bits"
	"math/rand"
)

const (
	maxLevel = 16
)

type Element struct {
	Member string
	Score  float64
}

type Level struct {
	forward *node // 前向指针，正向遍历
	span    int64 // span 用于rank命令
}

type node struct {
	Element
	backward *node // 后向指针，用于倒序遍历
	level    []*Level
}

type skipList struct {
	header *node
	tail   *node
	length int64
	level  int16
}

func makeNode(level int16, score float64, member string) *node {
	n := &node{
		Element: Element{
			Score:  score,
			Member: member,
		},
		level: make([]*Level, level),
	}
	for i := range n.level {
		n.level[i] = new(Level)
	}
	return n
}

// 位操作取随机的level
func randomLevel() int16 {
	total := uint64(1)<<uint64(maxLevel) - 1
	// 取rand数的低16位
	k := rand.Uint64() % total
	return maxLevel - int16(bits.Len64(k+1)) + 1
}

func (skipList *skipList) insert(member string, score float64) *node {
	update := make([]*node, maxLevel)

	n := skipList.header
	for i := skipList.level - 1; i >= 0; i-- {
		if n.level[i] != nil {
			for n.level[i].forward != nil &&
				(n.level[i].forward.Score < score ||
					(n.level[i].forward.Score == score && n.level[i].forward.Member < member)) {
				n = n.level[i].forward
			}
		}
	}
}

func (skipList *skipList) getByRank(rank int64) *node {
	var i int64 = 0
	n := skipList.header

	for l := skipList.level - 1; l >= 0; l-- {
		for n.level[l].forward != nil && (i+n.level[l].span < rank) {
			i += n.level[l].span
			n = n.level[l].forward
		}

		if i == rank {
			return n
		}
	}
	return nil
}
