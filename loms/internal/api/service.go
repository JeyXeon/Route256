package loms

import (
	"route256/loms/internal/service"
	desc "route256/loms/pkg/loms"
)

type Implementation struct {
	desc.UnimplementedLomsServer

	lomsService *service.Service
}

func NewLoms(lomsService *service.Service) *Implementation {
	return &Implementation{
		desc.UnimplementedLomsServer{},
		lomsService,
	}
}
