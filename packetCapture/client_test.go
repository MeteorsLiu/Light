package main

/*
#cgo LDFLAGS: -lpcap
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <unistd.h>
#include <pcap/pcap.h>
#include <signal.h>
#include <linux/if_ether.h>
#include <netinet/ip.h>
#include <netinet/in.h>
#include <netinet/if_ether.h>
#include <arpa/inet.h>
#include <netinet/tcp.h>
#include <netinet/udp.h>
#include <time.h>
#include "light.h"
*/
import "C"
import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

func main() {
	dev := C.CString("ens18")
	filter := C.CString("tcp dst port 22")
	defer C.free(unsafe.Pointer(dev))
	defer C.free(unsafe.Pointer(filter))
	defer C.stop_capture(C.int(1))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go C.Init(dev, filter)

	<-sigCh

}
