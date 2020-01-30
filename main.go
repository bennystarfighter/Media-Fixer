package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var filesToRemove []string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("I need a filepath.")
		return
	}

	workingPath := filepath.Dir(os.Args[1])

	filename := os.Args[1]
	fmt.Println("Arg:", filename)

	// Extract Audio
	tmpAudioPath := workingPath + "/" + randStringRunes(12) + ".mkv"
	filesToRemove = append(filesToRemove, tmpAudioPath)

	cmd := exec.Command("ffmpeg", "-y", "-i", filename, "-q:a", "0", "-map", "0:a", "-c:a", "libmp3lame", tmpAudioPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		cleanup()
	}

	var extension = filepath.Ext(filename)
	var newName = filename[0:len(filename)-len(extension)] + ".BENJI.mkv"

	// Replace audio
	cmd = exec.Command("ffmpeg", "-y", "-i", filename, "-i", tmpAudioPath, "-map", "0:v", "-map", "1:a", "-map", "0:s?", "-c", "copy", "-c:s", "copy", newName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		cleanup()
	}
	cleanup()
}

func cleanup() {
	for _, e := range filesToRemove {
		err := os.Remove(e)
		if err != nil {
			log.Println(err)
		}
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
