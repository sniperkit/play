package main

import (
	"fmt"
	"os"
	"time"

	. "gist.github.com/6096872.git"

	"github.com/bradfitz/iter"
	"github.com/shurcooL/pipe"
)

func ChanCombinedOutput(outch ChanWriter, p pipe.Pipe) error {
	s := pipe.NewState(outch, outch)

	// Test interrupting the pipe after a few seconds.
	/*go func() {
		time.Sleep(5 * time.Second)
		s.Kill()
	}()*/

	err := p(s)
	if err == nil {
		err = s.RunTasks()
	}
	close(outch)
	return err
}

func main() {
	fmt.Println("Starting.\n")
	p := pipe.Script(
		pipe.Println("Building."),
		pipe.Exec("go", "build", "-o", "/Users/Dmitri/Desktop/pipe_bin", "/Users/Dmitri/Dropbox/Work/2013/GoLand/src/gist.github.com/7176504.git/main.go"),
		pipe.Println("Running."),
		//pipe.Exec("/Users/Dmitri/Desktop/pipe_bin"),
		pipe.CancellableTaskFunc(func(s, s2 *pipe.State) error {
			for i := 1; i <= 10 && !s2.IsCancelled(); i++ {
				time.Sleep(1000 * time.Millisecond)
				if i%3 != 0 {
					fmt.Fprintln(s.Stdout, i)
				} else {
					fmt.Fprintln(s.Stderr, i, "stderr")
				}
			}
			return nil
		}),
		pipe.Println("Done."),
	)

	for _ = range iter.N(2) {
		outch := make(ChanWriter)

		go func() {
			err := ChanCombinedOutput(outch, p)
			if err != nil {
				fmt.Printf("Error: %v\n\n", err)
			}
		}()

		for {
			b, ok := <-outch
			if !ok {
				break
			}
			os.Stdout.Write(b)
		}

		fmt.Println()
	}

	fmt.Println("Pipe done.")
}
