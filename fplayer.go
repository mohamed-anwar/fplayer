package main

import (
	"encoding/binary"
	"fmt"
	"github.com/hajimehoshi/oto"
	"io"
	"os"
)

type WAVHeader struct {
	/* RIFF header */
	ChunkID   string
	ChunkSize uint32
	Format    string

	/* Subchunk 1 */
	Subchunk1ID   string
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16

	/* Subchunk 2 */
	Subchunk2ID   string
	Subchunk2Size uint32
}

type WAVFile struct {
	file   *os.File
	header WAVHeader
}

func (wav *WAVFile) Load(file *os.File) error {
	wav.file = file

	header := make([]byte, 48)

	if n, err := wav.file.Read(header); err != nil {
		panic(err)
	} else {
		fmt.Printf("read %v bytes\n", n)
	}

	/* parse header */
	wav.header.ChunkID = string(header[0:4])
	wav.header.ChunkSize = binary.LittleEndian.Uint32(header[4:8])
	wav.header.Format = string(header[8:12])

	/* parse subchunk 1 */
	wav.header.Subchunk1ID = string(header[12:16])
	wav.header.Subchunk1Size = binary.LittleEndian.Uint32(header[16:20])
	wav.header.AudioFormat = binary.LittleEndian.Uint16(header[20:22])
	wav.header.NumChannels = binary.LittleEndian.Uint16(header[22:24])
	wav.header.SampleRate = binary.LittleEndian.Uint32(header[24:28])
	wav.header.ByteRate = binary.LittleEndian.Uint32(header[28:32])
	wav.header.BlockAlign = binary.LittleEndian.Uint16(header[32:34])
	wav.header.BitsPerSample = binary.LittleEndian.Uint16(header[34:36])

	/* parse subchunk 2 */
	wav.header.Subchunk2ID = string(header[36:40])
	wav.header.Subchunk2Size = binary.LittleEndian.Uint32(header[40:44])

	fmt.Printf("ChunkID   %v\n", wav.header.ChunkID)
	fmt.Printf("ChunkSize %v\n", wav.header.ChunkSize)
	fmt.Printf("Format    %v\n", wav.header.Format)

	fmt.Println()

	fmt.Printf("Subchunk1ID   %v\n", wav.header.Subchunk1ID)
	fmt.Printf("Subchunk1Size %v\n", wav.header.Subchunk1Size)
	fmt.Printf("AudioFormat   %v\n", wav.header.AudioFormat)
	fmt.Printf("NumChannels   %v\n", wav.header.NumChannels)
	fmt.Printf("SampleRate    %v\n", wav.header.SampleRate)
	fmt.Printf("ByteRate      %v\n", wav.header.ByteRate)
	fmt.Printf("BlockAlign    %v\n", wav.header.BlockAlign)
	fmt.Printf("BitsPerSample %v\n", wav.header.BitsPerSample)

	fmt.Println()

	fmt.Printf("Subchunk2ID   %v\n", wav.header.Subchunk2ID)
	fmt.Printf("Subchunk2Size %v\n", wav.header.Subchunk2Size)

	return nil
}

func (wav WAVFile) Read(buf []byte) (int, error) {
	return wav.file.Read(buf)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: fplayer filename")
		os.Exit(-1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	wav := WAVFile{}

	err = wav.Load(file)
	if err != nil {
		panic(err)
	}

	player, err := oto.NewPlayer(int(wav.header.SampleRate), int(wav.header.NumChannels), int(wav.header.BitsPerSample/8), 4096)
	if err != nil {
		panic(err)
	}

	for {
		written, err := io.Copy(player, wav)

		if err != nil {
			panic(err)
		}

		if written == 0 {
			break
		}
	}
}
