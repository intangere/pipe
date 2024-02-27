package pipe

import (
	"fmt"
	"reflect"
	"errors"
)

// errType is the type of error interface.
var errType = reflect.TypeOf((*error)(nil)).Elem()

type Pipe_[T any] struct {
	args []any
	extra_args [][]reflect.Value
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

		true_index := p.executionIndex

		if len(p.extra_args[true_index]) > 0 {
			p.inputs = append(p.inputs, p.extra_args[true_index]...)
		}

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

func (p *Pipe_[T]) Result() (T, error) {
	return p.DoN(len(p.fs)).Unwrap()
}

func (p *Pipe_[T]) Flow(f any, args ...any) (p_ *Pipe_[T]) {
	p.fs = append(p.fs, f)

	extra_args := []reflect.Value{}
	for _, arg := range args {
		extra_args = append(extra_args, reflect.ValueOf(arg))
	}

	p.extra_args = append(p.extra_args, extra_args)
	return p.DoN(1)
}

func (p *Pipe_[T]) Next(f any, args ...any) *Pipe_[T] {
	p.fs = append(p.fs, f)

	extra_args := []reflect.Value{}
	for _, arg := range args {
		extra_args = append(extra_args, reflect.ValueOf(arg))
	}

	p.extra_args = append(p.extra_args, extra_args)
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
