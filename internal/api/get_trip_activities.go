package api

import (
	"github.com/google/uuid"
	"net/http"
	"nlw-journey/internal/api/spec"
)

func (api API) GetTripsTripIDActivities(_ http.ResponseWriter, r *http.Request, _tripID string) *spec.Response {
	tripID, err := uuid.Parse(_tripID)
	if err != nil {
		return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Id de viagem inválido."})
	}

	activities, err := api.repository.GetTripActivities(r.Context(), tripID)
	if err != nil {
		return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Algo deu errado, tente novamente mais tarde."})
	}

	parsedActivities := make([]spec.GetTripActivitiesInner, len(activities))

	for i, activity := range activities {
		parsedActivities[i] = spec.GetTripActivitiesInner{
			ID:       activity.ID.String(),
			OccursAt: activity.OccursAt.Time,
			Title:    activity.Title,
		}
	}

	return spec.GetTripsTripIDActivitiesJSON200Response(spec.GetTripActivitiesResponse{
		Activities: parsedActivities,
	})
}
