package journal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ItemMemoryRepository struct {
	data []*JournalItem
}

func NewMemoryRepo() *ItemMemoryRepository {
	return &ItemMemoryRepository{
		data: make([]*JournalItem, 0, 10),
	}
}

func (repo *ItemMemoryRepository) AddItemInRepo(r *http.Request) (*JournalItem, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("package journal, AddItem err : %v", err.Error())
	}
	r.Body.Close()

	newItem := &JournalItem{}
	err = json.Unmarshal(body, newItem)
	if err != nil {
		return nil, fmt.Errorf("package journal, AddItem err : %v", err.Error())
	}
	if newItem.Sender == "" || newItem.Recipient == "" {
		return nil, errors.New("bad JSON field")
	}

	repo.data = append(repo.data, newItem)
	return newItem, nil
}

func (repo *ItemMemoryRepository) CreateGrapgh() (*Graph, error) {
	nodesMap := make(map[string][]string)
	activityMap := make(map[string]int32)

	for _, item := range repo.data {
		val, ok := nodesMap[item.Sender]

		if !ok {
			nodesMap[item.Sender] = append(nodesMap[item.Sender], item.Recipient)
			activityMap[item.Sender] += 1
			nodesMap[item.Recipient] = append(nodesMap[item.Recipient], item.Sender)
			activityMap[item.Recipient] += 1
			continue
		}

		if !contains(val, item.Recipient) {
			nodesMap[item.Sender] = append(nodesMap[item.Sender], item.Recipient)
			activityMap[item.Sender] += 1
			nodesMap[item.Recipient] = append(nodesMap[item.Recipient], item.Sender)
			activityMap[item.Recipient] += 1
		}
	}

	maxActivity, err := findMaxActivity(activityMap)
	if err != nil {
		return nil, fmt.Errorf("package journal, GetGrapgh, findMaxActivity : %v", err.Error())
	}

	newGraph := &Graph{}
	newGraph.Nodes = nodesMap
	newGraph.ActivityIndex = getActivityIndex(activityMap, maxActivity)

	return newGraph, nil
}

func contains(sl []string, user string) bool {
	for _, item := range sl {
		if item == user {
			return true
		}
	}
	return false
}

func findMaxActivity(activity map[string]int32) (int32, error) {
	maxValue := int32(-1)
	for _, value := range activity {
		if value > maxValue {
			maxValue = value
		}
	}
	if maxValue == -1 {
		return 0, fmt.Errorf("empty activity map")
	}
	return maxValue, nil
}

func getActivityIndex(activity map[string]int32, max_value int32) map[string]string {
	activityIndex := make(map[string]string)

	for key, value := range activity {
		tmp := float32(value) / float32(max_value)
		if tmp < 0.34 {
			activityIndex[key] = "min"
		} else if tmp > 0.66 {
			activityIndex[key] = "max"
		} else {
			activityIndex[key] = "mid"
		}
	}

	return activityIndex
}

func AddStaticJournal(repo *ItemMemoryRepository) {
	newJournal := []*JournalItem{
		{Recipient: "Cat", Sender: "Alex"},
		{Recipient: "Mark", Sender: "Tom"},
		{Recipient: "Lancy", Sender: "Mark"},
		{Recipient: "Mark", Sender: "Alex"},
	}

	repo.data = newJournal
	fmt.Println("\t Communication journal")
	for _, item := range repo.data {
		fmt.Println("Recipient: ", item.Recipient, "\t", "Sender: ", item.Sender)
	}
}
