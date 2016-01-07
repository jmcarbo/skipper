package main

type insertCommand struct {
	loader   loader
	inserter inserter
}

func (a *args) validateSelectInsert() (*medium, *medium, error) {
	if len(a.media) == 0 || (len(a.media) == 1 && a.media[0].typ == innkeeper) {
		return nil, nil, missingInput
	}

	if len(a.media) > 2 {
		return nil, nil, tooManyInputs
	}

	if len(a.media) == 1 {
		return nil, nil, invalidInputType
	}

	var input, output *medium
	if a.media[0].typ == innkeeper {
		input = a.media[1]
		output = a.media[0]
	} else if a.media[1].typ == innkeeper {
		input = a.media[0]
		output = a.media[1]
	} else {
		return nil, nil, invalidInputType
	}

	switch input.typ {
	case stdin, file, inline:
		return input, output, nil
	default:
		return nil, nil, invalidInputType
	}
}

func newInsert(a *args) (command, error) {
	input, output, err := a.validateSelectInsert()
	if err != nil {
		return nil, err
	}

	inputClient, err := createClient(input)
	if err != nil {
		return nil, err
	}

	outputClient, err := createClient(output)
	if err != nil {
		return nil, err
	}

	return &insertCommand{
		inputClient.(loader),
		outputClient.(inserter)}, nil
}

func (ic *insertCommand) execute() error {
	routes, err := ic.loader.LoadAll()
	if err != nil {
		return err
	}

	for _, r := range routes {
		if err = ic.inserter.Insert(r); err != nil {
			return err
		}
	}

	return nil
}
