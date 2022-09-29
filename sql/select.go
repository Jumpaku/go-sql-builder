package sql

type Param struct {
	name  *string
	value any
}

func (p Param) Name() (string, bool) {
	if p.name == nil {
		return "", false
	}
	return *p.name, true
}

func (p Param) Value() any {
	return p.value
}

type Params struct {
	params []Param
}

func (p *Params) Append(values ...any) {
	for _, value := range values {
		p.params = append(p.params, Param{value: value})
	}
}
func (p *Params) AppendWithName(name string, value any) {
	p.params = append(p.params, Param{name: &name, value: value})
}
func (p *Params) Merge(params Params) {
	p.params = append(p.params, params.params...)
}

type SelectStmt struct {
	Template string
	Params   Params
}
type FromBuilder interface {
	From(items ...FromItem) WhereBuilder
}
type WhereBuilder interface {
	GroupByBuilder
	Where(expr boolExpr, exprs ...boolExpr) GroupByBuilder
}
type GroupByBuilder interface {
	OrderByBuilder
	GroupBy(expr Expr) OrderByBuilder
	GroupByHaving(expr Expr, cond boolExpr) OrderByBuilder
}
type OrderByBuilder interface {
	LimitOffsetBuilder
	OrderBy(expr Expr) LimitOffsetBuilder
	OrderByDesc(expr Expr) LimitOffsetBuilder
}
type LimitOffsetBuilder interface {
	SelectBuilder
	Limit(limit int) SelectBuilder
	LimitOffset(limit int, offset int) SelectBuilder
}
type SelectBuilder interface {
	Build() SelectStmt
}

type GroupBy struct {
	Expr   Expr
	Having boolExpr
}

type Order string

const (
	Order_ASC  Order = "ASC"
	Order_DESC Order = "DESC"
)

type OrderBy struct {
	Expr  Expr
	Order Order
}
type LimitOffset struct {
	Limit  int
	Offset int
}
type selectBuilder struct {
	columns     []string
	from        []FromItem
	where       boolExpr
	groupBy     *GroupBy
	orderBy     *OrderBy
	limitOffset *LimitOffset
}

func Select(column string, columns ...string) FromBuilder {
	return &selectBuilder{
		columns: append([]string{column}, columns...),
	}
}
func SelectAll() FromBuilder {
	return &selectBuilder{
		columns: ([]string{"*"}),
	}
}
func (stmt *selectBuilder) From(items ...FromItem) WhereBuilder {
	stmt.from = items
	return stmt
}
func (stmt *selectBuilder) Where(expr boolExpr, exprs ...boolExpr) GroupByBuilder {
	stmt.where = expr
	return stmt
}
func (stmt *selectBuilder) GroupBy(expr Expr) OrderByBuilder {
	return stmt.GroupByHaving(expr, nil)
}
func (stmt *selectBuilder) GroupByHaving(expr Expr, cond boolExpr) OrderByBuilder {
	stmt.groupBy = &GroupBy{Expr: expr, Having: cond}
	return stmt
}
func (stmt *selectBuilder) OrderBy(expr Expr) LimitOffsetBuilder {
	stmt.orderBy = &OrderBy{Expr: expr, Order: Order_ASC}
	return stmt
}
func (stmt *selectBuilder) OrderByDesc(expr Expr) LimitOffsetBuilder {
	stmt.orderBy = &OrderBy{Expr: expr, Order: Order_DESC}
	return stmt
}
func (stmt *selectBuilder) Limit(limit int) SelectBuilder {
	return stmt.LimitOffset(limit, 0)
}
func (stmt *selectBuilder) LimitOffset(limit int, offset int) SelectBuilder {
	stmt.limitOffset = &LimitOffset{Limit: limit, Offset: offset}
	return stmt
}
func (stmt *selectBuilder) Build() SelectStmt { return SelectStmt{} }
