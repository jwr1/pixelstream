# pixelstream

Stream videos to your awtrix clock with ease

### Usage

pixelstream is a command line utility, see the following section for examples of usage. [ffmpeg](https://ffmpeg.org/) is required in order for conversions to work.

```bash
pixelstream play video.mp4 http://192.168.1.170
pixelstream play video.mp4 http://192.168.1.170 --frame-rate 10
pixelstream play video.pxlstrm http://192.168.1.170
```

Basically, when using the `play` command, you need to specify the path to the video file and the clock's host (ip address and protocol). If you have not played a specific file before, pixelstream will use ffmpeg to convert it to a usable format, this could take between a few minutes to half an hour depending on the duration and resolution of the original file. Once the conversion is complete, it will save to a new file with `.pxlstrm` at the end, this will make it so it doesn't have to convert the same file again in the future.

For more commands and options, use the `--help` flag when running `pixelstream`.

If you'd like to try a pre-converted video (maybe you don't have ffmpeg installed), then download one of the `.pxlstrm` files from the `samples` directory of this repository.
