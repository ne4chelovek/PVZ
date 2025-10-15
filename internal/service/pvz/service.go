package pvz

import (
	"PVZ/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PVZService struct {
	pvzRepo repository.PVZRepository
	recRepo repository.ReceptionRepository
	outBox  repository.EventRepository
	dbPool  *pgxpool.Pool
}

func NewPVZService(repo repository.PVZRepository, recRepo repository.ReceptionRepository, dbPool *pgxpool.Pool, outBox repository.EventRepository) *PVZService {
	return &PVZService{pvzRepo: repo, recRepo: recRepo, dbPool: dbPool, outBox: outBox}
}
