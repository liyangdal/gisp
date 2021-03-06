package gisp

import (
	"fmt"
)

type TypeSignError struct {
	Type  Type
	Value interface{}
}

func (err TypeSignError) Error() string {
	return fmt.Sprintf("%v can't match %v", err.Value, err.Type)
}

type ParsexSignError struct {
	message string
	args    []interface{}
}

func ParsexSignErrorf(message string, args ...interface{}) ParsexSignError {
	return ParsexSignError{message, args}
}

func (err ParsexSignError) Error() string {
	return fmt.Sprintf(err.message, err.args...)
}

type Func interface {
	Functor
	Name() string
	Overload(functor Functor) error
	Content() []Functor
}

type TaskBox struct {
	task func(env Env) (interface{}, error)
}

func (tb TaskBox) Eval(env Env) (interface{}, error) {
	return tb.task(env)
}

type Function struct {
	atom    Atom
	Global  Env
	content []Functor
}

func NewFunction(name string, global Env, functor Functor) *Function {
	return &Function{
		atom:    Atom{name, Type{ANY, false}},
		Global:  global,
		content: []Functor{functor},
	}
}

func (fun Function) Name() string {
	return fun.atom.Name
}

func (fun Function) Task(env Env, args ...interface{}) (Lisp, error) {
	for _, functor := range fun.content {
		task, err := functor.Task(env, args...)
		if err == nil {
			return task, nil
		}
	}

	if f, ok := fun.Global.Global(fun.Name()); ok {
		switch foo := f.(type) {
		case Functor:
			return foo.Task(env, args...)
		case TaskExpr:
			task, err := foo(env, args...)
			if err != nil {
				return nil, err
			}
			return TaskBox{task}, nil
		case LispExpr:
			lisp, err := foo(env, args...)
			if err != nil {
				return nil, err
			}
			return lisp, nil
		}
	}
	return nil, fmt.Errorf("not found args type sign for %v", args)
}

func (fun *Function) Overload(functor Functor) error {
	fun.content = append([]Functor{functor}, fun.content...)
	return nil
}

func (fun Function) Content() []Functor {
	return fun.content
}

func DefunExpr(env Env, args ...interface{}) (Tasker, error) {
	funName := args[0].(Atom)
	_args := args[1].(List)
	lambda, err := DeclareLambda(env, _args, args[2:]...)
	if err != nil {
		return nil, err
	}
	if f, ok := env.Local(funName.Name); ok {
		if fun, ok := f.(Function); ok {
			err := fun.Overload(*lambda)
			if err == nil {
				return Q(fun).Eval, nil
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("%v is defined as no Expr", funName.Name)
		}
	} else {
		err := env.Defun(funName.Name, *lambda)
		if err == nil {
			return nil, nil
		} else {
			return nil, err
		}
	}
}
