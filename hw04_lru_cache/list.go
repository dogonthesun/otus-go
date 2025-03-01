package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any

	Prev *ListItem
	Next *ListItem
}

type DoubleLinkedList struct {
	len         int
	front, back *ListItem
}

func NewList() *DoubleLinkedList {
	return &DoubleLinkedList{}
}

func (d *DoubleLinkedList) Len() int {
	return d.len
}

func (d *DoubleLinkedList) Front() *ListItem {
	return d.front
}

func (d *DoubleLinkedList) Back() *ListItem {
	return d.back
}

func (d *DoubleLinkedList) PushFront(v any) *ListItem {
	return d.pushFront(&ListItem{v, nil, nil})
}

func (d *DoubleLinkedList) PushBack(v any) *ListItem {
	return d.pushBack(&ListItem{v, nil, nil})
}

func (d *DoubleLinkedList) pushFront(i *ListItem) *ListItem {
	i.Prev, i.Next = nil, d.front
	if d.front != nil {
		d.front.Prev = i
	} else {
		d.back = i
	}
	d.front = i
	d.len++
	return i
}

func (d *DoubleLinkedList) pushBack(i *ListItem) *ListItem {
	i.Prev, i.Next = d.back, nil
	if d.back != nil {
		d.back.Next = i
	} else {
		d.front = i
	}
	d.back = i
	d.len++
	return i
}

func (d *DoubleLinkedList) Remove(i *ListItem) {
	_ = d.unlist(i)
}

func (d *DoubleLinkedList) MoveToFront(i *ListItem) {
	d.pushFront(d.unlist(i))
}

func (d *DoubleLinkedList) unlist(i *ListItem) *ListItem {
	prev, next := i.Prev, i.Next
	if prev != nil {
		prev.Next = i.Next
	} else {
		d.front = i.Next
	}
	if next != nil {
		next.Prev = i.Prev
	} else {
		d.back = i.Prev
	}
	i.Next, i.Prev = nil, nil
	d.len--
	return i
}
