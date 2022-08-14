package main

/*
#cgo LDFLAGS: -lpcap
#include <stdlib.h>
#include <stdint.h>
#include <unistd.h>
#include <pcap/pcap.h>
#include "light.h"
extern void UPLOAD(uintptr_t rule, char * ip, char * rates);
*/
import "C"
import (
	"context"
	"runtime/cgo"
	"unsafe"
	"log"
	"github.com/MeteorsLiu/Light/interfaces"
	"github.com/MeteorsLiu/Light/queue"
)

type ElementC struct {
	queue  *queue.Queue
	signal context.Context
}

func (e *ElementC) Upload(ip, rates string) {
		if e.queue == nil {
		return
	}
	log.Println(e.queue.Push(interfaces.UploadPayload{
		IP:    ip,
		Rates: rates,
	}))
}

// Fuck Getter
func (e *ElementC) GetUpdate() *queue.Queue {
	return e.queue
}

//export UPLOAD
func UPLOAD(rule C.uintptr_t, ip, rates *C.char) {
	handle := cgo.Handle(rule)
	CALL := handle.Value().(func(string, string))
	_ip, _rates := C.GoString(ip), C.GoString(rates)
	CALL(_ip, _rates)
}
func NewClient(ctx context.Context, devName, filterRule string, returnPtr *ElementC) {
	dev := C.CString(devName)
	filter := C.CString(filterRule)
	defer C.free(unsafe.Pointer(dev))
	defer C.free(unsafe.Pointer(filter))
	defer C.stop_capture(C.int(1))

	rule := &ElementC{
		queue:  queue.New(),
		signal: ctx,
	}

	funcPtr := C.uintptr_t(cgo.NewHandle(rule.Upload))

	*returnPtr = *rule

	// Start the daemon process
	go C.Init(funcPtr, dev, filter)
	// Wait until signal
	<-ctx.Done()
}
