package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ClearConsole() {
	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
	case "windows":
		clearCmd := exec.Command("cmd", "/c", "cls")
		clearCmd.Stdout = os.Stdout
		if err := clearCmd.Run(); err != nil {
			fmt.Println(err)
		}
	}
}
