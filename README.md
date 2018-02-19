## videoScale
#### Template for using ffmpeg and ffprobe to scale videofiles and output as .mp4

* CC-BY Roy Dybing, Feb. 2018
* github: rDybing
* slack: rdybing

Written in Go 1.9.2

------

#### v1.2 Update Feb. 19th:

Added Scale and Crop in one go, which saves one the quality degration of encoding twice if first cropping and then scaling.

Whilst on the subject, optimized the scaling routine slightly. First check if even or odd, if odd make even. Second, if not divisible by 16, add two.

------

#### v1.1 Update Feb. 16th:

Testing on Win10-WSL was not too happy about the aac encoder, as the WSL Ubuntu 16.04 distro of ffmpeg apparently is a bit dated and considers this 'experimental' unlike the ffmpeg distro in Ubuntu 17.10. Changed ffmpeg to force using this encoder however. Works fine.

------

Need ffmpeg installed - this includes necessary encoders and ffprobe utility. To install on Linux or Win10-WSL (Debian/Ubuntu):

> ~$ sudo apt-get install ffmpeg

* To run from included Linux AMD64 binary on Linux or Win10-WSL navigate to the videoScale folder and then:

> ~$ ./vidScale

* Included in this repo is one video-file for demonstration-purposes, 'in.mp4'

------

App intended to demonstrate how to use ffmpeg to scale and keep resolution of video divisible by 16 for compatibility purposes. I probably could write this in C#, but it'd take twice as long (not my 'native' language) and it's getting late :)

Incoming video will be first checked to see if it's resolution is 512x512 - which all new videos produced by the client app should be. If not however, as is the case with our legacy database of videos, video will be made to be 512px high, and a new width is calculated to keep the aspect of the original video more or less intact. All legacy videos are in portrait mode, hence why height will determine width and not the other way around (or check which dimension to use).

It should be able to handle most video-formats as input, including .mov (though not if QT encoded IIRC due to Apple tight-assedness) but I've not tested. Throw in a few video-files in different formats into the project folder, and go wild :)

In the end a new .mp4 file is produced with the new resolution, and the app exits. Assuming no errors.

-----

