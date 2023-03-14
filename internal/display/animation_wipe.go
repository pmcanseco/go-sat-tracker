package display

import "time"

type WipeAnimation struct {
	device Device
}

func NewWipeAnimation(device Device) WipeAnimation {
	return WipeAnimation{device: device}
}

func (wa *WipeAnimation) Run() {
	wa.device.Clear()
	wa.doRun(on)
	time.Sleep(500 * time.Millisecond)
	wa.doRun(off)
	time.Sleep(500 * time.Millisecond)
}

func (wa *WipeAnimation) doRun(p PixelState) {
	for i := int16(0); i < 128; i++ {
		for j := int16(0); j < 32; j++ {
			wa.device.SetPixel(i, j, p)
		}
		_ = wa.device.Display()
		time.Sleep(1 * time.Millisecond)
	}
}
