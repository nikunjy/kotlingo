package kotlingo

import (
	"fmt"
	"strings"

	"github.com/nikunjy/kotlingo/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type KotlinPackageName struct {
	RawValue string

	ImportsAll bool
}

func (k KotlinPackageName) Parts() []string {
	return strings.Split(k.RawValue, ".")
}

type importListener struct {
	*parser.BaseKotlinParserListener
	logger Logger

	Imports []KotlinPackageName
}

// EnterImportList is called when production importList is entered.
func (s *importListener) EnterImportList(ctx *parser.ImportListContext) {
	count := ctx.GetChildCount()
	s.logger.Info("Found imports length %d\n", count)
	s.Imports = make([]KotlinPackageName, 0, count)
}

func (s *importListener) EnterImportHeader(ctx *parser.ImportHeaderContext) {
	count := ctx.GetChildCount()
	var ans KotlinPackageName
	for i := 0; i < count; i++ {
		child := ctx.GetChild(i)
		switch val := child.(type) {
		case *parser.IdentifierContext:
			ans.RawValue += val.GetText()
		case *antlr.TerminalNodeImpl:
			if val.GetSymbol().GetTokenType() == parser.KotlinLexerIMPORT {
				continue
			}
			if val.GetSymbol().GetTokenType() == parser.KotlinLexerMULT {
				ans.ImportsAll = true
				continue
			}
			ans.RawValue += val.GetText()
		}
	}
	s.Imports = append(s.Imports, ans)
}

// ExitImportList is called when production importList is exited.
func (s *importListener) ExitImportList(ctx *parser.ImportListContext) {}

type errorListener struct {
	err error
}

func (el *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	if e != nil {
		el.err = fmt.Errorf("syntax error %w", e.GetMessage())
	}
}

func (el *errorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}

func (el *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {

}

func (el *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
}

func (p *Processor) GetImports() ([]KotlinPackageName, error) {
	p.el.err = nil
	logger := p.cfg.logger
	listener := &importListener{
		logger: logger,
	}
	antlr.ParseTreeWalkerDefault.Walk(listener, p.file)
	return listener.Imports, p.el.err
}
