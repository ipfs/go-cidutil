package cidutil

import (
	cid "github.com/ipfs/go-cid"
	mhash "github.com/multiformats/go-multihash"
)

// Inliner is a cid.Builder that will use the id multihash when the
// size of the content is no more than limit
type Inliner struct {
	cid.Builder
	Limit int
}

// WithCodec implements the cid.Builder interface
func (p Inliner) WithCodec(c uint64) cid.Builder {
	return Inliner{p.Builder.WithCodec(c), p.Limit}
}

// Sum implements the cid.Builder interface
func (p Inliner) Sum(data []byte) (*cid.Cid, error) {
	if len(data) > p.Limit {
		return p.Builder.Sum(data)
	}
	return cid.V1Builder{Codec: p.GetCodec(), MhType: mhash.ID}.Sum(data)
}
