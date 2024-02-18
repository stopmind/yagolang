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
	Name        string
	Value       string
	Type        ContextType
	Definition  bool
	SubContexts []*BuilderContext
}

func (c *BuilderContext) collapse() {
	switch c.Type {
	case ValueContext:
		c.FinalString = c.Value
	case CallContext:
		args := []string{c.Name}

		for _, subContext := range c.SubContexts {
			if subContext.Type == ValueContext {
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
		c.FinalString = fmt.Sprint("$", c.Name)
		if c.Definition {
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

func (b *Builder) baseBegin(context *BuilderContext) {
	b.Contexts = append(b.Contexts, context)
}

func (b *Builder) beginCall(call string) {
	b.baseBegin(&BuilderContext{
		Name: call,
		Type: CallContext,
	})
}

func (b *Builder) end() {
	if len(b.Contexts) == 0 {
		return
	}

	context := b.Contexts[len(b.Contexts)-1]
	b.Contexts = b.Contexts[0 : len(b.Contexts)-1]

	context.collapse()
	if len(b.Contexts) == 0 {
		b.Result += fmt.Sprintf("{{%s}}", context.FinalString)
	} else {
		parent := b.Contexts[len(b.Contexts)-1]
		parent.SubContexts = append(parent.SubContexts, context)
	}
}

func (b *Builder) addValue(value string) {
	b.baseBegin(&BuilderContext{
		Value: value,
		Type:  ValueContext,
	})
	b.end()
}

func (b *Builder) getVariable(variable string) {
	b.baseBegin(&BuilderContext{
		Value: fmt.Sprint("$", variable),
		Type:  ValueContext,
	})
	b.end()
}

func (b *Builder) beginVariable(variable string, definition bool) {
	if len(b.Contexts) > 0 {
		panic("Invalid construction")
	}

	b.baseBegin(&BuilderContext{
		Name:       variable,
		Definition: definition,
		Type:       VariableContext,
	})
}

func newBuilder() *Builder {
	return &Builder{}
}
