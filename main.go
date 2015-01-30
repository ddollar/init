package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func reap() {
	syscall.Wait4(-1, nil, syscall.WNOHANG, &syscall.Rusage{})
}

func terminate() {
	syscall.Kill(-1, syscall.SIGTERM)
	time.Sleep(5 * time.Second)
	syscall.Kill(-1, syscall.SIGKILL)
	os.Exit(0)
}

func handleSignal(sig os.Signal, handler func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig)

	for _ = range ch {
		handler()
	}
}

func main() {
	go handleSignal(syscall.SIGCHLD, reap)
	go handleSignal(syscall.SIGTERM, terminate)
	go handleSignal(syscall.SIGKILL, terminate)

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("usage: init <command>")
		os.Exit(1)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}
