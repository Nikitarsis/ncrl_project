package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestCyclicReading(t *testing.T) {
	stringChan := make(chan *string, 1000)
	fmt.Println("STARTED READING")
	go CyclicReading(
		false,
		stringChan,
		func(s string) {},
		func(s string) {},
		func() {},
		"data/test",
	)
	i := 1
	for {
		str, opened := <-stringChan
		if !opened {
			fmt.Println("ENDED READING")
			return
		} else {
			fmt.Print(i)
			fmt.Print(":")
			i++
			fmt.Println(*str)
		}
	}
}

func TestLoop(t *testing.T) {
	stringChan := make(chan *string, 10)
	byteChan := make(chan *[]byte, 10)
	strs := []string{
		"Мы всѣ учились понемногу,",
		"Чему нибудь и какъ нибудь:",
		"Такъ воспитаньемъ, слава Богу,",
		"У насъ немудрено блеснуть.",
		"Онѣгинъ былъ, по мнѣнью многихъ",
		"(Судей рѣшительныхъ и строгихъ),",
		"Ученый малый, но педантъ.",
		"Имѣлъ онъ счастливый талантъ",
		"Безъ принужденья въ разговорѣ",
		"Коснуться до всего слегка,",
		"Съ ученымъ видомъ знатока",
		"Хранить молчанье въ важномъ спорѣ,",
		"И возбуждать улыбку дамъ",
		"Огнемъ нежданыхъ эпиграммъ",
	}
	analyzeFunc := func(s *string) (*[]byte, bool) {
		ret := []byte(*s)
		return &ret, true
	}
	go func() {
		for _, str := range strs {
			stringChan <- &str
		}
		close(stringChan)
	}()
	go loopRoutine(
		byteChan,
		stringChan,
		func(s string) {},
		func(s string) {},
		func() bool { return false },
		func() {},
		analyzeFunc,
	)
	for {
		out, opened := <-byteChan
		if !opened {
			return
		}
		fmt.Print(*out)
		fmt.Println()
	}
}

func TestBasicFunction(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Recovered: %v\n", r)
		}
	}()
	args := strings.Split("--output_file data/parsed --input_file data/test -tmsI", " ")
	basicFunction(args)
}
