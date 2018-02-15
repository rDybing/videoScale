/*****************************
 * piFanControl
 * CCBY Roy Dybing, Feb. 2017
 * github.com/rDybing
 *****************************/
package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type state_t struct {
	tempC    int
	fanOn    bool
	limitOn  int
	limitOff int
}

var fanPin int

func main() {
	origW, origH := getDimensions()
	fmt.Printf("w: %3d :: h: %3d\n", origW, origH)
}

func getDimensions() (int, int) {
	fmt.Println("Getting Dimensions")
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "stream=width,height", "-of", "default=noprint_wrappers=1", "in.mp4")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running ffprobe: %v\n", err)
	}
	tStr := string(cmdOutput.Bytes())
	width, height := cleanString(tStr)
	return width, height
}

func cleanString(s string) (int, int) {
	s = strings.Replace(s, "width=", "", -1)
	s = strings.Replace(s, "height=", "", -1)
	//s = strings.Replace(s, "\n", "", -1)
	result := strings.Split(s, "\n")
	width, err := strconv.Atoi(result[0])
	height, err := strconv.Atoi(result[1])
	if err != nil {
		log.Fatalf("Error converting string: %v\n", err)
	}
	return width, height
}
