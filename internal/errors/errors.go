package errors

type NotFound struct {
	Err string
}

func (e NotFound) Error() string {
	return e.Err
}

type Unauthorized struct {
	Err string
}

func (e Unauthorized) Error() string {
	return e.Err
}

type BadRequest struct {
	Err string
}

func (e BadRequest) Error() string {
	return e.Err
}
