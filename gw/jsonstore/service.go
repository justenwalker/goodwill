// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package jsonstore

import (
	"context"
	"go.justen.tech/goodwill/gw/taskcontext"
	"go.justen.tech/goodwill/gw/values"
	"google.golang.org/grpc"
)

// Store to access a JSONStore
type Store struct {
	orgName string
	store   string
	ex      *taskcontext.Service
}

func NewService(orgName string, store string, conn grpc.ClientConnInterface) *Store {
	return &Store{
		orgName: orgName,
		store:   store,
		ex:      taskcontext.NewService(conn),
	}
}

const (
	jsonStoreOrg         = "_goodwill_js_org"
	jsonStoreName        = "_goodwill_js_name"
	jsonStoreItemPath    = "_goodwill_js_item_path"
	jsonStoreItemData    = "_goodwill_js_item_data"
	jsonStoreQueryName   = "_goodwill_js_query_name"
	jsonStoreQueryText   = "_goodwill_js_query_text"
	jsonStoreQueryParams = "_goodwill_js_query_params"
	jsonStoreItemExists  = `${jsonStore.isExists(` + jsonStoreOrg + `, ` + jsonStoreName + `, ` + jsonStoreItemPath + `)}`
	jsonStorePut         = `${jsonStore.put(` + jsonStoreOrg + `, ` + jsonStoreName + `, ` + jsonStoreItemPath + `, ` + jsonStoreItemData + `)}`
	jsonStoreGet         = `${jsonStore.get(` + jsonStoreOrg + `, ` + jsonStoreName + `, ` + jsonStoreItemPath + `)}`
	jsonStoreDelete      = `${jsonStore.delete(` + jsonStoreOrg + `, ` + jsonStoreName + `, ` + jsonStoreItemPath + `)}`
	jsonStoreUpsert      = `${jsonStore.upsert(` + jsonStoreOrg + `, ` + jsonStoreName + `, ` + jsonStoreItemPath + `, ` + jsonStoreItemData + `)}`
	jsonStoreUpsertQuery = `${jsonStore.upsertQuery(` + jsonStoreOrg + `, ` + jsonStoreName + `, ` + jsonStoreQueryName + `, ` + jsonStoreQueryText + `)}`
	jsonStoreExecQuery   = `${jsonStore.executeQuery(` + jsonStoreOrg + `, ` + jsonStoreName + `, ` + jsonStoreQueryName + `, ` + jsonStoreQueryParams + `)}`
)

// ItemExists tests if an item exists in the store
func (c *Store) ItemExists(ctx context.Context, itemPath string) (bool, error) {
	var result values.Bool
	if err := c.ex.EvaluateParams(ctx, jsonStoreItemExists, &result, map[string]values.Value{
		jsonStoreOrg:      values.String(c.orgName),
		jsonStoreName:     values.String(c.store),
		jsonStoreItemPath: values.String(itemPath),
	}); err != nil {
		return false, err
	}
	return bool(result), nil
}

// Get gets the item in the json store
func (c *Store) Get(ctx context.Context, itemPath string, out values.ValueOut) error {
	return c.ex.EvaluateParams(ctx, jsonStoreGet, out, map[string]values.Value{
		jsonStoreOrg:      values.String(c.orgName),
		jsonStoreName:     values.String(c.store),
		jsonStoreItemPath: values.String(itemPath),
	})
}

// Delete an item in the json store.
// returns true if the item was deleted, otherwise false.
func (c *Store) Delete(ctx context.Context, itemPath string) (bool, error) {
	var value values.Bool
	if err := c.ex.EvaluateParams(ctx, jsonStoreDelete, &value, map[string]values.Value{
		jsonStoreOrg:      values.String(c.orgName),
		jsonStoreName:     values.String(c.store),
		jsonStoreItemPath: values.String(itemPath),
	}); err != nil {
		return false, err
	}
	return bool(value), nil
}

// Put puts the item in the json store
func (c *Store) Put(ctx context.Context, itemPath string, data values.Value) error {
	return c.ex.EvaluateParams(ctx, jsonStorePut, values.Discard, map[string]values.Value{
		jsonStoreOrg:      values.String(c.orgName),
		jsonStoreName:     values.String(c.store),
		jsonStoreItemPath: values.String(itemPath),
		jsonStoreItemData: data,
	})
}

// Upsert inserts or updates an item
func (c *Store) Upsert(ctx context.Context, itemPath string, data values.Value) error {
	return c.ex.EvaluateParams(ctx, jsonStoreUpsert, values.Discard, map[string]values.Value{
		jsonStoreOrg:      values.String(c.orgName),
		jsonStoreName:     values.String(c.store),
		jsonStoreItemPath: values.String(itemPath),
		jsonStoreItemData: data,
	})
}

// UpsertQuery creates or updates a named query
func (c *Store) UpsertQuery(ctx context.Context, queryName string, queryText string) error {
	return c.ex.EvaluateParams(ctx, jsonStoreUpsertQuery, values.Discard, map[string]values.Value{
		jsonStoreOrg:       values.String(c.orgName),
		jsonStoreName:      values.String(c.store),
		jsonStoreQueryName: values.String(queryName),
		jsonStoreQueryText: values.String(queryText),
	})
}

// ExecuteQuery executes a named query
func (c *Store) ExecuteQuery(ctx context.Context, queryName string, params map[string]interface{}, out values.ValueOut) error {
	return c.ex.EvaluateParams(ctx, jsonStoreExecQuery, out, map[string]values.Value{
		jsonStoreOrg:         values.String(c.orgName),
		jsonStoreName:        values.String(c.store),
		jsonStoreQueryName:   values.String(queryName),
		jsonStoreQueryParams: values.JSON{Value: params},
	})
}
