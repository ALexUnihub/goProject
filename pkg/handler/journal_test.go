package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	journal "github.com/ALexUnihub/goProject/pkg/jornal_repo"
	"github.com/golang/mock/gomock"
)

func TestAddRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := journal.NewMockJournalRepo(ctrl)
	jrnlHandler := &JournalHandler{
		JournalRepo: repo,
	}

	// good req
	newItem := &journal.JournalItem{
		Sender:    "Tom Clancy",
		Recipient: "Deep Thought",
	}

	bodyReader := strings.NewReader(`{"sender": "Tom Clancy", "recipient": "Deep Thought"}`)
	req := httptest.NewRequest("GET", "/api/record/", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	repo.EXPECT().AddItemInRepo(req).Return(newItem, nil)

	jrnlHandler.AddRecord(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err %v", err.Error())
		return
	}

	checkItem := &journal.JournalItem{}
	err = json.Unmarshal(body, checkItem)
	if err != nil {
		t.Errorf("json.Unmarshal err: %v", err.Error())
		return
	}

	if !reflect.DeepEqual(checkItem, newItem) {
		t.Errorf("result doesn't match")
		return
	}

	// bad json data
	bodyReader = strings.NewReader(`{"sr": "Tom Clancy", "recipient": "Deep Thought"}`)
	req = httptest.NewRequest("GET", "/api/record/", bodyReader)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	repo.EXPECT().AddItemInRepo(req).Return(nil, errors.New("bad JSON field"))

	jrnlHandler.AddRecord(w, req)

	resp = w.Result()

	if resp.StatusCode != 500 {
		t.Errorf("expected status code 500, got: %v", resp.StatusCode)
		return
	}
}

func TestGetGraph(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := journal.NewMockJournalRepo(ctrl)
	jrnlHandler := &JournalHandler{
		JournalRepo: repo,
	}

	expectedGraph := &journal.Graph{
		Nodes:         map[string][]string{},
		ActivityIndex: map[string]string{},
	}

	// testing that the function works correctly
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

	req := httptest.NewRequest("GET", "/api/graph/", nil)
	w := httptest.NewRecorder()

	repo.EXPECT().CreateGrapgh().Return(expectedGraph, nil)
	jrnlHandler.GetGraph(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll err %v", err.Error())
		return
	}

	gotGraph := &journal.Graph{}

	err = json.Unmarshal(body, gotGraph)
	if err != nil {
		t.Errorf("json.Unmarshal err: %v", err.Error())
		return
	}

	if !reflect.DeepEqual(expectedGraph, gotGraph) {
		t.Errorf("TestGetGraph: result doesn't match")
		return
	}

	// err in CreateGrapgh() case
	w = httptest.NewRecorder()

	repo.EXPECT().CreateGrapgh().Return(nil, errors.New("pckg handler, GetGraph err"))
	jrnlHandler.GetGraph(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("epxected status code 500, got: %v", resp.StatusCode)
		return
	}
}
