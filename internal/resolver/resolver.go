package resolver

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"template-wisp/internal/ast"
	"template-wisp/internal/scope"
)

// Resolver handles variable resolution and member access.
type Resolver struct {
	scope *scope.Scope
}

// NewResolver creates a new resolver with the given scope.
func NewResolver(s *scope.Scope) *Resolver {
	return &Resolver{scope: s}
}

// ResolveExpression resolves an expression to a value.
func (r *Resolver) ResolveExpression(expr ast.Expression) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.Identifier:
		return r.ResolveIdentifier(e)
	case *ast.DotExpression:
		return r.ResolveDotExpression(e)
	case *ast.PipeExpression:
		return r.ResolvePipeExpression(e)
	case *ast.IntegerLiteral:
		return e.Value, nil
	case *ast.Boolean:
		return e.Value, nil
	case *ast.StringLiteral:
		return e.Value, nil
	case *ast.InfixExpression:
		return r.ResolveInfixExpression(e)
	case *ast.PrefixExpression:
		return r.ResolvePrefixExpression(e)
	case *ast.IndexExpression:
		return r.ResolveIndexExpression(e)
	case *ast.ArrayLiteral:
		return r.ResolveArrayLiteral(e)
	case *ast.HashLiteral:
		return r.ResolveHashLiteral(e)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// ResolveIdentifier resolves an identifier to a value.
func (r *Resolver) ResolveIdentifier(ident *ast.Identifier) (interface{}, error) {
	val, ok := r.scope.Get(ident.Value)
	if !ok {
		return nil, fmt.Errorf("undefined variable: %s", ident.Value)
	}
	return val, nil
}

// ResolveDotExpression resolves a dot expression (member access).
func (r *Resolver) ResolveDotExpression(dot *ast.DotExpression) (interface{}, error) {
	// First, resolve the base field
	baseVal, err := r.ResolveIdentifier(dot.Field)
	if err != nil {
		return nil, err
	}

	// Then, traverse the chain
	current := baseVal
	for _, field := range dot.Chain {
		current, err = r.AccessMember(current, field.Value)
		if err != nil {
			return nil, err
		}
	}

	return current, nil
}

// AccessMember accesses a member of an object (struct field, map key, etc.).
func (r *Resolver) AccessMember(obj interface{}, member string) (interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("cannot access member %s of nil", member)
	}

	val := reflect.ValueOf(obj)

	// Handle pointer types
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("cannot access member %s of nil pointer", member)
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		// Try to access struct field
		field := val.FieldByName(member)
		if !field.IsValid() {
			// Try case-insensitive match
			for i := 0; i < val.NumField(); i++ {
				if strings.EqualFold(val.Type().Field(i).Name, member) {
					field = val.Field(i)
					break
				}
			}
		}
		if !field.IsValid() {
			return nil, fmt.Errorf("struct has no field %s", member)
		}
		if !field.CanInterface() {
			return nil, fmt.Errorf("cannot access unexported field %s", member)
		}
		return field.Interface(), nil

	case reflect.Map:
		// Try to access map key
		key := reflect.ValueOf(member)
		if !key.IsValid() {
			return nil, fmt.Errorf("invalid map key: %s", member)
		}
		if !key.Type().AssignableTo(val.Type().Key()) {
			// Try to convert key type
			key = key.Convert(val.Type().Key())
		}
		mapVal := val.MapIndex(key)
		if !mapVal.IsValid() {
			return nil, fmt.Errorf("map has no key %s", member)
		}
		return mapVal.Interface(), nil

	default:
		return nil, fmt.Errorf("cannot access member %s of type %T", member, obj)
	}
}

// ResolvePipeExpression resolves a pipe expression (function call).
func (r *Resolver) ResolvePipeExpression(pipe *ast.PipeExpression) (interface{}, error) {
	// Get the function
	fn, ok := r.scope.GetFunction(pipe.Function.Value)
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", pipe.Function.Value)
	}

	// Resolve arguments
	args := make([]interface{}, len(pipe.Arguments))
	for i, arg := range pipe.Arguments {
		val, err := r.ResolveExpression(arg)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	// Call the function
	return r.CallFunction(fn, args)
}

