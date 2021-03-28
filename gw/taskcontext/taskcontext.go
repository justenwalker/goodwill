// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package taskcontext

import (
	"context"
	"fmt"
	"go.justen.tech/goodwill/gw/values"
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/grpc"
)

type Service struct {
	client pb.ContextServiceClient
}

func NewService(conn grpc.ClientConnInterface) *Service {
	return &Service{
		client: pb.NewContextServiceClient(conn),
	}
}

// SetVariable sets the tasks variable to the given value
func (c *Service) SetVariable(ctx context.Context, name string, value values.Value) error {
	variable, err := newVariable(name, value)
	if err != nil {
		return err
	}
	_, err = c.client.SetVariable(ctx, variable)
	return err
}

// GetVariable gets the tasks variable to the given value
func (c *Service) GetVariable(ctx context.Context, name string, out values.ValueOut) error {
	value, err := c.client.GetVariable(ctx, &pb.VariableName{Name: name})
	if err != nil {
		return err
	}
	return values.Unmarshal(value, out)
}

// GetVariableNames gets the list of variable names currently set in the task
func (c *Service) GetVariableNames(ctx context.Context) ([]string, error) {
	names, err := c.client.GetVariableNames(ctx, &pb.GetVariableNameParams{})
	if err != nil {
		return nil, err
	}
	return names.Name, nil
}

// GetVariables gets all variable currently set in the task
func (c *Service) GetVariables(ctx context.Context) (map[string]interface{}, error) {
	vars, err := c.client.GetVariables(ctx, &pb.GetVariablesRequest{})
	if err != nil {
		return nil, err
	}
	variables := make(map[string]interface{})
	for key, value := range vars.Value {
		v, err := values.Interface(value)
		if err != nil {
			return nil, err
		}
		variables[key] = v
	}
	return variables, nil
}

// EvaluateParams evaluates the given expression, and returns the result into the output value.
// The given map of parameters are set as variables before the expression is evaluated,
// which approximates a parameterized query; allowing a safer expression evaluation compared to string concatenation.
func (c *Service) EvaluateParams(ctx context.Context, expr string, out values.ValueOut, params map[string]values.Value, ) error {
	var parameters []*pb.Variable
	for key, val := range params {
		mv, err := values.Marshal(val)
		if err != nil {
			return fmt.Errorf("failed to marshal parameter %q: %w", key, err)
		}
		parameters = append(parameters, &pb.Variable{
			Name:  key,
			Value: mv,
		})
	}
	v, err := c.client.Evaluate(ctx, &pb.EvaluateRequest{
		Expression: expr,
		Parameters: parameters,
		Type:       out.Type(),
	})
	if err != nil {
		return err
	}
	return values.Unmarshal(v, out)
}

// EvaluateParams evaluates the given expression, and returns the result into the output value
func (c *Service) Evaluate(ctx context.Context, expr string, out values.ValueOut) error {
	return c.EvaluateParams(ctx, expr, out, nil)
}

func newVariable(key string, value values.Value) (*pb.Variable, error) {
	v, err := values.Marshal(value)
	if err != nil {
		return nil, err
	}
	return &pb.Variable{
		Name:  key,
		Value: v,
	}, nil
}