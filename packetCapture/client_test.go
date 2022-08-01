package packagecapture

import (
	"context"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	t.Run("test1", func(t *testing.T){
		ctx, cancel := context.WithTimeout(context.Backgroud(), 1*time.Minute)
		defer cancel()

		var r ElementC
		NewClient(ctx, "ens18", "tcp port 9999", &r)
		update := r.GetUpdate()
		for {
		select {
		case <-ctx.Done():
			return
		case elem := <-update:
			t.Log(elem)
		}
		
		}
	})
}