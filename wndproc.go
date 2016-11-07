package gform

import (
	"log"
	"unsafe"

	"github.com/darkautism/w32"
)

func genPoint(p uintptr) (x, y int) {
	x = int(w32.LOWORD(uint32(p)))
	y = int(w32.HIWORD(uint32(p)))
	return
}

func genMouseEventArg(wparam, lparam uintptr) *MouseEventData {
	var data MouseEventData
	data.Button = int(wparam)
	data.X, data.Y = genPoint(lparam)

	return &data
}

func genDropFilesEventArg(wparam uintptr) *DropFilesEventData {
	hDrop := w32.HDROP(wparam)

	var data DropFilesEventData
	_, fileCount := w32.DragQueryFile(hDrop, 0xFFFFFFFF)
	data.Files = make([]string, fileCount)

	var i uint32
	for i = 0; i < fileCount; i++ {
		data.Files[i], _ = w32.DragQueryFile(hDrop, i)
	}

	data.X, data.Y, _ = w32.DragQueryPoint(hDrop)

	w32.DragFinish(hDrop)

	return &data
}

func generalWndProc(hwnd w32.HWND, msg uint, wparam, lparam uintptr) uintptr {
	if msg == w32.WM_INITDIALOG && gDialogWaiting != nil {
		gDialogWaiting.hwnd = hwnd
		RegMsgHandler(gDialogWaiting)
	}
	log.Println("hwnd:", hwnd, "msg:", msg)
	return w32.DefWindowProc(hwnd, uint32(msg), wparam, lparam)
	if controller := GetMsgHandler(hwnd); controller != nil {
		var ret uintptr
		switch msg {
		case w32.WM_NOTIFY: //Reflect notification to control
			ret = controller.WndProc(msg, wparam, lparam)
			nm := (*w32.NMHDR)(unsafe.Pointer(lparam))
			if controller := GetMsgHandler(nm.HwndFrom); controller != nil {
				ret := controller.WndProc(msg, wparam, lparam)
				if ret != 0 {
					w32.SetWindowLong(hwnd, w32.DWL_MSGRESULT, uint32(ret))
					return w32.TRUE
				}
			}
		case w32.WM_COMMAND:
			if lparam != 0 { //Reflect message to control
				ret = controller.WndProc(msg, wparam, lparam)
				h := w32.HWND(lparam)
				if controller := GetMsgHandler(h); controller != nil {
					ret := controller.WndProc(msg, wparam, lparam)
					if ret != 0 {
						w32.SetWindowLong(hwnd, w32.DWL_MSGRESULT, uint32(ret))
						return w32.TRUE
					}
				}
			}
		case w32.WM_CLOSE:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnClose().Fire(NewEventArg(controller, nil))
		case w32.WM_KILLFOCUS:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnKillFocus().Fire(NewEventArg(controller, nil))
		case w32.WM_SETFOCUS:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnSetFocus().Fire(NewEventArg(controller, nil))
		case w32.WM_DROPFILES:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnDropFiles().Fire(NewEventArg(controller, genDropFilesEventArg(wparam)))
		case w32.WM_LBUTTONDOWN:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnLBDown().Fire(NewEventArg(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_LBUTTONUP:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnLBUp().Fire(NewEventArg(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_MBUTTONDOWN:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnMBDown().Fire(NewEventArg(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_MBUTTONUP:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnMBUp().Fire(NewEventArg(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_RBUTTONDOWN:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnRBDown().Fire(NewEventArg(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_RBUTTONUP:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnRBUp().Fire(NewEventArg(controller, genMouseEventArg(wparam, lparam)))
		case w32.WM_PAINT:
			ret = controller.WndProc(msg, wparam, lparam)
			canvas := NewCanvasFromHwnd(hwnd)
			defer canvas.Dispose()
			controller.OnPaint().Fire(NewEventArg(controller, &PaintEventData{Canvas: canvas}))
		case w32.WM_KEYUP:
			ret = controller.WndProc(msg, wparam, lparam)
			controller.OnKeyUp().Fire(NewEventArg(controller, &KeyUpEventData{int(wparam), int(lparam)}))
		case w32.WM_SIZE:
			ret = controller.WndProc(msg, wparam, lparam)
			x, y := genPoint(lparam)
			controller.OnSize().Fire(NewEventArg(controller, &SizeEventData{uint(wparam), x, y}))
		default:
			if handler, ok := controller.BindedHandler(msg); ok {
				ret = controller.WndProc(msg, wparam, lparam)
				handler(NewEventArg(controller, &RawMsg{hwnd, msg, wparam, lparam}))
			} else {
				//ret = w32.DefWindowProc(hwnd, uint32(msg), wparam, lparam)
			}
		}
		_ = ret
		return w32.DefWindowProc(hwnd, uint32(msg), wparam, lparam)
	}

	return w32.DefWindowProc(hwnd, msg, wparam, lparam)
}
