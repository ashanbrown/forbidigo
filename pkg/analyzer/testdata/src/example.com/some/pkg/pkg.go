package pkg

func Forbidden() {
}

func NewCustom() CustomType {
	return CustomType{}
}

type CustomType struct {
	ForbiddenField int
}


type Result struct {
	Value int
}

func (c CustomType) AlsoForbidden() *Result {
	return nil
}

type CustomInterface interface {
	StillForbidden()
}
