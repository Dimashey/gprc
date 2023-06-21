//go:generate mockgen -destination=rocket_mocks_test.go -package=rocket github.com/Dimashey/gprc/internal/rocket Store

package rocket

import "context"

type Rocket struct {
	ID      string
	Name    string
	Type    string
	Flights int
}

type Store interface {
	GetRocketById(id string) (Rocket, error)
	InsertRocket(rocket Rocket) (Rocket, error)
	DeleteRocket(id string) error
}

type Service struct {
	Store Store
}

func New(store Store) Service {
	return Service{Store: store}
}

func (s Service) GetRocketByID(ctx context.Context, id string) (Rocket, error) {
	rocket, err := s.Store.GetRocketById(id)
	if err != nil {
		return Rocket{}, err
	}

	return rocket, nil
}

func (s Service) InsertRocket(ctx context.Context, rocket Rocket) (Rocket, error) {
	rkt, err := s.Store.InsertRocket(rocket)
	if err != nil {
		return Rocket{}, err
	}

	return rkt, nil
}

func (s Service) DeleteRocket(ctx context.Context, id string) error {
	err := s.Store.DeleteRocket(id)
	if err != nil {
		return err
	}

	return nil
}
