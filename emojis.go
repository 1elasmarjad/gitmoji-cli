package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	// check all args
	if len(os.Args) <= 1 {
		fmt.Println("Usage: emojis <cmd>")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "commit":

		// the index of the "-m" text in the args
		var targetIdx = -1

		for i, v := range os.Args {
			if v == "-m" {
				targetIdx = i
				break
			}
		}

		if targetIdx == -1 {
			fmt.Println("Usage: emojis commit -m <msg> (ERR1)")
			os.Exit(1)
		}

		// does something exist after?
		if len(os.Args) <= targetIdx+1 {
			fmt.Println("Usage: emojis commit -m <msg> (ERR2)")
			os.Exit(1)
		}

		// get the message
		msgIdx := targetIdx + 1

		InputArgs := os.Args

		var msg *string = &InputArgs[msgIdx]

		if strings.HasPrefix(*msg, "-") {
			fmt.Println("Usage: emojis commit -m <msg> (ERR3)")
			os.Exit(1)
		}

		newMsg := "🚀 " + *msg
		msg = &newMsg

		InputArgs[msgIdx] = `"` + *msg + `"`

		// remove first arg
		InputArgs = InputArgs[1:]

		fmt.Println("git", strings.Join(InputArgs, " "))

	default:
		fmt.Println("Invalid command")
		os.Exit(1)
	}

}
