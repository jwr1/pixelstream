package internal

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const pixelstreamFormatIdentifier = "PXLSTRM"
const frameSize = 256 * 3

type PixelStream struct {
	Version   uint8
	FrameRate uint8
	Frames    [][256][3]uint8
}

func (ps *PixelStream) GetDuration() time.Duration {
	return time.Duration((int64(len(ps.Frames)) * 1000 / int64(ps.FrameRate)) * int64(time.Millisecond))
}

func (ps *PixelStream) Stream(host string) {
	notifyUrl := host + "/api/notify"

	currentFrame := 0

	ticker := time.NewTicker(time.Duration((1000 / int64(ps.FrameRate)) * int64(time.Millisecond)))
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				if currentFrame < len(ps.Frames) {
					fmt.Println("Frame", currentFrame, "Time", t)

					sendFrame(notifyUrl, ps.Frames[currentFrame])

					currentFrame++
				} else {
					ticker.Stop()
					done <- true
					fmt.Println("Stream Done!")
				}
			}
		}
	}()

	<-done
}

func sendFrame(url string, frame [256][3]uint8) {
	colorVal := make([]string, 256)

	for i, pixel := range frame {
		colorVal[i] = fmt.Sprint((uint32(pixel[0]) << 16) | (uint32(pixel[1]) << 8) | (uint32(pixel[2]) << 0))
	}

	jsonVal := fmt.Sprintf("{\"stack\":false,\"draw\":[{\"db\":[0,0,32,8,[%s]]}]}", strings.Join(colorVal, ","))

	_, e := http.Post(url, "application/json", strings.NewReader(jsonVal))
	if e != nil {
		fmt.Fprintln(os.Stderr, "Error:", e)
	}
}
