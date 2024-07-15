package api

import (
	"github.com/discord-gophers/goapi-gen/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
)

func (api API) GetTripsTripIDParticipants(_ http.ResponseWriter, r *http.Request, _tripID string) *spec.Response {
	tripID, err := uuid.Parse(_tripID)
	if err != nil {
		return spec.GetTripsTripIDParticipantsJSON400Response(spec.Error{Message: "ID de viagem inválido."})
	}

	participants, err := api.repository.GetParticipants(r.Context(), tripID)
	if err != nil {
		api.logger.Error("failed to get trip's participants", zap.Error(err), zap.String("tripID", _tripID))
		return spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: "Algo deu errdao, tente novamente."})
	}

	mappedParticipants := make([]spec.Participant, len(participants))
	for i, participant := range participants {
		mappedParticipants[i] = spec.Participant{
			Email:       types.Email(participant.Email),
			ID:          participant.ID.String(),
			IsConfirmed: participant.IsConfirmed,
			Name:        nil,
		}
	}

	return spec.GetTripsTripIDParticipantsJSON200Response(struct {
		Participants []spec.Participant `json:"participants"`
	}{
		Participants: mappedParticipants,
	})
}
