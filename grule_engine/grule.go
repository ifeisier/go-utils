package grule_engine

import (
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

// NewKnowledgeLibrary 创建 KnowledgeLibrary
func NewKnowledgeLibrary() *ast.KnowledgeLibrary {
	return ast.NewKnowledgeLibrary()
}

// RegisterRule 注册规则
func RegisterRule(kl *ast.KnowledgeLibrary, jsonBytes []byte, name, version string) error {
	ruleBuilder := builder.NewRuleBuilder(kl)
	rs, _ := pkg.NewJSONResourceFromResource(pkg.NewBytesResource(jsonBytes))
	err := ruleBuilder.BuildRuleFromResource(name, version, rs)
	if err != nil {
		return err
	}
	return nil
}

// RemoveRuleEntry 从规则引擎中移除规则
func RemoveRuleEntry(kl *ast.KnowledgeLibrary, ruleName, name string, version string) {
	kl.RemoveRuleEntry(ruleName, name, version)
}

// Execute 执行规则
//
// dataCtx 可以使用 ast.NewDataContext() 创建
func Execute(kl *ast.KnowledgeLibrary, dataCtx ast.IDataContext, name, version string) (*ast.KnowledgeBase, error) {
	knowledgeBase, err := kl.NewKnowledgeBaseInstance(name, version)
	if err != nil {
		return knowledgeBase, err
	}

	e := engine.NewGruleEngine()
	err = e.Execute(dataCtx, knowledgeBase)
	if err != nil {
		return knowledgeBase, err
	}

	return knowledgeBase, nil
}
