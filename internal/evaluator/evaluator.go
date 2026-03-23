package evaluator

import (
	"fmt"
	"html"
	"reflect"
	"strings"

	"template-wisp/internal/ast"
	"template-wisp/internal/lexer"
	"template-wisp/internal/parser"
	"template-wisp/internal/resolver"
	"template-wisp/internal/scope"
)

// Loop signal types
const (
	signalNone     = 0
	signalBreak    = 1
	signalContinue = 2
)

// SafeString is a type that bypasses HTML auto-escaping.
type SafeString struct {
	Value string
}

// Evaluator handles evaluation of AST nodes.
type Evaluator struct {
	resolver        *resolver.Resolver
	scope           *scope.Scope
	output          strings.Builder
	signal          int                               // loop control signal (break/continue)
	templateFn      func(name string) (string, error) // function to load templates
	blocks          map[string]*ast.Program           // block overrides from child templates
	autoEscape      bool                              // whether to auto-escape HTML in output
	maxIter         int                               // max loop iterations (0 = unlimited)
	iterCount       int                               // current iteration count
	includeChain    map[string]bool                   // tracks active includes to detect circular references
	maxIncludeDepth int                               // max include depth (0 = 100 default)
	includeDepth    int                               // current include depth
}

const defaultMaxIncludeDepth = 100

// NewEvaluator creates a new evaluator with the given scope.
func NewEvaluator(s *scope.Scope) *Evaluator {
	return &Evaluator{
		resolver:        resolver.NewResolver(s),
		scope:           s,
		autoEscape:      true, // auto-escape by default for safety
		includeChain:    make(map[string]bool),
		maxIncludeDepth: defaultMaxIncludeDepth,
	}
}

// SetAutoEscape enables or disables HTML auto-escaping.
func (e *Evaluator) SetAutoEscape(enabled bool) {
	e.autoEscape = enabled
}

// SetMaxIterations sets the maximum number of loop iterations allowed.
func (e *Evaluator) SetMaxIterations(max int) {
	e.maxIter = max
}

// Evaluate evaluates a program and returns the output.
func (e *Evaluator) Evaluate(program *ast.Program) (string, error) {
	// Check if the program contains an extends statement
	for _, stmt := range program.Statements {
		if extendsStmt, ok := stmt.(*ast.ExtendsStatement); ok {
			err := e.EvaluateWithExtends(program, extendsStmt)
			if err != nil {
				return "", err
			}
			return e.output.String(), nil
		}
	}

	// Evaluate statements, handling blocks specially
	for i := 0; i < len(program.Statements); i++ {
		stmt := program.Statements[i]

		// Handle block tags specially - need to skip default body when overridden
		if blockStmt, ok := stmt.(*ast.BlockTagStatement); ok && blockStmt.Name != nil {
			blockName := blockStmt.Name.Value

			// Check if there's an override
			if e.blocks != nil {
				if override, ok := e.blocks[blockName]; ok {
					// Evaluate the override
					childEval := NewEvaluator(e.scope)
					childEval.templateFn = e.templateFn
					childEval.blocks = e.blocks
					result, err := childEval.Evaluate(override)
					if err != nil {
						return "", err
					}
					e.output.WriteString(result)

					// Skip default body until EndStatement
					for i+1 < len(program.Statements) {
						i++
						if _, isEnd := program.Statements[i].(*ast.EndStatement); isEnd {
							break
						}
					}
					continue
				}
			}

			// No override - evaluate default body until EndStatement
			for i+1 < len(program.Statements) {
				i++
				if _, isEnd := program.Statements[i].(*ast.EndStatement); isEnd {
					break
				}
				err := e.EvaluateStatement(program.Statements[i])
				if err != nil {
					return "", err
				}
			}
			continue
		}

		// Handle content tags specially
		if _, isContent := stmt.(*ast.ContentStatement); isContent {
			if e.blocks != nil {
				if override, ok := e.blocks["content"]; ok {
					childEval := NewEvaluator(e.scope)
					childEval.templateFn = e.templateFn
					childEval.blocks = e.blocks
					result, err := childEval.Evaluate(override)
					if err != nil {
						return "", err
					}
					e.output.WriteString(result)
					continue
				}
			}
			// No content override - skip
			continue
		}

		err := e.EvaluateStatement(stmt)
		if err != nil {
			return "", err
		}
	}
	return e.output.String(), nil
}

