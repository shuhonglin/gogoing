package structure

import "reflect"

type Node struct {
	next, prev *Node

	Value interface{}
}

func (node *Node) Next() *Node {
	return node.next
}

func (node *Node) Prev() *Node {
	return node.prev
}


type DoubleLinkList struct {
	first, last *Node
	len int
}

func (dlist *DoubleLinkList) Init() *DoubleLinkList {
	dlist.len = 0
	return dlist
}

func NewDoubleLinkList() *DoubleLinkList {
	return new(DoubleLinkList).Init()
}

// insert node before at, increments dlist.len, and returns node
func (dlist *DoubleLinkList) linkBefore(node, at *Node) *Node {
	pred := at.prev
	at.prev = node
	node.next = at
	node.prev = pred
	if pred==nil {
		dlist.first = node
	} else {
		pred.next = node
	}
	dlist.len++
	return node
}

// add node to head
func (dlist *DoubleLinkList) linkFirst(node *Node) *Node {
	first := dlist.first
	node.prev = nil
	node.next = first
	dlist.first = node
	if first==nil {
		dlist.last = node
	} else {
		first.prev = node
	}
	dlist.len++
	return node
}

// add node to tail
func (dlist *DoubleLinkList) linkLast(node *Node) *Node  {
	last := dlist.last
	node.prev = last
	node.next = nil
	dlist.last = node
	if last==nil {
		dlist.first = node
	} else {
		last.next=node
	}
	dlist.len++
	return node
}

// remove first node
func (dlist *DoubleLinkList) unlinkFirst() *Node {
	if dlist.first == nil {
		return nil
	}
	first := dlist.first
	next := first.next
	dlist.first = next
	if next == nil {
		 dlist.last = nil
	} else {
		first.next = nil // help gc
		next.prev = nil
	}
	dlist.len--
	return first
}


// remove last node
func (dlist *DoubleLinkList) unlinkLast() *Node {
	if dlist.last == nil {
		return nil
	}
	last := dlist.last
	prev := last.prev
	dlist.last = prev
	if prev==nil {
		dlist.first = nil
	} else {
		last.prev = nil // help gc
		prev.next = nil
	}
	dlist.len--
	return last
}

// remove node
func (dlist *DoubleLinkList) unlink(node *Node) *Node {
	next:=node.next
	prev:=node.prev
	if prev==nil {
		dlist.first = next
	} else {
		prev.next = next
		node.prev = nil
	}

	if next==nil {
		dlist.last = prev
	} else {
		next.prev = prev
		node.next = nil
	}
	dlist.len--
	return node
}

func (dlist *DoubleLinkList) getFirst() *Node {
	return dlist.first
}

func (dlist *DoubleLinkList) getLast() *Node {
	return dlist.last
}

// swap with next node(move current node to next node)
func (dlist *DoubleLinkList) moveNext(node *Node) bool {
	prev := node.prev
	next := node.next
	if next==nil {
		return false
	}

	if prev==nil {
		node.next = next.next
		next.next.prev = node
		next.next = node
		next.prev = nil
		node.prev = next
		dlist.first = next
		if reflect.DeepEqual(dlist.last, next) == true {
			node.next= nil
			dlist.last = node
		}
		return true
	}
	prev.next = next
	if next.next != nil {
		next.next.prev = node
	}
	node.next = next.next
	node.prev = next

	next.prev = prev
	next.next = node

	if reflect.DeepEqual(dlist.last, next)==true {
		node.next = nil
		dlist.last = node
	}
	return true
}

func (dlist *DoubleLinkList) moveToNext(node,newNextNode *Node) {
	prev := node.prev
	next := node.next

	node.next = newNextNode
	node.prev = newNextNode.prev
	newNextNode.prev.next = node
	newNextNode.prev = node

	next.prev = prev
	if prev==nil {
		dlist.first = next
	} else {
		prev.next = next
	}

}

