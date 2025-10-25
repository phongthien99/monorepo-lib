package core

import (
	"testing"
)

type TestConfig struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
	}
	Features map[string]bool
	Tags     []string
}

func TestDefaultMerge_SimpleOverride(t *testing.T) {
	dst := &TestConfig{}
	dst.Server.Host = "localhost"
	dst.Server.Port = 8080

	src := &TestConfig{}
	src.Server.Port = 9090

	if err := DefaultMerge(dst, src); err != nil {
		t.Fatalf("DefaultMerge failed: %v", err)
	}

	// Host should remain (src has zero value)
	if dst.Server.Host != "localhost" {
		t.Errorf("Expected host=localhost, got %s", dst.Server.Host)
	}

	// Port should be overridden
	if dst.Server.Port != 9090 {
		t.Errorf("Expected port=9090, got %d", dst.Server.Port)
	}
}

func TestDefaultMerge_DeepMerge(t *testing.T) {
	dst := &TestConfig{}
	dst.Server.Host = "localhost"
	dst.Server.Port = 8080
	dst.Database.Host = "dbhost"
	dst.Database.Port = 5432

	src := &TestConfig{}
	src.Server.Port = 9090
	src.Database.Username = "admin"
	src.Database.Password = "secret"

	if err := DefaultMerge(dst, src); err != nil {
		t.Fatalf("DefaultMerge failed: %v", err)
	}

	// Server.Host should remain
	if dst.Server.Host != "localhost" {
		t.Errorf("Expected server.host=localhost, got %s", dst.Server.Host)
	}

	// Server.Port should be overridden
	if dst.Server.Port != 9090 {
		t.Errorf("Expected server.port=9090, got %d", dst.Server.Port)
	}

	// Database.Host should remain
	if dst.Database.Host != "dbhost" {
		t.Errorf("Expected database.host=dbhost, got %s", dst.Database.Host)
	}

	// Database.Username should be merged
	if dst.Database.Username != "admin" {
		t.Errorf("Expected database.username=admin, got %s", dst.Database.Username)
	}

	// Database.Password should be merged
	if dst.Database.Password != "secret" {
		t.Errorf("Expected database.password=secret, got %s", dst.Database.Password)
	}
}

func TestDefaultMerge_MapMerge(t *testing.T) {
	dst := &TestConfig{}
	dst.Features = map[string]bool{
		"feature1": true,
		"feature2": false,
	}

	src := &TestConfig{}
	src.Features = map[string]bool{
		"feature2": true,
		"feature3": true,
	}

	if err := DefaultMerge(dst, src); err != nil {
		t.Fatalf("DefaultMerge failed: %v", err)
	}

	// feature1 should remain
	if !dst.Features["feature1"] {
		t.Error("Expected feature1=true")
	}

	// feature2 should be overridden
	if !dst.Features["feature2"] {
		t.Error("Expected feature2=true (overridden)")
	}

	// feature3 should be added
	if !dst.Features["feature3"] {
		t.Error("Expected feature3=true (added)")
	}
}

func TestDefaultMerge_SliceOverride(t *testing.T) {
	dst := &TestConfig{}
	dst.Tags = []string{"tag1", "tag2"}

	src := &TestConfig{}
	src.Tags = []string{"tag3", "tag4"}

	if err := DefaultMerge(dst, src); err != nil {
		t.Fatalf("DefaultMerge failed: %v", err)
	}

	// Slice should be completely overridden
	if len(dst.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(dst.Tags))
	}

	if dst.Tags[0] != "tag3" || dst.Tags[1] != "tag4" {
		t.Errorf("Expected [tag3, tag4], got %v", dst.Tags)
	}
}

func TestDefaultMerge_ZeroValueNotOverride(t *testing.T) {
	dst := &TestConfig{}
	dst.Server.Host = "localhost"
	dst.Server.Port = 8080

	src := &TestConfig{}
	// src has all zero values

	if err := DefaultMerge(dst, src); err != nil {
		t.Fatalf("DefaultMerge failed: %v", err)
	}

	// Values should remain (src has zero values)
	if dst.Server.Host != "localhost" {
		t.Errorf("Expected host=localhost, got %s", dst.Server.Host)
	}

	if dst.Server.Port != 8080 {
		t.Errorf("Expected port=8080, got %d", dst.Server.Port)
	}
}

func TestDefaultMerge_EmptySliceNotOverride(t *testing.T) {
	dst := &TestConfig{}
	dst.Tags = []string{"tag1", "tag2"}

	src := &TestConfig{}
	src.Tags = []string{} // Empty slice

	if err := DefaultMerge(dst, src); err != nil {
		t.Fatalf("DefaultMerge failed: %v", err)
	}

	// Tags should remain (src slice is empty)
	if len(dst.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(dst.Tags))
	}
}

func TestShallowMerge(t *testing.T) {
	dst := &TestConfig{}
	dst.Server.Host = "localhost"
	dst.Server.Port = 8080
	dst.Database.Host = "dbhost"

	src := &TestConfig{}
	src.Server.Host = "newhost"
	src.Server.Port = 9090
	// src.Database is zero value

	if err := ShallowMerge(dst, src); err != nil {
		t.Fatalf("ShallowMerge failed: %v", err)
	}

	// Entire struct should be replaced
	if dst.Server.Host != "newhost" {
		t.Errorf("Expected server.host=newhost, got %s", dst.Server.Host)
	}

	if dst.Server.Port != 9090 {
		t.Errorf("Expected server.port=9090, got %d", dst.Server.Port)
	}

	// Database should be zero (because src.Database is zero)
	if dst.Database.Host != "" {
		t.Errorf("Expected database.host to be empty, got %s", dst.Database.Host)
	}
}

func TestDefaultMerge_PointerFields(t *testing.T) {
	type ConfigWithPointer struct {
		Name  string
		Value *int
	}

	dst := &ConfigWithPointer{
		Name: "test",
	}
	dstVal := 100
	dst.Value = &dstVal

	src := &ConfigWithPointer{}
	srcVal := 200
	src.Value = &srcVal

	if err := DefaultMerge(dst, src); err != nil {
		t.Fatalf("DefaultMerge failed: %v", err)
	}

	// Name should remain
	if dst.Name != "test" {
		t.Errorf("Expected name=test, got %s", dst.Name)
	}

	// Value should be overridden
	if dst.Value == nil || *dst.Value != 200 {
		t.Errorf("Expected value=200, got %v", dst.Value)
	}
}
