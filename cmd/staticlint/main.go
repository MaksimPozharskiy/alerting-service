package main

import (
	"alerting-service/cmd/staticlint/noosexit"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers,
		printf.Analyzer,
		shadow.Analyzer,
		loopclosure.Analyzer,
	)

	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[:2] != "SA" {
			analyzers = append(analyzers, a.Analyzer)
			break
		}
	}

	analyzers = append(analyzers, staticcheck.Analyzers[0].Analyzer)

	analyzers = append(analyzers, noosexit.Analyzer)

	multichecker.Main(analyzers...)
}
