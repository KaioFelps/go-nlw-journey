package api

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
	"nlw-journey/internal/pgstore"
)

func (api API) PutTripsTripID(_ http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var body spec.PutTripsTripIDJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PutTripsTripIDJSON400Response(spec.Error{
			Message: "JSON inválido: " + err.Error(),
		})
	}

	if err := api.validator.Struct(body); err != nil {
		errors.As(err, &err)
		return spec.PostTripsJSON400Response(spec.Error{
			Message: "Input inválido: " + err.Error(),
		})
	}

	var parsedTripID, err = uuid.Parse(tripID)
	if err != nil {
		return spec.PutTripsTripIDJSON400Response(spec.Error{
			Message: "Id de viagem inválido.",
		})
	}

	if err := api.repository.UpdateTrip(r.Context(), pgstore.UpdateTripParams{
		Destination: body.Destination,
		EndsAt: pgtype.Timestamp{
			Time:  body.EndsAt,
			Valid: true,
		},
		StartsAt: pgtype.Timestamp{
			Time:  body.EndsAt,
			Valid: true,
		},
		ID: parsedTripID,
	}); err != nil {
		api.logger.Error("failed to update trip", zap.String("tripID", tripID), zap.Any("body", body))

		return spec.PutTripsTripIDJSON400Response(spec.Error{
			Message: "Algo deu errado, tente novamente mais tarde.",
		})
	}

	return spec.PutTripsTripIDJSON204Response(struct{}{})
}
