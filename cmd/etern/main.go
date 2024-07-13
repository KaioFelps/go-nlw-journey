package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"os/exec"
)

func main() {
	dotenvErr := godotenv.Load()

	if dotenvErr != nil {
		fmt.Println("Error loading .env file")
	}

	args := os.Args[1:]

	cmd := exec.Command("tern", args...)
	cmd.Env = append(os.Environ())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command and capture the output
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
