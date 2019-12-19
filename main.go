package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hokaccha/go-prettyjson"
)

var globalConfig *Config
var schema = NewSchema()

func main() {
	var err error
	if globalConfig, err = GetConfig(); err != nil && globalConfig != nil {
		fmt.Println("Missing config.json file!")
		return
	}

	if err = schema.StreamData(globalConfig); err != nil {
		fmt.Println("Failed to stream data from file, Missing data file")
		return
	}

	fmt.Println("Type help to view list command!")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = RunCommand(cmdString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// RunCommand process command string
func RunCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	if len(arrCommandStr) > 3 {
		arrCommandStr[2] = strings.Join(arrCommandStr[2:], " ")
		arrCommandStr = arrCommandStr[:3]
	} else if len(arrCommandStr) == 0 {
		return nil
	}
	switch arrCommandStr[0] {
	case "exit":
		os.Exit(0)
	case "help":
		ShowLegend()
		return nil
	case "describe":
		return DescribeProcess(schema, arrCommandStr)
	case "table":
		if result, err := SearchProcess(schema, arrCommandStr); err != nil {
			return err
		} else {
			RenderTabularView(arrCommandStr[1], result)
			return nil
		}
	case "search":
		result, err := SearchProcess(schema, arrCommandStr)
		s, _ := prettyjson.Marshal(result)
		fmt.Println(string(s))
		return err
	default:
		return fmt.Errorf("Unsupported command")
	}
	return nil
}

// ShowLegend show help legend
func ShowLegend() {
	fmt.Println("List command: ")
	fmt.Println("+=================================================+")
	fmt.Println("\t- describe [table-name]: show list table search field")
	fmt.Println("\t\tEx: describe user")
	fmt.Println("\t- search [table-name] [field-name]=[field-value]: search table by field name, response by JSON")
	fmt.Println("\t\tEx: search ticket status=pending")
	fmt.Println("\t- table [table-name] [field-name]=[field-value]: search table by field name, return tabular view")
	fmt.Println("\t\tEx: table organization name=Enthaze")
	fmt.Println("+=================================================+")
}
