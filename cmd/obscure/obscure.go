package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"vimagination.zapto.org/obscure"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (err error) {
	var (
		key, input, output string
		decode             bool
		inputFile          = os.Stdin
		outputFile         = os.Stdout
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\nThis program generates a cipher, using the key as a random seed, and encodes the given file to the output.\n\nOnly letters and numbers are encoded, punctuation is not modified.\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&key, "k", "", "key to generate seed from")
	flag.StringVar(&input, "i", "-", "input file")
	flag.StringVar(&output, "o", "-", "output file")
	flag.BoolVar(&decode, "d", false, "decode file")
	flag.Parse()

	if key == "" {
		return ErrNoKey
	}

	if input == "" {
		return ErrNoInput
	} else if input != "-" {
		inputFile, err = os.Open(input)
		if err != nil {
			return fmt.Errorf("error opening input file: %w", err)
		}

		defer inputFile.Close()
	}

	if output == "" {
		return ErrNoOutput
	} else if output != "-" {
		outputFile, err = os.Create(output)
		if err != nil {
			return fmt.Errorf("error opening output file: %w", err)
		}

		defer func() {
			if errr := outputFile.Close(); err == nil {
				err = errr
			}
		}()
	}

	_, err = io.Copy(outputFile, obscure.NewEncoder([]byte(key), inputFile, decode))

	return err
}

var (
	ErrNoKey    = errors.New("no key specified")
	ErrNoInput  = errors.New("no input file specified")
	ErrNoOutput = errors.New("no output file specified")
)
