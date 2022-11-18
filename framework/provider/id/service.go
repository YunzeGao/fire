package id

import "github.com/rs/xid"

type FireIDService struct {
}

func NewFireIDService(params ...interface{}) (interface{}, error) {
	return &FireIDService{}, nil
}

func (s *FireIDService) NewID() string {
	return xid.New().String()
}
