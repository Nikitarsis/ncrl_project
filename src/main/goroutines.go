package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func CyclicReading(
	readFromStdIn bool,
	out chan<- *string,
	warn func(string),
	log func(string),
	fileNames ...string,
) {
	files := make([]*os.File, len(fileNames))
	for i, name := range fileNames {
		file, err := os.Open(name)
		if err != nil {
			warn(err.Error())
			continue
		}
		log(fmt.Sprintf("File %s opened", name))
		files[i] = file
		defer file.Close()
	}
	for _, file := range files {
		reader := bufio.NewReader(file)
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			warn(err.Error())
			continue
		}
		out <- &line
	}
	if !readFromStdIn {
		log("Finished without STDIN")
		return
	}
	log("Starting read from STDIN")
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			warn(err.Error())
			continue
		}
		out <- &line
	}
	log("Finished STDIN")
}

func CyclicWriting(
	writeToPipeline bool,
	out <-chan *string,
	warn func(s string),
	log func(s string),
	fileNames ...string,
) {
	writers := make([]*bufio.Writer, len(fileNames))
	for i, name := range fileNames {
		file, err := os.Open(name)
		if err != nil {
			warn(err.Error())
			continue
		}
		log(fmt.Sprintf("File %s opened", name))
		writers[i] = bufio.NewWriter(file)
		defer file.Close()
	}
	if writeToPipeline {
		writers = append(writers, bufio.NewWriter(os.Stdin))
	}
	for {
		line := <-out
		for _, writer := range writers {
			nn, err := writer.WriteString(*line)
			if err != nil {
				warn(err.Error())
				continue
			}
			if nn != len(*line) {
				warn("Writing wasn't completed")
				continue
			}
			writer.WriteByte('\n')
		}
	}
}

func loopRoutine(out chan<- *string, in <-chan *string, warn func(string), function func(*string) (*[]byte, bool)) {
	for {
		str := <-in
		st, check := function(str)
		if !check {
			warn("String wasn't processed")
			continue
		}
		ret := string(*st)
		out <- &ret
	}
}
