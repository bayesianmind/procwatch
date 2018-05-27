package watch

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

// thank you https://stackoverflow.com/questions/22949444/using-golang-to-get-windows-idle-time-getlastinputinfo-or-similar/22950406#22950406
type lastInputInfo struct {
	cbSize uint32
	dwTime uint32
}

var cbSize = uint32(unsafe.Sizeof(lastInputInfo{}))
var user32 = syscall.MustLoadDLL("user32.dll")
var kernal32 = syscall.MustLoadDLL("Kernel32.dll")
var getLastInputInfo = user32.MustFindProc("GetLastInputInfo")
var getTickCount = kernal32.MustFindProc("GetTickCount")

// GetIdleTime gets the time since last user input
func GetIdleTime() (time.Duration, error) {
	// err is always non-nil calling these
	tickCount, _, err := getTickCount.Call()
	if tickCount == 0 {
		return 0, fmt.Errorf("could not get tick count: %v", err)
	}
	lii := lastInputInfo{cbSize: cbSize}
	r0, _, err := getLastInputInfo.Call(uintptr(unsafe.Pointer(&lii)))
	if r0 == 0 {
		return 0, fmt.Errorf("could not get last input info: %v", err)
	}
	idle_i64 := int64(tickCount) - int64(lii.dwTime)
	return time.Duration(idle_i64) * time.Millisecond, nil
}
