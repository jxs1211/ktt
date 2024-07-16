package services

import (
	. "ktt/backend/storage"
	_ "ktt/backend/utils/proxy"
)

type Configuration struct {
	content string
}

func NewConfiguration(content string) *Configuration {
	return &Configuration{
		content: content,
	}
}

func (c *Configuration) Load() bool {

	return true
}

func (c *Configuration) Validate() bool {
	return true
}
