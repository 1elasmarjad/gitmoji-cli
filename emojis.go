package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// check all args
	if len(os.Args) <= 1 {
		fmt.Print("Usage: gitmoji <cmd>")
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
			fmt.Print("Usage: gitmoji commit -m <msg> (ERR1)")
			os.Exit(1)
		}

		// does something exist after?
		if len(os.Args) <= targetIdx+1 {
			fmt.Print("Usage: gitmoji commit -m <msg> (ERR2)")
			os.Exit(1)
		}

		// get the message
		msgIdx := targetIdx + 1

		InputArgs := os.Args

		var msg *string = &InputArgs[msgIdx]

		if strings.HasPrefix(*msg, "-") {
			fmt.Print("Usage: gitmoji commit -m <msg> (ERR3)")
			os.Exit(1)
		}

		emoji, err := GetEmoji(*msg)

		if err != nil {
			fmt.Print("Error getting emoji: ", err)
			os.Exit(1)
		}

		newMsg := emoji + " " + *msg
		msg = &newMsg

		InputArgs[msgIdx] = *msg

		// remove first arg
		InputArgs = InputArgs[1:]

		cmd := exec.Command("git", InputArgs...)

		output, err := cmd.Output()

		if err != nil {
			fmt.Print("Error running git command:", err)
			os.Exit(1)
		}

		fmt.Println("Running git command: ", cmd.String())
		fmt.Print(string(output))
		os.Exit(0)

	case "help":
		fmt.Print("Usage: gitmoji <cmd>\n\n")
		fmt.Print("Commands:\n")
		fmt.Print("  commit -m <msg>  -  Add an emoji to a git commit message\n")
		fmt.Print("  help             -  Show this help message\n\n")
		os.Exit(0)
	default:
		fmt.Print("Invalid command")
		os.Exit(1)
	}

}

type ChatGptResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}
	} `json:"choices"`
}

func GetEmoji(commitMessage string) (string, error) {

	postBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "The user will give you a git-commit message, and you will provide an emoji to prepend to the message to make it more fun! Only send one emoji as the response, and make sure it's appropriate for a commit message.",
			},
			{
				"role":    "user",
				"content": commitMessage,
			},
		},
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

	var chatGptResponse ChatGptResponse

	err = json.Unmarshal(body, &chatGptResponse)

	if err != nil {
		return "", err
	}

	if (len(chatGptResponse.Choices) == 0) || (len(chatGptResponse.Choices[0].Message.Content) == 0) {
		return "", fmt.Errorf("no emoji returned")
	}

	// return the emoji
	return chatGptResponse.Choices[0].Message.Content, nil
}
