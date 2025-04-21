package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidStructType     = errors.New("value is not structure")
	ErrUnsupportedType       = errors.New("unsupported type")
	ErrInvalidValidationRule = errors.New("invalid validation rule")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errList := make([]string, 0, len(v))
	for _, err := range v {
		errList = append(errList, fmt.Sprintf("%s: %v", err.Field, err.Err))
	}
	return "Invalid fields: " + strings.Join(errList, ", ")
}

func (v ValidationErrors) Unwrap() []error {
	errs := make([]error, 0, len(v))
	for _, err := range v {
		errs = append(errs, err.Err)
	}
	return errs
}

func Validate(s any) (err error) {
	v, t := reflect.ValueOf(s), reflect.TypeOf(s)

	if v.Kind() != reflect.Struct {
		err = ErrInvalidStructType
		return
	}

	var errs ValidationErrors

	for i := range v.NumField() {
		fv := v.Field(i)
		ft := t.Field(i)

		rule := ft.Tag.Get("validate")
		if len(rule) == 0 {
			continue
		}

		match, err := newMatcher(fv, rule)
		if err != nil {
			return err // rule definition error
		}

		if err := match(fv); err != nil {
			errs = append(errs, ValidationError{ft.Name, err}) // validation errors
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

var (
	ErrValidationMaxValue    = errors.New("more than maximum value")
	ErrValidationMinValue    = errors.New("less than minimum value")
	ErrValidationInValue     = errors.New("not included in the set")
	ErrValidationLength      = errors.New("invalid length of string")
	ErrValidationRegexpValue = errors.New("not match regular expression")
)

type matchFunc func(v reflect.Value) error

func splitRules(rule string) []string {
	return strings.Split(rule, "|")
}

func parseRule(rule string) (id string, condition string, err error) {
	q := strings.SplitN(rule, ":", 2)
	if len(q) != 2 {
		err = ErrInvalidValidationRule
	} else {
		id, condition = q[0], q[1]
	}
	return
}

// chainMatchers returns match function which calls
// matchers consequentially and returns nil iff
// all functions return nil. Otherwise, the first
// error returned by matchers is returned.
func chainMatchers(matchers []matchFunc) matchFunc {
	return func(v reflect.Value) error {
		for _, matcher := range matchers {
			if err := matcher(v); err != nil {
				return err
			}
		}
		return nil
	}
}

func newMatcher(v reflect.Value, rule string) (matchFunc, error) {
	switch v.Kind() { //nolint:exhaustive
	case reflect.String, reflect.Int:
		return newMatcherForKind(v.Kind(), rule)
	case reflect.Slice:
		match, err := newMatcherForKind(v.Type().Elem().Kind(), rule)
		if err != nil {
			return nil, err
		}
		return func(v reflect.Value) error {
			for i := range v.Len() {
				if err := match(v.Index(i)); err != nil {
					return fmt.Errorf("[%d] %w", i, err)
				}
			}
			return nil
		}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

// newMatcherForKind returns match function for reflect.Kind type.
func newMatcherForKind(k reflect.Kind, rule string) (matchFunc, error) {
	rules := splitRules(rule)
	matchers := make([]matchFunc, 0, len(rules))

	for _, r := range rules {
		id, requirement, err := parseRule(r)
		if err != nil {
			return nil, err
		}

		var matcher matchFunc
		switch k { //nolint:exhaustive
		case reflect.Int:
			matcher, err = newIntMatcher(id, requirement)
		case reflect.String:
			matcher, err = newStrMatcher(id, requirement)
		default:
			err = ErrUnsupportedType
		}
		if err != nil {
			return nil, err
		}

		matchers = append(matchers, matcher)
	}

	return chainMatchers(matchers), nil
}

// newIntMatcher returns match function for integer type
// built according to condition.
func newIntMatcher(id string, condition string) (matchFunc, error) {
	switch id {
	case "in":
		vals := make(map[int64]struct{})
		for _, v := range strings.Split(condition, ",") {
			val, err := strconv.Atoi(v)
			if err != nil {
				return nil, ErrInvalidValidationRule
			}
			vals[int64(val)] = struct{}{}
		}
		return func(v reflect.Value) error {
			if _, ok := vals[v.Int()]; !ok {
				return ErrValidationInValue
			}
			return nil
		}, nil
	case "max":
		maxValue, err := strconv.Atoi(condition)
		if err != nil {
			return nil, ErrInvalidValidationRule
		}
		return func(v reflect.Value) error {
			if v.Int() > int64(maxValue) {
				return ErrValidationMaxValue
			}
			return nil
		}, nil
	case "min":
		minValue, err := strconv.Atoi(condition)
		if err != nil {
			return nil, ErrInvalidValidationRule
		}
		return func(v reflect.Value) error {
			if v.Int() < int64(minValue) {
				return ErrValidationMinValue
			}
			return nil
		}, nil
	default:
		return nil, ErrInvalidValidationRule
	}
}

// newStrMatcher returns match function for string type
// built according to condtition.
func newStrMatcher(id string, condition string) (matchFunc, error) {
	switch id {
	case "in":
		vals := make(map[string]struct{})
		for _, v := range strings.Split(condition, ",") {
			vals[v] = struct{}{}
		}
		return func(v reflect.Value) error {
			if _, ok := vals[v.String()]; !ok {
				return ErrValidationInValue
			}
			return nil
		}, nil
	case "len":
		length, err := strconv.Atoi(condition)
		if err != nil {
			return nil, ErrInvalidValidationRule
		}
		return func(v reflect.Value) error {
			if len(v.String()) != length {
				return ErrValidationLength
			}
			return nil
		}, nil
	case "regexp":
		re, err := regexp.Compile(condition)
		if err != nil {
			return nil, ErrInvalidValidationRule
		}
		return func(v reflect.Value) error {
			if !re.MatchString(v.String()) {
				return ErrValidationRegexpValue
			}
			return nil
		}, nil
	default:
		return nil, ErrInvalidValidationRule
	}
}
