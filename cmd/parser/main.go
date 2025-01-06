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
	sender   string
	message  string
	dateTime string
}

/*
 * isNewMessage is needed because of multiline messages
 * isNewMessage checks if a line has the following format
 * 7/2/22, 6:50 AM - User Name: message
 * by removing all spaces and testing against `isMessageRe` regex
 */
func isNewMessage(line string) bool {
	spacesRe := regexp.MustCompile(`\p{Zs}`)
	line = spacesRe.ReplaceAllString(line, "")
	var isMessageRe = `(?m)^\d{1,2}\/\d{1,2}\/\d{2},\d{1,2}:\d{2}\s?(AM|PM)?-[\w\s]+:.+$`
	match, _ := regexp.MatchString(isMessageRe, line)
	return match
}

/*
 * extractMessageFromLine parses a line
 * and returns a Message struct
 */
func extractMessageFromLine(line string) Message {
	parts := strings.Split(line, " - ")
	dateTime := parts[0]
	senderAndMessage := strings.Split(parts[1], ": ")
	sender := senderAndMessage[0]
	message := senderAndMessage[1]
	return Message{sender, message, dateTime}
}

func main() {
	flag.Usage = func() {
		fmt.Println("WhatsApp chat parser.")
		fmt.Println()
		fmt.Println("-file                     Provide the path to the file to process")
	}
	filePath := flag.String("file", "", "Path to the file to process")
	flag.Parse()
	if *filePath == "" {
		fmt.Print("Error: -file flag is required\n")
		os.Exit(1)
	}
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Failed to open file at", *filePath)
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewReader(file)
	var currentMessage Message
	var currentText string
	messageHandler := func() {
		fmt.Printf("Message: %s\n", currentMessage)
	}
	performActionOnEachMessage := func() {
		messageHandler()
		if currentText != "" {
			currentMessage.message = strings.TrimSpace(currentText)
			currentText = ""
		}
	}
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				performActionOnEachMessage()
				fmt.Print("Parsed the whole file.\n")
				os.Exit(1)
			}
			os.Exit(1)
		}
		if isNewMessage(line) {
			performActionOnEachMessage()
			currentMessage = extractMessageFromLine(line)
			currentText = currentMessage.message
			continue
		}
		currentText += line
	}
}
