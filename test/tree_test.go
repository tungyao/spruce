package test

import (
	"fmt"
	"testing"
)

func Reset(trx *TreeNode) {
	if trx.Color == 1 && trx.Left != nil && trx.Left.Color == 1 { // 去右旋
		// 看有没有 父父节点有没有 uncle元素
		fmt.Println(trx)
		fmt.Println(trx.Parent)
		fmt.Println(trx.Left)
		fmt.Println(trx.Right)
		// 判断有没有uncle节点
		if trx.Parent.Right == nil { // 没有直接旋转+变色
			t := *trx.Parent
			*trx.Parent = *trx
			*trx = t
		}
	}
}
func TestTree(t *testing.T) {
	tr := initTree()

	tr.insertRBTree(tr.Node, 123, "123", nil)
	tr.insertRBTree(tr.Node, 10, "2", nil)
	tr.insertRBTree(tr.Node, 8, "a", nil)
	t.Log(maxDeep(tr.Node))
	// tr.insertRBTree(3, "3", nil)
	// tr.insertRBTree(426, "516", nil)
	// tr.insertRBTree(4, "4", nil)
	// t.Log(maxDeep(tr.Node))
	// fmt.Println(searchNode(tr, 426))
	// invertTree(tr)
}

// 升级结构到二叉树
// 搜索二叉树
type Tree struct {
	Init bool
	Node *TreeNode
}
type TreeNode struct {
	Data     map[int]string
	Left     *TreeNode
	Right    *TreeNode
	Position int
	Parent   *TreeNode // 不要这个节点 是真的找不到 父节点
	Color    int       // 1是红色
}

func initTree() *Tree {
	// 取中间值 100 是最大插槽 这里手动 100 个插槽
	return &Tree{
		Init: false,
		Node: nil,
	}
}

// 应该在插入的时候只能是红色
func (tr *Tree) insertRBTree(trx *TreeNode, position int, value string, parent *TreeNode) *TreeNode {
	if !tr.Init {
		tr.Init = true
		tr.Node = &TreeNode{
			Data:     map[int]string{position: value},
			Left:     nil,
			Right:    nil,
			Position: position,
			Parent:   nil,
			Color:    0,
		}
		return tr.Node
	}
	// 前期插入过程不变
	// 重点在插入过后的颜色调整
	// 两个红色不能相连!! 调整也是这 左右旋也是根据在 那一片叶子上
	if trx == nil {
		return &TreeNode{
			Data:     map[int]string{position: value},
			Left:     nil,
			Right:    nil,
			Position: position,
			Parent:   parent,
			Color:    1,
		}
	}
	if position < trx.Position {
		trx.Left = tr.insertRBTree(trx.Left, position, value, trx)
	} else {
		trx.Right = tr.insertRBTree(trx.Right, position, value, trx)
	}
	Reset(trx)
	return trx
}

//
func insertNode(tr *TreeNode, position int, value string) *TreeNode {
	if tr == nil {
		return &TreeNode{
			Data:     map[int]string{position: value},
			Left:     nil,
			Right:    nil,
			Position: position,
		}
	}
	// 比position小放到左边
	if position < tr.Position {
		tr.Left = insertNode(tr.Left, position, value)
		// 比position大放到右边
	} else {
		tr.Right = insertNode(tr.Right, position, value)
	}
	Reset(tr)
	return tr
}

// 妈的递归查找真的好使
func searchNode(tr *TreeNode, position int) map[int]string {
	if tr.Position == position {
		return tr.Data
	}
	if position < tr.Position {
		return searchNode(tr.Left, position)
		// 比position大放到右边
	} else {
		return searchNode(tr.Right, position)
	}
}

// 最大深度
func maxDeep(tr *TreeNode) int {
	if tr == nil {
		return 0
	}
	var maxDegree = 0
	left := maxDeep(tr.Left)
	right := maxDeep(tr.Right)
	if left > right {
		maxDegree = left
	} else {
		maxDegree = right
	}
	return maxDegree + 1
}

// 最大宽度 层级遍历
func maxWidth(tr *TreeNode) int {
	return 0
}

// 反转二叉树 反转的是 节点下的叶子
func invertTree(tr *TreeNode) *TreeNode {
	if tr == nil {
		return nil
	}
	tr.Left, tr.Right = invertTree(tr.Right), invertTree(tr.Left)
	return tr
}
