package filter

// https://github.com/go-audio/transforms/blob/master/bit_crush.go

import (
	"github.com/go-audio/audio"
	"math"
)

var (
	crusherStepSize  = 0.000001
	CrusherMinFactor = 1.0
	CrusherMaxFactor = 2097152.0
)

// BitCrusher reduces the resolution of the sample to the target bit depth
// Note that bit crusher effects are usually made of this feature + a decimator
type BitCrusher struct {
	stepSize float64
	Factor   float64
}

func NewBitCrusher(cr *BitCrusher) func(buf *audio.FloatBuffer) {
	if cr.Factor < CrusherMinFactor {
		cr.Factor = 1.0
	}

	if cr.Factor > CrusherMaxFactor {
		cr.Factor = CrusherMaxFactor
	}

	return func(buf *audio.FloatBuffer) {
		stepSize := crusherStepSize * cr.Factor
		for i := 0; i < len(buf.Data); i++ {
			frac, exp := math.Frexp(buf.Data[i])
			frac = signum(frac) * math.Floor(math.Abs(frac)/stepSize+0.5) * stepSize
			buf.Data[i] = math.Ldexp(frac, exp)
		}
	}
}

func signum(v float64) float64 {
	if v >= 0.0 {
		return 1.0
	}
	return -1.0
}
