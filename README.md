# wpp-chat-parser

Parses the text file generated when you export a whatsapp chat.

## How to run
```console
go run cmd/parser/main.go -file /path/to/file
```

## Description
What a message looks like on the exported file:
```
8/8/24, 8:32 PM - John Doe: I use vim btw.
     TIME          Sender      Message
```

It is extracted to a Message struct:
```go
type Message struct {
	sender   string
	message  string
	dateTime string
}
```

Which will look something like:
```go
Message: {
    John Doe
    I use vim btw.
    8/8/24, 8:32 PM
 }
 ```
 Perform any actions you wish on each message on the **messageHandler**
 ```go
messageHandler := func() {
    fmt.Printf("Message: %s\n", currentMessage)
}
```
> default just prints each Message struct created
