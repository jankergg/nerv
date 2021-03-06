package service

import (
	"encoding/json"
	"fmt"
	"github.com/ChaosXu/nerv/lib/yml"
	"io/ioutil"
	"os"
)

const (
	FileBeatIndexDir = "../../data/filebeat"
	FileBeatIndex    = "../../data/filebeat/config_index.json"
	FileBeatConfig   = "../../log/config/filebeat.yml"
)

// LogConfigService merge all filebeat's configs into one
type LogConfigService struct {
	AppService *AppService `inject:"App"`
}

func (p *LogConfigService) Init() error {
	p.AppService.On("Update", p)
	return nil
}

func (p *LogConfigService) Handle(event string, data interface{}) {
	fmt.Printf("%s:%+v\n", event, data)
}

// Add a filebeat config of app
func (p *LogConfigService) Add(file string) error {
	configs := []string{}
	if _, err := os.Stat(FileBeatIndex); err == nil {
		buf, err := ioutil.ReadFile(FileBeatIndex)
		if err != nil {
			fmt.Printf("read log config failed when add. %s\n", err.Error())
		}

		json.Unmarshal(buf, &configs)
	} else {
		configs = append(configs, FileBeatConfig)
	}

	for _, v := range configs {
		if v == file {
			return nil
		}
	}
	configs = append(configs, file)
	buf, err := json.Marshal(configs)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(FileBeatIndexDir, os.ModeDir|os.ModePerm); err != nil {
		fmt.Printf("write log config failed when add. file: %s\n", err.Error())
		return err
	}
	err = ioutil.WriteFile(FileBeatIndex, buf, os.ModePerm)
	if err != nil {
		fmt.Printf("write log config failed when add. file: %s\n", err.Error())
		return err
	}

	m := yml.NewMerger(FileBeatConfig)
	for _, item := range configs {
		m.Add(item)
	}
	return m.Merge()
}

// Remove a filebeat config of app
func (p *LogConfigService) Remove(file string) error {
	if _, err := os.Stat(FileBeatIndex); err != nil {
		return nil
	}
	buf, err := ioutil.ReadFile("../data/filebeat/config_index.json")
	if err != nil {
		fmt.Printf("read log config failed when remove. %s \n", err.Error())
	}
	configs := []string{}
	json.Unmarshal(buf, &configs)
	for i, v := range configs {
		if v == file {
			configs = append(configs[:i], configs[i+1:]...)
			break
		}
	}

	buf, err = json.Marshal(configs)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(FileBeatIndexDir, os.ModeDir|os.ModePerm); err != nil {
		fmt.Printf("write log config failed when remove. file: %s\n", err.Error())
		return err
	}
	err = ioutil.WriteFile(FileBeatIndex, buf, os.ModePerm)
	if err != nil {
		return err
	}

	m := yml.NewMerger(FileBeatConfig)
	for _, item := range configs {
		fmt.Printf("write log config failed when remove. %s\n", err.Error())
		m.Add(item)
	}
	return m.Merge()
}
