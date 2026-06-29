package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServFlowDAGExecution(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/workflows/define", handleDefine)
	mux.HandleFunc("/api/workflows/execute", handleExecute)
	mux.HandleFunc("/api/workflows/instances/", handleGetInstance)

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	// 1. Define DAG Workflow: Task A -> Task B (depends on A)
	defPayload := WorkflowDef{
		ID: "onboarding-flow",
		Tasks: []Task{
			{Name: "CreateUser", DependsOn: nil, Action: "success"},
			{Name: "SendWelcomeEmail", DependsOn: []string{"CreateUser"}, Action: "success"},
		},
	}
	body, _ := json.Marshal(defPayload)
	resp, err := http.Post(testServer.URL+"/api/workflows/define", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to define workflow: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected StatusCreated, got %d", resp.StatusCode)
	}

	// 2. Execute Workflow Instance
	execPayload := map[string]string{"workflow_id": "onboarding-flow"}
	execBody, _ := json.Marshal(execPayload)
	execResp, err := http.Post(testServer.URL+"/api/workflows/execute", "application/json", bytes.NewReader(execBody))
	if err != nil {
		t.Fatalf("failed to execute workflow: %v", err)
	}

	var inst WorkflowInstance
	json.NewDecoder(execResp.Body).Decode(&inst)
	execResp.Body.Close()

	if inst.Status != "running" {
		t.Fatalf("expected running workflow status, got %q", inst.Status)
	}

	// Wait briefly for background execution to complete
	time.Sleep(100 * time.Millisecond)

	// 3. Query Instance Status
	getResp, err := http.Get(testServer.URL + "/api/workflows/instances/" + inst.ID)
	if err != nil {
		t.Fatalf("failed to query instance: %v", err)
	}
	defer getResp.Body.Close()

	var finalInst WorkflowInstance
	json.NewDecoder(getResp.Body).Decode(&finalInst)

	if finalInst.Status != "completed" {
		t.Errorf("expected workflow completion, got %q. Logs: %v", finalInst.Status, finalInst.Logs)
	}

	if finalInst.TaskStates["SendWelcomeEmail"].Status != "completed" {
		t.Errorf("expected SendWelcomeEmail to be completed, got %q", finalInst.TaskStates["SendWelcomeEmail"].Status)
	}
}
