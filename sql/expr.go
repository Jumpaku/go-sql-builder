package sql

type expr interface {
	BuildParams() Params
	BuildTemplate() string
}
type simpleTerm struct {
	template string
	params   Params
}

func (e *simpleTerm) BuildParams() Params   { return e.params }
func (e *simpleTerm) BuildTemplate() string { return e.template }

type predOp string

const (
	predOp_AND predOp = "AND"
	predOp_OR  predOp = "OR"
)

type pred interface {
	BuildParams() Params
	BuildTemplate() string
}
type predImpl struct {
	negated bool
	pred    pred
}

func (e *predImpl) BuildParams() Params   { return e.pred.BuildParams() }
func (e *predImpl) BuildTemplate() string { return "NOT ( " + e.pred.BuildTemplate() + ")" }

type predTerms struct {
	predOp predOp
	terms  []pred
}

func (e *predTerms) BuildParams() Params {
	var params Params
	for _, term := range e.terms {
		params.Merge(term.BuildParams())
	}
	return params
}
func (e *predTerms) BuildTemplate() string {
	template := ""
	for i, term := range e.terms {
		if i > 0 {
			template += " " + string(e.predOp) + " "
		}
		template += term.BuildTemplate()
	}
	return template
}

type simplePred struct {
	template string
	params   Params
}

func (e *simplePred) BuildParams() Params   { return e.params }
func (e *simplePred) BuildTemplate() string { return e.template }

func Pred(template string, params ...Param) pred {
	return &simplePred{template: template, params: Params{params: params}}
}
func True() pred         { return Pred("TRUE") }
func False() pred        { return Pred("FALSE") }
func Not(pred pred) pred { return &predImpl{negated: true, pred: pred} }
func And(left pred, right pred, rest ...pred) pred {
	return &predTerms{predOp: predOp_AND, terms: append([]pred{left, right}, rest...)}
}
func Or(left pred, right pred, rest ...pred) pred {
	return &predTerms{predOp: predOp_OR, terms: append([]pred{left, right}, rest...)}
}
func Eq(target string, param any) pred {
	return compareTerm(compareOp_EQ, target, param)
}
func Neq(target string, param any) pred {
	return compareTerm(compareOp_NEQ, target, param)
}
func Lt(target string, param any) pred {
	return compareTerm(compareOp_LT, target, param)
}
func Leq(target string, param any) pred {
	return compareTerm(compareOp_LEQ, target, param)
}
func Gt(target string, param any) pred {
	return compareTerm(compareOp_GT, target, param)
}
func Geq(target string, param any) pred {
	return compareTerm(compareOp_GEQ, target, param)
}
func Like(text string, pattern string) pred {
	return Pred(text + " LIKE " + pattern)
}
func IsNull(target string) pred {
	return Pred(target + " IS NULL")
}
func IsNotNull(target string) pred {
	return Pred(target + " IS NOT NULL")
}
func In(target string, values ...any) pred {
	template := target + " IN ("
	params := []Param{}
	for i, v := range values {
		if i > 0 {
			template += ","
		}
		template += "?"
		params = append(params, Param{value: v})

	}
	return Pred(template, params...)
}
func Between(target string, fromValue any, toValue any) pred {
	return Pred(target + " IS NOT NULL")
}

type simpleExpr struct {
	template string
	params   Params
}

func (e *simpleExpr) BuildParams() Params   { return e.params }
func (e *simpleExpr) BuildTemplate() string { return e.template }

type compareOp int

const (
	compareOp_EQ compareOp = iota
	compareOp_NEQ
	compareOp_LT
	compareOp_LEQ
	compareOp_GT
	compareOp_GEQ
)

func compareTerm(op compareOp, target string, param any) pred {
	template := target
	switch op {
	case compareOp_EQ:
		template += " = "
	case compareOp_NEQ:
		template += " <> "
	case compareOp_LT:
		template += " < "
	case compareOp_LEQ:
		template += " <= "
	case compareOp_GT:
		template += " > "
	case compareOp_GEQ:
		template += " >= "
	}
	template += "?"
	return Pred(template, Param{value: param})
}
