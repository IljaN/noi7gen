package encoding

import (
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"os"
)

// WriteWav renders a buffer to a wave-file
func WriteWAV(fileName string, out *audio.FloatBuffer, bitRate int) error {
	wavFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer wavFile.Close()

	enc := wav.NewEncoder(wavFile, out.PCMFormat().SampleRate, bitRate, out.PCMFormat().NumChannels, 1)
	if err := enc.Write(out.AsIntBuffer()); err != nil {
		return err
	}

	defer enc.Close()
	return nil
}
