package generator

import (
	"math"
	"math/rand"
)

const (
	TwoPi = float64(2 * math.Pi)
)

const (
	SineB = 4.0 / math.Pi
	SineC = -4.0 / (math.Pi * math.Pi)
	Q     = 0.775
	SineP = 0.225
)

type WaveFunc = func(x float64) float64

// Sine takes an input value from -Pi to Pi
// and returns a value between -1 and 1
func Sine(x float64) float64 {
	xs := float64(x)
	ys := SineB*xs + SineC*xs*(math.Abs(xs))
	ys = SineP*(ys*(math.Abs(ys))-ys) + ys
	return float64(ys)
}

const TringleA = 2.0 / math.Pi

// Triangle takes an input value from -Pi to Pi
// and returns a value between -1 and 1
func Triangle(x float64) float64 {
	return float64(TringleA*x) - 1.0
}

// Square takes an input value from -Pi to Pi
// and returns -1 or 1
func Square(x float64) float64 {
	if x >= 0.0 {
		return 1
	}
	return -1.0
}

const SawtoothA = 1.0 / math.Pi

// Triangle takes an input value from -Pi to Pi
// and returns a value between -1 and 1
func Sawtooth(x float64) float64 {
	return SawtoothA * x
}

func WhiteNoise(x float64) float64 {
	return (rand.Float64()*2 - 1)
}
