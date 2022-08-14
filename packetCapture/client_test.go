package main

import (
	"context"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	t.Run("test1", func(t *testing.T){
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		var r ElementC
		NewClient(ctx, "ens18", "tcp port 22", &r)
		for {
		select {
		case <-ctx.Done():
			return
		case elem := <-r.GetUpdate().Queue:
			t.Log(elem)
		}
		
		}
	})
}
