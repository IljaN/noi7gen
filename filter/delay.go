package filter

import (
	"container/ring"
	"github.com/go-audio/audio"
)

type DelayFilter struct {
	LeftTime     float64
	LeftFactor   float64
	LeftFeedback float64

	RightTime     float64
	RightFactor   float64
	RightFeedback float64
	rightDelayed  *ring.Ring
	leftDelayed   *ring.Ring
}

func NewDelayFilter(f *DelayFilter) func(buf *audio.FloatBuffer) {
	return func(buf *audio.FloatBuffer) {
		isStereo := buf.Format.NumChannels == 2
		sampleRate := buf.Format.SampleRate
		leftDelayTime := int(float64(sampleRate) * f.LeftTime)

		if f.leftDelayed == nil {
			f.leftDelayed = ring.New(leftDelayTime)
		}

		if isStereo {
			rightDelayTime := int(float64(sampleRate) * f.LeftTime)
			if f.rightDelayed == nil {
				f.rightDelayed = ring.New(rightDelayTime)
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

			s := Delay(buf.Data[ix], f.LeftFactor, f.LeftFeedback, f.leftDelayed)
			f.leftDelayed = f.leftDelayed.Next()
			buf.Data[ix] = s

			if isStereo {
				s := Delay(buf.Data[ix+1], f.RightFactor, f.RightFeedback, f.rightDelayed)
				f.rightDelayed = f.rightDelayed.Next()
				buf.Data[ix+1] = s
			}
		}
	}
}

func Delay(s, factor, feedback float64, ring *ring.Ring) float64 {
	if ring.Value != nil {
		prev := ring.Value.(float64)
		ring.Value = s
		s += prev * factor
		ring.Value = ring.Value.(float64) + feedback*s
	} else {
		ring.Value = s
	}
	return s
}
