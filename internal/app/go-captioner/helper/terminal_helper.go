package helper

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GetCommandPrintable(cmd *exec.Cmd) string {
	return strings.Join(cmd.Args, " ")
}

// WaitForOptionalEdits pauses for up to 60 seconds asking the user if they want to make edits.
// - If the user answers "n"/"no": it returns immediately (skip edits).
// - If the user answers "y"/"yes": it waits until the user presses Enter again to continue.
// - If no answer is provided within 60 seconds: it returns (skip edits).
// It returns true if user chose to make edits, false otherwise, and an error on I/O issues.
func WaitForOptionalEdits() (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Pause 60s: Make edits now? [y/N]: ")
	answerCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		text, err := reader.ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		answerCh <- strings.TrimSpace(strings.ToLower(text))
	}()

	const secondsInMin = 60

	select {
	case <-time.After(secondsInMin * time.Second):
		fmt.Println() // move to next line after timeout
		return false, nil
	case err := <-errCh:
		return false, err
	case ans := <-answerCh:
		if ans == "y" || ans == "yes" {
			fmt.Print("Editing... Press Enter to continue: ")
			// Block until an empty line (Enter) is submitted.
			_, err := reader.ReadString('\n')
			if err != nil {
				return true, err
			}
			return true, nil
		}
		return false, nil
	}
}
