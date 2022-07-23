package journal

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestAddItem(t *testing.T) {
	repo := NewMemoryRepo()

	newItem := &JournalItem{
		Sender:    "Tom Clancy",
		Recipient: "Deep Thought",
	}

	// correct input
	bodyReader := strings.NewReader(`{"sender": "Tom Clancy", "recipient": "Deep Thought"}`)
	req := httptest.NewRequest("GET", "/api/addRecord", bodyReader)
	req.Header.Set("Content-Type", "application/json")

	checkItem, err := repo.AddItemInRepo(req)
	last_idx := len(repo.data) - 1
	if err != nil {
		t.Errorf(err.Error(), "correct input test: unexpected error")
		return
	}
	if !reflect.DeepEqual(checkItem, newItem) || !reflect.DeepEqual(repo.data[last_idx], newItem) {
		t.Errorf("correct input test: result doesn't match")
		return
	}

	// bad JSON
	bodyReader = strings.NewReader(`"}`)
	req = httptest.NewRequest("GET", "/api/addRecord", bodyReader)
	req.Header.Set("Content-Type", "application/json")

	_, err = repo.AddItemInRepo(req)
	if err == nil {
		t.Errorf(err.Error(), "bad JSON test")
		return
	}

	// bad JSON field (sender)
	bodyReader = strings.NewReader(`{"sendr": "Tom Clancy", "recipient": "Deep Thought"}`)
	req = httptest.NewRequest("GET", "/api/addRecord", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	_, err = repo.AddItemInRepo(req)

	if err != nil && err.Error() != "bad JSON field" {
		t.Errorf(err.Error(), "bad JSON field (sender) test")
		return
	}

	// bad JSON field (recipient)
	bodyReader = strings.NewReader(`{"sender": "Tom Clancy", "recipnt": "Deep Thought"}`)
	req = httptest.NewRequest("GET", "/api/addRecord", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	_, err = repo.AddItemInRepo(req)

	if err != nil && err.Error() != "bad JSON field" {
		t.Errorf(err.Error(), "bad JSON field (recipient) test")
		return
	}
}

func TestGetGrapgh(t *testing.T) {
	// testing that the function works correctly
	repo := NewMemoryRepo()
	repo.data = []*JournalItem{
		{Recipient: "Cat", Sender: "Alex"},
		{Recipient: "Mark", Sender: "Tom"},
		{Recipient: "Lancy", Sender: "Mark"},
		{Recipient: "Mark", Sender: "Alex"},
		{Recipient: "Alex", Sender: "Cat"},
		{Recipient: "Alex", Sender: "Sonya"},
		{Recipient: "Alex", Sender: "Bob"},
		{Recipient: "Joe", Sender: "Dave"},
		{Recipient: "Joe", Sender: "Dave"},
		{Recipient: "Cat", Sender: "Bo"},
	}

	expectedGraph := &Graph{
		Nodes:         map[string][]string{},
		ActivityIndex: map[string]string{},
	}
	expectedGraph.Nodes["Tom"] = []string{"Mark"}
	expectedGraph.Nodes["Joe"] = []string{"Dave"}
	expectedGraph.Nodes["Bo"] = []string{"Cat"}
	expectedGraph.Nodes["Bob"] = []string{"Alex"}
	expectedGraph.Nodes["Dave"] = []string{"Joe"}
	expectedGraph.Nodes["Alex"] = []string{"Cat", "Mark", "Sonya", "Bob"}
	expectedGraph.Nodes["Cat"] = []string{"Alex", "Bo"}
	expectedGraph.Nodes["Mark"] = []string{"Tom", "Lancy", "Alex"}
	expectedGraph.Nodes["Lancy"] = []string{"Mark"}
	expectedGraph.Nodes["Sonya"] = []string{"Alex"}

	expectedGraph.ActivityIndex["Tom"] = "min"
	expectedGraph.ActivityIndex["Joe"] = "min"
	expectedGraph.ActivityIndex["Bo"] = "min"
	expectedGraph.ActivityIndex["Bob"] = "min"
	expectedGraph.ActivityIndex["Dave"] = "min"
	expectedGraph.ActivityIndex["Alex"] = "max"
	expectedGraph.ActivityIndex["Cat"] = "mid"
	expectedGraph.ActivityIndex["Mark"] = "max"
	expectedGraph.ActivityIndex["Lancy"] = "min"
	expectedGraph.ActivityIndex["Sonya"] = "min"

	newGraph, err := repo.CreateGrapgh()
	if err != nil {
		t.Error("TestGetGrapgh", err.Error())
		return
	}
	if !reflect.DeepEqual(newGraph, expectedGraph) {
		t.Errorf("TestGetGrapgh: result doesn't match")
	}
}
