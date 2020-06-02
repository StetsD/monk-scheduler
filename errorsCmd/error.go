package errorsCmd

type ErrorCmd string

func (e ErrorCmd) Error() string {
	return string(e)
}
