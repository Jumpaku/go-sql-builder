package sql

import "strings"

func Table(name string) fromItem { return &fromItemTable{name: name} }
func TableAs(name string, alias string) fromItem {
	return &fromItemTable{name: name, fromItemBase: fromItemBase{alias: &alias}}
}
func Sub(stmt SelectStmt) fromItem {
	return &fromItemSub{stmt: stmt}
}
func SubAs(stmt SelectStmt, alias string) fromItem {
	return &fromItemSub{stmt: stmt, fromItemBase: fromItemBase{alias: &alias}}
}

type fromItem interface {
	BuildParams() Params
	BuildTemplate() string
	JoinOn(item fromItem, cond boolExpr) fromItem
	JoinUsing(item fromItem, column string, columns ...string) fromItem
	FullJoinOn(item fromItem, cond boolExpr) fromItem
	FullJoinUsing(item fromItem, column string, columns ...string) fromItem
	LeftJoinOn(item fromItem, cond boolExpr) fromItem
	LeftJoinUsing(item fromItem, column string, columns ...string) fromItem
	RightJoinOn(item fromItem, cond boolExpr) fromItem
	RightJoinUsing(item fromItem, column string, columns ...string) fromItem
	CrossJoin(item fromItem) fromItem
}

type fromItemBase struct {
	alias *string
}

func (t *fromItemBase) getAlias() (alias string, exists bool) {
	if t.alias == nil {
		return "", false
	}
	return *t.alias, true
}

type fromItemTable struct {
	fromItemBase
	name string
}

func (t *fromItemTable) BuildParams() Params { return Params{} }
func (t *fromItemTable) BuildTemplate() string {
	template := t.name
	if alias, ok := t.getAlias(); ok {
		template += " AS " + alias
	}
	return template
}

type fromItemSub struct {
	fromItemBase
	stmt SelectStmt
}

func (t *fromItemSub) BuildParams() Params { return t.stmt.Params }
func (t *fromItemSub) BuildTemplate() string {
	template := "(" + t.stmt.Template + ")"
	if alias, ok := t.getAlias(); ok {
		template += " AS " + alias
	}
	return template
}

type fromItemJoin struct {
	fromItemBase
	joinOp       JoinOp
	leftItem     fromItem
	rightItem    fromItem
	onCondition  boolExpr
	usingColumns []string
}

func (t *fromItemJoin) BuildParams() Params {
	var params Params
	params.Merge(t.leftItem.BuildParams())
	params.Merge(t.rightItem.BuildParams())
	if t.joinOp != JoinOp_CROSS && t.onCondition != nil {
		params.Merge(t.onCondition.BuildParams())
	}
	return params
}
func (t *fromItemJoin) BuildTemplate() string {
	template := t.leftItem.BuildTemplate() + " " + string(t.joinOp) + " JOIN " + t.rightItem.BuildTemplate()
	if t.joinOp != JoinOp_CROSS && t.onCondition != nil {
		template += " ON " + t.onCondition.BuildTemplate()
	}
	if t.joinOp != JoinOp_CROSS && len(t.usingColumns) > 0 {
		template += " USING (" + strings.Join(t.usingColumns, ", ") + ")"
	}
	template = "(" + template + ")"
	if alias, ok := t.getAlias(); ok {
		template += " AS " + alias
	}
	return template
}

type JoinOp string

const (
	JoinOp_CROSS       JoinOp = "CROSS"
	JoinOp_INNER       JoinOp = "INNER"
	JoinOp_FULL_OUTER  JoinOp = "FULL OUTER"
	JoinOp_LEFT_OUTER  JoinOp = "LEFT OUTER"
	JoinOp_RIGHT_OUTER JoinOp = "RIGHT OUTER"
)

func (i *fromItemTable) JoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_INNER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemTable) JoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_INNER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemTable) FullJoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_FULL_OUTER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemTable) FullJoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_FULL_OUTER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemTable) LeftJoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_LEFT_OUTER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemTable) LeftJoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_LEFT_OUTER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemTable) RightJoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_RIGHT_OUTER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemTable) RightJoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_RIGHT_OUTER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemTable) CrossJoin(item fromItem) fromItem {
	return &fromItemJoin{joinOp: JoinOp_CROSS, leftItem: i, rightItem: item}
}

func (i *fromItemSub) JoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_INNER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemSub) JoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_INNER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemSub) FullJoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_FULL_OUTER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemSub) FullJoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_FULL_OUTER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemSub) LeftJoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_LEFT_OUTER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemSub) LeftJoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_LEFT_OUTER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemSub) RightJoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_RIGHT_OUTER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemSub) RightJoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_RIGHT_OUTER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemSub) CrossJoin(item fromItem) fromItem {
	return &fromItemJoin{joinOp: JoinOp_CROSS, leftItem: i, rightItem: item}
}

func (i *fromItemJoin) JoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_INNER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemJoin) JoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_INNER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemJoin) FullJoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_FULL_OUTER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemJoin) FullJoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_FULL_OUTER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemJoin) LeftJoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_LEFT_OUTER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemJoin) LeftJoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_LEFT_OUTER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemJoin) RightJoinOn(item fromItem, cond boolExpr) fromItem {
	return &fromItemJoin{joinOp: JoinOp_RIGHT_OUTER, leftItem: i, rightItem: item, onCondition: cond}
}
func (i *fromItemJoin) RightJoinUsing(item fromItem, column string, columns ...string) fromItem {
	return &fromItemJoin{joinOp: JoinOp_RIGHT_OUTER, leftItem: i, rightItem: item, usingColumns: append([]string{column}, columns...)}
}
func (i *fromItemJoin) CrossJoin(item fromItem) fromItem {
	return &fromItemJoin{joinOp: JoinOp_CROSS, leftItem: i, rightItem: item}
}
