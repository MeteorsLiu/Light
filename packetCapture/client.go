package main

/*
#cgo LDFLAGS: -lpcap
#include <stdlib.h>
#include <stdint.h>
#include <unistd.h>
#include <pcap/pcap.h>
#include "light.h"
*/
import "C"
import (
	"context"
	"unsafe"

	"github.com/MeteorsLiu/Light/interfaces"
	"github.com/MeteorsLiu/Light/queue"
)

var ResultQueue queue.Queue = nil

// C Interface upload
func upload(ip, rates *C.char) {
	if ResultQueue == nil {
		return
	}
	ResultQueue.Push(interfaces.UploadPayload{
		IP:    C.GoString(ip),
		Rates: C.GoString(rates),
	})
}
func _init(ctx context.Context, q queue.Queue, devName, filterRule string) {
	ResultQueue = q

	dev := C.CString(devName)
	filter := C.CString(filterRule)
	defer C.free(unsafe.Pointer(dev))
	defer C.free(unsafe.Pointer(filter))
	defer C.stop_capture(C.int(1))

	go C.Init(dev, filter)

	<-ctx.Done()
	ResultQueue = nil
}
