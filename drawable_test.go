package main

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// MockDrawable is a mock implementation of the Drawable interface for testing purposes.
type MockDrawable struct {
	zIndex int
	drawn  bool
}

func (m *MockDrawable) Draw(screen *ebiten.Image) {
	m.drawn = true
}

func (m *MockDrawable) ZIndex() int {
	return m.zIndex
}

func TestDrawHandlerAdd(t *testing.T) {
	handler := &DrawHandler{}

	// Add objects with different ZIndexes
	obj1 := &MockDrawable{zIndex: 1}
	obj2 := &MockDrawable{zIndex: 2}
	obj3 := &MockDrawable{zIndex: 3}

	handler.Add(obj1)
	handler.Add(obj2)
	handler.Add(obj3)

	// Check the order
	if len(handler.drawable) != 3 {
		t.Errorf("expected 3 objects, got %d", len(handler.drawable))
	}
	if handler.drawable[0] != obj1 || handler.drawable[1] != obj2 || handler.drawable[2] != obj3 {
		t.Errorf("objects are not sorted correctly by ZIndex")
	}
}

func TestDrawHandlerRemove(t *testing.T) {
	handler := &DrawHandler{}

	obj1 := &MockDrawable{zIndex: 1}
	obj2 := &MockDrawable{zIndex: 2}

	handler.Add(obj1)
	handler.Add(obj2)
	handler.Remove(obj1)

	// Check the remaining objects
	if len(handler.drawable) != 1 {
		t.Errorf("expected 1 object, got %d", len(handler.drawable))
	}
	if handler.drawable[0] != obj2 {
		t.Errorf("unexpected object in the list")
	}
}

func TestDrawHandlerDraw(t *testing.T) {
	handler := &DrawHandler{}
	screen := ebiten.NewImage(100, 100)

	obj1 := &MockDrawable{zIndex: 1}
	obj2 := &MockDrawable{zIndex: 2}

	handler.Add(obj1)
	handler.Add(obj2)

	handler.Draw(screen)

	// Check if both objects were drawn
	if !obj1.drawn || !obj2.drawn {
		t.Errorf("expected objects to be drawn")
	}
}

func TestDrawHandlerClear(t *testing.T) {
	handler := &DrawHandler{}

	obj1 := &MockDrawable{zIndex: 1}
	obj2 := &MockDrawable{zIndex: 2}

	handler.Add(obj1)
	handler.Add(obj2)
	handler.Clear()

	// Check if the list is empty
	if len(handler.drawable) != 0 {
		t.Errorf("expected 0 objects, got %d", len(handler.drawable))
	}
}