// EvaluateStatement evaluates a statement.
func (e *Evaluator) EvaluateStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		return e.EvaluateExpressionStatement(s)
	case *ast.LetStatement:
		return e.EvaluateLetStatement(s)
	case *ast.AssignStatement:
		return e.EvaluateAssignStatement(s)
	case *ast.IfStatement:
		return e.EvaluateIfStatement(s)
	case *ast.UnlessStatement:
		return e.EvaluateUnlessStatement(s)
	case *ast.ForStatement:
		return e.EvaluateForStatement(s)
	case *ast.WhileStatement:
		return e.EvaluateWhileStatement(s)
	case *ast.RangeStatement:
		return e.EvaluateRangeStatement(s)
	case *ast.WithStatement:
		return e.EvaluateWithStatement(s)
	case *ast.CaseStatement:
		return e.EvaluateCaseStatement(s)
	case *ast.EndStatement:
		return e.EvaluateEndStatement(s)
	case *ast.ElseStatement:
		return e.EvaluateElseStatement(s)
	case *ast.ElsifStatement:
		return e.EvaluateElsifStatement(s)
	case *ast.WhenStatement:
		return e.EvaluateWhenStatement(s)
	case *ast.BreakStatement:
		return e.EvaluateBreakStatement(s)
	case *ast.ContinueStatement:
		return e.EvaluateContinueStatement(s)
	case *ast.CycleStatement:
		return e.EvaluateCycleStatement(s)
	case *ast.IncrementStatement:
		return e.EvaluateIncrementStatement(s)
	case *ast.DecrementStatement:
		return e.EvaluateDecrementStatement(s)
	case *ast.IncludeStatement:
		return e.EvaluateIncludeStatement(s)
	case *ast.RenderStatement:
		return e.EvaluateRenderStatement(s)
	case *ast.ComponentStatement:
		return e.EvaluateComponentStatement(s)
	case *ast.ExtendsStatement:
		return e.EvaluateExtendsStatement(s)
	case *ast.BlockTagStatement:
		return e.EvaluateBlockTagStatement(s)
	case *ast.ContentStatement:
		return e.EvaluateContentStatement(s)
	case *ast.RawStatement:
		return e.EvaluateRawStatement(s)
	case *ast.CommentStatement:
		return e.EvaluateCommentStatement(s)
	case *ast.TextContent:
		return e.EvaluateTextContent(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

// escapeOutput escapes HTML in a value before outputting.
func (e *Evaluator) escapeOutput(val interface{}) string {
	if !e.autoEscape {
		return fmt.Sprintf("%v", val)
	}
	// SafeString bypasses escaping
	if safe, ok := val.(SafeString); ok {
		return safe.Value
	}
	str := fmt.Sprintf("%v", val)
	return html.EscapeString(str)
}

// EvaluateExpressionStatement evaluates an expression statement.
func (e *Evaluator) EvaluateExpressionStatement(stmt *ast.ExpressionStatement) error {
	val, err := e.resolver.ResolveExpression(stmt.Expression)
	if err != nil {
		return err
	}

	// Convert value to string with auto-escaping
	str := e.escapeOutput(val)

	e.output.WriteString(str)
	return nil
}

// EvaluateLetStatement evaluates a let statement.
func (e *Evaluator) EvaluateLetStatement(stmt *ast.LetStatement) error {
	// Evaluate the value
	val, err := e.resolver.ResolveExpression(stmt.Value)
	if err != nil {
		return err
	}

	// Set the variable in the current scope
	e.scope.Set(stmt.Name.Value, val)
	return nil
}

// EvaluateAssignStatement evaluates an assign statement.
func (e *Evaluator) EvaluateAssignStatement(stmt *ast.AssignStatement) error {
	// Evaluate the value
	val, err := e.resolver.ResolveExpression(stmt.Value)
	if err != nil {
		return err
	}

	// Set the variable in the current scope
	e.scope.Set(stmt.Name.Value, val)
	return nil
}

// EvaluateIfStatement evaluates an if statement.
func (e *Evaluator) EvaluateIfStatement(stmt *ast.IfStatement) error {
	// Evaluate the condition
	condVal, err := e.resolver.ResolveExpression(stmt.Condition)
	if err != nil {
		return err
	}

	// Convert to boolean
	condBool, err := e.resolver.ToBool(condVal)
	if err != nil {
		return err
	}

	// Create a new scope for the if block
	e.scope = scope.NewChildScope(e.scope)
	e.resolver.SetScope(e.scope)

	// Execute the appropriate branch
	if condBool {
		if stmt.Consequence != nil {
			for _, s := range stmt.Consequence.Statements {
				err := e.EvaluateStatement(s)
				if err != nil {
					return err
				}
			}
		}
	} else if stmt.Alternative != nil {
		for _, s := range stmt.Alternative.Statements {
			err := e.EvaluateStatement(s)
			if err != nil {
				return err
			}
		}
	}

	// Restore the parent scope
	e.scope = e.scope.Parent()
	e.resolver.SetScope(e.scope)

	return nil
}

// EvaluateUnlessStatement evaluates an unless statement.
func (e *Evaluator) EvaluateUnlessStatement(stmt *ast.UnlessStatement) error {
	// Evaluate the condition
	condVal, err := e.resolver.ResolveExpression(stmt.Condition)
	if err != nil {
		return err
	}

	// Convert to boolean
	condBool, err := e.resolver.ToBool(condVal)
	if err != nil {
		return err
	}

	// Create a new scope for the unless block
	e.scope = scope.NewChildScope(e.scope)
	e.resolver.SetScope(e.scope)

	// Execute if condition is false
	if !condBool {
		if stmt.Consequence != nil {
			for _, s := range stmt.Consequence.Statements {
				err := e.EvaluateStatement(s)
				if err != nil {
					return err
				}
			}
		}
	} else if stmt.Alternative != nil {
		for _, s := range stmt.Alternative.Statements {
			err := e.EvaluateStatement(s)
			if err != nil {
				return err
			}
		}
	}

	// Restore the parent scope
	e.scope = e.scope.Parent()
	e.resolver.SetScope(e.scope)

	return nil
}

// EvaluateForStatement evaluates a for statement.
func (e *Evaluator) EvaluateForStatement(stmt *ast.ForStatement) error {
	// Evaluate the collection
	collectionVal, err := e.resolver.ResolveExpression(stmt.Collection)
	if err != nil {
		return err
	}

	// Convert to slice
	collection, err := e.ToSlice(collectionVal)
	if err != nil {
		return err
	}

	// Create a new scope for the loop
	e.scope = scope.NewChildScope(e.scope)
	e.resolver.SetScope(e.scope)

	// Iterate over the collection
	for i, item := range collection {
		// Set the loop variable
		if stmt.LoopVar != nil {
			e.scope.Set(stmt.LoopVar.Value, item)
		}

		// Set the index variable if present
		if stmt.IndexVar != nil {
			e.scope.Set(stmt.IndexVar.Value, i)
		}

		// Execute the loop body
		if stmt.Body != nil {
			for _, s := range stmt.Body.Statements {
				err := e.EvaluateStatement(s)
				if err != nil {
					return err
				}
				// Check for break signal
				if e.signal == signalBreak {
					e.signal = signalNone
					goto endLoop
				}
				// Check for continue signal
				if e.signal == signalContinue {
					e.signal = signalNone
					break
				}
			}
		}
	}
endLoop:

	// Restore the parent scope
	e.scope = e.scope.Parent()
	e.resolver.SetScope(e.scope)

	return nil
}

// EvaluateWhileStatement evaluates a while statement.
func (e *Evaluator) EvaluateWhileStatement(stmt *ast.WhileStatement) error {
	// Create a new scope for the loop
	e.scope = scope.NewChildScope(e.scope)
	e.resolver.SetScope(e.scope)

	iterCount := 0
	// Loop while condition is true
	for {
		// Check iteration limit
		iterCount++
		if e.maxIter > 0 && iterCount > e.maxIter {
			e.scope = e.scope.Parent()
			e.resolver.SetScope(e.scope)
			return fmt.Errorf("loop iteration limit exceeded (%d iterations)", e.maxIter)
		}

		// Evaluate the condition
		condVal, err := e.resolver.ResolveExpression(stmt.Condition)
		if err != nil {
			return err
		}

		// Convert to boolean
		condBool, err := e.resolver.ToBool(condVal)
		if err != nil {
			return err
		}

		if !condBool {
			break
		}

		// Execute the loop body
		if stmt.Body != nil {
			for _, s := range stmt.Body.Statements {
				err := e.EvaluateStatement(s)
				if err != nil {
					return err
				}
				// Check for break signal
				if e.signal == signalBreak {
					e.signal = signalNone
					goto endLoop
				}
				// Check for continue signal
				if e.signal == signalContinue {
					e.signal = signalNone
					break
				}
			}
		}
	}
endLoop:

	// Restore the parent scope
	e.scope = e.scope.Parent()
	e.resolver.SetScope(e.scope)

	return nil
}

// EvaluateRangeStatement evaluates a range statement.
func (e *Evaluator) EvaluateRangeStatement(stmt *ast.RangeStatement) error {
	// Evaluate the start value
	startVal, err := e.resolver.ResolveExpression(stmt.Start)
	if err != nil {
		return err
	}

	// Evaluate the end value
	endVal, err := e.resolver.ResolveExpression(stmt.End)
	if err != nil {
		return err
	}

	// Convert to integers
	start, err := e.resolver.ToInt(startVal)
	if err != nil {
		return err
	}

	end, err := e.resolver.ToInt(endVal)
	if err != nil {
		return err
	}

	// Create a new scope for the loop
	e.scope = scope.NewChildScope(e.scope)
	e.resolver.SetScope(e.scope)

	// Iterate over the range
	for i := start; i <= end; i++ {
		// Check iteration limit
		if e.maxIter > 0 && (i-start) >= e.maxIter {
			e.scope = e.scope.Parent()
			e.resolver.SetScope(e.scope)
			return fmt.Errorf("loop iteration limit exceeded (%d iterations)", e.maxIter)
		}

		// Set the loop variable
		e.scope.Set("i", i)

		// Execute the loop body
		if stmt.Body != nil {
			for _, s := range stmt.Body.Statements {
				err := e.EvaluateStatement(s)
				if err != nil {
					return err
				}
				// Check for break signal
				if e.signal == signalBreak {
					e.signal = signalNone
					goto endLoop
				}
				// Check for continue signal
				if e.signal == signalContinue {
					e.signal = signalNone
					break
				}
			}
		}
	}
endLoop:

	// Restore the parent scope
	e.scope = e.scope.Parent()
	e.resolver.SetScope(e.scope)

	return nil
}

// EvaluateWithStatement evaluates a with statement.
func (e *Evaluator) EvaluateWithStatement(stmt *ast.WithStatement) error {
	// Evaluate the source expression
	sourceVal, err := e.resolver.ResolveExpression(stmt.Source)
	if err != nil {
		return err
	}

	// Create a new scope for the with block
	e.scope = scope.NewChildScope(e.scope)
	e.resolver.SetScope(e.scope)

	// Set the target variable
	if stmt.Target != nil {
		e.scope.Set(stmt.Target.Value, sourceVal)
	}

	// Execute the with body
	if stmt.Body != nil {
		for _, s := range stmt.Body.Statements {
			err := e.EvaluateStatement(s)
			if err != nil {
				return err
			}
		}
	}

	// Restore the parent scope
	e.scope = e.scope.Parent()
	e.resolver.SetScope(e.scope)

	return nil
}

// EvaluateCaseStatement evaluates a case statement.
func (e *Evaluator) EvaluateCaseStatement(stmt *ast.CaseStatement) error {
	// Evaluate the case value
	caseVal, err := e.resolver.ResolveExpression(stmt.Value)
	if err != nil {
		return err
	}

	// Convert case value to string for comparison
	caseStr, err := e.resolver.ToString(caseVal)
	if err != nil {
		return err
	}

	// Create a new scope for the case block
	e.scope = scope.NewChildScope(e.scope)
	e.resolver.SetScope(e.scope)
	defer func() {
		e.scope = e.scope.Parent()
		e.resolver.SetScope(e.scope)
	}()

	// Walk through body statements, finding matching when clauses
	matched := false
	defaultBranch := (*ast.BlockStatement)(nil)

	i := 0
	for i < len(stmt.Body.Statements) {
		s := stmt.Body.Statements[i]

		if whenStmt, ok := s.(*ast.WhenStatement); ok {
			// Evaluate the when value
			whenVal, err := e.resolver.ResolveExpression(whenStmt.Value)
			if err != nil {
				return err
			}
			whenStr, err := e.resolver.ToString(whenVal)
			if err != nil {
				return err
			}

			if whenStr == caseStr {
				// Match found - collect body until next when/else
				matched = true
				var bodyStmts []ast.Statement
				i++
				for i < len(stmt.Body.Statements) {
					nextStmt := stmt.Body.Statements[i]
					if _, isWhen := nextStmt.(*ast.WhenStatement); isWhen {
						break
					}
					if _, isElse := nextStmt.(*ast.ElseStatement); isElse {
						break
					}
					bodyStmts = append(bodyStmts, nextStmt)
					i++
				}
				// Evaluate the matched branch
				for _, bs := range bodyStmts {
					if err := e.EvaluateStatement(bs); err != nil {
						return err
					}
				}
				continue
			} else {
				// No match - skip this branch
				i++
				for i < len(stmt.Body.Statements) {
					nextStmt := stmt.Body.Statements[i]
					if _, isWhen := nextStmt.(*ast.WhenStatement); isWhen {
						break
					}
					if _, isElse := nextStmt.(*ast.ElseStatement); isElse {
						break
					}
					i++
				}
				continue
			}
		} else if _, isElse := s.(*ast.ElseStatement); isElse {
			if !matched {
				// No when matched - collect else body
				elseBody := &ast.BlockStatement{}
				i++
				for i < len(stmt.Body.Statements) {
					elseBody.Statements = append(elseBody.Statements, stmt.Body.Statements[i])
					i++
				}
				defaultBranch = elseBody
			}
			break
		} else {
			i++
		}
	}

	// Evaluate default branch if no when matched
	if !matched && defaultBranch != nil {
		for _, s := range defaultBranch.Statements {
			if err := e.EvaluateStatement(s); err != nil {
				return err
			}
		}
	}

	return nil
}

// EvaluateEndStatement evaluates an end statement.
func (e *Evaluator) EvaluateEndStatement(stmt *ast.EndStatement) error {
	// End statements are handled by the parser
	// No action needed during evaluation
	return nil
}

// EvaluateElseStatement evaluates an else statement.
func (e *Evaluator) EvaluateElseStatement(stmt *ast.ElseStatement) error {
	// Else statements are handled by the if/unless evaluation
	// No action needed during evaluation
	return nil
}

// EvaluateElsifStatement evaluates an elsif statement.
func (e *Evaluator) EvaluateElsifStatement(stmt *ast.ElsifStatement) error {
	// Elsif statements are handled by the if evaluation
	// No action needed during evaluation
	return nil
}

// EvaluateWhenStatement evaluates a when statement.
func (e *Evaluator) EvaluateWhenStatement(stmt *ast.WhenStatement) error {
	// When statements are handled by the case evaluation
	// No action needed during evaluation
	return nil
}

// EvaluateBreakStatement evaluates a break statement.
func (e *Evaluator) EvaluateBreakStatement(stmt *ast.BreakStatement) error {
	e.signal = signalBreak
	return nil
}

// EvaluateContinueStatement evaluates a continue statement.
func (e *Evaluator) EvaluateContinueStatement(stmt *ast.ContinueStatement) error {
	e.signal = signalContinue
	return nil
}

// EvaluateCycleStatement evaluates a cycle statement.
func (e *Evaluator) EvaluateCycleStatement(stmt *ast.CycleStatement) error {
	if len(stmt.Values) == 0 {
		return nil
	}

	// Get the cycle counter
	counter, ok := e.scope.Get("__cycle_counter__")
	if !ok {
		counter = 0
	}

	// Get the value at the current counter position
	counterInt, _ := e.resolver.ToInt(counter)
	val, err := e.resolver.ResolveExpression(stmt.Values[counterInt])
	if err != nil {
		return err
	}

	// Increment the counter for next use
	counterInt++
	if counterInt >= len(stmt.Values) {
		counterInt = 0
	}
	e.scope.Set("__cycle_counter__", counterInt)

	// Convert to string and output
	str, err := e.resolver.ToString(val)
	if err != nil {
		return err
	}

	e.output.WriteString(str)
	return nil
}

// EvaluateIncrementStatement evaluates an increment statement.
func (e *Evaluator) EvaluateIncrementStatement(stmt *ast.IncrementStatement) error {
	// Get the current value
	val, ok := e.scope.Get(stmt.Variable.Value)
	if !ok {
		val = 0
	}

	// Increment the value
	valInt, _ := e.resolver.ToInt(val)
	valInt++

	// Set the new value
	e.scope.Set(stmt.Variable.Value, valInt)
	return nil
}

// EvaluateDecrementStatement evaluates a decrement statement.
func (e *Evaluator) EvaluateDecrementStatement(stmt *ast.DecrementStatement) error {
	// Get the current value
	val, ok := e.scope.Get(stmt.Variable.Value)
	if !ok {
		val = 0
	}

	// Decrement the value
	valInt, _ := e.resolver.ToInt(val)
	valInt--

	// Set the new value
	e.scope.Set(stmt.Variable.Value, valInt)
	return nil
}

// EvaluateIncludeStatement evaluates an include statement.
func (e *Evaluator) EvaluateIncludeStatement(stmt *ast.IncludeStatement) error {
	if e.templateFn == nil {
		return fmt.Errorf("include statement requires a template loader")
	}

	// Get the template name
	if stmt.Template == nil {
		return fmt.Errorf("include statement requires a template name")
	}
	templateName := stmt.Template.Value

	// Check for circular include
	if err := e.pushInclude(templateName); err != nil {
		return err
	}
	defer e.popInclude(templateName)

	// Load the template
	templateContent, err := e.templateFn(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template %s: %w", templateName, err)
	}

	// Parse and evaluate the template
	l := lexer.NewLexer(templateContent)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse errors in template %s: %v", templateName, p.Errors())
	}

	// Create a child evaluator with the current scope
	childEval := e.childEvaluator(e.scope)
	result, err := childEval.Evaluate(program)
	if err != nil {
		return err
	}

	e.output.WriteString(result)
	return nil
}

