package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vasilisp/aivoke/internal/openai"
	"github.com/vasilisp/aivoke/internal/prompt"
	"github.com/vasilisp/aivoke/internal/util"
)

func Main() {
	key, err := util.GetOpenAIKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get OpenAI key: %v\n", err)
		os.Exit(1)
	}

	client := openai.NewClient(key)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <prompt-id> <message>\n", os.Args[0])
		os.Exit(1)
	}

	params, args := util.ParseArgs(os.Args[2:])

	show := false
	if value, ok := params["show"]; ok {
		show = (value != "false" && value != "0")
		delete(params, "show")
	}

	stdin := false
	if len(args) == 0 {
		stdin = true
	}

	pr, err := prompt.Build(os.Args[1], params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get prompt: %v\n", err)
		os.Exit(1)
	}

	if show {
		fmt.Print(string(pr.Content))
		os.Exit(0)
	}

	var userPrompt []byte
	if stdin {
		userPrompt, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading stdin:", err)
			os.Exit(1)
		}
	} else {
		userPrompt = []byte(strings.Join(args, " "))
	}

	response, err := client.AskGPT(string(pr.Content), string(userPrompt))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get response: %v\n", err)
		os.Exit(1)
	}

	if pr.Config.Postprocess {
		response = prompt.Postprocess(response)
	}

	fmt.Println(response)
}
