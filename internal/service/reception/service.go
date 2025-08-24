package reception

import "PVZ/internal/repository"

type ReceptionService struct {
	Repo repository.ReceptionRepository
}

func NewReceptionService(repo repository.ReceptionRepository) *ReceptionService {
	return &ReceptionService{Repo: repo}
}
