package bidsservice

import eapi "tender-manager/internal/app/employeeservice"
import "log/slog"

type BidsClient struct {
	storage bidsStorage
	eapi    eapi.Client
	logger  *slog.Logger
}

func New(b bidsStorage, e eapi.Client, l *slog.Logger) *BidsClient {
	return &BidsClient{
		storage: b,
		eapi:    e,
		logger:  l}
}