// EvaluateRenderStatement evaluates a render statement.
func (e *Evaluator) EvaluateRenderStatement(stmt *ast.RenderStatement) error {
	if e.templateFn == nil {
		return fmt.Errorf("render statement requires a template loader")
	}

	// Get the template name
	if stmt.Template == nil {
		return fmt.Errorf("render statement requires a template name")
	}
	templateName := stmt.Template.Value

	// Check for circular includes
	if err := e.pushInclude(templateName); err != nil {
		return err
	}
	defer e.popInclude(templateName)

	// Load the template
	templateContent, err := e.templateFn(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template %s: %w", templateName, err)
	}

	// Parse and evaluate the template
	l := lexer.NewLexer(templateContent)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse errors in template %s: %v", templateName, p.Errors())
	}

	// Create a child evaluator with an isolated scope
	childScope := scope.NewIsolatedScope()
	defer childScope.Release()
	childEval := e.childEvaluator(childScope)
	result, err := childEval.Evaluate(program)
	if err != nil {
		return err
	}

	e.output.WriteString(result)
	return nil
}

// EvaluateComponentStatement evaluates a component statement.
// A component is like include but with named props passed to an isolated scope.
func (e *Evaluator) EvaluateComponentStatement(stmt *ast.ComponentStatement) error {
	if e.templateFn == nil {
		return fmt.Errorf("component statement requires a template loader")
	}

	if stmt.Name == nil {
		return fmt.Errorf("component statement requires a component name")
	}
	templateName := stmt.Name.Value

	// Check for circular includes
	if err := e.pushInclude(templateName); err != nil {
		return err
	}
	defer e.popInclude(templateName)

	// Load the component template
	templateContent, err := e.templateFn(templateName)
	if err != nil {
		return fmt.Errorf("failed to load component %s: %w", templateName, err)
	}

	// Parse the component template
	l := lexer.NewLexer(templateContent)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse errors in component %s: %v", templateName, p.Errors())
	}

	// Create an isolated scope for the component
	childScope := scope.NewIsolatedScope()
	defer childScope.Release()

	// Evaluate props and set them in the component scope
	for _, prop := range stmt.Props {
		val, err := e.resolver.ResolveExpression(prop)
		if err != nil {
			return fmt.Errorf("failed to resolve component prop: %w", err)
		}
		// If the prop is an identifier, use their name as the variable name
		if ident, ok := prop.(*ast.Identifier); ok {
			childScope.Set(ident.Value, val)
		} else if dot, ok := prop.(*ast.DotExpression); ok {
			// For dot expressions like .title, use the field name
			childScope.Set(dot.Field.Value, val)
		}
	}

	// Evaluate the component in the isolated scope
	childEval := e.childEvaluator(childScope)
	result, err := childEval.Evaluate(program)
	if err != nil {
		return err
	}

	e.output.WriteString(result)
	return nil
}

