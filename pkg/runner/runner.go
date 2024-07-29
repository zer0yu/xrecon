package runner

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chainreactors/fingers"
	"github.com/chainreactors/utils/httputils"
)

type Options struct {
	URL          string
	URLFile      string
	OutputFormat string
	OutputFile   string
}

type Runner struct {
	options     *Options
	engine      *fingers.Engine
	input       *InputProcessor
	output      OutputProcessor
	cdnDetector *CDNDetector
}

func NewRunner(options *Options) (*Runner, error) {
	engine, err := fingers.NewEngine()
	if err != nil {
		return nil, fmt.Errorf("error creating fingers engine: %v", err)
	}

	inputProcessor, err := NewInputProcessor(options.URL, options.URLFile)
	if err != nil {
		return nil, fmt.Errorf("error creating input processor: %v", err)
	}

	var outputProcessor OutputProcessor
	switch options.OutputFormat {
	case "csv":
		outputProcessor = NewCSVOutput(options.OutputFile + ".csv")
	case "txt":
		outputProcessor = NewTXTOutput(options.OutputFile + ".txt")
	default:
		outputProcessor = NewTerminalOutput()
	}

	cdnDetector := NewCDNDetector()

	return &Runner{
		options:     options,
		engine:      engine,
		input:       inputProcessor,
		output:      outputProcessor,
		cdnDetector: cdnDetector,
	}, nil
}

func (r *Runner) Run() error {
	defer r.output.Close()

	for url := range r.input.URLs() {
		fingerprint, cdnInfo, err := r.processURL(url)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", url, err)
			continue
		}
		r.output.Write(url, fingerprint, cdnInfo)
	}

	return nil
}

func (r *Runner) processURL(url string) (string, string, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	// 检测 CDN
	isCDN, provider, itemType, err := r.cdnDetector.Detect(strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://"))
	cdnInfo := fmt.Sprintf("CDN: %v, Provider: %s, Type: %s", isCDN, provider, itemType)
	if err != nil {
		cdnInfo = fmt.Sprintf("CDN detection error: %v", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", cdnInfo, err
	}
	defer resp.Body.Close()

	content := httputils.ReadRaw(resp)
	frames, err := r.engine.DetectContent(content)
	if err != nil {
		return "", cdnInfo, err
	}

	return frames.String(), cdnInfo, nil
}
