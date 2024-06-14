package main

import (
	"sort"
)

type Clickable interface {
	OnClick()
	IsClicked(x, y int) bool
	ZIndex() int
}

type OnClickHandler struct {
	// List of clickable objects
	clickableObjects []Clickable
}

func (o *OnClickHandler) Add(obj Clickable) {
	o.clickableObjects = append(o.clickableObjects, obj)

	// ZIndex でソートしておく
	// [0] にもっとも ZIndex が大きいオブジェクトがくるようにする
	sort.Slice(o.clickableObjects, func(i, j int) bool {
		return o.clickableObjects[i].ZIndex() > o.clickableObjects[j].ZIndex()
	})
}

func (o *OnClickHandler) Remove(obj Clickable) {
	for i, v := range o.clickableObjects {
		if v == obj {
			o.clickableObjects = append(o.clickableObjects[:i], o.clickableObjects[i+1:]...)
			return
		}
	}
}

func (o *OnClickHandler) HandleClick(x, y int) {
	for _, obj := range o.clickableObjects {
		if obj.IsClicked(x, y) {
			obj.OnClick()

			// 一つクリックされたら終了
			// オブジェクトが重なっている場合、一番手前のオブジェクトのみがクリックされるという仕様にしておく
			return
		}
	}
}
