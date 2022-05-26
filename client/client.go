package main

// #cgo LDFLAGS: -lpcap
// #include "../packetCapture/light.h"

import "C"
import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

func main() {
	dev := C.CString("ens5")
	filter := C.CString("tcp dst port 22")
	defer C.free(unsafe.Pointer(dev))
	defer C.free(unsafe.Pointer(filter))
	defer C.stop_capture(C.int(1))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go C.Init(dev, filter)

	<-sigCh

}
