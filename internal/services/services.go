package services

import "github.com/namkatcedrickjumtock/travel-planner/persistence"


type Planner interface{

}

type TravelPlannerServiceImpl struct {
	repo persistence.Repository
}

var _ Planner = (*TravelPlannerServiceImpl)(nil)

func NewTravelPlannerService(repo persistence.Repository) (*TravelPlannerServiceImpl, error) {
	return &TravelPlannerServiceImpl{repo: repo}, nil
}

func (t *TravelPlannerServiceImpl) BookHotelAppointment()  {
	
}