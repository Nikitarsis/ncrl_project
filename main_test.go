package main

import (
	"fmt"
	"strings"
	"testing"
)

func loopTest(analyzeFunc func(*string) (*[]byte, bool), strs []string) {
	stringChan := make(chan *string, 10)
	byteChan := make(chan *[]byte, 10)
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
	i := 0
	for {
		i++
		out, opened := <-byteChan
		if !opened {
			return
		}
		fmt.Printf("%d:%s\n", i, string(*out))
	}
}

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
	loopTest(analyzeFunc, strs)
}

func TestLoopRealAnalyze(t *testing.T) {
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
	reformed := []string{
		"Любить иных — тяжелый крест,",
		"А ты прекрасна без извилин,",
		"И прелести твоей секрет",
		"Разгадке жизни равносилен.",
		"Весною слышен шорох снов",
		"И шелест новостей и истин.",
		"Ты из семьи таких основ.",
		"Твой смысл, как воздух, бескорыстен.",
		"Легко проснуться и прозреть,",
		"Словесный сор из сердца вытрясть",
		"И жить, не засоряясь впредь,",
		"Все это — небольшая хитрость.",
	}
	analyzer := getStringAnalyzer(true, true)
	analyzeFunc := func(s *string) (*[]byte, bool) {
		ret, err := analyzer.AnalyzeString(s).GetJson()
		if err != nil {
			t.Error(err)
			return nil, false
		}
		return ret, true
	}
	loopTest(analyzeFunc, strs)
	fmt.Printf("\n\n%s\n\n", "REFORMED")
	loopTest(analyzeFunc, reformed)
}

func TestCyclicWriting(t *testing.T) {
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
	byteChan := make(chan *[]byte, 10)
	go func() {
		for _, str := range strs {
			bytes := []byte(str)
			byteChan <- &bytes
		}
		close(byteChan)
	}()
	CyclicWriting(
		true,
		byteChan,
		func(s string) {},
		func(s string) {},
		func() bool { return false },
	)
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
