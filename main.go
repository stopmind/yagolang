package main

func main() {
	translator := NewTranslator("./")
	result, err := translator.TranslateFile("test.yg")
	if err != nil {
		panic(err)
	}

	print(result)
}
