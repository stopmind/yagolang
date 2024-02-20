package main

import (
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

func (c *BuilderContext) Collapse() {
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
			panic("Invalid subcontexts count")
		}
		c.FinalString = fmt.Sprint("$", c.name)
		if c.definition {
			c.FinalString += " := "
		} else {
			c.FinalString += " = "
		}

		c.FinalString += c.SubContexts[0].FinalString

	default:
		panic("Unknown context type!")
	}

	c.SubContexts = nil
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

func (b *Builder) End() {
	if len(b.Contexts) == 0 {
		return
	}

	context := b.Contexts[len(b.Contexts)-1]
	b.Contexts = b.Contexts[0 : len(b.Contexts)-1]

	context.Collapse()
	if len(b.Contexts) == 0 {
		b.Result += fmt.Sprintf("{{%s}}", context.FinalString)
	} else {
		parent := b.Contexts[len(b.Contexts)-1]
		parent.SubContexts = append(parent.SubContexts, context)
	}
}

func (b *Builder) AddValue(value string) {
	b.BaseBegin(&BuilderContext{
		value:       value,
		contextType: ValueContext,
	})
	b.End()
}

func (b *Builder) BeginVariable(variable string, definition bool) {
	if len(b.Contexts) > 0 {
		panic("Invalid construction")
	}

	b.BaseBegin(&BuilderContext{
		name:        variable,
		definition:  definition,
		contextType: VariableContext,
	})
}

func NewBuilder() *Builder {
	return &Builder{}
}
