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

type listener struct {
	*parser.BaseKotlinParserListener
	p *Processor
}

func (l *listener) EnterPackageHeader(ctx *parser.PackageHeaderContext) {
	p := l.p
	for i := 0; i < ctx.GetChildCount(); i++ {
		child := ctx.GetChild(i)
		switch val := child.(type) {
		case *parser.IdentifierContext:
			p.metadata.packageName.RawValue += val.GetText()
		}
	}
}

// EnterImportList is called when production importList is entered.
func (l *listener) EnterImportList(ctx *parser.ImportListContext) {
	p := l.p
	logger := p.cfg.logger
	count := ctx.GetChildCount()
	logger.Info("Found imports length %d\n", count)
	p.metadata.imports = make([]KotlinPackageName, 0, count)
}

func (l *listener) EnterImportHeader(ctx *parser.ImportHeaderContext) {
	p := l.p
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
	p.metadata.imports = append(p.metadata.imports, ans)
}

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

func (p *Processor) GetImports() []KotlinPackageName {
	return p.metadata.imports
}

func (p *Processor) GetPackageName() KotlinPackageName {
	return p.metadata.packageName
}
