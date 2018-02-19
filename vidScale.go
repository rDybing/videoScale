/*****************************
 * videoScale
 * CC-BY Roy Dybing, Feb. 2018
 * github: rDybing
 * slack:  rdybing
 *****************************/

/*****************************
 * v1.2 Update Feb. 19th:
 *
 * Added Scale and Crop in one go, which saves one the quality
 * degration of encoding twice if first cropping and then
 * scaling.
 *
 * Whilst on the subject, optimized the scaling routine
 * slightly. First check if even or odd, if odd make even.
 * Second, if not divisible by 16, add two.
 *****************************/

/*****************************
 * v1.1 Update Feb. 16th:
 *
 * Testing on Win10-WSL was not too happy about the aac encoder,
 * as the WSL Ubuntu 16.04 distro of ffmpeg apparently is a bit
 * dated and considers this 'experimental' unlike the ffmpeg
 * distro in Ubuntu 17.10. Changed ffmpeg to force using this
 * encoder however. Works fine.
 *****************************/
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type video_t struct {
	width  int
	height int
}

func main() {
	var ok bool
	var cropOK bool
	var inFile string

	for ok == false {
		inFile = getInput("Name of file to scale:")
		if _, err := os.Stat(inFile); err != nil {
			fmt.Println("No such file, try again...")
			ok = false
		} else {
			ok = true
		}
	}

	cropInput := getInput("Crop file to square aspect (y/Y):")
	if cropInput == "y" || cropInput == "Y" {
		cropOK = true
	} else {
		cropOK = false
	}

	oldVid := getDimensions(inFile)
	fmt.Printf("old - w: %4d :: h: %4d\n", oldVid.width, oldVid.height)

	if cropOK {
		outFile := getInput("Save as (.mp4 will be added):")
		outFile += ".mp4"
		scaleAndCropNewFile(inFile, outFile)
	} else {
		newVid := calcNewSize(oldVid)
		fmt.Printf("new - w: %4d :: h: %4d\n", newVid.width, newVid.height)
		outFile := getInput("Save as (.mp4 will be added):")
		outFile += ".mp4"
		scaleNewFile(inFile, outFile, newVid)
	}
}

func getDimensions(inFile string) video_t {
	fmt.Println("Getting Dimensions")
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "stream=width,height", "-of", "default=noprint_wrappers=1", inFile)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running ffprobe: %v\n", err)
	}
	tStr := string(cmdOutput.Bytes())
	return cleanString(tStr)
}

func scaleAndCropNewFile(inFile string, outFile string) {
	var scaleCmd = `scale=(iw*sar)*max(512/(iw*sar)\,512/ih):ih*max(512/(iw*sar)\,512/ih), crop=512:512`
	cmd := exec.Command("ffmpeg", "-i", inFile, "-f", "mp4", "-c:v", "libx264", "-r", "30", "-vf", scaleCmd, "-c:a", "aac", "-strict", "-2", "-b:a", "128k", "-ar", "44100", outFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running ffmpeg (scale and crop): %v\n%s\n", err, stderr.String())
	}
	fmt.Println("Success Scale and Crop!")
}

func scaleNewFile(inFile string, outFile string, vid video_t) {
	outSize := fmt.Sprintf("%d:%d", vid.width, vid.height)
	cmd := exec.Command("ffmpeg", "-i", inFile, "-f", "mp4", "-c:v", "libx264", "-r", "30", "-s:v", outSize, "-c:a", "aac", "-strict", "-2", "-b:a", "128k", "-ar", "44100", outFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running ffmpeg (scale): %v\n%s\n", err, stderr.String())
	}
	fmt.Println("Success Scale!")
}

func calcNewSize(in video_t) video_t {
	var out video_t
	out.height = 512
	scaleValue := float32(out.height) / float32(in.height)
	out.width = int(float32(in.width) * scaleValue)
	if out.width%2 != 0 {
		out.width++
	}
	for out.width%16 != 0 {
		out.width += 2
	}
	return out
}

func getInput(helpText string) string {
	var input string
	fmt.Println(helpText)
	fmt.Scanf("%s\n", &input)
	return input
}

func cleanString(s string) video_t {
	var vid video_t
	s = strings.Replace(s, "width=", "", -1)
	s = strings.Replace(s, "height=", "", -1)
	result := strings.Split(s, "\n")
	w, err := strconv.Atoi(result[0])
	if err != nil {
		log.Fatalf("Error converting string line 1: %v\n", err)
	}
	h, err := strconv.Atoi(result[1])
	if err != nil {
		log.Fatalf("Error converting string line 2: %v\n", err)
	}
	vid.width = w
	vid.height = h
	return vid
}
