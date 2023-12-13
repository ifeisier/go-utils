package grule_engine

// LogicalOperator 逻辑运算符
type LogicalOperator struct {
	EQ  []*Operator `json:"eq,omitempty"`
	NOT []*Operator `json:"not,omitempty"`
	GT  []*Operator `json:"gt,omitempty"`
	GTE []*Operator `json:"gte,omitempty"`
	LT  []*Operator `json:"lt,omitempty"`
	LTE []*Operator `json:"lte,omitempty"`
}

// When 针对 pkg.GruleJSON 中的 when 数据结构
type When struct {
	And []*LogicalOperator `json:"and,omitempty"`
	Or  []*LogicalOperator `json:"or,omitempty"`
}

// NewWhenAnd 针对 when and
func NewWhenAnd() *When {
	return &When{And: make([]*LogicalOperator, 0)}
}

// NewWhenOr 针对 when or
func NewWhenOr() *When {
	return &When{Or: make([]*LogicalOperator, 0)}
}

// //////////////// 添加逻辑运算符 //////////////////

// AddLogicalOperator 直接将 LogicalOperator 添加到 When
func (build *When) AddLogicalOperator(lo *LogicalOperator) *When {
	return build.set(lo)
}

// AddEq 添加 when 中的 Eq(==) 逻辑运算
func (build *When) AddEq(lo1, lo2 *Operator) *When {
	return build.set(&LogicalOperator{EQ: []*Operator{lo1, lo2}})
}

// AddNot 添加 when 中的 Not(!=) 逻辑运算
func (build *When) AddNot(lo1, lo2 *Operator) *When {
	return build.set(&LogicalOperator{NOT: []*Operator{lo1, lo2}})
}

// AddGt 添加 when 中的 Gt(>) 逻辑运算
func (build *When) AddGt(lo1, lo2 *Operator) *When {
	return build.set(&LogicalOperator{GT: []*Operator{lo1, lo2}})
}

// AddGte 添加 when 中的 Gte(>=) 逻辑运算
func (build *When) AddGte(lo1, lo2 *Operator) *When {
	return build.set(&LogicalOperator{GTE: []*Operator{lo1, lo2}})
}

// AddLt 添加 when 中的 Lt(<) 逻辑运算
func (build *When) AddLt(lo1, lo2 *Operator) *When {
	return build.set(&LogicalOperator{LT: []*Operator{lo1, lo2}})
}

// AddLte 添加 when 中的 Lte(<=) 逻辑运算
func (build *When) AddLte(lo1, lo2 *Operator) *When {
	return build.set(&LogicalOperator{LTE: []*Operator{lo1, lo2}})
}

func (build *When) set(lo *LogicalOperator) *When {
	if build.And != nil {
		build.And = append(build.And, lo)
	} else {
		build.Or = append(build.Or, lo)
	}
	return build
}
