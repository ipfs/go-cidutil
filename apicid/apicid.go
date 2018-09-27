package apicid

import (
	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-cidutil/cidenc"
	mbase "github.com/multiformats/go-multibase"
)

// JSONBase is the base to use when Encoding into JSON.
var JSONBase mbase.Encoder = mbase.MustNewEncoder(mbase.Base58BTC)

// apicid.Hash is a type to respesnt a CID in the API which marshals
// as a string
type Hash struct {
	str string
}

// FromCid creates an APICid from a Cid
func FromCid(c cid.Cid) Hash {
	return Hash{c.Encode(JSONBase)}
}

// Cid converts an APICid to a CID
func (c Hash) Cid() (cid.Cid, error) {
	return cid.Decode(c.str)
}

func (c Hash) String() string {
	return c.Encode(cidenc.Default)
}

func (c Hash) Encode(enc cidenc.Interface) string {
	if c.str == "" {
		return ""
	}
	str, err := enc.Recode(c.str)
	if err != nil {
		return c.str
	}
	return str
}

func (c *Hash) UnmarshalText(b []byte) error {
	c.str = string(b)
	return nil
}

func (c Hash) MarshalText() ([]byte, error) {
	return []byte(c.str), nil
}

// Cid is type to represent a normal CID in the API which marshals
// like a normal CID i.e. ({"/": <HASH>}) but may uses cidenc.Default
// for the String() to optionally upgrade a version 0 CID to version 1
type Cid struct {
	cid.Cid
}

func (c Cid) String() string {
	return cidenc.Default.Encode(c.Cid)
}