// CallFunction calls a function with the given arguments.
func (r *Resolver) CallFunction(fn interface{}, args []interface{}) (interface{}, error) {
	fnVal := reflect.ValueOf(fn)
	if fnVal.Kind() != reflect.Func {
		return nil, fmt.Errorf("not a function: %T", fn)
	}

	fnType := fnVal.Type()

	// Convert arguments to reflect values
	argVals := make([]reflect.Value, len(args))
	for i, arg := range args {
		if arg == nil {
			// Determine the expected type for this parameter position
			var expectedType reflect.Type
			if fnType.IsVariadic() && i >= fnType.NumIn()-1 {
				// Variadic element type
				expectedType = fnType.In(fnType.NumIn() - 1).Elem()
			} else if i < fnType.NumIn() {
				expectedType = fnType.In(i)
			}
			if expectedType != nil {
				argVals[i] = reflect.Zero(expectedType)
			} else {
				argVals[i] = reflect.Zero(reflect.TypeOf((*interface{})(nil)).Elem())
			}
		} else {
			argVals[i] = reflect.ValueOf(arg)
		}
	}

	// Call the function
	results := fnVal.Call(argVals)

	// Handle return values
	if len(results) == 0 {
		return nil, nil
	} else if len(results) == 1 {
		return results[0].Interface(), nil
	} else {
		// Multiple return values - return as slice
		resultSlice := make([]interface{}, len(results))
		for i, result := range results {
			resultSlice[i] = result.Interface()
		}
		return resultSlice, nil
	}
}

// ResolveInfixExpression resolves an infix expression.
func (r *Resolver) ResolveInfixExpression(infix *ast.InfixExpression) (interface{}, error) {
	left, err := r.ResolveExpression(infix.Left)
	if err != nil {
		return nil, err
	}

	right, err := r.ResolveExpression(infix.Right)
	if err != nil {
		return nil, err
	}

	return r.ApplyOperator(infix.Operator, left, right)
}

// ResolvePrefixExpression resolves a prefix expression.
func (r *Resolver) ResolvePrefixExpression(prefix *ast.PrefixExpression) (interface{}, error) {
	right, err := r.ResolveExpression(prefix.Right)
	if err != nil {
		return nil, err
	}

	switch prefix.Operator {
	case "!":
		return r.ApplyNotOperator(right)
	case "-":
		return r.ApplyNegateOperator(right)
	default:
		return nil, fmt.Errorf("unknown prefix operator: %s", prefix.Operator)
	}
}

// ResolveIndexExpression resolves an index expression.
func (r *Resolver) ResolveIndexExpression(index *ast.IndexExpression) (interface{}, error) {
	left, err := r.ResolveExpression(index.Left)
	if err != nil {
		return nil, err
	}

	idx, err := r.ResolveExpression(index.Index)
	if err != nil {
		return nil, err
	}

	return r.AccessIndex(left, idx)
}

// AccessIndex accesses an element by index (array index, map key, etc.).
func (r *Resolver) AccessIndex(obj interface{}, index interface{}) (interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("cannot index nil")
	}

	val := reflect.ValueOf(obj)

	// Handle pointer types
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("cannot index nil pointer")
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		// Convert index to int
		idx, err := r.ToInt(index)
		if err != nil {
			return nil, fmt.Errorf("invalid index type: %T", index)
		}
		if idx < 0 || idx >= val.Len() {
			return nil, fmt.Errorf("index out of range: %d", idx)
		}
		return val.Index(idx).Interface(), nil

	case reflect.Map:
		// Use index as map key
		key := reflect.ValueOf(index)
		if !key.IsValid() {
			return nil, fmt.Errorf("invalid map key: %v", index)
		}
		if !key.Type().AssignableTo(val.Type().Key()) {
			// Try to convert key type
			key = key.Convert(val.Type().Key())
		}
		mapVal := val.MapIndex(key)
		if !mapVal.IsValid() {
			return nil, fmt.Errorf("map has no key: %v", index)
		}
		return mapVal.Interface(), nil

	case reflect.String:
		// String indexing
		idx, err := r.ToInt(index)
		if err != nil {
			return nil, fmt.Errorf("invalid index type: %T", index)
		}
		str := val.String()
		if idx < 0 || idx >= len(str) {
			return nil, fmt.Errorf("index out of range: %d", idx)
		}
		return string(str[idx]), nil

	default:
		return nil, fmt.Errorf("cannot index type %T", obj)
	}
}

// ResolveArrayLiteral resolves an array literal.
func (r *Resolver) ResolveArrayLiteral(arr *ast.ArrayLiteral) (interface{}, error) {
	elements := make([]interface{}, len(arr.Elements))
	for i, elem := range arr.Elements {
		val, err := r.ResolveExpression(elem)
		if err != nil {
			return nil, err
		}
		elements[i] = val
	}
	return elements, nil
}

