package cidenc

import (
	"context"
	"testing"

	mbase "github.com/multiformats/go-multibase"
)

func TestContext(t *testing.T) {
	enc := Encoder{Base: mbase.MustNewEncoder(mbase.Base64)}
	ctx := context.Background()
	ctx = Enable(ctx, enc)
	e, ok := Get(ctx).(Encoder)
	if !ok || e.Base.Encoding() != mbase.Base64 {
		t.Fatal("Failed to retrive encoder from context")
	}
}
