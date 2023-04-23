// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

// +build goodwill

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.justen.tech/goodwill/gw"
	"go.justen.tech/goodwill/gw/docker"
	"go.justen.tech/goodwill/gw/secret"
	"go.justen.tech/goodwill/gw/value"
	"io/ioutil"
	"os"
	"time"
)

// Default is a flow that prints the working directory
func Default(ctx context.Context, ts *gw.Task) error {
	fmt.Println("MinServerVersion:", gw.MinServerVersion)
	fmt.Println("ServerVersion:", ts.ServerVersion)
	fmt.Println("====== Get Task Config")
	cfg, err := ts.Config().Configuration(ctx)
	if err != nil {
		return err
	}
	bs, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println("Configuration:")
	fmt.Println(string(bs))
	_ = ts.Log(ctx, "Log Message: %d", 42)
	return nil
}

// SetVariables gets and sets variables
func SetVariables(ctx context.Context, ts *gw.Task) (map[string]value.Value, error) {
	var t value.Time
	tc := ts.Context()
	fmt.Println("====== Evaluate Expression")
	if err := tc.Evaluate(ctx, "${datetime.current()}", &t); err != nil {
		return nil, fmt.Errorf("evaluate expression failed: %w", err)
	}
	fmt.Printf("datetime.current: %v\n", t.Format(time.RFC3339Nano))

	fmt.Println("====== Set/Get Variables")
	fmt.Println("Set timeStamp:", t)
	if err := tc.SetVariable(ctx, "timeStamp", t); err != nil {
		return nil, fmt.Errorf("set timeStamp var failed: %w", err)
	}
	fmt.Println("Get timeStamp")
	if err := tc.GetVariable(ctx, "timeStamp", &t); err != nil {
		return nil, fmt.Errorf("get timeStamp var failed: %w", err)
	}
	fmt.Printf("timeStamp: %v\n", t.Format(time.RFC3339Nano))

	fmt.Println("====== Dump Variables")
	vars, err := tc.GetVariables(ctx)
	if err != nil {
		return nil, fmt.Errorf("ger variables failed: %w", err)
	}
	fmt.Println("Variables:")
	bs, _ := json.MarshalIndent(vars, "", "  ")
	fmt.Println(string(bs))
	return map[string]value.Value{
		"foo": value.String("bar"),
		"baz": value.Int64(1000),
	}, nil
}

// Crypto does crypto operations on secrets
func Crypto(ts *gw.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	crypt := ts.Secret()
	fmt.Println("====== Encrypt String")
	enc, err := crypt.EncryptString(ctx, "foo")
	if err != nil {
		return fmt.Errorf("encrypt string failed: %w", err)
	}
	fmt.Println("Encrypted:", enc)
	fmt.Println("====== Decrypt String")
	dec, err := crypt.DecryptString(ctx, enc)
	if err != nil {
		return fmt.Errorf("decrypt string failed: %w", err)
	}
	fmt.Println("Decrypted:", dec)
	gkpname := fmt.Sprintf("genkp-%s", ts.ProcessID)
	nkpname := fmt.Sprintf("newkp-%s", ts.ProcessID)
	fmt.Println("====== Generate Key Pair")
	kp, err := crypt.GenerateKeyPair(ctx, gkpname, secret.GeneratePassword)
	if err != nil {
		return fmt.Errorf("generate keypair %q failed: %w", "gen-keypair", err)
	}
	fmt.Println("ID:", kp.ID)
	fmt.Println("PublicKey:", kp.PublicKey)
	fmt.Println("StorePassword:", kp.StorePassword)
	fmt.Println("====== Export Key Pair")
	keypair, err := crypt.KeyPairFiles(ctx, gkpname, secret.StorePassword(kp.StorePassword))
	if err != nil {
		return fmt.Errorf("export keypair %q files failed: %w", gkpname, err)
	}
	fmt.Println("Public Key File: ", keypair.PublicKeyFile)
	fmt.Println("Private Key File: ", keypair.PrivateKeyFile)
	pub, err := ioutil.ReadFile(keypair.PublicKeyFile)
	if err != nil {
		return fmt.Errorf("read public key %q  failed: %w", keypair.PublicKeyFile, err)
	}
	priv, err := ioutil.ReadFile(keypair.PrivateKeyFile)
	if err != nil {
		return fmt.Errorf("read private key %q  failed: %w", keypair.PrivateKeyFile, err)
	}
	fmt.Println("====== Create Key Pair")
	nkp, err := crypt.CreateKeyPair(ctx, nkpname, secret.KeyPair{
		PublicKey:  pub,
		PrivateKey: priv,
	})
	if err != nil {
		return fmt.Errorf("create keypair failed: %w", err)
	}
	fmt.Println("ID:", nkp.ID)
	fmt.Println("PublicKey:", nkp.PublicKey)
	fmt.Println("====== Update Access Level")
	if err := crypt.UpdateAccess(ctx, gkpname, []secret.AccessEntry{
		{
			TeamName: gw.DefaultTeamName,
			OrgName:  gw.DefaultOrgName,
			Level:    secret.Owner,
		},
	}); err != nil {
		return fmt.Errorf("update gen-keypair access failed: %w", err)
	}
	fmt.Println("====== Get Access Level")
	access, err := crypt.ListAccess(ctx, gkpname)
	if err != nil {
		return fmt.Errorf("list gen-keypair access failed: %w", err)
	}
	for _, a := range access {
		fmt.Printf("- %+v\n", a)
	}
	fmt.Println("====== Delete Secrets")
	if err := crypt.Delete(ctx, gkpname); err != nil {
		return fmt.Errorf("delete gen-keypair failed: %w", err)
	}
	if err := crypt.Delete(ctx, nkpname); err != nil {
		return fmt.Errorf("delete create-keypair failed: %w", err)
	}
	return nil
}

