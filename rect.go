package gform

import (
	"github.com/darkautism/w32"
)

type Rect struct {
	rect w32.RECT
}

func NewEmptyRect() *Rect {
	var newRect Rect
	w32.SetRectEmpty(&newRect.rect)

	return &newRect
}

func NewRect(left, top, right, bottom int) *Rect {
	var newRect Rect
	w32.SetRectEmpty(&newRect.rect)
	newRect.Set(left, top, right, bottom)

	return &newRect
}

func NewRect32(left, top, right, bottom int32) *Rect {
	var newRect Rect
	w32.SetRectEmpty(&newRect.rect)
	newRect.Set32(left, top, right, bottom)

	return &newRect
}

func (this *Rect) Data() (left, top, right, bottom int) {
	left = int(this.rect.Left)
	top = int(this.rect.Top)
	right = int(this.rect.Right)
	bottom = int(this.rect.Bottom)
	return
}

func (this *Rect) GetW32Rect() *w32.RECT {
	return &this.rect
}

func (this *Rect) Set(left, top, right, bottom int) {
	w32.SetRect(&this.rect, int32(left), int32(top), int32(right), int32(bottom))
}

func (this *Rect) Set32(left, top, right, bottom int32) {
	w32.SetRect(&this.rect, left, top, right, bottom)
}

func (this *Rect) IsEqual(rect *Rect) bool {
	return w32.EqualRect(&this.rect, &rect.rect)
}

func (this *Rect) Inflate(x, y int) {
	w32.InflateRect(&this.rect, int32(x), int32(y))
}

func (this *Rect) Intersect(src *Rect) {
	w32.IntersectRect(&this.rect, &this.rect, &src.rect)
}

func (this *Rect) IsEmpty() bool {
	return w32.IsRectEmpty(&this.rect)
}

func (this *Rect) Offset(x, y int) {
	w32.OffsetRect(&this.rect, int32(x), int32(y))
}

func (this *Rect) IsPointIn(x, y int) bool {
	return w32.PtInRect(&this.rect, int32(x), int32(y))
}

func (this *Rect) Substract(src *Rect) {
	w32.SubtractRect(&this.rect, &this.rect, &src.rect)
}

func (this *Rect) Union(src *Rect) {
	w32.UnionRect(&this.rect, &this.rect, &src.rect)
}
