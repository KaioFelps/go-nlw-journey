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
	db       Database
	logger   *zap.Logger
	port     int
	host     string
	username string
	password string
}

func NewMailPit(pool *pgxpool.Pool, logger *zap.Logger) MailPit {
	port, _ := strconv.Atoi(os.Getenv("MAILER_PORT"))

	return MailPit{
		db:       pgstore.New(pool),
		logger:   logger,
		port:     port,
		host:     os.Getenv("MAILER_HOST"),
		username: os.Getenv("MAILER_USERNAME"),
		password: os.Getenv("MAILER_PASSWORD"),
	}
}

func (mailPit MailPit) SendConfirmedTripNotificationEmail(trip pgstore.Trip) error {
	var ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	msg, err := mailPit.GenerateMsg("mailpit@jorney.com", trip.OwnerEmail, fmt.Sprintf("Confirme a sua viagem para %s.", trip.Destination))
	if err != nil {
		return err
	}

	client, err := mailPit.GenerateClient()
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles("internal/mail/mailpit/confirm_notification.tmpl")
	if err != nil {
		return fmt.Errorf("MailPit: failed to render template: %w", err)
	}

	if err := msg.SetBodyHTMLTemplate(tmpl, trip); err != nil {
		return fmt.Errorf("MailPit: failed to set 'body' html template: %w", err)
	}

	if err := client.DialAndSend(msg); err != nil {
		return fmt.Errorf("MailPit: failed to send mail: %w", err)
	}

	mailPit.logger.Info(fmt.Sprintf("MailPit: successfully sent notification e-mail to %s.", trip.OwnerEmail))

	return nil

}

func (mailPit MailPit) SendConfirmTripEmailToTripOwner(tripID uuid.UUID) error {
	var ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var trip, err = mailPit.db.GetTrip(ctx, tripID)

	if err != nil {
		return fmt.Errorf("MailPit failed to senc fonrimation trip email to id %s : %w", tripID.String(), err)
	}

	msg, err := mailPit.GenerateMsg("mailpit@jorney.com", trip.OwnerEmail, fmt.Sprintf("Confirme a sua viagem para %s.", trip.Destination))
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles("internal/mail/mailpit/confirm.tmpl")
	if err != nil {
		return fmt.Errorf("MailPit: failed to render template: %w", err)
	}

	if err := msg.SetBodyHTMLTemplate(tmpl, trip); err != nil {
		return fmt.Errorf("MailPit: failed to set 'body' html template: %w", err)
	}

	client, err := mailPit.GenerateClient()
	if err != nil {
		return err
	}

	if err := client.DialAndSend(msg); err != nil {
		return fmt.Errorf("MailPit: failed to send mail: %w", err)
	}

	mailPit.logger.Info(fmt.Sprintf("MailPit: successfully sent e-mail to %s.", trip.OwnerEmail))

	return nil
}

func (mailPit MailPit) GenerateMsg(from string, to string, subject string) (*mail.Msg, error) {
	msg := mail.NewMsg()
	if err := msg.From(from); err != nil {
		return nil, fmt.Errorf("MailPit: failed to set 'from' email: %w", err)
	}

	if err := msg.To(to); err != nil {
		return nil, fmt.Errorf("MailPit: failed to set 'to' email: %w", err)
	}

	msg.Subject(subject)

	return msg, nil
}

func (mailPit MailPit) GenerateClient() (*mail.Client, error) {
	client, err := mail.NewClient(
		mailPit.host,
		mail.WithPort(mailPit.port),
		mail.WithUsername(mailPit.username),
		mail.WithPassword(mailPit.password),
		mail.WithTLSPortPolicy(mail.NoTLS),
	)

	if err != nil {
		return nil, fmt.Errorf("MailPit: failed to create mail client: %w", err)
	}

	return client, nil
}
