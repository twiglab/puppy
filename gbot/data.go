package gbot

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Attr map[string]any

type Data struct {
	Projects []Project
	Attr     Attr
}

type Project struct {
	ID    string
	Name  string
	Areas []Area
	Attr  Attr
}

type Area struct {
	ID      string
	Name    string
	Cameras []string
	Attr    Attr
}

type YamlDataLoad struct {
	yamlFile string
	data     Data
}

func (l *YamlDataLoad) Load() error {
	f, err := os.Open(l.yamlFile)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	return dec.Decode(&l.yamlFile)
}

func (l *YamlDataLoad) Each(proj, area string, f func(Project, Area) error) error {
	for _, pr := range l.data.Projects {
		if proj == "*" || pr.ID == proj {
			for _, a := range pr.Areas {
				if area == "*" || a.ID == area {
					if err := f(pr, a); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
