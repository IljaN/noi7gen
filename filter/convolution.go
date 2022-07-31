package filter

import (
	"container/ring"
	"github.com/go-audio/audio"
)

// The simple convolution filter is efficient when the number
// of coefficients is small (on my desktop: <300), and works
// by multiplying the coefficients with the values in the
// delay line.
//
// This filter should be appropriate for implementing low pass,
// high pass, band pass and band stop filters.
//
type SimpleConvolutionFilter struct {
	Coefficients   []float64
	LeftDelayLine  *ring.Ring
	RightDelayLine *ring.Ring
}

func NewSimpleConvolutionFilter(s *SimpleConvolutionFilter) func(buf *audio.FloatBuffer) {
	return func(buf *audio.FloatBuffer) {
		isStereo := true
		order := len(s.Coefficients)
		if s.LeftDelayLine == nil {
			s.LeftDelayLine = ring.New(order)
		}
		if isStereo {
			if s.RightDelayLine == nil {
				s.RightDelayLine = ring.New(order)
			}
		}

		n := len(buf.Data)
		if isStereo {
			n = n / 2
		}
		for i := 0; i < n; i++ {
			ix := i
			if isStereo {
				ix *= 2
			}

			s.LeftDelayLine.Value = buf.Data[ix]
			sample := 0.0
			left := s.LeftDelayLine
			for j := 0; j < order; j++ {
				v := left.Value
				left = left.Prev()
				value := 0.0
				if v != nil {
					value = v.(float64)
					sample += s.Coefficients[j] * value
				}
			}
			buf.Data[ix] = sample
			s.LeftDelayLine = s.LeftDelayLine.Next()

			if isStereo {
				ix += 1
				s.RightDelayLine.Value = buf.Data[ix]
				sample = 0.0
				right := s.RightDelayLine
				for j := 0; j < order; j++ {
					v := right.Value
					right = right.Prev()
					value := 0.0
					if v != nil {
						value = v.(float64)
						sample += s.Coefficients[j] * value
					}
				}
				buf.Data[ix] = sample
				s.RightDelayLine = s.RightDelayLine.Next()
			}
		}
	}

}
