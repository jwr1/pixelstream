package internal

import (
	"fmt"
	"time"
)

var Host string

type PixelStream struct {
	Version   uint8
	FrameRate uint8
	Frames    []Frame
}

func (ps *PixelStream) GetTotalDuration() time.Duration {
	return time.Duration((int64(len(ps.Frames)) * 1000 / int64(ps.FrameRate)) * int64(time.Millisecond))
}

func (ps *PixelStream) GetFrame(d time.Duration) *Frame {
	frameNum := float64(ps.FrameRate) * d.Seconds()

	return &ps.Frames[int64(frameNum)]
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

					ps.Frames[currentFrame].SendFrame(notifyUrl)

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
