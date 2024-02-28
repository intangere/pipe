package pipe

import (
	"reflect"
	"fmt"
)

// errType is the type of error interface.
var errType = reflect.TypeOf((*error)(nil)).Elem()

type PipeFast_[R any, F func(R) R] struct {
	//args []any
	extra_args [][]R
	fs []any
	err error
	inputs []R
}

func (p *PipeFast_[R, F]) Errored() bool {
     return p.err != nil
}

func (p *PipeFast_[R, F]) Flow(f F, args ...R) (p_ *PipeFast_[R, F]) {
	defer func() {
		if r := recover(); r != nil {
			p.err = fmt.Errorf("panicked: %w", r)
			p_ = p
		}
	}()

	p.inputs = append(p.inputs, args...)
	output := f(p.inputs[len(p.inputs)-1])
	p.inputs = append(p.inputs, output)
	return p
}

func Pipe[R any, F func (R) R](args ...R) *PipeFast_[R, F] {
	return &PipeFast_[R, F]{
		inputs: args,
	}
}

func (p *PipeFast_[R, F]) Result() (R, error) {
	if len(p.inputs) > 0 {
		return p.inputs[len(p.inputs)-1], p.err
	}
	return *new(R), p.err
}