// EvaluateExtendsStatement evaluates an extends statement.
// This implements template inheritance: the child template's blocks override
// the parent layout's blocks.
func (e *Evaluator) EvaluateExtendsStatement(stmt *ast.ExtendsStatement) error {
	if e.templateFn == nil {
		return fmt.Errorf("extends statement requires a template loader")
	}

	if stmt.Layout == nil {
		return fmt.Errorf("extends statement requires a layout name")
	}

	return fmt.Errorf("extends must be handled via EvaluateWithExtends")
}

// EvaluateBlockTagStatement evaluates a block tag statement.
// If the evaluator has block overrides from a child template, the override is used.
// Otherwise, the default content between {% block %} and {% end %} is evaluated.
func (e *Evaluator) EvaluateBlockTagStatement(stmt *ast.BlockTagStatement) error {
	if stmt.Name == nil {
		return fmt.Errorf("block statement requires a block name")
	}

	blockName := stmt.Name.Value

	// Check if there's an override for this block from a child template
	if e.blocks != nil {
		if override, ok := e.blocks[blockName]; ok {
			// Evaluate the child's block content
			childEval := NewEvaluator(e.scope)
			childEval.templateFn = e.templateFn
			childEval.blocks = e.blocks
			result, err := childEval.Evaluate(override)
			if err != nil {
				return err
			}
			e.output.WriteString(result)
			return nil
		}
	}

	// No override - evaluate default content until {% end %}
	// The default content is the statements between this block tag and the next end tag
	// This is handled by the main Evaluate loop encountering the EndStatement
	return nil
}

