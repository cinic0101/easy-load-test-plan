package ui

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const StoragePath = "ui/plans/plans.yaml"

type Storage struct {
	File TestPlans
}

type TestPlans struct {
	Plans []TestPlan
}

type TestPlan struct {
	ID string
	Name string
	Plan []Test
}

type Test struct {
	Name string
	Config string
}

func (s *Storage) CreateIfNotExists() {
	if _, err := os.Stat(StoragePath); err == nil {
		return
	}

	p := &TestPlans{
	}

	if data, err := yaml.Marshal(p); err == nil {
		ioutil.WriteFile(StoragePath, data, 0755)
	}
}

func (s *Storage) Load() error {
	testPlans := &TestPlans{}

	buf, err := ioutil.ReadFile(StoragePath)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, &testPlans); err != nil {
		return err
	}

	s.File = *testPlans

	return err
}

func (s *Storage) Save() {
	if data, err := yaml.Marshal(s.File); err == nil {
		ioutil.WriteFile(StoragePath, data, 0777)
	}
}

func (s *Storage) AddNewTestPlan(id string, name string)  {
	s.Load()

	s.File.Plans = append(s.File.Plans, TestPlan{
		ID: id,
		Name: name,
	})

	s.Save()
}