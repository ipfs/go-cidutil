package cidutil

import (
	"math/rand"
	"testing"

	cid "github.com/ipfs/go-cid"
	mhash "github.com/multiformats/go-multihash"
)

func TestInlinerSmallValue(t *testing.T) {
	builder := Inliner{cid.V0Builder{}, 64}
	c, err := builder.Sum([]byte("Hello World"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Prefix().MhType != mhash.ID {
		t.Fatal("Inliner builder failed to use ID Multihash on small values")
	}
}

func TestInlinerLargeValue(t *testing.T) {
	builder := Inliner{cid.V0Builder{}, 64}
	data := make([]byte, 512)
	rand.Read(data)
	c, err := builder.Sum(data)
	if err != nil {
		t.Fatal(err)
	}
	if c.Prefix().MhType == mhash.ID {
		t.Fatal("Inliner builder used ID Multihash on large values")
	}
}
