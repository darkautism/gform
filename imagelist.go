package gform

import (
	"github.com/darkautism/w32"
)

type ImageList struct {
	handle w32.HIMAGELIST
}

func NewImageList(cx, cy int, flags uint, cInitial, cGrow int) *ImageList {
	imgl := new(ImageList)
	imgl.handle = w32.ImageList_Create(int32(cx), int32(cy), uint32(flags), int32(cInitial), int32(cGrow))

	return imgl
}

func (this *ImageList) Handle() w32.HIMAGELIST {
	return this.handle
}

func (this *ImageList) Destroy() bool {
	return w32.ImageList_Destroy(this.handle)
}

func (this *ImageList) SetImageCount(uNewCount uint) bool {
	return w32.ImageList_SetImageCount(this.handle, uint32(uNewCount))
}

func (this *ImageList) ImageCount() int32 {
	return w32.ImageList_GetImageCount(this.handle)
}

func (this *ImageList) AddIcon(icon *Icon) int32 {
	return w32.ImageList_AddIcon(this.handle, icon.Handle())
}

func (this *ImageList) RemoveAll() bool {
	return w32.ImageList_RemoveAll(this.handle)
}

func (this *ImageList) Remove(i int) bool {
	return w32.ImageList_Remove(this.handle, int32(i))
}
