package internal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

const currentPixelStreamFormatVersion uint8 = 1

func SaveFile(output_file string, data *PixelStream) error {
	buf := bytes.NewBuffer(make([]byte, 0, len(pixelstreamFormatIdentifier)+2+(len(data.Frames)*frameSize)))

	_, err := buf.WriteString(pixelstreamFormatIdentifier)
	if err != nil {
		return err
	}

	_, err = buf.Write([]byte{data.Version, data.FrameRate})
	if err != nil {
		return err
	}

	for _, frame := range data.Frames {
		for _, pixel := range frame {
			_, err = buf.Write(pixel[:])
			if err != nil {
				return err
			}
		}
	}

	err = os.WriteFile(output_file, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadFile(source_path string) (*PixelStream, error) {
	file, err := os.ReadFile(source_path)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(pixelstreamFormatIdentifier, string(file[0:len(pixelstreamFormatIdentifier)])) {
		return nil, errors.New("invalid/corrupt pxlstrm file")
	}

	fileVersion := file[len(pixelstreamFormatIdentifier)]

	if fileVersion != currentPixelStreamFormatVersion {
		return nil, fmt.Errorf("unsupported pxlstrm format version: found %d, expected %d", fileVersion, currentPixelStreamFormatVersion)
	}

	frameCount := (len(file) - (len(pixelstreamFormatIdentifier) + 2)) / frameSize

	output := &PixelStream{
		Version:   fileVersion,
		FrameRate: file[len(pixelstreamFormatIdentifier)+1],
		Frames:    make([][256][3]uint8, 0, frameCount),
	}

	for i := len(pixelstreamFormatIdentifier) + 2; i < len(file); i += frameSize {
		var frame [256][3]uint8

		for j := 0; j < frameSize; j += 3 {
			frame[j/3] = [3]uint8{file[i+j], file[i+j+1], file[i+j+2]}
		}

		output.Frames = append(output.Frames, frame)
	}

	return output, nil
}
