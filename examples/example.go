package main

import (
	"log"
	"fmt"
	"strings"
	"errors"

	. "github.com/intangere/pipe"
)

func main() {

    // instead of doing this
    example := "Hello World!"
    example = strings.ToLower(example)
    splits := strings.Split(example, " ")
    new := "Bye " + splits[1]
    log.Println("Result:", new)

    // with pipe you can do
    res, err := Pipe[string]("Hello World!").
	Flow(strings.ToLower).
        Flow(strings.Split, " ").
        Flow(func (parts []string) string {
		return "Bye " + parts[1]
	}).
	Unwrap()

    log.Println("Result:", res, "Error:", err)

    // or with deferred execution
    res, err = Pipe[string]("Hello World!").
	Next(strings.ToLower).
        Next(strings.Split, " ").
        Next(func (parts []string) string {
		return "Bye " + parts[1]
	}).
	Do().
	Unwrap()

    log.Println("Result:", res, "Error:", err)

    // .Result() is shorthand for .Do().Unwrap()
    // It can also be used with .Flow() as .Do() will be a no-op

    res, err = Pipe[string]("Hello World!").
	Next(strings.ToLower).
        Next(strings.Split, " ").
        Next(func (parts []string) string {
		return "Bye " + parts[1]
	}).
	Result()

    log.Println("Result:", res, "Error:", err)

    // add function calls to the pipe with deferred execution ( `.Do()` )
    res, err = Pipe[string]("Hello").
	Next(strings.ToLower).
	Next(func (s string) string {
		return s + " world1!"
	}).
	Next(strings.Title).
	Do().
	Unwrap()

    log.Println("Result:", res, "Error:", err)


    // execute function calls as they are added to the pipe
    res, err = Pipe[string]("Hello").
	Flow(strings.ToLower).
	Flow(func (s string) string {
		return s + " world!"
	}).
	Flow(strings.Title).
	Unwrap()

    log.Println("Result:", res, "Error:", err)


    // execute calls one at a time and inspect the current state
    p := Pipe[string]("Hello").
	Next(strings.ToLower).
	Next(func (s string) string {
		return s + " world!"
	}).
	Next(strings.Title)

    for !p.Finished() {
	p.DoN(1)

        res, err = p.Unwrap()
        log.Println("Current state |", "Result:", res, "Err:", err, "Errored?", p.Errored())
    }

    res, err = p.Unwrap()
    log.Println("Final state: ", res, err)


    // automatic panic handling
    _, err = Pipe[string]("Hello").
	Flow(func (s string) string {
		panic("i panicked")
		return s + " world!"
	}).
	Unwrap()

    if err != nil {
       log.Println("Pipe panicked", err)
    }


    // mixed function types - i'm surprised this works.
    // you can use any amount of outputs as inputs to the next function if it expects it
    res, err = Pipe[string]("Hello").
	Flow(func (s string) (string, string) {
		return s, "world!"
	}).
	Flow(func (s string, s2 string) (string) {
		return s + " " + s2
	}).
	Unwrap()

    log.Println("Merged result:", res, err)


    // mixed functions that may fail
    _, err = Pipe[string]("Hello").
	Flow(func (s string) (string, error) {
		return "testing", nil
	}).
	Flow(func (s string) (string) {
		return s
	}).
	Flow(func (s string) (string, error) {
		return "", errors.New("I errored")
	}).
	Unwrap()

    if err != nil {
       log.Println("This pipe always fails", err)
    }

    // expand an existing pipeline
    p = Pipe[string]("Hello").
	Flow(func (s string) (string, error) {
		return "testing", nil
	})

    res, err = p.Unwrap()
    log.Println("Current res:", res)

    p.Flow(func (s string) (string) {
       return s + " #2"
    })

    res, err = p.Unwrap()
    log.Println("New res:", res)


    //reusable pipe
    p = Pipe[string]().
	Next(func (s string) string {
           return "With argument: " + s
        }).
        Next(log.Println)

    for i := 0; i < 3; i++ {
        // no result
        _, err = p.Result(fmt.Sprintf("%d", i))
    }

    // will print:
    // With argument: 0
    // With argument: 1
    // With argument: 2


    //reusable pipe cannot be used with Flow(). It will throw an input error
    p = Pipe[string]().
	Flow(func (s string) (string, error) {
		return "testing", nil
	})

    _, err = p.Unwrap()
    log.Print("Errored: ", p.Errored(), err)
}
