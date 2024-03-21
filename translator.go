package main

import (
	"errors"
	"fmt"
	"os"
)

type Translator struct {
	projectPath string
	files       []string
	entities    []IEntity
}

func (t *Translator) AddFile(filePath string) {
	for _, file := range t.files {
		if file == filePath {
			return
		}
	}

	t.files = append(t.files, filePath)
}

func (t *Translator) ProcessFile(filePath string) error {
	data, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	tokens := Tokenize(string(data))

	fmt.Printf("File \"%v\" was reading.\n", filePath)
	for i, token := range tokens {
		fmt.Printf("Token (index: %v)\n"+
			" Type: %v (index: %v)\n"+
			" Line: %v\n"+
			" Row: %v\n"+
			" Data: %v\n\n",
			i,
			TokenTypeName(token.Type), token.Type,
			token.Line,
			token.Row,
			token.Data)
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case NameToken:
			i++
			if i >= len(tokens) {
				return errors.New("Unexpected file end")
			}
			nextToken := tokens[i]

			switch nextToken.Type {
			case OpenBracketToken:
				i++
				if i >= len(tokens) || tokens[i].Type != CloseBracketToken {
					return errors.New("this call not closed or use args")
				}

				t.entities = append(t.entities, CallEntity{funcName: token.Data})
			default:
				return errors.New("this construction not supported or incorrect")
			}
		}

	}

	return nil
}

func (t *Translator) ProcessFiles() error {
	for _, file := range t.files {
		err := t.ProcessFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Translator) Build(builder *Builder) error {
	ctx := EntitiesContext{}

	for _, entity := range t.entities {
		err := entity.Registration(&ctx)
		if err != nil {
			return err
		}
	}

	for _, entity := range t.entities {
		err := entity.Check(&ctx)
		if err != nil {
			return err
		}
	}

	for _, entity := range ctx.RegisteredEntities {
		if entity.UsesCount > 0 {
			err := entity.Entity.Build(&ctx, builder)
			if err != nil {
				return err
			}
		}
	}

	for _, entity := range t.entities {
		if entity.Code() {
			err := entity.Build(&ctx, builder)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Translator) TranslateFile(filePath string) (string, error) {
	t.AddFile(filePath)
	err := t.ProcessFiles()
	if err != nil {
		return "", err
	}

	builder := NewBuilder()
	err = t.Build(builder)
	if err != nil {
		return "", err
	}

	return builder.Result, nil
}

func NewTranslator(projectPath string) Translator {
	return Translator{
		projectPath: projectPath,
	}
}
