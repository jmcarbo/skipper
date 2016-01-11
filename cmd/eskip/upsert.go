package main

type upsertCommand struct {
	loader   loader
	upserter upserter
}

func isOutputMedium(m *medium) bool {
    switch m.typ {
    case /* innkeeper, */etcd:
        return true
    default:
        return false
    }
}

func validateSelectUpsert(a *args) (*medium, *medium, error) {
    if len(a.media) == 0 || (len(a.media) == 1 && isOutputMedium(a.media[0])) {
		return nil, nil, missingInput
    }

	if len(a.media) == 1 {
		return nil, nil, missingOutput
	}

	if len(a.media) > 2 {
		return nil, nil, tooManyInputs
	}

	var input, output *medium
	if isOutputMedium(a.media[0]) {
		input = a.media[1]
		output = a.media[0]
	} else if isOutputMedium(a.media[1]) {
		input = a.media[0]
		output = a.media[1]
	} else {
		return nil, nil, missingOutput
	}

	switch input.typ {
	case stdin, file, inline:
		return input, output, nil
	default:
		return nil, nil, invalidInputType
	}
}

func newUpsert(a *args) (command, error) {
	input, output, err := validateSelectUpsert(a)
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

	return &upsertCommand{
		inputClient.(loader),
		outputClient.(upserter)}, nil
}

func (uc *upsertCommand) execute() error {
	routes, err := uc.loader.LoadAll()
	if err != nil {
		return err
	}

	for _, r := range routes {
		if err = uc.upserter.Upsert(r); err != nil {
			return err
		}
	}

	return nil
}
