package sql

func Not(expr boolExpr) boolExpr { return &boolExprImpl{negate: true, terms: []boolExpr{expr}} }
func And(expr ...boolExpr) boolExpr {
	return &boolExprImpl{negate: false, boolOp: boolOp_AND, terms: expr}
}
func Or(expr ...boolExpr) boolExpr {
	return &boolExprImpl{negate: false, boolOp: boolOp_OR, terms: expr}
}
func Eq(target string, param any) boolExpr {
	return newBoolCompareTerm(compareOp_EQ, target, param)
}
func Neq(target string, param any) boolExpr {
	return newBoolCompareTerm(compareOp_NEQ, target, param)
}
func Lt(target string, param any) boolExpr {
	return newBoolCompareTerm(compareOp_LT, target, param)
}
func Leq(target string, param any) boolExpr {
	return newBoolCompareTerm(compareOp_LEQ, target, param)
}
func Gt(target string, param any) boolExpr {
	return newBoolCompareTerm(compareOp_GT, target, param)
}
func Geq(target string, param any) boolExpr {
	return newBoolCompareTerm(compareOp_GEQ, target, param)
}
func Like(text string, pattern string) boolExpr {
	return &likeTerm{text: text, pattern: pattern}
}

type Expr interface {
	BuildParams() Params
	BuildTemplate() string
}
type simpleExpr struct {
	template string
	params   Params
}

func (e *simpleExpr) BuildParams() Params   { return e.params }
func (e *simpleExpr) BuildTemplate() string { return e.template }

type boolExpr interface {
	BuildParams() Params
	BuildTemplate() string
}
type boolOp string

const (
	boolOp_AND boolOp = "AND"
	boolOp_OR  boolOp = "OR"
)

type boolExprImpl struct {
	negate bool
	boolOp boolOp
	terms  []boolExpr
}

func (e *boolExprImpl) BuildParams() Params {
	var params Params
	for _, term := range e.terms {
		params.Merge(term.BuildParams())
	}
	return params
}
func (e *boolExprImpl) BuildTemplate() string {
	template := ""
	for i, term := range e.terms {
		if i > 0 {
			template += " " + string(e.boolOp) + " "
		}
		template += term.BuildTemplate()
	}
	if e.negate {
		template = "NOT (" + template + ")"
	}
	return template
}

type compareOp int

const (
	compareOp_EQ compareOp = iota
	compareOp_NEQ
	compareOp_LT
	compareOp_LEQ
	compareOp_GT
	compareOp_GEQ
)

type boolCompareTerm struct {
	operator compareOp
	left     Expr
	right    Expr
}

func (e *boolCompareTerm) BuildParams() Params {
	var params Params
	params.Append(e.left.BuildParams())
	params.Append(e.right.BuildParams())
	return params
}
func (e *boolCompareTerm) BuildTemplate() string {
	var template string
	template += e.left.BuildTemplate()
	switch e.operator {
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
	template += e.right.BuildTemplate()
	return template
}

func newBoolCompareTerm(op compareOp, target string, param any) boolExpr {
	return &boolCompareTerm{
		operator: op,
		left:     &simpleExpr{template: target},
		right:    &simpleExpr{template: "?", params: Params{params: []Param{{value: param}}}},
	}
}

type likeTerm struct {
	text    string
	pattern string
}

func (e *likeTerm) BuildParams() Params { return Params{} }
func (e *likeTerm) BuildTemplate() string {
	return e.text + " LIKE(" + e.pattern + ")"
}
