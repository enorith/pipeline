package action

type Action interface {
	Handle(params ...ActionParam) ([]ActionParam, error)
	InputTypes() []string
	OutputTypes() []string
}
