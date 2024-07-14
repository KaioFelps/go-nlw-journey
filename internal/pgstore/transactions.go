package pgstore

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"nlw-journey/internal/api/spec"
)

func (selfQueries *Queries) CreateTrip(ctx context.Context, pool *pgxpool.Pool, params spec.PostTripsJSONBody) (uuid.UUID, error) {
	tx, err := pool.Begin(ctx)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("pgstore: failed to begin trx for CreateTrip: %w", err)
	}

	/*
		Rollback the commit on the function finishing.
		If it is successfully committed, there will be nothing to be rollback.
		However, doing this will guarantee we never leave a blocked transaction in our database pool,
		since it will roll back the transaction on any return before we commit
	*/
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// starts a sql transaction
	selfWithTransaction := selfQueries.WithTx(tx)

	// insert a trip to the database
	tripID, err := selfWithTransaction.InsertTrip(ctx, InsertTripParams{
		Destination: params.Destination,
		OwnerEmail:  string(params.OwnerEmail),
		OwnerName:   params.OwnerName,
		StartsAt: pgtype.Timestamp{
			Time:  params.StartsAt,
			Valid: true,
		},
		EndsAt: pgtype.Timestamp{
			Time:  params.EndsAt,
			Valid: true,
		},
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("pgstore: failed to insert trip: %w", err)
	}

	// emails mapped to InviteParticipantsToTripParams array
	var participants = make([]InviteParticipantsToTripParams, len(params.EmailsToInvite))
	for i, email := range params.EmailsToInvite {
		participants[i] = InviteParticipantsToTripParams{
			TripID: tripID,
			Email:  string(email),
		}
	}

	// insert all the emails to the database
	if _, err := selfWithTransaction.InviteParticipantsToTrip(ctx, participants); err != nil {
		return uuid.UUID{}, fmt.Errorf("pgstore: failed to invite participants to trip: %w", err)
	}

	// try to commit (guarantee both queries have been successfully done)
	if err := tx.Commit(ctx); err != nil {
		return uuid.UUID{}, fmt.Errorf("pgstore: failed to commit CreateTrip: %w", err)
	}

	return tripID, nil
}
