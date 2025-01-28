package main

import (
	"strconv"
)

type FLAGS string

const (
	SAVE_STRING FLAGS = "should_save_string"
	COMBO       FLAGS = "should_save_combo"
	NO_INPUT    FLAGS = "stop_input_pipeline"
	NO_OUTPUT   FLAGS = "stop_output_pipeline"
	LOG_TRASH   FLAGS = "log_unnecessary_information"
)

type Config struct {
	flags           map[FLAGS]bool
	readingFiles    []string
	outputFiles     []string
	sizeOfChan      int
	numOfGoroutines int
}

func GetConfig() Config {
	return Config{
		map[FLAGS]bool{
			SAVE_STRING: false,
			COMBO:       false,
			NO_INPUT:    false,
			NO_OUTPUT:   false,
			LOG_TRASH:   false,
		},
		make([]string, 0),
		make([]string, 0),
		10000,
		12,
	}
}

func (c *Config) flagUp(flag FLAGS) {
	c.flags[flag] = true
}

func (c Config) checkFlag(flag FLAGS) bool {
	return c.flags[flag]
}

func (c *Config) SetReadingFiles(s ...string) {
	c.readingFiles = s
}

func (c *Config) SetOutputFiles(s ...string) {
	c.outputFiles = s
}

func (c *Config) SetSizeOfChan(s ...string) {
	if len(s) != 1 {
		panic("Incorrect array")
	}
	ret, err := strconv.Atoi(s[0])
	if err != nil {
		panic(err.Error())
	}
	if ret < 0 {
		panic("Size of chan cannot be negative")
	}
	c.sizeOfChan = ret
}

func (c *Config) SetNumOfGoroutines(s ...string) {
	if len(s) != 1 {
		panic("Incorrect array")
	}
	ret, err := strconv.Atoi(s[0])
	if err != nil {
		panic(err.Error())
	}
	if ret < 0 {
		panic("Size of chan cannot be negative")
	}
	c.numOfGoroutines = ret
}

func (c Config) GetNumOfGoroutines() int {
	return c.numOfGoroutines
}

func (c Config) GetSizeOfChan() int {
	return c.sizeOfChan
}

func (c Config) GetReadingFiles() []string {
	return c.readingFiles
}

func (c Config) GetOutputFiles() []string {
	return c.outputFiles
}
