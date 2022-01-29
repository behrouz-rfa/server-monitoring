package adminservice

import "server-monitoring/domain/visits"

var (
	VisitedService visitedServiceInterface = &visitedService{}
)

type visitedServiceInterface interface {
	Insert(visit *visits.Visit) error
}

type visitedService struct {
}

func (v visitedService) Insert(visit *visits.Visit) error {

	return visit.Insert()
}
