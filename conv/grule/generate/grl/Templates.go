package grl

import (
	"text/template"

	"github.com/global-soft-ba/decisionTable/conv/grule/data"
	"github.com/global-soft-ba/decisionTable/data/hitPolicy"
)

func GenerateTemplates(hp hitPolicy.HitPolicy, interference bool) (*template.Template, error) {
	var t *template.Template

	t, err := template.New(RULE).Funcs(
		template.FuncMap{
			"getFormat": func() data.OutputFormat {
				return data.GRL
			},
		},
	).Parse(RULE)
	if err != nil {
		return &template.Template{}, err
	}

	_, err = t.New(RULENAME).Parse(RULENAME)
	if err != nil {
		return &template.Template{}, err
	}

	_, err = t.New(WHEN).Parse(WHEN)
	if err != nil {
		return &template.Template{}, err
	}

	_, err = t.New(THEN).Parse(THEN)
	if err != nil {
		return &template.Template{}, err
	}

	_, err = t.New(ENTRY).Parse(ENTRY)
	if err != nil {
		return &template.Template{}, err
	}

	_, err = t.New(SALIENCE).Parse(generatePolicyTemplate(hp))
	if err != nil {
		return &template.Template{}, err
	}

	_, err = t.New(INTERFERENCE).Parse(generateInterferenceTemplate(interference))
	if err != nil {
		return &template.Template{}, err
	}

	return t, nil
}

func generatePolicyTemplate(hp hitPolicy.HitPolicy) string {
	switch hp {
	case hitPolicy.Unique:
		return HitPolicyUnique
	case hitPolicy.First:
		return HitPolicyFirst
	case hitPolicy.Priority:
		return HitPolicyPriority
	default:
		return HitPolicyDefault
	}
}

func generateInterferenceTemplate(interference bool) string {
	switch interference {
	case true:
		return InterferenceExists
	default:
		return InterferenceNotExists
	}
}

/*
We use nested templates, so we walk along the internal grule data model


	rule <RuleName> <RuleDescription> [salience <priority>] {
    when
        <boolean expression>
    then
        <assignment or operation expression>
}

*/

const RULE = `rule {{ template "RULENAME" . }} {{template "SALIENCE" . }} {
 when {{ template "WHEN" .}}
 then {{ template "THEN" .}}
 {{ template "INTERFERENCE" }}
}`

const RULENAME = `{{define "RULENAME"}}{{.TableKey}}_row_{{ .Rule.Name }} "{{ .Rule.Annotation }}"{{end}}`

const WHEN = `{{define "WHEN" }}
{{- range $index, $val := .Rule.Expressions }}
{{- if eq $index 0}}
	{{template "ENTRY" $val}}
	{{- else}}
	&& {{template "ENTRY" $val}}
{{- end}}
{{- end}}
{{- end}}`

const THEN = `{{define "THEN"}}
 {{- range $index, $val := .Rule.Assignments }}
	{{template "ENTRY" $val}};
 {{- end}}
{{- end}}`

const ENTRY = `{{define "ENTRY"}}{{.Expression.Convert getFormat }}{{end}}`

// HitPolicies
/* We assume that the table has unique non overlapping conditions. So maximal hit conditions can be one (or nothing).
   We cannot check the uniqueness. In case that the table is not unique, multi hits are solved with a priority hit policy
   (highest priority wins).
*/
/* Would be correct for real unique => No Salience needed */
/* const UNIQUE = `{{define "SALIENCE"}} {{end}}` */

// UNIQUE Case with fallback in case of overlapping expressions
const (
	SALIENCE          = `SALIENCE`
	HitPolicyDefault  = `{{define "SALIENCE"}}{{end}}`
	HitPolicyUnique   = `{{define "SALIENCE"}}salience {{.Rule.Salience}}{{end}}`
	HitPolicyFirst    = `{{define "SALIENCE"}}salience {{.Rule.InvSalience}}{{end}}`
	HitPolicyPriority = `{{define "SALIENCE"}}salience {{.Rule.Salience}}{{end}}`
)

const (
	INTERFERENCE          = `INTERFERENCE`
	InterferenceExists    = `{{define "INTERFERENCE"}}{{end}}`
	InterferenceNotExists = `{{define "INTERFERENCE"}}Complete();{{end}}`
)
