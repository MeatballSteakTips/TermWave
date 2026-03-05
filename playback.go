package main

import (
	"fmt"
	"os/exec"
)

var currentPlayer *exec.Cmd

func PlayStream(streamURL string) error {
	StopStream()

	currentPlayer = exec.Command("mpv", "--no-video", streamURL)

	err := currentPlayer.Start()
	if err != nil {
		return fmt.Errorf("MPV failed to start: %w", err)
	}
	return nil
}

func StopStream() {
	if currentPlayer != nil && currentPlayer.Process != nil {
		currentPlayer.Process.Kill()
		currentPlayer.Process.Wait()
		currentPlayer = nil
	}
}
