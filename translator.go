package main

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
