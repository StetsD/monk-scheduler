package errorsApp

type ErrorApp string

func (ea ErrorApp) Error() string {
	return string(ea)
}

func Error(errText string) string {
	return ErrorApp(errText).Error()
}
