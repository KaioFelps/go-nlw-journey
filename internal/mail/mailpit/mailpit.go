package mailpit

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wneessen/go-mail"
	"go.uber.org/zap"
	"html/template"
	"nlw-journey/internal/pgstore"
	"os"
	"strconv"
	"time"
)

type Database interface {
	GetTrip(ctx context.Context, tripID uuid.UUID) (pgstore.Trip, error)
}

type MailPit struct {
	db     Database
	logger *zap.Logger
}

func NewMailPit(pool *pgxpool.Pool, logger *zap.Logger) MailPit {
	return MailPit{
		db:     pgstore.New(pool),
		logger: logger,
	}
}

func (mailPit MailPit) SendConfirmTripEmailToTripOwner(tripID uuid.UUID) error {
	var ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var trip, err = mailPit.db.GetTrip(ctx, tripID)

	if err != nil {
		return fmt.Errorf("MailPit failed to senc fonrimation trip email to id %s : %w", tripID.String(), err)
	}

	msg := mail.NewMsg()
	if err := msg.From("mailpit@jorney.com"); err != nil {
		return fmt.Errorf("MailPit: failed to set 'from' email: %w", err)
	}

	if err := msg.To(trip.OwnerEmail); err != nil {
		return fmt.Errorf("MailPit: failed to set 'to' email: %w", err)
	}

	msg.Subject(fmt.Sprintf("Confirme a sua viagem para %s.", trip.Destination))

	tmpl, err := template.ParseFiles("internal/mail/mailpit/confirmar.tmpl")
	if err != nil {
		return fmt.Errorf("MailPit: failed to render template: %w", err)
	}

	if err := msg.SetBodyHTMLTemplate(tmpl, trip); err != nil {
		return fmt.Errorf("MailPit: failed to set 'body' html template: %w", err)
	}

	port, _ := strconv.Atoi(os.Getenv("MAILER_PORT"))
	host := os.Getenv("MAILER_HOST")
	username := os.Getenv("MAILER_USERNAME")
	password := os.Getenv("MAILER_PASSWORD")

	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTLSPortPolicy(mail.NoTLS),
	)

	if err != nil {
		return fmt.Errorf("MailPit: failed to create mail client: %w", err)
	}

	if err := client.DialAndSend(msg); err != nil {
		return fmt.Errorf("MailPit: failed to send mail: %w", err)
	}

	mailPit.logger.Info(fmt.Sprintf("MailPit: successfully sent e-mail to %s.", trip.OwnerEmail))

	return nil
}
