# File to Video
This is inspired by the [Infinite-Storage-Glitch](https://github.com/DvorakDwarf/Infinite-Storage-Glitch) by DvorakDwarf. The idea is to encode binary data into a video to then upload it to basically take advantage of video-hosters as a kind of file host. This is just an experiment at replicating this in go. This repo isn't meant to actually store files on video platforms. I would strongly recommend against it because:
1. Your files are not safe. They may be deleted at any time.
2. You are most-likely violating some kind of TOS or laws.
3. It's really unpractical. Your video-files will be magnitudes bigger than your source file and up- and downloading your files will be painfully slow.

Keep in mind that this is an unfinished experiment and is not meant to actually store meaningful data in videos.

# Prerequisites
To use this, ffmpeg and ffprobe must be installed in the system path.

For v1 and v2, you have to install [vidio](https://pkg.go.dev/github.com/AlexEidt/Vidio)

For v1, you have to install [gozxing](https://github.com/makiuchi-d/gozxing)

# Data-Matrix Codes
I originally started this to make use of the error-correction found in qr-codes and data-martix codes.
I also had the idea of "overlaying" 3 data-matrix codes for each RGB-channel to increase data density, which looks like this:
This is implemented in the v1.go. You have to modify the main-function if you want to try it yourself.
| RGB-Combined |  Red Channel | Green Channel | Blue Channel |
:-:|:-:|:-:|:--:
![](https://github.com/Alex23582/FileToVideo/assets/117467716/9de91014-e662-4d4c-9c27-e2392c5115d8) | ![](https://github.com/Alex23582/FileToVideo/assets/117467716/dc2dc7a0-e1b1-447a-8eb5-d293aeb6b606) | ![](https://github.com/Alex23582/FileToVideo/assets/117467716/9644e645-2140-4a87-89b9-9d1c0559b7f5) | ![](https://github.com/Alex23582/FileToVideo/assets/117467716/211e2894-64c3-4370-85f3-af7372c0aac7)





With many frames combined into a video, the end-result looks like this:
<video src="https://github.com/Alex23582/FileToVideo/assets/117467716/b3b84f0a-e5bf-4bbd-ad0b-607e99259e62"/>
# Binary bits
The data-matrix codes had a significant performance impact. Probably because it's meant to read unaligned camera-pictures. For my purpose that wasn't necessary, so i just wrote the bits "as is" to the images with the idea of maybe utilising this [Reed-Solomon](https://github.com/klauspost/reedsolomon) go library to add the error-correction which is now missing.
You can try this with v2.go using the command-line: ```/goqrfile <encode/decode> <filename> <videofilename>```
<video src="https://github.com/Alex23582/FileToVideo/assets/117467716/8ef645e8-0d6d-4afb-99e2-c59d78709b80"/>
# Video-Compression
As DvorakDwarf already explained in his repo, the compression of youtube makes this program almost impossible to use. It's definitely not suitable for actual use of file storage. This may better with the implementation of Reed-Solomon codes in the future.

# Known Issues
 - The output-file is padded with 0-bytes, which mostly isn't a problem
