package pkg

func Forbidden() {
}

type CustomType struct {
	ForbiddenField int
}

func (c CustomType) AlsoForbidden() {}

type CustomInterface interface {
	StillForbidden()
}
