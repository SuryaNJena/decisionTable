package generate

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	grule "github.com/global-soft-ba/decisionTable/conv/grule/data"
	"github.com/global-soft-ba/decisionTable/conv/grule/generate/grl"
	"github.com/global-soft-ba/decisionTable/conv/grule/generate/json"
	"github.com/global-soft-ba/decisionTable/data/standard"
)

var (
	ErrGruleOutputFormatNotSupported = errors.New("output format not supported")
)

func CreateGruleGenerator() GruleGenerator {
	return GruleGenerator{}
}

type GruleGenerator struct {
	templates *template.Template
	format    grule.OutputFormat
}

type ruleTemplateData struct {
	Rule      grule.Rule // The current rule being processed
	TableKey  string     // The name of the decision table
}

func (g *GruleGenerator) Generate(rules grule.RuleSet, targetFormat string) (interface{}, error) {
	switch targetFormat {
	case string(standard.GRULE):
		tmpl, err := grl.GenerateTemplates(rules.HitPolicy, rules.Interference)
		if err != nil {
			return nil, err
		}
		g.templates = tmpl
		g.format = grule.GRL
		return g.generate(rules)

	case string(grule.JSON):
		tmpl, err := json.GenerateTemplates(rules.HitPolicy, rules.Interference)
		if err != nil {
			return nil, err
		}
		g.templates = tmpl
		g.format = grule.JSON
		return g.generate(rules)

	}
	return nil, ErrGruleOutputFormatNotSupported
}

func (g *GruleGenerator) generate(ruleSet grule.RuleSet) ([]string, error) {
	var result []string
	for _, v := range ruleSet.Rules {
		var tpl bytes.Buffer

		// Create the combined data structure
		templateData := ruleTemplateData{
			Rule:      v,
			TableKey: ruleSet.Key, // Get the table name from the RuleSet
		}

		err := g.templates.Execute(&tpl, templateData)
		if err != nil {
			return []string{}, fmt.Errorf("template execution failed for rule %s: %w", v.Name, err)
		}
		result = append(result, tpl.String())
	}
	return result, nil
}
