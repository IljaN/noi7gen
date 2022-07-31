package filter

import (
	"github.com/go-audio/audio"
	"math"
)

type FlangerFilter struct {
	Time    float64
	Factor  float64
	LFORate float64

	leftDelayed  []float64
	rightDelayed []float64
	phase        int
}

func NewFlangerFilter(f *FlangerFilter) func(buf *audio.FloatBuffer) {
	return func(buf *audio.FloatBuffer) {
		sampleRate := buf.Format.SampleRate
		isStereo := buf.Format.NumChannels == 2
		time := int(float64(sampleRate) * f.Time)
		if f.leftDelayed == nil {
			f.leftDelayed = make([]float64, time)
			f.rightDelayed = make([]float64, time)
		}

		n := len(buf.Data)
		if isStereo {
			n = n / 2
		}
		stepSize := (f.LFORate * math.Pi) / float64(sampleRate)
		for i := 0; i < n; i++ {
			ix := i
			if isStereo {
				ix *= 2
			}

			currentDelay := time - int(math.Ceil(float64(time)*math.Abs(math.Sin(float64(f.phase)*stepSize))))
			delayedIx := f.phase - currentDelay
			if delayedIx < 0 {
				delayedIx += time
			}

			f.leftDelayed[f.phase] = buf.Data[ix]
			if isStereo {
				f.rightDelayed[f.phase] = buf.Data[ix+1]
			}

			buf.Data[ix] = (1.0-f.Factor)*buf.Data[ix] + f.Factor*f.leftDelayed[delayedIx]
			if isStereo {
				buf.Data[ix+1] = (1.0-f.Factor)*buf.Data[ix+1] + f.Factor*f.rightDelayed[delayedIx]
			}

			f.phase += 1
			if f.phase >= len(f.leftDelayed) {
				f.phase = 0
			}
		}
	}
}
