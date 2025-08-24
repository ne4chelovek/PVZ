package pvz

import "PVZ/internal/repository"

type PVZService struct {
	Repo    repository.PVZRepository
	RecRepo repository.ReceptionRepository
}

func NewPVZService(repo repository.PVZRepository, recRepo repository.ReceptionRepository) *PVZService {
	return &PVZService{Repo: repo, RecRepo: recRepo}
}