// EvaluateContentStatement evaluates a content statement.
// Content is a special block that represents the main content area in a layout.
func (e *Evaluator) EvaluateContentStatement(stmt *ast.ContentStatement) error {
	// Check if there's a "content" block override
	if e.blocks != nil {
		if override, ok := e.blocks["content"]; ok {
			childEval := NewEvaluator(e.scope)
			childEval.templateFn = e.templateFn
			childEval.blocks = e.blocks
			result, err := childEval.Evaluate(override)
			if err != nil {
				return err
			}
			e.output.WriteString(result)
			return nil
		}
	}

	// No content override - output nothing
	return nil
}

// EvaluateWithExtends evaluates a program that contains an extends statement.
// It first collects block overrides from the child template, then evaluates
// the parent layout with those overrides.
func (e *Evaluator) EvaluateWithExtends(program *ast.Program, extendsStmt *ast.ExtendsStatement) error {
	if e.templateFn == nil {
		return fmt.Errorf("extends requires a template loader")
	}

	// Collect block overrides from the child template
	blockOverrides := make(map[string]*ast.Program)
	e.collectBlockOverrides(program, blockOverrides)

	// Load the parent layout
	layoutName := extendsStmt.Layout.Value
	layoutContent, err := e.templateFn(layoutName)
	if err != nil {
		return fmt.Errorf("failed to load layout %s: %w", layoutName, err)
	}

	// Parse the parent layout
	l := lexer.NewLexer(layoutContent)
	p := parser.NewParser(l)
	parentProgram := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse errors in layout %s: %v", layoutName, p.Errors())
	}

	// Evaluate the parent layout with block overrides
	e.blocks = blockOverrides
	result, err := e.Evaluate(parentProgram)
	if err != nil {
		return err
	}
	e.output.Reset()
	e.output.WriteString(result)
	return nil
}

