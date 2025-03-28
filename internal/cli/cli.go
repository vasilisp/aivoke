package cli

import (
	"fmt"
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

	prompt, err := prompt.Build(os.Args[1], params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get prompt: %v\n", err)
		os.Exit(1)
	}

	response, err := client.AskGPT(string(prompt), strings.Join(args, " "))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(response)
}
