package main

import (
	"os/exec"
)

func i3lock(color string) {
	cmd := exec.Command("i3lock", "-c", color, "-p", "win", "-u", "-n")
	go cmd.Run()
}
