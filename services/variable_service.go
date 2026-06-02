package services

import (
	"regexp"

	"paw-api/repositories"
)

var varPattern = regexp.MustCompile(`\{\{(\w+)\}\}`)

type VariableService struct {
	VariableRepo *repositories.EnvVariableRepository
}

func NewVariableService(variableRepo *repositories.EnvVariableRepository) *VariableService {
	return &VariableService{VariableRepo: variableRepo}
}

func (s *VariableService) loadLookup(envID int64) (map[string]string, error) {
	vars, err := s.VariableRepo.ListByEnvironment(envID)
	if err != nil {
		return nil, err
	}
	lookup := make(map[string]string, len(vars))
	for _, v := range vars {
		if v.Enabled {
			lookup[v.Key] = v.Value
		}
	}
	return lookup, nil
}

func (s *VariableService) Resolve(input string, envID int64) (string, error) {
	lookup, err := s.loadLookup(envID)
	if err != nil {
		return input, err
	}
	return resolveWithLookup(input, lookup), nil
}

func (s *VariableService) ResolveMap(input map[string]string, envID int64) (map[string]string, error) {
	lookup, err := s.loadLookup(envID)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(input))
	for k, v := range input {
		result[k] = resolveWithLookup(v, lookup)
	}
	return result, nil
}

func resolveWithLookup(input string, lookup map[string]string) string {
	result := input
	for i := 0; i < 10; i++ {
		prev := result
		result = varPattern.ReplaceAllStringFunc(result, func(match string) string {
			key := match[2 : len(match)-2]
			if val, ok := lookup[key]; ok {
				return val
			}
			return match
		})
		if result == prev {
			break
		}
	}
	return result
}
