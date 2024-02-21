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
