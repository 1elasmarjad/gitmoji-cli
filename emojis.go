package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

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

		emoji, err := GetEmoji(*msg)

		if err != nil {
			fmt.Println("Error getting emoji: ", err)
			os.Exit(1)
		}

		newMsg := emoji + " " + *msg
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

func GetEmoji(commitMessage string) (string, error) {

	postBody, err := json.Marshal(map[string]string{
		"model": "gpt-4o-mini",
	})

	if err != nil {
		return "", err
	}

	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", responseBody)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("EMOJIS_OPENAI_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	sbdy := string(body)

	// TODO
	fmt.Println(sbdy)

	return "🚀", nil
}
