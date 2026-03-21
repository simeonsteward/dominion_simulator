package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/simeonsteward/dominion_simulator/internal/core"
	"github.com/simeonsteward/dominion_simulator/internal/strategy"
)

func prompt(reader *bufio.Reader, msg string) string {
	fmt.Print(msg)
	t, _ := reader.ReadString('\n')
	return strings.TrimSpace(t)
}

func interactiveCreate(name string) error {
	reader := bufio.NewReader(os.Stdin)
	items := []strategy.NameCondition{}
	for {
		card := strings.ToLower(prompt(reader, "card name (or done): "))
		if card == "done" {
			break
		}
		if !core.IsCard(card) {
			fmt.Println("invalid card name")
			continue
		}
		items = append(items, strategy.NameCondition{Card: card, Condition: strategy.ConditionType{"True": []byte("{}")}})
	}
	return strategy.SavePriorityList(items, name)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: strategy-cli [create|validate|migrate] [flags]")
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "create":
		fs := flag.NewFlagSet("create", flag.ExitOnError)
		name := fs.String("name", "new_strategy", "strategy filename")
		input := fs.String("input", "", "optional source json path")
		_ = fs.Parse(os.Args[2:])
		if *input != "" {
			b, err := os.ReadFile(*input)
			if err != nil {
				panic(err)
			}
			if err := os.WriteFile("strategies/"+*name+".json", b, 0o644); err != nil {
				panic(err)
			}
			fmt.Println("strategy copied")
			return
		}
		if err := interactiveCreate(*name); err != nil {
			panic(err)
		}
		fmt.Println("strategy created")
	case "validate":
		fs := flag.NewFlagSet("validate", flag.ExitOnError)
		name := fs.String("name", "", "strategy filename without .json")
		_ = fs.Parse(os.Args[2:])
		if *name == "" {
			panic("name required")
		}
		if _, err := strategy.GetPriorityList(*name); err != nil {
			panic(err)
		}
		fmt.Println("valid")
	case "migrate":
		fs := flag.NewFlagSet("migrate", flag.ExitOnError)
		name := fs.String("name", "", "strategy filename without .json")
		_ = fs.Parse(os.Args[2:])
		if *name == "" {
			panic("name required")
		}
		items, err := strategy.LoadPriorityList(*name)
		if err != nil {
			panic(err)
		}
		if err := strategy.SavePriorityList(items, *name); err != nil {
			panic(err)
		}
		fmt.Println("migrated to v1 formatting")
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}
}
