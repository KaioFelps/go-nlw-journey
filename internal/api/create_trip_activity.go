package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
	"nlw-journey/internal/pgstore"
)

func (api API) PostTripsTripIDActivities(_ http.ResponseWriter, r *http.Request, _tripID string) *spec.Response {
	var tripID, err = uuid.Parse(_tripID)
	if err != nil {
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Id de viagem inválido."})
	}

	var body spec.PostTripsTripIDActivitiesJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: "JSON inválido: " + err.Error()})
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{
			Message: "Input inválido: " + err.Error(),
		})
	}

	activityID, err := api.repository.CreateActivity(r.Context(), pgstore.CreateActivityParams{
		TripID: tripID,
		Title:  body.Title,
		OccursAt: pgtype.Timestamp{
			Time:  body.OccursAt,
			Valid: true,
		},
	})
	if err != nil {
		api.logger.Error("failed to create a trip's activity", zap.Error(err), zap.Any("body", body))

		return spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{
			Message: "Algo deu errado, tente novamente mais tarde.",
		})
	}

	return spec.PostTripsTripIDActivitiesJSON201Response(spec.CreateTripActivitiesResponse{
		ActivityID: activityID.String(),
	})
}
