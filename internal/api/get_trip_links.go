package api

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"nlw-journey/internal/api/spec"
)

func (api API) GetTripsTripIDLinks(_ http.ResponseWriter, r *http.Request, _tripID string) *spec.Response {
	tripID, err := uuid.Parse(_tripID)
	if err != nil {
		return spec.GetTripsTripIDLinksJSON400Response(spec.Error{Message: "Id de viagem inválido."})
	}

	links, err := api.repository.GetTripLinks(r.Context(), tripID)
	if err != nil {
		api.logger.Error("failed to fetch trip links", zap.Error(err))
		return spec.GetTripsTripIDLinksJSON400Response(spec.Error{Message: "Algo deu errado, tente mais tarde."})
	}

	parsedLinks := make([]spec.Link, len(links))
	for i, link := range links {
		parsedLinks[i] = spec.Link{
			ID:    link.ID.String(),
			Title: link.Title,
			URL:   link.Url,
		}
	}

	return spec.GetTripsTripIDLinksJSON200Response(spec.GetTripLinksResponse{Links: parsedLinks})
}
