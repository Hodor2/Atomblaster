package audio

import (
	"bytes"
	"encoding/binary"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// generateWavData creates a valid WAV file in memory with a sine-wave PCM data.
func generateWavData(frequency float32, duration float32, sampleRate int, volume float32) []byte {
	channels := 1
	bitsPerSample := 16
	
	sampleCount := int(duration * float32(sampleRate))
	dataSize := sampleCount * channels * (bitsPerSample / 8)
	chunkSize := 36 + dataSize

	buf := new(bytes.Buffer)
	// Write the RIFF header.
	buf.WriteString("RIFF")
	binary.Write(buf, binary.LittleEndian, uint32(chunkSize))
	buf.WriteString("WAVE")

	// Write the fmt subchunk.
	buf.WriteString("fmt ")
	binary.Write(buf, binary.LittleEndian, uint32(16))         // Subchunk1Size for PCM
	binary.Write(buf, binary.LittleEndian, uint16(1))          // AudioFormat 1=PCM
	binary.Write(buf, binary.LittleEndian, uint16(channels))   // NumChannels
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate)) // SampleRate
	byteRate := sampleRate * channels * bitsPerSample / 8
	binary.Write(buf, binary.LittleEndian, uint32(byteRate))      // ByteRate
	blockAlign := channels * bitsPerSample / 8
	binary.Write(buf, binary.LittleEndian, uint16(blockAlign))    // BlockAlign
	binary.Write(buf, binary.LittleEndian, uint16(bitsPerSample)) // BitsPerSample

	// Write the data subchunk header.
	buf.WriteString("data")
	binary.Write(buf, binary.LittleEndian, uint32(dataSize))

	// Apply an envelope to avoid clicks (attack and release)
	attackTime := float32(0.05) // 50ms attack
	releaseTime := float32(0.1) // 100ms release
	attackSamples := int(attackTime * float32(sampleRate))
	releaseSamples := int(releaseTime * float32(sampleRate))
	
	// Generate and write the sample data (int16 samples).
	for i := 0; i < sampleCount; i++ {
		// Apply envelope
		amplitude := volume
		if i < attackSamples {
			amplitude *= float32(i) / float32(attackSamples)
		} else if i > sampleCount - releaseSamples {
			amplitude *= float32(sampleCount - i) / float32(releaseSamples)
		}
		
		// Sine wave formula for PCM16 with envelope
		sinValue := float32(math.Sin(2*math.Pi*float64(frequency)*float64(i)/float64(sampleRate)))
		sample := int16(32767 * amplitude * sinValue)
		binary.Write(buf, binary.LittleEndian, sample)
	}

	return buf.Bytes()
}

// GenerateSound creates a sound effect by generating a valid WAV file in memory
func GenerateSound(frequency float32, duration float32, volume float32) rl.Sound {
	sampleRate := 44100
	
	// Generate the WAV file as a byte slice with volume control
	wavData := generateWavData(frequency, duration, sampleRate, volume)
	
	// Load the wave from memory.
	wave := rl.LoadWaveFromMemory(".wav", wavData, int32(len(wavData)))
	
	// Convert the Wave to a Sound.
	sound := rl.LoadSoundFromWave(wave)
	rl.UnloadWave(wave)
	return sound
}
