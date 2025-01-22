package usecase

import (
	"backend/internal/entities"
	"backend/internal/repository"
)

type TrackUsecase interface {
	CreateTrack(track *entities.Track) error
	DeleteTrackByID(trackID int) error
	GetTrackByID(trackID int) (*entities.Track, error)
	GetTrackLinkByTrackID(trackID int) (*entities.Track, error)
	GetAllTracks() ([]entities.Track, error)
}

type trackUsecase struct {
	repo repository.TrackRepository
}

func NewTrackUsecase(repo repository.TrackRepository) TrackUsecase {
	return &trackUsecase{
		repo: repo,
	}
}

func (u *trackUsecase) CreateTrack(track *entities.Track) error {
	return u.repo.CreateTrack(track)
}

func (u *trackUsecase) DeleteTrackByID(trackID int) error {
	return u.repo.DeleteTrackByID(trackID)
}

func (u *trackUsecase) GetTrackByID(trackID int) (*entities.Track, error) {
	return u.repo.GetTrackByID(trackID)
}

func (u *trackUsecase) GetTrackLinkByTrackID(trackID int) (*entities.Track, error) {
	return u.repo.GetTrackLinkByTrackID(trackID)
}

func (u *trackUsecase) GetAllTracks() ([]entities.Track, error) {
	return u.repo.GetAllTracks()
}
