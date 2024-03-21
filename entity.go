package main

type RegisteredEntity struct {
	UsesCount int
	Entity    IEntity
}

type EntitiesContext struct {
	RegisteredEntities []RegisteredEntity
}

type IEntity interface {
	Code() bool
	Registration(ctx *EntitiesContext) error
	Check(ctx *EntitiesContext) error
	Build(ctx *EntitiesContext, builder *Builder) error
}

type CallEntity struct {
	funcName string
}

func (c CallEntity) Code() bool {
	return true
}

func (c CallEntity) Registration(ctx *EntitiesContext) error {
	return nil
}

func (c CallEntity) Check(ctx *EntitiesContext) error {
	return nil
}

func (c CallEntity) Build(ctx *EntitiesContext, builder *Builder) error {

	builder.BeginCall(c.funcName)
	if err := builder.End(); err != nil {
		return err
	}

	return nil
}
