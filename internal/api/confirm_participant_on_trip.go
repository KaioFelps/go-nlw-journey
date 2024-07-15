package api

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
	"nlw-journey/internal/pgstore"
)

func (api API) PatchParticipantsParticipantIDConfirm(_ http.ResponseWriter, r *http.Request, participantID string) *spec.Response {
	id, err := uuid.Parse(participantID)

	if err != nil {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(
			spec.Error{Message: "Uuid inválido."},
		)
	}

	participant, err := api.repository.GetParticipant(r.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{
				Message: "Participante não encontrado.",
			})
		}

		api.logger.Error("Failed to get participant", zap.Error(err), zap.String("participant's ID", participantID))
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{
			Message: "Alguma coisa deu errado... Tente novamente mais tarde.",
		})
	}

	if participant.IsConfirmed {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{
			Message: "Participante já confirmado.",
		})
	}

	if err := api.repository.ConfirmParticipant(r.Context(), participant.ID); err != nil {
		api.logger.Error("Failed to confirm participant", zap.Error(err), zap.String("participant's ID", participantID))

		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{
			Message: "Alguma coisa deu errado... Tente novamente mais tarde.",
		})
	}

	return spec.PatchParticipantsParticipantIDConfirmJSON204Response(
		struct{ participant pgstore.Participant }{participant},
	)
}
