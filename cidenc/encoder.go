package cidenc

import (
	cid "github.com/ipfs/go-cid"
	path "github.com/ipfs/go-path"
	mbase "github.com/multiformats/go-multibase"
)

// Encoder is a basic Encoder that will encode Cid's using
// a specifed base, optionally upgrading a Cid if is Version 0
type Encoder struct {
	Base    mbase.Encoder
	Upgrade bool
}

// Interface is a generic interface to the Encoder functionally.
type Interface interface {
	Encode(c cid.Cid) string
	Recode(v string) (string, error)
}

// Default is the
var Default = Encoder{
	Base:    mbase.MustNewEncoder(mbase.Base58BTC),
	Upgrade: false,
}

func (enc Encoder) Encode(c cid.Cid) string {
	if enc.Upgrade && c.Version() == 0 {
		c = cid.NewCidV1(c.Type(), c.Hash())
	}
	return c.Encode(enc.Base)
}

func (enc Encoder) Recode(v string) (string, error) {
	skip, err := enc.noopRecode(v)
	if skip || err != nil {
		return v, err
	}

	c, err := cid.Decode(v)
	if err != nil {
		return v, err
	}

	return enc.Encode(c), nil
}

func (enc Encoder) noopRecode(v string) (bool, error) {
	if len(v) < 2 {
		return false, cid.ErrCidTooShort
	}
	ver := cidVer(v)
	skip := ver == 0 && !enc.Upgrade || ver == 1 && v[0] == byte(enc.Base.Encoding())
	return skip, nil
}

func cidVer(v string) int {
	if len(v) == 46 && v[:2] == "Qm" {
		return 0
	} else {
		return 1
	}
}

// FromPath creates a new encoder that is influenced from the encoded
// Cid in a Path.  For CidV0 the multibase from the base encoder is
// used and automatic upgrades are disabled.  For CidV1 the multibase
// from the CID is used and upgrades are eneabled.  On error the base
// encoder is returned.  If you don't care about the error condiation
// it is safe to ignore the error returned.
func FromPath(enc Encoder, p string) (Encoder, error) {
	v := extractCidString(p)
	if cidVer(v) == 0 {
		return Encoder{enc.Base, false}, nil
	}
	e, err := mbase.NewEncoder(mbase.Encoding(v[0]))
	if err != nil {
		return enc, err
	}
	return Encoder{e, true}, nil
}

func extractCidString(p string) string {
	segs := path.FromString(p).Segments()
	v := segs[0]
	if v == "ipfs" && len(segs) > 0 {
		v = segs[1]
	}
	return v
}