// collectBlockOverrides walks through the program's statements and collects
// block definitions into the overrides map.
func (e *Evaluator) collectBlockOverrides(program *ast.Program, overrides map[string]*ast.Program) {
	i := 0
	for i < len(program.Statements) {
		stmt := program.Statements[i]

		if blockStmt, ok := stmt.(*ast.BlockTagStatement); ok && blockStmt.Name != nil {
			// Collect statements until we find the matching EndStatement
			blockName := blockStmt.Name.Value
			var bodyStmts []ast.Statement
			i++ // skip the block tag

			for i < len(program.Statements) {
				if _, isEnd := program.Statements[i].(*ast.EndStatement); isEnd {
					i++ // skip the end tag
					break
				}
				bodyStmts = append(bodyStmts, program.Statements[i])
				i++
			}

			overrides[blockName] = &ast.Program{Statements: bodyStmts}
		} else if _, isContent := stmt.(*ast.ContentStatement); isContent {
			// Content statement in child template: collect body until EndStatement
			var bodyStmts []ast.Statement
			i++ // skip the content tag

			for i < len(program.Statements) {
				if _, isEnd := program.Statements[i].(*ast.EndStatement); isEnd {
					i++ // skip the end tag
					break
				}
				bodyStmts = append(bodyStmts, program.Statements[i])
				i++
			}

			overrides["content"] = &ast.Program{Statements: bodyStmts}
		} else {
			i++
		}
	}
}

