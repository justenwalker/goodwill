package mage

import (
	"os"
	"sync"
)

func LoadOnce(file string) func() []byte {
	var once sync.Once
	var err error
	var data []byte
	return func() []byte {
		once.Do(func() {
			data, err = os.ReadFile(file)
		})
		if err != nil {
			debug.Fatalln("ERROR", err)
		}
		return data
	}
}

func DirOnce(dir string, mode os.FileMode) func() string {
	var once sync.Once
	var err error
	return func() string {
		once.Do(func() {
			err = os.MkdirAll(dir, mode)
		})
		if err != nil {
			debug.Fatalln("ERROR", err)
		}
		return dir
	}
}

func StringOnce(fn func() (string, error)) func() string {
	var once sync.Once
	var err error
	var data string
	return func() string {
		once.Do(func() {
			data, err = fn()
		})
		if err != nil {
			debug.Fatalln("ERROR", err)
		}
		return data
	}
}