// ResolveHashLiteral resolves a hash literal.
func (r *Resolver) ResolveHashLiteral(hash *ast.HashLiteral) (interface{}, error) {
	result := make(map[interface{}]interface{})
	for keyExpr, valExpr := range hash.Pairs {
		key, err := r.ResolveExpression(keyExpr)
		if err != nil {
			return nil, err
		}
		val, err := r.ResolveExpression(valExpr)
		if err != nil {
			return nil, err
		}
		result[key] = val
	}
	return result, nil
}

// ApplyOperator applies a binary operator to two values.
func (r *Resolver) ApplyOperator(op string, left, right interface{}) (interface{}, error) {
	switch op {
	case "+":
		return r.ApplyAddOperator(left, right)
	case "-":
		return r.ApplySubtractOperator(left, right)
	case "*":
		return r.ApplyMultiplyOperator(left, right)
	case "/":
		return r.ApplyDivideOperator(left, right)
	case "==":
		return r.ApplyEqualOperator(left, right)
	case "!=":
		return r.ApplyNotEqualOperator(left, right)
	case "<":
		return r.ApplyLessThanOperator(left, right)
	case "<=":
		return r.ApplyLessEqualOperator(left, right)
	case ">":
		return r.ApplyGreaterThanOperator(left, right)
	case ">=":
		return r.ApplyGreaterEqualOperator(left, right)
	case "in":
		return r.ApplyInOperator(left, right)
	default:
		return nil, fmt.Errorf("unknown operator: %s", op)
	}
}

// ApplyAddOperator applies the + operator.
func (r *Resolver) ApplyAddOperator(left, right interface{}) (interface{}, error) {
	leftNum, leftErr := r.ToNumber(left)
	rightNum, rightErr := r.ToNumber(right)

	if leftErr == nil && rightErr == nil {
		// Both are numbers
		if r.IsFloat(left) || r.IsFloat(right) {
			return leftNum + rightNum, nil
		}
		return int64(leftNum) + int64(rightNum), nil
	}

	// String concatenation
	leftStr, leftStrErr := r.ToString(left)
	rightStr, rightStrErr := r.ToString(right)
	if leftStrErr == nil && rightStrErr == nil {
		return leftStr + rightStr, nil
	}

	return nil, fmt.Errorf("cannot add %T and %T", left, right)
}

// ApplySubtractOperator applies the - operator.
func (r *Resolver) ApplySubtractOperator(left, right interface{}) (interface{}, error) {
	leftNum, err := r.ToNumber(left)
	if err != nil {
		return nil, err
	}
	rightNum, err := r.ToNumber(right)
	if err != nil {
		return nil, err
	}

	if r.IsFloat(left) || r.IsFloat(right) {
		return leftNum - rightNum, nil
	}
	return int64(leftNum) - int64(rightNum), nil
}

// ApplyMultiplyOperator applies the * operator.
func (r *Resolver) ApplyMultiplyOperator(left, right interface{}) (interface{}, error) {
	leftNum, err := r.ToNumber(left)
	if err != nil {
		return nil, err
	}
	rightNum, err := r.ToNumber(right)
	if err != nil {
		return nil, err
	}

	if r.IsFloat(left) || r.IsFloat(right) {
		return leftNum * rightNum, nil
	}
	return int64(leftNum) * int64(rightNum), nil
}

// ApplyDivideOperator applies the / operator.
func (r *Resolver) ApplyDivideOperator(left, right interface{}) (interface{}, error) {
	leftNum, err := r.ToNumber(left)
	if err != nil {
		return nil, err
	}
	rightNum, err := r.ToNumber(right)
	if err != nil {
		return nil, err
	}

	if rightNum == 0 {
		return nil, fmt.Errorf("division by zero")
	}

	return leftNum / rightNum, nil
}

// ApplyEqualOperator applies the == operator.
func (r *Resolver) ApplyEqualOperator(left, right interface{}) (interface{}, error) {
	// Try numeric comparison first to handle int/int64/float64 cross-type equality.
	leftNum, leftErr := r.ToNumber(left)
	rightNum, rightErr := r.ToNumber(right)
	if leftErr == nil && rightErr == nil {
		return leftNum == rightNum, nil
	}
	// Fall back to string comparison for non-numeric types.
	leftStr, leftStrErr := r.ToString(left)
	rightStr, rightStrErr := r.ToString(right)
	if leftStrErr == nil && rightStrErr == nil {
		return leftStr == rightStr, nil
	}
	return reflect.DeepEqual(left, right), nil
}