// swap tow node. note that nodeA must before nodeB
func (dlist *DoubleLinkList) swapNode(nodeA, nodeB *Node) bool {
	if reflect.DeepEqual(nodeA, nodeB) == true {
		return false
	}
	if reflect.DeepEqual(nodeA.next, nodeB) == true { // brother node
		return dlist.moveNext(nodeA)
	}

	if reflect.DeepEqual(nodeA, dlist.first) == true {
		if reflect.DeepEqual(nodeB, dlist.last) == true {  // swap last node for first node
			nodeBPre := nodeB.prev

			nodeA.next.prev = nodeB
			nodeB.prev = nil
			nodeB.next = nodeA.next
			dlist.first = nodeB

			nodeBPre.next = nodeA
			nodeA.prev = nodeBPre
			nodeA.next = nil
			dlist.last = nodeA
		} else { // swap common node for first node
			nodeANext := nodeA.next

			nodeB.prev.next = nodeA
			nodeB.next.prev = nodeA
			nodeA.prev = nodeB.prev
			nodeA.next = nodeB.next

			nodeANext.prev = nodeB
			nodeB.next = nodeANext
			nodeB.prev = nil
			dlist.first = nodeB
		}
	} else {
		if reflect.DeepEqual(nodeB, dlist.last) == true {  // swap last node for commom node
			nodeBPre := nodeB.prev

			nodeA.prev.next = nodeB
			nodeA.next.prev = nodeB
			nodeB.prev = nodeA.prev
			nodeB.next = nodeA.next

			nodeBPre.next = nodeA
			nodeA.prev = nodeBPre
			nodeA.next = nil
			dlist.last = nodeA


		} else { // swap common node for commom node
			nodeBPre := nodeB.prev
			nodeBNext := nodeB.next

			nodeA.prev.next = nodeB
			nodeA.next.prev = nodeB
			nodeB.prev = nodeA.prev
			nodeB.next = nodeA.next

			nodeBPre.next = nodeA
			nodeBNext.prev = nodeA
			nodeA.prev = nodeBPre
			nodeA.next = nodeBNext

		}
	}
	return true
}

func (dlist *DoubleLinkList) movePrev(node *Node) bool {
	prev := node.prev
	next := node.next
	if prev == nil {
		return false
	}
	if next == nil {
		node.prev = prev.prev
		node.prev.next = node
		prev.prev = node
		prev.next = nil
		node.next = prev
		dlist.last = prev
		if reflect.DeepEqual(dlist.first, prev) == true {
			node.prev = nil
			dlist.first = prev
		}
		return true
	}
	next.prev = prev
	if prev.prev != nil {
		prev.prev.next = node
	}
	node.next = prev
	node.prev = prev.prev
	prev.next = next
	prev.prev = node
	if reflect.DeepEqual(dlist.first, prev) == true {
		node.prev = nil
		dlist.first = node
	}
	return true
}

func (dlist *DoubleLinkList) moveToPrev(node, newPrevNode *Node) {
	prev := node.prev
	next := node.next
	node.prev = newPrevNode
	node.next = newPrevNode.next
	newPrevNode.next.prev = node
	newPrevNode.next = node

	prev.next = next

	if next==nil {
		dlist.last = prev
	} else {
		next.prev = prev
	}
}

func (dlist *DoubleLinkList) moveToFirst(node *Node)  {
	prev := node.prev
	next := node.next
	node.next = dlist.first
	dlist.first.prev = node
	node.prev = nil

	prev.next = next
	if next != nil {
		next.prev = prev
	} else {
		dlist.last = prev
	}
	dlist.first = node
}

func (dlist *DoubleLinkList) LinkFirst(node *Node) *Node {
	return dlist.linkFirst(node)
}

func (dlist *DoubleLinkList) Len() int  {
	return dlist.len
}

func (dlist *DoubleLinkList) GetFirst() *Node {
	return dlist.getFirst()
}

func (dlist *DoubleLinkList) GetLast() *Node {
	return dlist.getLast()
}