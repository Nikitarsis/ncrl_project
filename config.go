package main

import "strconv"

type Config struct {
	shouldSaveString   bool
	shouldCountCombo   bool
	stopInputPipeline  bool
	stopOutputPipeline bool
	readingFiles       []string
	outputFiles        []string
	sizeOfChan         int
	numOfGoroutines    int
}

func GetConfig() Config {
	return Config{false, false, false, false, make([]string, 0), make([]string, 0), 10000, 12}
}

func (c *Config) SaveString() {
	c.shouldSaveString = true
}

func (c *Config) CountCombo() {
	c.shouldCountCombo = true
}

func (c *Config) TurnOffOutputPipeline() {
	c.stopOutputPipeline = true
}

func (c *Config) TurnOffInputPipeline() {
	c.stopInputPipeline = true
}

func (c *Config) SetReadingFiles(s ...string) {
	c.readingFiles = s
}

func (c *Config) SetOutputFiles(s ...string) {
	c.outputFiles = s
}

func (c Config) ShouldSaveString() bool {
	return c.shouldSaveString
}

func (c Config) ShouldCountCombo() bool {
	return c.shouldCountCombo
}

func (c Config) ShouldStopInPipeline() bool {
	return c.stopInputPipeline
}

func (c Config) ShouldStopOutPipeline() bool {
	return c.stopOutputPipeline
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
