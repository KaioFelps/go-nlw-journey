package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
	"nlw-journey/internal/pgstore"
)

type Database interface {
	GetParticipant(ctx context.Context, participantID uuid.UUID) (pgstore.Participant, error)
	ConfirmParticipant(ctx context.Context, participantID uuid.UUID) error
	CreateTrip(context.Context, *pgxpool.Pool, spec.PostTripsJSONBody) (uuid.UUID, error)
}

type Mailer interface {
	SendConfirmTripEmailToTripOwner(tripID uuid.UUID) error
}

type API struct {
	db        Database
	pool      *pgxpool.Pool
	logger    *zap.Logger
	validator *validator.Validate
	mailer    Mailer
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger, mailer Mailer) API {
	_validator := validator.New(validator.WithRequiredStructEnabled())

	return API{
		pgstore.New(pool),
		pool,
		logger,
		_validator,
		mailer,
	}
}

func (api API) PatchParticipantsParticipantIDConfirm(w http.ResponseWriter, r *http.Request, participantID string) *spec.Response {
	id, err := uuid.Parse(participantID)

	if err != nil {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(
			spec.Error{Message: "Uuid inválido."},
		)
	}

	participant, err := api.db.GetParticipant(r.Context(), id)

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

	if err := api.db.ConfirmParticipant(r.Context(), participant.ID); err != nil {
		api.logger.Error("Failed to confirm participant", zap.Error(err), zap.String("participant's ID", participantID))

		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{
			Message: "Alguma coisa deu errado... Tente novamente mais tarde.",
		})
	}

	return spec.PatchParticipantsParticipantIDConfirmJSON204Response(
		struct{ participant pgstore.Participant }{participant},
	)
}

func (api API) PostTrips(w http.ResponseWriter, r *http.Request) *spec.Response {
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

	tripId, err := api.db.CreateTrip(r.Context(), api.pool, body)

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

func (api API) GetTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}

func (api API) PutTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}

func (api API) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}

func (api API) PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}

func (api API) GetTripsTripIDConfirm(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}

func (api API) PostTripsTripIDInvites(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}

func (api API) GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}

func (api API) PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}

func (api API) GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}
