package quizlet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	baseQuizletURL             = "https://quizlet.com/webapi/3.4"
	defaultUserAgent           = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36" //nolint:lll
	fetchModuleAttempts        = 10
	sleepDurationBeforeAttempt = 200 * time.Millisecond
)

type Parser struct {
	client *http.Client
}

func NewParser() (*Parser, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Parser{
		client: &http.Client{
			Jar: jar,
		},
	}, nil
}

func (p *Parser) fetchStudiableItems(
	ctx context.Context,
	moduleID string,
	attempts int,
) (*studiableItemsResponse, error) {
	query := fmt.Sprintf(
		"filters[studiableContainerId]=%s&filters[studiableContainerType]=1&perPage=1000&page=1",
		moduleID,
	)
	studiableItemsURL := fmt.Sprintf("%s/studiable-item-documents?%s", baseQuizletURL, query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, studiableItemsURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusForbidden {
		if attempts <= 0 {
			return nil, &ModuleFetchingError{ID: moduleID}
		}

		resp.Body.Close()
		time.Sleep(sleepDurationBeforeAttempt)

		return p.fetchStudiableItems(ctx, moduleID, attempts-1)
	}

	defer resp.Body.Close()

	var studiableItemsResponse studiableItemsResponse

	if err := json.NewDecoder(resp.Body).Decode(&studiableItemsResponse); err != nil {
		return nil, err
	}

	return &studiableItemsResponse, nil
}

func (p *Parser) Parse(ctx context.Context, moduleID string) ([]Card, error) {
	var cards []Card

	studiableItemsResponse, err := p.fetchStudiableItems(ctx, moduleID, fetchModuleAttempts)
	if err != nil {
		return cards, err
	}

	if len(studiableItemsResponse.Responses) == 0 {
		return cards, &ModuleParsingError{ID: moduleID}
	}

	resp := studiableItemsResponse.Responses[0]

	for _, item := range resp.Models.StudiableItem {
		card := Card{}

		for _, side := range item.CardSides {
			if len(side.Media) == 0 {
				continue
			}

			media := side.Media[0]

			if side.Label == wordLabel {
				card.Front = media.PlainText
			} else if side.Label == definitionLabel {
				card.Back = media.PlainText
			}
		}

		if card.Front != "" && card.Back != "" {
			cards = append(cards, card)
		}
	}

	return cards, nil
}
