package cidenc

import (
	"context"

	cidutil "github.com/ipfs/go-cidutil"

	cid "github.com/ipfs/go-cid"
	path "github.com/ipfs/go-path"
	mbase "github.com/multiformats/go-multibase"
)

// Encoder is a basic Encoder that will encode Cid's using a specifed
// base and optionally upgrade a CidV0 to CidV1
type Encoder struct {
	Base    mbase.Encoder
	Upgrade bool
}

// Interface is a generic interface to the Encoder functionally.
type Interface interface {
	Encode(c cid.Cid) string
	Recode(v string) (string, error)
}

// Default is the default encoder
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

// Recode reencodes the cid string to match the paramaters of the
// encoder
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

// WithOverride is like Encoder but also contains a override map to
// preserve the original encoding of select CIDs
type WithOverride struct {
	base     Encoder
	override map[cid.Cid]string
}

func (enc WithOverride) Encoder() Encoder {
	return enc.base
}

func (enc WithOverride) Map() map[cid.Cid]string {
	return enc.override
}

func NewOverride(enc Encoder) WithOverride {
	return WithOverride{base: enc, override: map[cid.Cid]string{}}
}

// Add adds a Cid to the override map if it will be encoded
// differently than the base encoder
func (enc WithOverride) Add(cids ...string) {
	for _, p := range cids {
		v := p
		c, err := cid.Decode(v)
		if err != nil {
			continue
		}
		if enc.base.Encode(c) != v {
			enc.override[c] = v
		}
		c2 := cidutil.TryOtherCidVersion(c)
		if c2.Defined() && enc.base.Encode(c2) != v {
			enc.override[c2] = v
		}
	}
}

func (enc WithOverride) Encode(c cid.Cid) string {
	v, ok := enc.override[c]
	if ok {
		return v
	}
	return enc.base.Encode(c)
}

func (enc WithOverride) Recode(v string) (string, error) {
	if len(enc.override) == 0 {
		return enc.base.Recode(v)
	}
	c, err := cid.Decode(v)
	if err != nil {
		return v, err
	}

	return enc.Encode(c), nil
}

type encoderKey struct{}

// Enable "enables" the encoder in the context using WithValue
func Enable(ctx context.Context, enc Interface) context.Context {
	return context.WithValue(ctx, encoderKey{}, enc)
}

// Get gets an encoder from the context if it exists, otherwise the
// default context is called.
func Get(ctx context.Context) Interface {
	enc, ok := ctx.Value(encoderKey{}).(Interface)
	if !ok {
		// FIXME: Warning?
		enc = Default
	}
	return enc
}
