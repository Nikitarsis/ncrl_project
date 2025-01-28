package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	args "github.com/Nikitarsis/golang_args"
	stringanalyzer "github.com/Nikitarsis/string_analyzer"
)

//CLASS := regexp.MustCompile(`[ѢѣІіѲѳѴѵ]|([ВКСфкцнгшщзхфвпрлджчсмтб]ъ[ ,.;:?!\-"'])`)
//REFORM := regexp.MustCompile(`([иИ][яеёоыеиюэ])|([ВКСфкцнгшщзхфвпрлджчсмтб][ ,.;:?!\-"'])`)
//TRASH := regexp.MustCompile(`.{,5}`)

var configSingleton Config

func getParser() *args.ArgsParser {
	builder := args.InitParserBuilder()
	builder.AddElementAtLeast(func(s ...string) { configSingleton.SetReadingFiles(s...) }, 1, "input_file", false, "i")
	builder.AddElementAtLeast(func(s ...string) { configSingleton.SetOutputFiles(s...) }, 1, "output_file", false, "o")
	builder.AddElementAtMost(func(s ...string) { configSingleton.flagUp(SAVE_STRING) }, 0, "save_strings", false, "s")
	builder.AddElementAtMost(func(s ...string) { configSingleton.flagUp(COMBO) }, 0, "count_combinations", false, "m")
	builder.AddElementAtMost(func(s ...string) { configSingleton.flagUp(NO_INPUT) }, 0, "no_in_pipeline", false, "I")
	builder.AddElementAtMost(func(s ...string) { configSingleton.flagUp(NO_OUTPUT) }, 0, "no_out_pipeline", false, "O")
	builder.AddElementAtMost(func(s ...string) { configSingleton.flagUp(LOG_TRASH) }, 0, "trace", false, "t")
	builder.AddElimentSingle(configSingleton.SetNumOfGoroutines, 1, "number_goroutines", false, "n")
	builder.AddElimentSingle(configSingleton.SetSizeOfChan, 1, "chan_size", false, "c")
	ret, err := builder.Construct()
	if err != nil {
		panic(err)
	}
	return ret
}

func getStringAnalyzer(saveStr bool, countComb bool) *stringanalyzer.StringAnalyzer {
	CLASS := regexp.MustCompile(`[ѢѣІіѲѳѴѵ]|([ВКСфкцнгшщзхфвпрлджчсмтб]ъ[ ,.;:?!\-"'])`)
	REFORM := regexp.MustCompile(`([иИ][яеёоыеиюэ])|([ВКСфкцнгшщзхфвпрлджчсмтб][ ,.;:?!\-"'])`)
	TRASH := regexp.MustCompile(`.{,5}`)

	builder := stringanalyzer.CreateSABuilder()
	if countComb {
		builder.SaveCombinations()
	}
	if saveStr {
		builder.SaveOriginalString()
	}
	builder.AddChecker("isYoficated", func(s *string) bool { return strings.ContainsAny(*s, "ёЁ") })
	builder.AddChecker("isClassical", func(s *string) bool { return CLASS.MatchString(*s) })
	builder.AddChecker("isReformed", func(s *string) bool { return REFORM.MatchString(*s) })
	builder.AddChecker("isTrash", func(s *string) bool { return TRASH.MatchString(*s) })
	return builder.Construct()
}

func basicFunction(args []string) {
	configSingleton = GetConfig() //Setup config
	logTrash := configSingleton.checkFlag(LOG_TRASH)
	log := func(str string) {
		if logTrash {
			fmt.Fprint(os.Stderr, str)
		}
	}
	warn := func(str string) {
		fmt.Fprint(os.Stderr, str)
	}
	log("started")
	parser := getParser()                                                                                   //Create parser
	parser.ParseArgs(args...)                                                                               //Parsing arguments
	analyzer := getStringAnalyzer(configSingleton.checkFlag(SAVE_STRING), configSingleton.checkFlag(COMBO)) //Constructs string analyzer
	//function that analyzes string
	analyzeFunc := func(s *string) (*[]byte, bool) {
		ret, err := analyzer.AnalyzeString(s).GetJson()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, false
		}
		return ret, true
	}
	size := configSingleton.GetSizeOfChan()       //size of chan
	numGo := configSingleton.GetNumOfGoroutines() //number of goroutines
	stringChan := make(chan *string, size)        //input chan
	byteChan := make(chan *[]byte, size)          //output chan

	//Launch reading
	go CyclicReading(
		!configSingleton.checkFlag(NO_INPUT),
		stringChan,
		warn,
		log,
		func() {},
		configSingleton.GetReadingFiles()...,
	)
	//loop function
	loop := func() {
		loopRoutine(
			byteChan,
			stringChan,
			warn,
			log,
			func() bool { return false },
			func() {},
			analyzeFunc,
		)
	}
	for i := 0; i < numGo; i++ {
		go loop()
	}
	//Launch writing
	CyclicWriting(
		!configSingleton.checkFlag(NO_OUTPUT),
		byteChan,
		warn,
		log,
		func() bool { return false },
		configSingleton.GetOutputFiles()...,
	)
}

func main() {
	basicFunction(os.Args[1:])
}
