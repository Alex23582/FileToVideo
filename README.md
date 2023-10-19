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
![](https://github.com/Alex23582/FileToVideo/assets/117467716/fc2b08fd-60ae-448b-b767-288636469de6) | ![](https://github.com/Alex23582/FileToVideo/assets/117467716/6898df38-c1d3-4951-b348-20aba34d1b9f) | ![](https://github.com/Alex23582/FileToVideo/assets/117467716/b3b72185-2bf9-4530-90ea-5717b84f2367) | ![](https://github.com/Alex23582/FileToVideo/assets/117467716/9a20de0e-85f0-4d2d-97ad-1a9563d7bde3)

With many frames combined into a video, the end-result looks like this:
<video src="https://github.com/Alex23582/FileToVideo/assets/117467716/f1878acc-f9de-44fa-8e32-36b0b7c24812"/>
# Binary bits
The data-matrix codes had a significant performance impact. Probably because it's meant to read unaligned camera-pictures. For my purpose that wasn't necessary, so i just wrote the bits "as is" to the images with the idea of maybe utilising this [Reed-Solomon](https://github.com/klauspost/reedsolomon) go library to add the error-correction which is now missing.
You can try this with v2.go using the command-line: ```/goqrfile <encode/decode> <filename> <videofilename>```
<video src="https://github.com/Alex23582/FileToVideo/assets/117467716/c8a1016d-a05b-4883-9899-62436285a222"/>
# Video-Compression
As DvorakDwarf already explained in his repo, the compression of youtube makes this program almost impossible to use. It's definitely not suitable for actual use of file storage. This may better with the implementation of Reed-Solomon codes in the future.

# Known Issues
 - The output-file is padded with 0-bytes, which mostly isn't a problem
