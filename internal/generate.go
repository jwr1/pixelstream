package internal

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"os/exec"
	"path"

	"golang.org/x/image/bmp"
)

func GeneratePixelStream(sourceFile FileLocation, frameRate uint8) (*PixelStream, error) {
	if sourceFile.System != OSFS {
		return nil, fmt.Errorf("GeneratePixelStream source file must be from the OS FS")
	}

	dirPath, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(dirPath)

	cmd := exec.Command("ffmpeg", "-i", "/"+sourceFile.Path, "-filter:v", fmt.Sprintf("fps=%d,scale=32:8", frameRate), "-c:a", "copy", path.Join(dirPath, "%d.bmp"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	dirRead, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	frameCount := len(dirRead)

	pixelstream := &PixelStream{
		Version:   pixelstreamFormatVersion,
		FrameRate: frameRate,
		Frames:    make([]Frame, frameCount),
	}

	for i := 0; i < frameCount; i++ {
		file, err := os.Open(path.Join(dirPath, fmt.Sprint(i+1)+".bmp"))
		if err != nil {
			return nil, err
		}

		defer file.Close()

		img, err := bmp.Decode(file)
		if err != nil {
			return nil, err
		}

		rect := img.Bounds()
		rgba := image.NewRGBA(rect)
		draw.Draw(rgba, rect, img, rect.Min, draw.Src)

		var frame [256][3]uint8

		for j := 0; j < len(rgba.Pix); j += 4 {
			frame[j/4] = [3]uint8{rgba.Pix[j], rgba.Pix[j+1], rgba.Pix[j+2]}
		}

		pixelstream.Frames[i] = frame
	}

	return pixelstream, nil
}
