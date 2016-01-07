package main

const (
	check  = "check"
	print  = "print"
    insert = "insert"
	upsert = "upsert"
	reset  = "reset"
	delete = "delete"
)

// map command string to command constructors
var commands = map[string]func(*args) (command, error) {
	// check:  newCheck,
	// print:  newPrint,
    insert: newInsert,
	// upsert: newUpsert,
	// reset:  newReset,
	// delete: newDelete}
    }

type command interface {
    execute() error
}

type loader interface {
    LoadAll() ([]*eskip.Route, error)
}

type loaderParser interface {
    LoadAndParseAll() ([]eskip.RouteInfo, error)
}

type upserter interface {
    Upsert(r *eskip.Route) error
}

type deleter interface {
    Delete(id string) error
}

type upsertCommand struct {
    loader loader
    upserter upserter
}

type checkCommand struct {
    loaderParser loaderParser
}

type printCommand struct {
    loaderParser loaderParser
}

type loaderWrapper struct {
    loader loader
}
