package main

import "context"

func main() {
	config := readConf()
	ctx, cancel := context.WithCancel()
	gRPCClient, conn := NewGRPCClient(config)
	defer conn.Close()
	defer cancel()

	go NewPollEventDaemon(ctx)
}
