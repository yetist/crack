package crack

import (
	"time"
)

type Cracker interface {
	GenEntry(path string, entry chan<- string)
	CrackIt(pass string) (bool, error)
}

func getCrack(data Cracker, in <-chan string, end chan<- bool) {
	for {
		select {
		case <-time.After(time.Second * 2):
			end <- true
			break
		case pass := <-in:
			if ret, _ := data.CrackIt(pass); ret {
				end <- true
				break
			}
		}
	}
}

func Crack(files []string, data Cracker) {
	ch := make(chan string)
	defer close(ch)

	end := make(chan bool)
	defer close(end)

	for _, path := range files {
		go data.GenEntry(path, ch)
	}
	go getCrack(data, ch, end)
	<-end
}
