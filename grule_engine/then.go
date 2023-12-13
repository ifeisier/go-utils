package grule_engine

type ThenSet struct {
	Set []*ArithmeticOperator `json:"set,omitempty"`
}

type ThenCall struct {
	Call []any `json:"call,omitempty"`
}

// ArithmeticOperator 算数运算符
type ArithmeticOperator struct {
	// obj: set 的目标属性
	Field string      `json:"obj,omitempty"`
	Value any         `json:"const,omitempty"`
	BOR   []*Operator `json:"bor,omitempty"`
	BAND  []*Operator `json:"band,omitempty"`
	PLUS  []*Operator `json:"plus,omitempty"`
	MINUS []*Operator `json:"minus,omitempty"`
	DIV   []*Operator `json:"div,omitempty"`
	MUL   []*Operator `json:"mul,omitempty"`
	MOD   []*Operator `json:"mod,omitempty"`
}

// Then 针对 pkg.GruleJSON 中的 then 数据结构
type Then struct {
	setOrCall []any
}

// NewThen 创建 Then
func NewThen() *Then {
	return &Then{}
}

func (build *Then) Build() []any {
	return build.setOrCall
}

////////////////// 设置 Call //////////////////

// AddCall 添加要调用的方法
//
// 根据情况在最后添加 Retract("规则名")
func (build *Then) AddCall(methodName string, parameter ...*Operator) *Then {
	anies := make([]any, 0)
	anies = append(anies, methodName)
	for _, operator := range parameter {
		anies = append(anies, operator)
	}
	build.setOrCall = append(build.setOrCall, ThenCall{anies})
	return build
}

////////////////// 设置 Set //////////////////

// AddSetValueObj 设置对象的某个属性
func (build *Then) AddSetValueObj(field string, lo1 string) *Then {
	build.setOrCall = append(build.setOrCall,
		ThenSet{Set: []*ArithmeticOperator{{Field: field}, {Field: lo1}}})
	return build
}

// AddSetValueConst 设置常量值
func (build *Then) AddSetValueConst(field string, lo1 any) *Then {
	build.setOrCall = append(build.setOrCall,
		ThenSet{Set: []*ArithmeticOperator{{Field: field}, {Value: lo1}}})
	return build
}

// AddSetBor 创建 | 运算
func (build *Then) AddSetBor(field string, lo1, lo2 *Operator) *Then {
	build.setOrCall = append(build.setOrCall,
		ThenSet{Set: []*ArithmeticOperator{{Field: field}, {BOR: []*Operator{lo1, lo2}}}})
	return build
}

// AddSetBand 创建 & 运算
func (build *Then) AddSetBand(field string, lo1, lo2 *Operator) *Then {
	build.setOrCall = append(build.setOrCall,
		ThenSet{Set: []*ArithmeticOperator{{Field: field}, {BAND: []*Operator{lo1, lo2}}}})
	return build
}

// AddSetPlus 创建 + 运算
func (build *Then) AddSetPlus(field string, lo1, lo2 *Operator) *Then {
	build.setOrCall = append(build.setOrCall,
		ThenSet{Set: []*ArithmeticOperator{{Field: field}, {PLUS: []*Operator{lo1, lo2}}}})
	return build
}

// AddSetMinus 创建 - 运算
func (build *Then) AddSetMinus(field string, lo1, lo2 *Operator) *Then {
	build.setOrCall = append(build.setOrCall,
		ThenSet{Set: []*ArithmeticOperator{{Field: field}, {MINUS: []*Operator{lo1, lo2}}}})
	return build
}

// AddSetDiv 创建 / 运算
func (build *Then) AddSetDiv(field string, lo1, lo2 *Operator) *Then {
	build.setOrCall = append(build.setOrCall,
		ThenSet{Set: []*ArithmeticOperator{{Field: field}, {DIV: []*Operator{lo1, lo2}}}})
	return build
}

// AddSetMul 创建 * 运算
func (build *Then) AddSetMul(field string, lo1, lo2 *Operator) *Then {
	build.setOrCall = append(build.setOrCall,
		ThenSet{Set: []*ArithmeticOperator{{Field: field}, {MUL: []*Operator{lo1, lo2}}}})
	return build
}

// AddSetMod 创建 % 运算
func (build *Then) AddSetMod(field string, lo1, lo2 *Operator) *Then {
	build.setOrCall = append(build.setOrCall,
		ThenSet{Set: []*ArithmeticOperator{{Field: field}, {MOD: []*Operator{lo1, lo2}}}})
	return build
}
