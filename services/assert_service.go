package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"paw-api/pkg/httpclient"
)

type AssertRule struct {
	Type   string
	Target string
	Value  string
}

type AssertResult struct {
	Rule   AssertRule
	Passed bool
	Actual string
	Error  string
}

type AssertService struct{}

func NewAssertService() *AssertService {
	return &AssertService{}
}

func (s *AssertService) Run(resp *httpclient.Response, rules []AssertRule) []AssertResult {
	var results []AssertResult
	for _, rule := range rules {
		result := s.evaluate(resp, rule)
		results = append(results, result)
	}
	return results
}

func (s *AssertService) evaluate(resp *httpclient.Response, rule AssertRule) AssertResult {
	switch rule.Type {
	case "status":
		return s.assertStatus(resp, rule)
	case "body_contains":
		return s.assertBodyContains(resp, rule)
	case "body_jsonpath":
		return s.assertJSONPath(resp, rule)
	case "header_equals":
		return s.assertHeader(resp, rule)
	case "duration_lt":
		return s.assertDuration(resp, rule)
	default:
		return AssertResult{Rule: rule, Passed: false, Error: "unknown assert type: " + rule.Type}
	}
}

func (s *AssertService) assertStatus(resp *httpclient.Response, rule AssertRule) AssertResult {
	expected, err := strconv.Atoi(rule.Value)
	if err != nil {
		return AssertResult{Rule: rule, Passed: false, Actual: fmt.Sprintf("%d", resp.Status), Error: "invalid status code: " + rule.Value}
	}
	parts := strings.SplitN(rule.Target, " ", 2)
	op := "eq"
	if len(parts) == 2 {
		op = parts[0]
		expected, _ = strconv.Atoi(parts[1])
	}
	var passed bool
	switch op {
	case "eq", "=", "==":
		passed = resp.Status == expected
	case "neq", "!=":
		passed = resp.Status != expected
	case "gt", ">":
		passed = resp.Status > expected
	case "lt", "<":
		passed = resp.Status < expected
	default:
		passed = resp.Status == expected
	}
	return AssertResult{Rule: rule, Passed: passed, Actual: fmt.Sprintf("%d", resp.Status)}
}

func (s *AssertService) assertBodyContains(resp *httpclient.Response, rule AssertRule) AssertResult {
	contains := strings.Contains(resp.Body, rule.Value)
	negate := strings.HasPrefix(rule.Target, "not ")
	if negate {
		return AssertResult{Rule: rule, Passed: !contains, Actual: fmt.Sprintf("%v", contains)}
	}
	return AssertResult{Rule: rule, Passed: contains, Actual: fmt.Sprintf("%v", contains)}
}

func (s *AssertService) assertJSONPath(resp *httpclient.Response, rule AssertRule) AssertResult {
	var data interface{}
	if err := json.Unmarshal([]byte(resp.Body), &data); err != nil {
		return AssertResult{Rule: rule, Passed: false, Error: "response is not valid JSON: " + err.Error()}
	}
	val, err := resolveJSONPath(data, rule.Target)
	if err != nil {
		return AssertResult{Rule: rule, Passed: false, Actual: "<not found>", Error: "JSONPath error: " + err.Error()}
	}
	actual := fmt.Sprintf("%v", val)
	if rule.Value == "" {
		return AssertResult{Rule: rule, Passed: true, Actual: actual}
	}
	return AssertResult{Rule: rule, Passed: actual == rule.Value, Actual: actual}
}

func resolveJSONPath(data interface{}, path string) (interface{}, error) {
	if path == "$" || path == "" {
		return data, nil
	}
	path = strings.TrimPrefix(path, "$.")
	parts := strings.Split(path, ".")
	current := data
	for _, part := range parts {
		if mapData, ok := current.(map[string]interface{}); ok {
			val, found := mapData[part]
			if !found {
				return nil, fmt.Errorf("key not found: %s", part)
			}
			current = val
		} else if arrData, ok := current.([]interface{}); ok {
			idx := 0
			if part[0] == '[' && part[len(part)-1] == ']' {
				fmt.Sscanf(part, "[%d]", &idx)
			} else {
				fmt.Sscanf(part, "%d", &idx)
			}
			if idx < 0 || idx >= len(arrData) {
				return nil, fmt.Errorf("index out of range: %d", idx)
			}
			current = arrData[idx]
		} else {
			return nil, fmt.Errorf("cannot navigate into %T", current)
		}
	}
	return current, nil
}

func (s *AssertService) assertHeader(resp *httpclient.Response, rule AssertRule) AssertResult {
	headerName := rule.Target
	values, ok := resp.Headers[headerName]
	actual := strings.Join(values, ", ")
	negate := strings.HasPrefix(rule.Value, "not ")
	if negate {
		return AssertResult{Rule: rule, Passed: !ok || actual != rule.Value[4:], Actual: actual}
	}
	return AssertResult{Rule: rule, Passed: ok && actual == rule.Value, Actual: actual}
}

func (s *AssertService) assertDuration(resp *httpclient.Response, rule AssertRule) AssertResult {
	expected, err := strconv.Atoi(rule.Value)
	if err != nil {
		return AssertResult{Rule: rule, Passed: false, Actual: fmt.Sprintf("%dms", resp.DurationMs), Error: "invalid duration: " + rule.Value}
	}
	passed := resp.DurationMs < int64(expected)
	return AssertResult{Rule: rule, Passed: passed, Actual: fmt.Sprintf("%dms", resp.DurationMs)}
}
