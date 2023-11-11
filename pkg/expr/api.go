package expr

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// EvalBooleanExpression evaluates a boolean expression using CEL (e.g. "1 == 1") and returns the result
func EvalBooleanExpression(expression string, context map[string]interface{}) (bool, error) {
	// empty expression always evaluates to false
	if expression == "" {
		return false, nil
	}

	// init cel go environment
	var exprDecl []*exprpb.Decl
	for key, value := range context {
		switch v := value.(type) {
		case int:
			exprDecl = append(exprDecl, decls.NewVar(key, decls.Int))
		case string:
			exprDecl = append(exprDecl, decls.NewVar(key, decls.String))
		case []string:
			exprDecl = append(exprDecl, decls.NewVar(key, decls.NewListType(decls.String)))
		case map[string]string:
			exprDecl = append(exprDecl, decls.NewVar(key, decls.NewMapType(decls.String, decls.String)))
		default:
			return false, fmt.Errorf("unsupported context value type: %T", v)
		}
	}

	// generate cel evaluation environment
	options := append([]cel.EnvOption{cel.Declarations(exprDecl...)}, additionalFunctions...)
	celConfig, err := cel.NewEnv(options...)
	if err != nil {
		return false, fmt.Errorf("failed to create cel environment: %w", err)
	}

	// prepare program for evaluation
	ast, issues := celConfig.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return false, fmt.Errorf("failed to compile expression: %w", issues.Err())
	}
	prg, err := celConfig.Program(ast)
	if err != nil {
		return false, fmt.Errorf("failed to construct program: %w", err)
	}

	// evaluate
	execOut, _, err := prg.Eval(context)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate expression. expr: %s, error: %w", expression, err)
	}

	// check result
	if execOut.Type() != types.BoolType {
		return false, fmt.Errorf("expression did not evaluate to boolean. expr: %s, type: %s", expression, execOut.Type())
	}

	return execOut.Value() == true, nil
}

// EvaluateRules will check all rules and returns the count of matching rules
func EvaluateRules(rules []string, evalContext map[string]interface{}) int {
	result := 0

	for _, rule := range rules {
		if EvaluateRule(rule, evalContext) {
			result++
		}
	}

	return result
}

// EvaluateRule will evaluate a WorkflowRule and return the result
func EvaluateRule(rule string, evalContext map[string]interface{}) bool {
	match, err := EvalBooleanExpression(rule, evalContext)
	if err != nil {
		return false
	}

	return match
}
