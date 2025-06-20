// Code generated by BobGen psql v0.38.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package factory

import (
	"context"
	"testing"
)

func TestCreateUser(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping test, no DSN provided")
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("Error starting transaction: %v", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			t.Fatalf("Error rolling back transaction: %v", err)
		}
	}()

	if _, err := New().NewUser(ctx).Create(ctx, tx); err != nil {
		t.Fatalf("Error creating User: %v", err)
	}
}

func TestCreateUserAuth(t *testing.T) {
	if testDB == nil {
		t.Skip("skipping test, no DSN provided")
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	tx, err := testDB.Begin(ctx)
	if err != nil {
		t.Fatalf("Error starting transaction: %v", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			t.Fatalf("Error rolling back transaction: %v", err)
		}
	}()

	if _, err := New().NewUserAuth(ctx).Create(ctx, tx); err != nil {
		t.Fatalf("Error creating UserAuth: %v", err)
	}
}
