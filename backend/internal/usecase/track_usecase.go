package usecase

import (
	"backend/internal/entities"
	"backend/internal/repository"
	"math/rand"
	"time"
)

type TrackUsecase interface {
	CreateTrack(track *entities.Track) error
	DeleteTrackByID(trackID int) error
	GetTrackByID(trackID int) (*entities.Track, error)
	GetTrackLinkByTrackID(trackID int) (*entities.Track, error)
	GetAllTracks() ([]entities.Track, error)
	GetRandomTrack() (*entities.Track, error)
	GetTrackByAuthorID(authorID int) ([]entities.Track, error)
	GetRandomTrackByGenre(genreID int) (*entities.Track, error)
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

func (u *trackUsecase) GetRandomTrack() (*entities.Track, error) {
	trackCount, err := u.repo.GetTrackCount()
	if err != nil {
		return nil, err
	}
	if trackCount == 0 {
		return nil, nil
	}
	rand.Seed(time.Now().UnixNano())
	randomOffset := rand.Intn(trackCount)

	track, err := u.repo.GetRandomTrack(randomOffset)
	if err != nil {
		return nil, err
	}
	return track, nil
}

func (u *trackUsecase) GetTrackByAuthorID(authorID int) ([]entities.Track, error) {
	return u.repo.GetTrackByAuthorID(authorID)
}

func (u *trackUsecase) GetRandomTrackByGenre(genreID int) (*entities.Track, error) {
	return u.repo.GetRandomTrackByGenre(genreID)
}
