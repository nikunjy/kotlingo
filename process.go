package kotlingo

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/nikunjy/kotlingo/parser"
)

type Processor struct {
	file parser.IKotlinFileContext

	el *errorListener

	cfg *Config
}

func NewProcessor(fileName string, opts ...Option) (*Processor, error) {
	cfg := defaultCommonConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	input, err := antlr.NewFileStream(fileName)
	if err != nil {
		return nil, err
	}
	lexer := parser.NewKotlinLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewKotlinParser(stream)
	// Parser prints a lot of errors that we dont' need for our purpose
	p.RemoveErrorListeners()
	el := &errorListener{}
	p.AddErrorListener(el)
	//p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	file := p.KotlinFile()
	return &Processor{
		file: file,
		el:   el,
		cfg:  &cfg,
	}, nil
}
