// Oh My Molik is an experimental, read-only coding assistant.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/MoLik-Labs/FCC/ohmyfox/agent"
)

func main() {
	endpoint := flag.String("endpoint", "", "Full HTTP(S) Chat Completions URL, e.g. https://gateway.example/v1/chat/completions")
	model := flag.String("model", "", "Model ID recognized by the endpoint")
	workspace := flag.String("workspace", ".", "Workspace directory to inspect (read-only)")
	maxTurns := flag.Int("max-turns", 8, "Maximum model turns (1-32)")
	promptStdin := flag.Bool("prompt-stdin", false, "Read prompt from standard input instead of command-line arguments")
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Oh My Molik (experimental read-only preview)")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: ohmyfox --endpoint URL --model MODEL [options] 'prompt'")
		fmt.Fprintln(flag.CommandLine.Output(), "Optional bearer token: OMF_API_KEY environment variable")
		flag.PrintDefaults()
	}
	flag.Parse()
	config := agent.Config{Endpoint: *endpoint, Model: *model, APIKey: os.Getenv("OMF_API_KEY"), Workspace: *workspace, MaxTurns: *maxTurns}
	fox, err := agent.New(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ohmyfox:", err)
		os.Exit(2)
	}
	prompt := strings.Join(flag.Args(), " ")
	if *promptStdin {
		if len(flag.Args()) != 0 {
			fmt.Fprintln(os.Stderr, "ohmyfox: do not combine --prompt-stdin with a command-line prompt")
			os.Exit(2)
		}
		data, err := io.ReadAll(io.LimitReader(os.Stdin, (64<<10)+1))
		if err != nil {
			fmt.Fprintln(os.Stderr, "ohmyfox: cannot read prompt")
			os.Exit(1)
		}
		prompt = string(data)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	result, err := fox.Run(ctx, prompt)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ohmyfox:", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
