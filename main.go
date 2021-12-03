package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/LoganFievez/yjq/config"
	"github.com/icza/dyno"

	"gopkg.in/yaml.v2"
)

const (
	NO_ERR         = iota
	NO_ARG         = iota
	FILE_NOT_EXIST = iota
	NOT_YML_YAML   = iota
	ONLY_ONE       = iota
)

/*
 * POUR LANCER UN DEBUG:
 * #1: dlv debug --headless --listen=:2345 --log --api-version=2 -- assets/*.yml
 * #2: Lancer le debug vscode qui est paramétré
 */

func main() {
	info, _ := os.Stdin.Stat()
	var (
		output []byte
		err    error
	)

	// fmt.Println(info.Mode())

	if (info.Mode() & os.ModeCharDevice) != 0 {
		config := handleArgs(false)
		// Pas en mode pipe
		output, err = os.ReadFile(config.Filename)
		if err != nil {
			panic(err)
		}

	} else {
		// en mode pipe
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadBytes('\n')
			output = append(output, input...)

			if err != nil && err == io.EOF {
				break
			}
		}
	}

	var out interface{}
	err = yaml.Unmarshal(output, &out)
	if err != nil {
		panic(err)
	}

	out = dyno.ConvertMapI2MapS(out)

	var converted []byte

	converted, err = json.Marshal(out)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", converted)
	os.Exit(NO_ERR)
}

func handleArgs(pipemode bool) config.Config {
	config := config.New()
	config.Pipemode = pipemode
	if len(os.Args) == 1 {
		printErrorMsg()
		os.Exit(NO_ARG)
	} else {
		if match, _ := regexp.MatchString(`^-{1}\bh\b$|^-{2}\bhelp\b$`, os.Args[1]); match {
			printErrorMsg()
			os.Exit(NO_ERR)
		}
		config.Filename = os.Args[len(os.Args)-1]
		if _, err := os.Stat(config.Filename); err != nil {
			if os.IsNotExist(err) {
				printErrorMsg(fmt.Sprintf("%s does not exist.", config.Filename))
				os.Exit(FILE_NOT_EXIST)
			}
		}
		if match, _ := regexp.MatchString(`\*+`, config.Filename); match {
			printErrorMsg("Can only handle one file at a time.")
			os.Exit(ONLY_ONE)
		}
		if match, _ := regexp.MatchString(`\.(?i)(yml|yaml)$`, config.Filename); !match {
			printErrorMsg("Not a YML file.", "Not a YAML file.")
			os.Exit(NOT_YML_YAML)
		}
	}
	return config
}

func printErrorMsg(errs ...string) {
	fmt.Fprintf(os.Stdout, "yjq is a tool for converting a YML or a YAML into JSON.\n")
	fmt.Fprintf(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "Usage:\n")
	fmt.Fprintf(os.Stdout, "\tyjq [options] <yml_file>\n")
	fmt.Fprintf(os.Stdout, "\t<yml_content> | yjq\n")
	fmt.Fprintf(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "The options are:\n")
	fmt.Fprintf(os.Stdout, "\t--help -h\tuse this to show the help\n")
	fmt.Fprintf(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "Notes:\n")
	fmt.Fprintf(os.Stdout, "\tIf multiple files are passed to yjq only the last one is used.\n")
	fmt.Fprintf(os.Stdout, "\n")
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "\033[31mError(s):\n\033[0m")
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "\033[31m\t%s\n\033[0m", err)
		}
	}
}
