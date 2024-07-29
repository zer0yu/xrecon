package runner

import (
	"encoding/csv"
	"fmt"
	"os"
)

type OutputProcessor interface {
	Write(url, fingerprint, cdnInfo string) error
	Close() error
}

type CSVOutput struct {
	writer *csv.Writer
	file   *os.File
}

func NewCSVOutput(filename string) *CSVOutput {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(file)
	writer.Write([]string{"URL", "Fingerprint", "CDN Info"})

	return &CSVOutput{
		writer: writer,
		file:   file,
	}
}

func (c *CSVOutput) Write(url, fingerprint, cdnInfo string) error {
	return c.writer.Write([]string{url, fingerprint, cdnInfo})
}

func (c *CSVOutput) Close() error {
	c.writer.Flush()
	return c.file.Close()
}

type TXTOutput struct {
	file *os.File
}

func NewTXTOutput(filename string) *TXTOutput {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	return &TXTOutput{
		file: file,
	}
}

func (t *TXTOutput) Write(url, fingerprint, cdnInfo string) error {
	_, err := fmt.Fprintf(t.file, "URL: %s\nFingerprint: %s\nCDN Info: %s\n\n", url, fingerprint, cdnInfo)
	return err
}

func (t *TXTOutput) Close() error {
	return t.file.Close()
}

type TerminalOutput struct{}

func NewTerminalOutput() *TerminalOutput {
	return &TerminalOutput{}
}

func (t *TerminalOutput) Write(url, fingerprint, cdnInfo string) error {
	fmt.Printf("URL: %s\nFingerprint: %s\nCDN Info: %s\n\n", url, fingerprint, cdnInfo)
	return nil
}

func (t *TerminalOutput) Close() error {
	return nil
}
