package api

import (
	"github.com/google/uuid"
	"net/http"
	"nlw-journey/internal/api/spec"
)

func (api API) GetTripsTripID(_ http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	parsedTripID, err := uuid.Parse(tripID)
	if err != nil {
		return spec.GetTripsTripIDJSON400Response(spec.Error{
			Message: "ID de viagem inválido.",
		})
	}

	trip, err := api.repository.GetTrip(r.Context(), parsedTripID)
	if err != nil {
		return spec.GetTripsTripIDJSON400Response(spec.Error{
			Message: "Viagem não encontrada.",
		})
	}
	return spec.GetTripsTripIDJSON200Response(struct {
		Trip spec.Trip `json:"trip"`
	}{
		Trip: spec.Trip{
			Destination: trip.Destination,
			EndsAt:      trip.EndsAt.Time,
			ID:          trip.ID.String(),
			IsConfirmed: trip.IsConfirmed,
			StartsAt:    trip.StartsAt.Time,
		},
	})
}
