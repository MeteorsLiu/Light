package packetcapture

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

	"github.com/MeteorsLiu/Light/interfaces"
	"github.com/MeteorsLiu/Light/queue"
)

type ElementC struct {
	queue  *queue.Queue
	signal context.Context
}

func (e *ElementC) Upload(ip, rates *C.char) {
	if e.queue == nil {
		return
	}
	e.queue.Push(interfaces.UploadPayload{
		IP:    C.GoString(ip),
		Rates: C.GoString(rates),
	})
}

// Fuck Getter
func (e *ElementC) GetUpdate() *queue.Queue {
	return e.queue
}

//export UPLOAD
func UPLOAD(rule C.uintptr_t, ip, rates *C.char) {
	handle := cgo.Handle(rule)
	defer handle.Delete()
	CALL := handle.Value().(func(*C.char, *C.char))
	CALL(ip, rates)
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
