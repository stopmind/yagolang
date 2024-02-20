package main

func main() {
	builder := NewBuilder()

	// args var
	builder.BeginVariable("args", true)

	// parseArgs
	builder.BeginCall("parseArgs")
	builder.AddValue("3")
	builder.AddValue("\"\"")

	// action
	builder.BeginCall("carg")
	builder.AddValue("\"string\"")
	builder.AddValue("\"action\"")
	builder.End()

	// user
	builder.BeginCall("carg")
	builder.AddValue("\"userid\"")
	builder.AddValue("\"user\"")
	builder.End()

	// key
	builder.BeginCall("carg")
	builder.AddValue("\"string\"")
	builder.AddValue("\"key\"")
	builder.End()

	// data
	builder.BeginCall("carg")
	builder.AddValue("\"string\"")
	builder.AddValue("\"value\"")
	builder.End()

	builder.End()

	builder.End()

	print(builder.Result)
}
