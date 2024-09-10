package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

const pixelstreamFileExt = ".pxlstrm"
const pixelstreamFormatIdentifier = "PXLSTRM"
const pixelstreamFormatVersion uint8 = 1

const frameSize = frameArea * 3

func (ps *PixelStream) SaveFile(fl FileLocation) error {
	buf := bytes.NewBuffer(make([]byte, 0, len(pixelstreamFormatIdentifier)+2+(len(ps.Frames)*frameSize)))

	_, err := buf.WriteString(pixelstreamFormatIdentifier)
	if err != nil {
		return err
	}

	_, err = buf.Write([]byte{ps.Version, ps.FrameRate})
	if err != nil {
		return err
	}

	for _, frame := range ps.Frames {
		for _, pixel := range frame {
			_, err = buf.Write(pixel[:])
			if err != nil {
				return err
			}
		}
	}

	err = fl.WriteFile(buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadFile(fl FileLocation) (*PixelStream, error) {
	file, err := fs.ReadFile(fl.System, fl.Path)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(pixelstreamFormatIdentifier, string(file[0:len(pixelstreamFormatIdentifier)])) {
		return nil, errors.New("invalid/corrupt pxlstrm file")
	}

	fileVersion := file[len(pixelstreamFormatIdentifier)]

	if fileVersion != pixelstreamFormatVersion {
		return nil, fmt.Errorf("unsupported pxlstrm format version: found %d, expected %d", fileVersion, pixelstreamFormatVersion)
	}

	frameCount := (len(file) - (len(pixelstreamFormatIdentifier) + 2)) / frameSize

	output := &PixelStream{
		Version:   fileVersion,
		FrameRate: file[len(pixelstreamFormatIdentifier)+1],
		Frames:    make([]Frame, 0, frameCount),
	}

	for i := len(pixelstreamFormatIdentifier) + 2; i < len(file); i += frameSize {
		var frame Frame

		for j := 0; j < frameSize; j += 3 {
			frame[j/3] = [3]uint8{file[i+j], file[i+j+1], file[i+j+2]}
		}

		output.Frames = append(output.Frames, frame)
	}

	return output, nil
}
