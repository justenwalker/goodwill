// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package jsonstore

import (
	"context"

	"google.golang.org/grpc"

	"go.justen.tech/goodwill/gw/taskcontext"
	"go.justen.tech/goodwill/gw/value"
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
	var result value.Bool
	if err := c.ex.EvaluateParams(ctx, jsonStoreItemExists, &result, map[string]value.Value{
		jsonStoreOrg:      value.String(c.orgName),
		jsonStoreName:     value.String(c.store),
		jsonStoreItemPath: value.String(itemPath),
	}); err != nil {
		return false, err
	}
	return bool(result), nil
}

// Get gets the item in the json store
func (c *Store) Get(ctx context.Context, itemPath string, out value.ValueOut) error {
	return c.ex.EvaluateParams(ctx, jsonStoreGet, out, map[string]value.Value{
		jsonStoreOrg:      value.String(c.orgName),
		jsonStoreName:     value.String(c.store),
		jsonStoreItemPath: value.String(itemPath),
	})
}

// Delete an item in the json store.
// returns true if the item was deleted, otherwise false.
func (c *Store) Delete(ctx context.Context, itemPath string) (bool, error) {
	var v value.Bool
	if err := c.ex.EvaluateParams(ctx, jsonStoreDelete, &v, map[string]value.Value{
		jsonStoreOrg:      value.String(c.orgName),
		jsonStoreName:     value.String(c.store),
		jsonStoreItemPath: value.String(itemPath),
	}); err != nil {
		return false, err
	}
	return bool(v), nil
}

// Put puts the item in the json store
func (c *Store) Put(ctx context.Context, itemPath string, data value.Value) error {
	return c.ex.EvaluateParams(ctx, jsonStorePut, value.Discard, map[string]value.Value{
		jsonStoreOrg:      value.String(c.orgName),
		jsonStoreName:     value.String(c.store),
		jsonStoreItemPath: value.String(itemPath),
		jsonStoreItemData: data,
	})
}

// Upsert inserts or updates an item
func (c *Store) Upsert(ctx context.Context, itemPath string, data value.Value) error {
	return c.ex.EvaluateParams(ctx, jsonStoreUpsert, value.Discard, map[string]value.Value{
		jsonStoreOrg:      value.String(c.orgName),
		jsonStoreName:     value.String(c.store),
		jsonStoreItemPath: value.String(itemPath),
		jsonStoreItemData: data,
	})
}

// UpsertQuery creates or updates a named query
func (c *Store) UpsertQuery(ctx context.Context, queryName string, queryText string) error {
	return c.ex.EvaluateParams(ctx, jsonStoreUpsertQuery, value.Discard, map[string]value.Value{
		jsonStoreOrg:       value.String(c.orgName),
		jsonStoreName:      value.String(c.store),
		jsonStoreQueryName: value.String(queryName),
		jsonStoreQueryText: value.String(queryText),
	})
}

// ExecuteQuery executes a named query
func (c *Store) ExecuteQuery(ctx context.Context, queryName string, params map[string]interface{}, out value.ValueOut) error {
	return c.ex.EvaluateParams(ctx, jsonStoreExecQuery, out, map[string]value.Value{
		jsonStoreOrg:         value.String(c.orgName),
		jsonStoreName:        value.String(c.store),
		jsonStoreQueryName:   value.String(queryName),
		jsonStoreQueryParams: value.JSON{Value: params},
	})
}
