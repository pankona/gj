package main

import (
	"testing"
)

type MockClickable struct {
	clicked       bool
	x, y          int
	width, height int
	zIndex        int
}

func (m *MockClickable) OnClick(x, y int) bool {
	m.clicked = true
	return false
}

func (m *MockClickable) IsClicked(x, y int) bool {
	return x >= m.x && x <= m.x+m.width && y >= m.y && y <= m.y+m.height
}

func (m *MockClickable) ZIndex() int {
	return m.zIndex
}

func TestOnClickHandlerAdd(t *testing.T) {
	handler := &OnClickHandler{}
	obj1 := &MockClickable{zIndex: 1}
	obj2 := &MockClickable{zIndex: 2}

	handler.Add(obj1)
	handler.Add(obj2)

	if len(handler.clickableObjects) != 2 {
		t.Errorf("Expected 2 objects, got %d", len(handler.clickableObjects))
	}

	if handler.clickableObjects[0] != obj2 {
		t.Errorf("Expected obj2 to be first due to higher ZIndex, but got obj1")
	}
}

func TestOnClickHandlerRemove(t *testing.T) {
	handler := &OnClickHandler{}
	obj1 := &MockClickable{zIndex: 1}
	obj2 := &MockClickable{zIndex: 2}

	handler.Add(obj1)
	handler.Add(obj2)
	handler.Remove(obj1)

	if len(handler.clickableObjects) != 1 {
		t.Errorf("Expected 1 object, got %d", len(handler.clickableObjects))
	}

	if handler.clickableObjects[0] != obj2 {
		t.Errorf("Expected obj2 to be remaining, but got obj1")
	}
}

func TestOnClickHandlerHandleClick(t *testing.T) {
	handler := &OnClickHandler{}
	obj1 := &MockClickable{x: 0, y: 0, width: 10, height: 10, zIndex: 1}
	obj2 := &MockClickable{x: 5, y: 5, width: 10, height: 10, zIndex: 2}

	handler.Add(obj1)
	handler.Add(obj2)

	handler.HandleClick(6, 6)

	if !obj2.clicked {
		t.Errorf("obj2 should not be clicked due to lower ZIndex")
	}

	if obj1.clicked {
		t.Errorf("obj1 should be clicked due to higher ZIndex")
	}
}

func TestOnClickHandlerHandleClickNoClick(t *testing.T) {
	handler := &OnClickHandler{}
	obj1 := &MockClickable{x: 0, y: 0, width: 10, height: 10, zIndex: 1}

	handler.Add(obj1)

	handler.HandleClick(20, 20)

	if obj1.clicked {
		t.Errorf("obj1 should not be clicked because click was outside the bounds")
	}
}
