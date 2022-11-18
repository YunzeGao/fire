package demo

import "github.com/YunzeGao/fire/framework"

type Service struct {
	container framework.IContainer
}

func NewService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.IContainer)
	return &Service{container: container}, nil
}

func (s *Service) GetAllStudent() []Student {
	return []Student{
		{
			ID:   1,
			Name: "foo",
		},
		{
			ID:   2,
			Name: "bar",
		},
	}
}
