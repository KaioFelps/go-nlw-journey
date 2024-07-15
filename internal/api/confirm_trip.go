package api

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
	"nlw-journey/internal/pgstore"
)

// GetTripsTripIDConfirm Confirm a trip and send e-mail invitations.
func (api API) GetTripsTripIDConfirm(_ http.ResponseWriter, r *http.Request, _tripID string) *spec.Response {
	tripID, err := uuid.Parse(_tripID)
	if err != nil {
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{
			Message: "Id de viagem inválido.",
		})
	}

	trip, err := api.repository.GetTrip(r.Context(), tripID)
	if err != nil {
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: fmt.Sprintf("Não foi possível encontrar a viagem de id %s", _tripID)})
	}

	if err := api.repository.UpdateTrip(r.Context(), pgstore.UpdateTripParams{
		Destination: trip.Destination,
		EndsAt:      trip.EndsAt,
		StartsAt:    trip.StartsAt,
		IsConfirmed: true,
		ID:          trip.ID,
	}); err != nil {
		api.logger.Error("failed to update trip", zap.Error(err), zap.String("tripID", _tripID), zap.Any("trip", trip))
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: "Algo deu errado agora, tente mais tarde."})
	}

	go func() {
		if err := api.mailer.SendConfirmedTripNotificationEmail(trip); err != nil {
			api.logger.Error("failed to send email on ConfirmTrip", zap.Error(err), zap.Any("trip", trip))
		}
	}()

	return spec.GetTripsTripIDConfirmJSON204Response(struct{}{})
}
