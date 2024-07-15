package api

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
	"nlw-journey/internal/pgstore"
)

type Repository interface {
	GetParticipant(context.Context, uuid.UUID) (pgstore.Participant, error)
	ConfirmParticipant(context.Context, uuid.UUID) error
	CreateTrip(context.Context, *pgxpool.Pool, spec.PostTripsJSONBody) (uuid.UUID, error)
	GetTrip(context.Context, uuid.UUID) (pgstore.Trip, error)
	UpdateTrip(context.Context, pgstore.UpdateTripParams) error
	GetTripActivities(context.Context, uuid.UUID) ([]pgstore.Activity, error)
	CreateTripLink(context.Context, pgstore.CreateTripLinkParams) (uuid.UUID, error)
	GetParticipants(context.Context, uuid.UUID) ([]pgstore.Participant, error)
	GetTripLinks(context.Context, uuid.UUID) ([]pgstore.Link, error)
	InsertTrip(context.Context, pgstore.InsertTripParams) (uuid.UUID, error)
	CreateActivity(context.Context, pgstore.CreateActivityParams) (uuid.UUID, error)
	InviteParticipantsToTrip(ctx context.Context, arg []pgstore.InviteParticipantsToTripParams) (int64, error)
}

type Mailer interface {
	SendConfirmTripEmailToTripOwner(tripID uuid.UUID) error
	SendConfirmedTripNotificationEmail(trip pgstore.Trip) error
}

type API struct {
	repository Repository
	pool       *pgxpool.Pool
	logger     *zap.Logger
	validator  *validator.Validate
	mailer     Mailer
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

func (api API) PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}

func (api API) GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	//TODO implement me
	panic("implement me")
}
