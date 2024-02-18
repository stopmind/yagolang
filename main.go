package main

func main() {
	builder := newBuilder()

	// args var
	builder.beginVariable("args", true)

	// parseArgs
	builder.beginCall("parseArgs")
	builder.addValue("3")
	builder.addValue("\"\"")

	// action
	builder.beginCall("carg")
	builder.addValue("\"string\"")
	builder.addValue("\"action\"")
	builder.end()

	// user
	builder.beginCall("carg")
	builder.addValue("\"userid\"")
	builder.addValue("\"user\"")
	builder.end()

	// key
	builder.beginCall("carg")
	builder.addValue("\"string\"")
	builder.addValue("\"key\"")
	builder.end()

	// data
	builder.beginCall("carg")
	builder.addValue("\"string\"")
	builder.addValue("\"value\"")
	builder.end()

	builder.end()

	builder.end()

	print(builder.Result)
}