// EvaluateRawStatement evaluates a raw statement.
func (e *Evaluator) EvaluateRawStatement(stmt *ast.RawStatement) error {
	// In raw mode, output content literally without any processing
	if stmt.Content != "" {
		e.output.WriteString(stmt.Content)
	}
	return nil
}

// EvaluateCommentStatement evaluates a comment statement.
func (e *Evaluator) EvaluateCommentStatement(stmt *ast.CommentStatement) error {
	// Comments are not output
	return nil
}

// EvaluateTextContent evaluates literal text content.
func (e *Evaluator) EvaluateTextContent(stmt *ast.TextContent) error {
	// Output the text content directly
	e.output.WriteString(stmt.Value)
	return nil
}

// ToSlice converts a value to a slice.
func (e *Evaluator) ToSlice(val interface{}) ([]interface{}, error) {
	if val == nil {
		return nil, fmt.Errorf("cannot convert nil to slice")
	}

	// If it's already a slice, return it
	if slice, ok := val.([]interface{}); ok {
		return slice, nil
	}

	// Use reflection to handle typed slices
	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Slice {
		result := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			result[i] = rv.Index(i).Interface()
		}
		return result, nil
	}

	return nil, fmt.Errorf("cannot convert %T to slice", val)
}

// GetOutput returns the current output.
func (e *Evaluator) GetOutput() string {
	return e.output.String()
}

