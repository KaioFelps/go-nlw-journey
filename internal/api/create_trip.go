package api

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
)

func (api API) PostTrips(_ http.ResponseWriter, r *http.Request) *spec.Response {
	var body spec.PostTripsJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PostTripsJSON400Response(spec.Error{
			Message: "JSON inválido: " + err.Error(),
		})
	}

	if err := api.validator.Struct(body); err != nil {
		errors.As(err, &err)
		return spec.PostTripsJSON400Response(spec.Error{
			Message: "Input inválido: " + err.Error(),
		})
	}

	tripId, err := api.repository.CreateTrip(r.Context(), api.pool, body)

	if err != nil {
		api.logger.Error("Failed to post trips", zap.Error(err), zap.Any("body", body))
		return spec.PostTripsJSON400Response(spec.Error{
			Message: "Alguma coisa deu errado. Tente novamente mais tarde.",
		})
	}

	go func() {
		if err := api.mailer.SendConfirmTripEmailToTripOwner(tripId); err != nil {
			api.logger.Error("failed to send email on PostTrips", zap.Error(err), zap.String("tripID", tripId.String()))
		}
	}()

	return spec.PostTripsJSON201Response(spec.CreateTripResponse{TripID: tripId.String()})
}
