package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Message struct {
	Sender   string
	Message  string
	DateTime string
}

// isNewMessage verifies if a line is the beginning of a new message.
func isNewMessage(line string) bool {
	spacesRegex := regexp.MustCompile(`\p{Zs}`)
	line = spacesRegex.ReplaceAllString(line, "")
	isMessageRegex := `(?m)^\d{1,2}\/\d{1,2}\/\d{2},\d{1,2}:\d{2}\s?(AM|PM)?-[\w\s]+:.+$`
	match, _ := regexp.MatchString(isMessageRegex, line)
	return match
}

// extractMessageFromLine analyzes a line and returns a Message struct
func extractMessageFromLine(line string) (Message, error) {
	parts := strings.SplitN(line, " - ", 2)
	if len(parts) < 2 {
		return Message{}, fmt.Errorf("malformed line: %s", line)
	}
	dateTime := parts[0]
	senderAndMessage := strings.SplitN(parts[1], ": ", 2)
	if len(senderAndMessage) < 2 {
		return Message{}, fmt.Errorf("malformed line (no sender or message): %s", line)
	}
	return Message{
		Sender:   senderAndMessage[0],
		Message:  senderAndMessage[1],
		DateTime: dateTime,
	}, nil
}

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "Path to the file to be processed")

	flag.Usage = func() {
		fmt.Println("WhatsApp chat parser.")
		fmt.Println()
		fmt.Println("-file                     Provide the path to the file to analyzed")
	}

	flag.Parse()
	if filePath == "" {
		fmt.Fprintln(os.Stderr, "Error: The -file parameter is required.")
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error trying to open file '%s': %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
	var currentMessage Message
	var currentText strings.Builder

	processMessage := func() {
		if currentText.Len() > 0 {
			currentMessage.Message = strings.TrimSpace(currentText.String())
			fmt.Printf("Date/Time: %s\nSender: %s\nMessage: %s\n\n", currentMessage.DateTime, currentMessage.Sender, currentMessage.Message)
			currentText.Reset()
		}
	}
	for {
		line, err := scanner.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				processMessage()
				break
			}
			fmt.Fprintf(os.Stderr, "Erro ao ler linha: %v\n", err)
			break
		}
		if isNewMessage(line) {
			processMessage()
			currentMessage, err = extractMessageFromLine(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error analyzing line: %v\n", err)
			}
			continue
		}
		currentText.WriteString(line + "\n")
	}
	fmt.Println("File processed successfully.")
}
