package model

import "testing"

func TestModelsUseExistingTableNames(t *testing.T) {
	if got := (User{}).TableName(); got != "user" {
		t.Fatalf("user table name = %q", got)
	}
	if got := (UserStats{}).TableName(); got != "player_stats" {
		t.Fatalf("stats table name = %q", got)
	}
}
