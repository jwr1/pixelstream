package internal

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/muesli/termenv"
)

type Frame [256][3]uint8

func (f *Frame) SendFrame(url string) error {
	colorVal := make([]string, 256)

	for i, pixel := range f {
		colorVal[i] = fmt.Sprint((uint32(pixel[0]) << 16) | (uint32(pixel[1]) << 8) | (uint32(pixel[2]) << 0))
	}

	jsonVal := fmt.Sprintf("{\"stack\":false,\"draw\":[{\"db\":[0,0,32,8,[%s]]}]}", strings.Join(colorVal, ","))

	_, err := http.Post(url, "application/json", strings.NewReader(jsonVal))
	if err != nil {
		return err
	}

	return nil
}

func (f *Frame) ReceiveFrame(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	for index, v := range strings.Split(strings.Trim(string(body), "[]"), ",") {
		num, err := strconv.ParseUint(v, 10, 24)
		if err != nil {
			return err
		}

		f[index] = [3]uint8{uint8((num & 0xFF0000) >> 16), uint8((num & 0x00FF00) >> 8), uint8((num & 0x0000FF) >> 0)}
	}

	return nil
}

func (f *Frame) View() string {
	var s strings.Builder

	for i := 0; i < 8; i++ {
		for j := 0; j < 32; j++ {
			s.WriteString(termenv.
				String("██").
				Foreground(termenv.ColorProfile().FromColor(color.RGBA{f[i*32+j][0], f[i*32+j][1], f[i*32+j][2], 255})).
				String(),
			)
		}
		s.WriteRune('\n')
	}

	return s.String()
}
