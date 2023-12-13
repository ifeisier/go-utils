package grule_engine

type Operator struct {
	Obj   string `json:"obj,omitempty"`
	Const any    `json:"const,omitempty"`
}

// NewOperatorsObj 创建对象属性
func NewOperatorsObj(s string) *Operator {
	return &Operator{Obj: s}
}

// NewOperatorsConst 创建常量
func NewOperatorsConst(c any) *Operator {
	return &Operator{Const: c}
}

//func newLogicalOperator(lo1, lo2 any) []*model.Operator {
//	operators := make([]*model.Operator, 0)
//
//	switch lo1.(type) {
//	case *model.Operator_Obj:
//		operators = append(operators, &model.Operator{Value: lo1.(*model.Operator_Obj)})
//	case *model.Operator_Const:
//		operators = append(operators, &model.Operator{Value: lo1.(*model.Operator_Const)})
//	}
//
//	switch lo2.(type) {
//	case *model.Operator_Obj:
//		operators = append(operators, &model.Operator{Value: lo2.(*model.Operator_Obj)})
//	case *model.Operator_Const:
//		operators = append(operators, &model.Operator{Value: lo2.(*model.Operator_Const)})
//	}
//
//	return operators
//}
