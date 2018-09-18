package cidutil

import (
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

func TryOtherCidVersion(c cid.Cid) cid.Cid {
	prefix := c.Prefix()
	if prefix.Codec != cid.DagProtobuf || prefix.MhType != mh.SHA2_256 || prefix.MhLength != 32 {
		return cid.Undef
	}
	var c1 cid.Cid
	if prefix.Version == 0 {
		c1 = cid.NewCidV1(cid.DagProtobuf, c.Hash())
	} else {
		c1 = cid.NewCidV0(c.Hash())
	}
	return c1
}
