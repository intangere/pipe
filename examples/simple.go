package main

import . "github.com/intangere/pipe/simple"
import "strings"
import "log"

// this is pretty much 1-1 with a |> operator input->output->(output as input)->output
// no magic here what so ever so it is fast, and it is constrained to 1 type Pipe[your input/output type]

func main() {
	res, err := Pipe[string]("Hello World!\n  ").
		Flow(strings.ToLower).
		Flow(strings.Title).
		Flow(strings.TrimSpace).
		Result()

	log.Println("Res: ", res, "Error:", err)

	// you can omit type
	res, err = Pipe("Hello World!\n  ").
		Flow(strings.ToLower).
		Flow(strings.Title).
		Flow(strings.TrimSpace).
		Result()

	log.Println("Res: ", res, "Error:", err)
}
