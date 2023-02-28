package loms

import (
	desc "route256/loms/pkg/loms"
)

type Implementation struct {
	desc.UnimplementedLomsServer
}

func NewLoms() *Implementation {
	return &Implementation{
		desc.UnimplementedLomsServer{},
	}
}
