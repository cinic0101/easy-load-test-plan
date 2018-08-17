package ui

import (
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const StoragePath = "ui/projects/projects.yaml"

type Storage struct {
	AllProjects Projects
}

type Projects struct {
	Projects []TestProject
}

type TestProject struct {
	ID string
	Name string
	Desc string
	Time string
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

	p := &Projects{
	}

	if data, err := yaml.Marshal(p); err == nil {
		ioutil.WriteFile(StoragePath, data, 0755)
	}
}

func (s *Storage) Load() error {
	testPlans := &Projects{}

	buf, err := ioutil.ReadFile(StoragePath)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, &testPlans); err != nil {
		return err
	}

	s.AllProjects = *testPlans

	return err
}

func (s *Storage) Save() {
	if data, err := yaml.Marshal(s.AllProjects); err == nil {
		ioutil.WriteFile(StoragePath, data, 0777)
	}
}

func (s *Storage) AddNewTestPlan(id string, name string, desc string) {
	s.Load()

	s.AllProjects.Projects = append([]TestProject{{
		ID: id,
		Name: name,
		Desc: desc,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}}, s.AllProjects.Projects...)

	s.Save()
}