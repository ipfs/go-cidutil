package main

import (
	"fmt"
	"testing"

	c "github.com/ipfs/go-cid"
)

func TestCidConv(t *testing.T) {
	cidv0 := "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"
	cidv1 := "bafybeiczsscdsbs7ffqz55asqdf3smv6klcw3gofszvwlyarci47bgf354"
	cid, err := c.Decode(cidv0)
	if err != nil {
		t.Fatal(err)
	}
	cid, err = toCidV1(cid)
	if err != nil {
		t.Fatal(err)
	}
	if cid.String() != cidv1 {
		t.Fatalf("conversion failure: %s != %s", cid, cidv1)
	}
	cid, err = toCidV0(cid)
	if err != nil {
		t.Fatal(err)
	}
	cidStr := cid.String()
	if cidStr != cidv0 {
		t.Errorf("conversion failure, expected: %s; but got: %s", cidv0, cidStr)
	}
}

func TestBadCidConv(t *testing.T) {
	// this cid is a raw leaf and should not be able to convert to cidv0
	cidv1 := "bafkreifit7vvfkf2cwwzvyycdczm5znbdbqx54ab6shbesvwgkwthdf77y"
	cid, err := c.Decode(cidv1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = toCidV0(cid)
	if err == nil {
		t.Fatal("expected failure")
	}
}
