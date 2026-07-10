package storage

import (
	"testing"
)

func TestNewServStoreWorkflowStore(t *testing.T) {
	s := NewServStoreWorkflowStore(nil)
	if s == nil {
		t.Fatal("expected NewServStoreWorkflowStore to return non-nil")
	}
	if s.GetClient() != nil {
		t.Error("expected Client to be nil")
	}
}

func TestLoadDefinitionsMissingClient(t *testing.T) {
	s := NewServStoreWorkflowStore(nil)
	_, err := s.LoadDefinitions()
	if err == nil {
		t.Error("expected error loading with nil client")
	}
}

func TestSaveDefinitionsMissingClient(t *testing.T) {
	s := NewServStoreWorkflowStore(nil)
	err := s.SaveDefinitions(nil)
	if err == nil {
		t.Error("expected error saving with nil client")
	}
}

func TestLoadInstancesMissingClient(t *testing.T) {
	s := NewServStoreWorkflowStore(nil)
	_, err := s.LoadInstances()
	if err == nil {
		t.Error("expected error loading with nil client")
	}
}

func TestSaveInstancesMissingClient(t *testing.T) {
	s := NewServStoreWorkflowStore(nil)
	err := s.SaveInstances(nil)
	if err == nil {
		t.Error("expected error saving with nil client")
	}
}
