package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
	"nlw-journey/internal/pgstore"
)

// PostTripsTripIDInvites Invite someone to the trip.
// (POST /trips/{tripId}/invites)
func (api API) PostTripsTripIDInvites(_ http.ResponseWriter, r *http.Request, _tripID string) *spec.Response {
	tripID, err := uuid.Parse(_tripID)
	if err != nil {
		return spec.PostTripsTripIDInvitesJSON400Response(spec.Error{Message: "Id de viagem inválido."})
	}

	var body spec.PostTripsTripIDInvitesJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PostTripsTripIDInvitesJSON400Response(spec.Error{Message: "JSON inválido: " + err.Error()})
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PostTripsTripIDInvitesJSON400Response(spec.Error{Message: "Input inválido: " + err.Error()})
	}

	_, err = api.repository.InviteParticipantsToTrip(r.Context(),
		[]pgstore.InviteParticipantsToTripParams{
			{
				TripID: tripID,
				Email:  string(body.Email),
			},
		})

	if err != nil {
		api.logger.Error("failed to add participant to trip", zap.Error(err), zap.String("email", string(body.Email)), zap.String("tripID", _tripID))
		return spec.PostTripsTripIDInvitesJSON400Response(spec.Error{Message: "Algo deu errado enquanto convidávamos o usuário. Tente novamente mais tarde."})
	}

	return spec.PostTripsTripIDInvitesJSON201Response(struct{}{})
}
