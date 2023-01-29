package pkg

func Forbidden() {
}

func NewCustom() CustomType {
	return CustomType{}
}

type CustomType struct {
	ForbiddenField int
}

func (c CustomType) AlsoForbidden() {}

type CustomInterface interface {
	StillForbidden()
}
