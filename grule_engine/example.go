package grule_engine

import (
	"encoding/json"
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type Tx struct {
	Test  bool
	Test2 string
}

// 这是一个跑不起来的例子，想要使用这个模块要先学会 Grule，这个模块只是一个工具类。
func main() {

	when := NewWhenOr()
	then := NewThen().
		AddCall("OBJ.Test", &Operator{Const: "参数1"}).
		AddCall("Retract", &Operator{Const: "规则名"})

	gruleJSON := pkg.GruleJSON{
		Name:        "规则名",
		Description: "规则说明",
		Salience:    10,
		When:        when,
		Then:        then.Build(),
	}

	gruleSlice := []string{"{\"eq\":[{\"obj\":\"Type\"},{\"const\":5}]}"}
	for _, grule := range gruleSlice {
		operator := &LogicalOperator{}
		_ = json.Unmarshal([]byte(grule), operator)
		when.AddLogicalOperator(operator)
	}

	marshal, _ := json.Marshal(gruleJSON)

	knowledgeLibrary := NewKnowledgeLibrary()
	err := RegisterRule(knowledgeLibrary, marshal, "规则名称", "1.0")
	fmt.Println(err)

	// 规则中使用的变量都要添加，不添加就不能执行。
	dataCtx := ast.NewDataContext()
	dataCtx.Add("OBJ", "这里填写结构体实例")
	dataCtx.Add("Type", 5)
	dataCtx.Add("brightness", 10)
	dataCtx.Add("height_temperature", 2)
	_, err = Execute(knowledgeLibrary, dataCtx, "规则名称", "1.0")
	fmt.Println(err)
}
