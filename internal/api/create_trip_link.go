package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
	"nlw-journey/internal/pgstore"
)

/*
PostTripsTripIDLinks
Create a trip link
*/
func (api API) PostTripsTripIDLinks(_ http.ResponseWriter, r *http.Request, _tripID string) *spec.Response {
	tripID, err := uuid.Parse(_tripID)
	if err != nil {
		return spec.PostTripsTripIDLinksJSON400Response(spec.Error{
			Message: "ID de viagem inválido.",
		})
	}

	var body spec.PostTripsTripIDLinksJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return spec.PostTripsTripIDLinksJSON400Response(spec.Error{Message: "JSON inválido: " + err.Error()})
	}

	linkID, err := api.repository.CreateTripLink(r.Context(), pgstore.CreateTripLinkParams{
		TripID: tripID,
		Title:  body.Title,
		Url:    body.URL,
	})
	if err != nil {
		api.logger.Error("failed to create trip link", zap.Error(err), zap.String("tripID", _tripID), zap.Any("body", body))
		return spec.PostTripsTripIDLinksJSON400Response(spec.Error{Message: "Alguma coisa deu errado. Tente mais tarde."})
	}

	return spec.PostTripsTripIDLinksJSON201Response(struct {
		LinkID string `json:"linkId"`
	}{
		LinkID: linkID.String(),
	})
}
