package lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return &list{0, nil, nil}
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newFront := ListItem{v, l.front, nil}

	return l.pushFrontItem(&newFront)
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBack := ListItem{v, nil, l.back}

	return l.pushBackItem(&newBack)
}

func (l *list) Remove(i *ListItem) {
	if l.Len() == 0 {
		return
	}

	if l.Len() == 1 {
		// если переданный элемент не совпадает, то ничего не делаем
		if *(l.front) == *(i) {
			newEmptyList := NewList().(*list)
			*l = *newEmptyList
		}

		return
	}

	switch {
	case i.Next == nil: // переданный элемент последний
		i.Prev.Next = nil
		l.back = i.Prev
	case i.Prev == nil: // переданный элемент первый
		i.Next.Prev = nil
		l.front = i.Next
	default: // переданный элемент посередине
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)

	*i = *l.pushFrontItem(i)
}

func (l *list) pushFrontItem(newFront *ListItem) *ListItem {
	if l.len == 0 {
		return l.initList(newFront.Value)
	}

	newFront.Prev = nil
	newFront.Next = l.front

	l.front.Prev = newFront
	l.front = newFront

	l.len++

	return l.Front()
}

func (l *list) pushBackItem(newBack *ListItem) *ListItem {
	if l.len == 0 {
		return l.initList(newBack.Value)
	}

	newBack.Prev = l.back
	newBack.Next = nil

	l.back.Next = newBack
	l.back = newBack

	l.len++

	return l.Back()
}

func (l *list) initList(v interface{}) *ListItem {
	firstValue := &ListItem{v, nil, nil}
	l.front = firstValue
	l.back = firstValue
	l.len++

	return firstValue
}