// ApplyNotEqualOperator applies the != operator.
func (r *Resolver) ApplyNotEqualOperator(left, right interface{}) (interface{}, error) {
	return !reflect.DeepEqual(left, right), nil
}

// ApplyLessThanOperator applies the < operator.
func (r *Resolver) ApplyLessThanOperator(left, right interface{}) (interface{}, error) {
	leftNum, err := r.ToNumber(left)
	if err != nil {
		return nil, err
	}
	rightNum, err := r.ToNumber(right)
	if err != nil {
		return nil, err
	}
	return leftNum < rightNum, nil
}

// ApplyLessEqualOperator applies the <= operator.
func (r *Resolver) ApplyLessEqualOperator(left, right interface{}) (interface{}, error) {
	leftNum, err := r.ToNumber(left)
	if err != nil {
		return nil, err
	}
	rightNum, err := r.ToNumber(right)
	if err != nil {
		return nil, err
	}
	return leftNum <= rightNum, nil
}

// ApplyGreaterThanOperator applies the > operator.
func (r *Resolver) ApplyGreaterThanOperator(left, right interface{}) (interface{}, error) {
	leftNum, err := r.ToNumber(left)
	if err != nil {
		return nil, err
	}
	rightNum, err := r.ToNumber(right)
	if err != nil {
		return nil, err
	}
	return leftNum > rightNum, nil
}

// ApplyGreaterEqualOperator applies the >= operator.
func (r *Resolver) ApplyGreaterEqualOperator(left, right interface{}) (interface{}, error) {
	leftNum, err := r.ToNumber(left)
	if err != nil {
		return nil, err
	}
	rightNum, err := r.ToNumber(right)
	if err != nil {
		return nil, err
	}
	return leftNum >= rightNum, nil
}

// ApplyInOperator checks if left is contained in right (slice or string).
func (r *Resolver) ApplyInOperator(left, right interface{}) (interface{}, error) {
	leftStr := fmt.Sprintf("%v", left)
	switch v := right.(type) {
	case []interface{}:
		for _, item := range v {
			if fmt.Sprintf("%v", item) == leftStr {
				return true, nil
			}
		}
		return false, nil
	case string:
		return strings.Contains(v, leftStr), nil
	default:
		return false, nil
	}
}

// ApplyNotOperator applies the ! operator.
func (r *Resolver) ApplyNotOperator(val interface{}) (interface{}, error) {
	b, err := r.ToBool(val)
	if err != nil {
		return nil, err
	}
	return !b, nil
}

// ApplyNegateOperator applies the - operator (unary).
func (r *Resolver) ApplyNegateOperator(val interface{}) (interface{}, error) {
	num, err := r.ToNumber(val)
	if err != nil {
		return nil, err
	}
	return -num, nil
}

// ToNumber converts a value to a float64.
func (r *Resolver) ToNumber(val interface{}) (float64, error) {
	if val == nil {
		return 0, fmt.Errorf("cannot convert nil to number")
	}

	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string %q to number", v)
		}
		return f, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to number", val)
	}
}

// ToInt converts a value to an int.
func (r *Resolver) ToInt(val interface{}) (int, error) {
	num, err := r.ToNumber(val)
	if err != nil {
		return 0, err
	}
	return int(num), nil
}

// ToString converts a value to a string.
func (r *Resolver) ToString(val interface{}) (string, error) {
	if val == nil {
		return "", nil
	}

	switch v := val.(type) {
	case string:
		return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", v), nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// ToBool converts a value to a boolean.
func (r *Resolver) ToBool(val interface{}) (bool, error) {
	if val == nil {
		return false, nil
	}

	switch v := val.(type) {
	case bool:
		return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		num, err := r.ToNumber(val)
		if err != nil {
			return false, err
		}
		return num != 0, nil
	case string:
		return v != "", nil
	default:
		return true, nil
	}
}

// IsFloat checks if a value is a float type.
func (r *Resolver) IsFloat(val interface{}) bool {
	switch val.(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}

// SetScope sets the resolver's scope.
func (r *Resolver) SetScope(s *scope.Scope) {
	r.scope = s
}

// GetScope returns the resolver's scope.
func (r *Resolver) GetScope() *scope.Scope {
	return r.scope
}
