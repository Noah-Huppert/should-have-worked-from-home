package stdin

import "bufio"
import "os"
import "fmt"
import "log"

// Read retrieves user input from standard in.
func Read() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading user input from stdin: %s",
			err.Error())
	}

	return text, nil
}

// Prompt displays a message and then reads user input.
func Prompt(logger *log.Logger, msg string) (string, error) {
	logger.Printf(msg)
	return Read()
}
