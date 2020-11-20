package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
)



const (
	RealFilename = "real.txt" // keep send package count number in this file

)

// increment value and return as string
func IncrementValue(pastValue string)(newValue string){
	newValueInt, _ := strconv.Atoi(pastValue)
	return strconv.Itoa(newValueInt + 1)
}

//read from file and write to file
func readAndWrite(filename string, m *sync.Mutex ) (str string ,err error){
	m.Lock()
	defer m.Unlock()

	initialValue := "1"
	someText , err:= readFile(filename)
	if err == nil {
		newValue := IncrementValue(string(someText))
		err = ioutil.WriteFile(filename, []byte(newValue), 0644)
		if err != nil {
			return "", err
		}
		return newValue, nil
	} else {
		err = ioutil.WriteFile(filename,[]byte(initialValue), 0644)
		if err != nil {
			return initialValue,err
		}
	}

	return "",nil
}

// read from file and return the text
func readFile(filename string) (str string ,err error){
	if _, err := os.Stat(filename); err == nil {
		someText, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Error reading")
			return "", err
		}
		return string(someText), nil
	} else {
		return "",err
	}
	return "",nil
}

