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
	finishCallback func(),
	fileNames ...string,
) {
	defer finishCallback()
	defer close(out)
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
	in <-chan *[]byte,
	warn func(s string),
	log func(s string),
	processingStoped func() bool,
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
		writers = append(writers, bufio.NewWriter(os.Stdout))
	}
	defer func() {
		for _, writer := range writers {
			err := writer.Flush()
			if err != nil {
				warn(fmt.Sprintf("Writer wasn't flushed: %s", err.Error()))
			}
		}
	}()
	for i := 0; ; i++ {
		line, opened := <-in
		if !opened {
			log("Channel closed: writing stopped")
			return
		}
		for _, writer := range writers {
			nn, err := writer.Write(*line)
			if err != nil {
				warn(err.Error())
				continue
			}
			if nn != len(*line) {
				warn("Writing wasn't completed")
				continue
			}
			writer.WriteByte('\n')
			if i >= 100 && i%100 == 0 {
				err := writer.Flush()
				if err != nil {
					warn(fmt.Sprintf("Writer wasn't flushed: %s", err.Error()))
				}
			}
		}
		if processingStoped() && len(in) == 0 {
			log("Writing stopped")
			return
		}
	}
}

func loopRoutine(
	out chan<- *[]byte,
	in <-chan *string,
	warn func(string),
	log func(string),
	shouldStop func() bool,
	finishCallback func(),
	analyzeFunction func(*string) (*[]byte, bool),
) {
	defer finishCallback()
	for {
		str, opened := <-in
		if !opened {
			log("Channel closed: processing stopped")
			return
		}
		ret, check := analyzeFunction(str)
		if !check {
			warn("String wasn't processed")
			continue
		}
		out <- ret
		if shouldStop() {
			log("Processing stopped")
			close(out)
			return
		}
	}
}
