package main

import (
	"errors"
	"fmt"
	"strings"
)

type ContextType uint

const (
	CallContext ContextType = iota
	ValueContext
	VariableContext
)

type BuilderContext struct {
	FinalString string
	name        string
	value       string
	contextType ContextType
	definition  bool
	SubContexts []*BuilderContext
}

func (c *BuilderContext) Collapse() error {
	switch c.contextType {
	case ValueContext:
		c.FinalString = c.value
	case CallContext:
		args := []string{c.name}

		for _, subContext := range c.SubContexts {
			if subContext.contextType == ValueContext {
				args = append(args, subContext.FinalString)
			} else {
				args = append(args, fmt.Sprintf("(%s)", subContext.FinalString))
			}
		}

		c.FinalString = strings.Join(args, " ")
	case VariableContext:
		if len(c.SubContexts) != 1 {
			return errors.New("invalid subcontexts count")
		}
		c.FinalString = fmt.Sprint("$", c.name)
		if c.definition {
			c.FinalString += " := "
		} else {
			c.FinalString += " = "
		}

		c.FinalString += c.SubContexts[0].FinalString

	default:
		return errors.New("unknown context type")
	}

	c.SubContexts = nil
	return nil
}

type Builder struct {
	Contexts []*BuilderContext
	Result   string
}

func (b *Builder) BaseBegin(context *BuilderContext) {
	b.Contexts = append(b.Contexts, context)
}

func (b *Builder) BeginCall(call string) {
	b.BaseBegin(&BuilderContext{
		name:        call,
		contextType: CallContext,
	})
}

func (b *Builder) End() error {
	if len(b.Contexts) == 0 {
		return nil
	}

	context := b.Contexts[len(b.Contexts)-1]
	b.Contexts = b.Contexts[0 : len(b.Contexts)-1]

	err := context.Collapse()
	if err != nil {
		return err
	}

	if len(b.Contexts) == 0 {
		b.Result += fmt.Sprintf("{{%s}}", context.FinalString)
	} else {
		parent := b.Contexts[len(b.Contexts)-1]
		parent.SubContexts = append(parent.SubContexts, context)
	}

	return nil
}

func (b *Builder) AddValue(value string) {
	b.BaseBegin(&BuilderContext{
		value:       value,
		contextType: ValueContext,
	})
	_ = b.End()
}

func (b *Builder) BeginVariable(variable string, definition bool) error {
	if len(b.Contexts) > 0 {
		return errors.New("invalid construction")
	}

	b.BaseBegin(&BuilderContext{
		name:        variable,
		definition:  definition,
		contextType: VariableContext,
	})

	return nil
}

func NewBuilder() *Builder {
	return &Builder{}
}
