package pipe

import (
	"fmt"
	"reflect"
	"strings"
	"log"
	"errors"
)

// errType is the type of error interface.
var errType = reflect.TypeOf((*error)(nil)).Elem()

type Pipe_[T any] struct {
	args []any
	fs []any
	err error
	inputs []reflect.Value
	executionIndex int
}

func (p *Pipe_[T]) Errored() bool {
     return p.err != nil
}

func (p *Pipe_[T]) Finished() bool {
     return p.executionIndex == len(p.fs) || p.Errored()
}

func (p *Pipe_[T]) DoN(calls_to_execute int) (p_ *Pipe_[T]) {

	if calls_to_execute == 0 {
		p.err = errors.New("Argument to DoN must be at least 1")
		return p
	}

	if p.executionIndex+calls_to_execute > len(p.fs) {
		calls_to_execute = len(p.fs) - p.executionIndex
	}

	defer func() {
		if r := recover(); r != nil {
			p.err = fmt.Errorf("pipeline panicked: %v", r)
			p_ = p
		}
	}()

	for _, arg := range p.args {
		p.inputs = append(p.inputs, reflect.ValueOf(arg))
	}

        p.args = []any{}

	for fIndex, f := range p.fs[p.executionIndex:p.executionIndex+calls_to_execute] {

		p.executionIndex++

		funcType := reflect.TypeOf(f)

		outputs := reflect.ValueOf(f).Call(p.inputs)
		p.inputs = p.inputs[:0]

		for oIndex, output := range outputs {
			if funcType.Out(oIndex).Implements(errType) {
				if !output.IsNil() {
					p.err = fmt.Errorf("%s func failed: %w", ord(fIndex), output.Interface().(error))
					return p
				}
			} else {
				p.inputs = append(p.inputs, output)
			}
		}
	}

	return p
}

func (p *Pipe_[T]) Do() (p_ *Pipe_[T]) {
	return p.DoN(len(p.fs))
}

func (p *Pipe_[T]) Flow(f any) (p_ *Pipe_[T]) {
	p.fs = append(p.fs, f)
	return p.DoN(1)
}

func (p *Pipe_[T]) Next(f any) *Pipe_[T] {
	p.fs = append(p.fs, f)
	return p
}

func Pipe[T any](args ...any) *Pipe_[T] {
	return &Pipe_[T]{
		args: args,
	}
}

func ord(index int) string {
	order := index + 1
	switch {
	case order > 10 && order < 20:
		return fmt.Sprintf("%dth", order)
	case order%10 == 1:
		return fmt.Sprintf("%dst", order)
	case order%10 == 2:
		return fmt.Sprintf("%dnd", order)
	case order%10 == 3:
		return fmt.Sprintf("%drd", order)
	default:
		return fmt.Sprintf("%dth", order)
	}
}

func (p *Pipe_[T]) Unwrap() (T, error) {
	if len(p.inputs) > 0 {
		return p.inputs[len(p.inputs)-1].Interface().(T), p.err
	}
	return *new(T), p.err
}

func main() {

    // add function calls to the pipe with deferred execution ( `.Do()` )
    res, err := Pipe[string]("Hello").
	Next(strings.ToLower).
	Next(func (s string) string {
		return s + " world!"
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
}
