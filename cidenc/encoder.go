package cidenc

import (
	"context"

	cid "github.com/ipfs/go-cid"
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
