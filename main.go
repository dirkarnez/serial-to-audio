package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/tarm/serial"
)

const (
	comPort      = "/dev/ttyUSB0" // Change this to your COM port
	baudRate     = 9600           // Change this to your desired baud rate
	sampleRate   = 44100          // Audio sample rate
	numChannels  = 1              // Number of audio channels
	bitsPerSample = 16             // Bits per sample
)

func main() {
	// Open the COM port
	port, err := serial.OpenPort(&serial.Config{
		Name: comPort,
		Baud: baudRate,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// Create a new WAV file
	file, err := os.Create("output.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new WAV encoder
	enc := wav.NewEncoder(file, sampleRate, bitsPerSample, numChannels, 1)

	// Create an audio buffer
	buf := &audio.IntBuffer{
		Format: &audio.Format{
			SampleRate:  sampleRate,
			NumChannels: numChannels,
		},
	}

	// Read ADC values from the COM port and save to the WAV file
	for {
		// Read a byte from the COM port
		buf := make([]byte, 1)
		_, err := port.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		// Convert the byte to an int16 value
		value := int16(buf[0])

		// Add the value to the audio buffer
		buf.Data = append(buf.Data, int(value))

		// If the buffer is full, write it to the WAV file
		if len(buf.Data) >= sampleRate {
			if err := enc.Write(buf); err != nil {
				log.Fatal(err)
			}
			buf.Data = nil
		}

		// Sleep for a short period to control the data rate
		time.Sleep(time.Millisecond)
	}

	// Flush any remaining data in the audio buffer
	if err := enc.Write(buf); err != nil {
		log.Fatal(err)
	}
}