// Lock and unlock a project lock
func Lock(ts *gw.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	fmt.Println("====== Project Lock")
	if err := ts.Lock().ProjectLock(ctx, "lockName"); err != nil {
		return fmt.Errorf("lock project failed: %w", err)
	}
	if err := ts.Lock().ProjectUnlock(ctx, "lockName"); err != nil {
		return fmt.Errorf("unlock project failed: %w", err)
	}
	return nil
}

// KeyValue tests the project key/value store
func KeyValue(ts *gw.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	fmt.Println("====== Key Value")
	kv := ts.KV()
	if err := kv.PutString(ctx, "str", "ok"); err != nil {
		return fmt.Errorf("put string failed: %w", err)
	}
	if err := kv.PutLong(ctx, "long", 100); err != nil {
		return fmt.Errorf("put long failed: %w", err)
	}
	l, err := kv.GetLong(ctx, "long")
	if err != nil {
		return fmt.Errorf("get long failed: %w", err)
	}
	fmt.Println("kv.long =", l)
	l, err = kv.Inc(ctx, "long")
	if err != nil {
		return fmt.Errorf("inc long failed: %w", err)
	}
	fmt.Println("inc kv.long =", l)
	str, err := kv.GetString(ctx, "str")
	if err != nil {
		return fmt.Errorf("get string failed: %w", err)
	}
	fmt.Println("kv.str =", str)
	if err := kv.Remove(ctx, "str"); err != nil {
		return fmt.Errorf("remove string failed: %w", err)
	}
	if err := kv.Remove(ctx, "long"); err != nil {
		return fmt.Errorf("remove long failed: %w", err)
	}
	return nil
}

// JSONStore tests the json store task
func JSONStore(ts *gw.Task) error {
	type serviceObject struct {
		Users   []string `json:"users"`
		Service string   `json:"service"`
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	fmt.Println("====== JSONStore")
	js := ts.JSONStore("TestStore")
	fmt.Println("Put service_c")
	if err := js.Put(ctx, "service_c", value.JSON{
		Value: serviceObject{
			Users:   []string{"mike", "alice"},
			Service: "service_c",
		},
	}); err != nil {
		return fmt.Errorf("failed insert service_c: %w", err)
	}
	fmt.Println("Execute lookupServiceByUser")
	var result value.JSON
	if err := js.ExecuteQuery(ctx, "lookupServiceByUser", map[string]interface{}{
		"users": []string{"mike"},
	}, &result); err != nil {
		return fmt.Errorf("failed lookupServiceByUser: %w", err)
	}
	fmt.Println("jsonStore.lookupServiceByUser:", result)

	fmt.Println("Remove service_c")
	removed, err := js.Delete(ctx, "service_c")
	if err != nil {
		return fmt.Errorf("failed delete service_c: %w", err)
	}
	fmt.Println("removed:", removed)
	fmt.Println("Remove service_c (again)")
	removed, err = js.Delete(ctx, "service_c")
	if err != nil {
		return fmt.Errorf("failed delete service_c: %w", err)
	}
	fmt.Println("removed:", removed)
	return nil
}

// Docker runs a docker container
func Docker(ts *gw.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	fmt.Println("====== Docker Container")
	return ts.Docker().RunContainer(ctx, "ubuntu:20.04",
		docker.Command("ls", "-la"),
		docker.Stdout(os.Stdout),
		docker.Stderr(os.Stderr),
	)
}