// ClearOutput clears the output.
func (e *Evaluator) ClearOutput() {
	e.output.Reset()
}

// GetScope returns the current scope.
func (e *Evaluator) GetScope() *scope.Scope {
	return e.scope
}

// SetTemplateFn sets the function used to load templates for include/render.
func (e *Evaluator) SetTemplateFn(fn func(name string) (string, error)) {
	e.templateFn = fn
}

// SetMaxIncludeDepth sets the maximum include nesting depth.
func (e *Evaluator) SetMaxIncludeDepth(max int) {
	e.maxIncludeDepth = max
}

// pushInclude checks for circular includes and tracks the include chain.
func (e *Evaluator) pushInclude(templateName string) error {
	if e.includeChain[templateName] {
		return fmt.Errorf("circular include detected: %s is already in the include chain", templateName)
	}
	if e.maxIncludeDepth > 0 && e.includeDepth >= e.maxIncludeDepth {
		return fmt.Errorf("include depth limit exceeded (%d includes)", e.maxIncludeDepth)
	}
	e.includeChain[templateName] = true
	e.includeDepth++
	return nil
}

// popInclude removes a template from the include chain.
func (e *Evaluator) popInclude(templateName string) {
	delete(e.includeChain, templateName)
	e.includeDepth--
}

// childEvaluator creates a child evaluator that inherits include tracking.
func (e *Evaluator) childEvaluator(childScope *scope.Scope) *Evaluator {
	child := NewEvaluator(childScope)
	child.templateFn = e.templateFn
	child.autoEscape = e.autoEscape
	child.maxIter = e.maxIter
	child.blocks = e.blocks
	child.maxIncludeDepth = e.maxIncludeDepth
	child.includeDepth = e.includeDepth
	child.includeChain = e.includeChain
	return child
}
func (e *Evaluator) SetScope(s *scope.Scope) {
	e.scope = s
	e.resolver.SetScope(s)
}
