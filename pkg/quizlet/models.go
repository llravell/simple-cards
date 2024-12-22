package quizlet

import "fmt"

type cardSideLabel string

const (
	wordLabel       cardSideLabel = "word"
	definitionLabel cardSideLabel = "definition"
)

type Card struct {
	Front string
	Back  string
}

type studiableItemsResponse struct {
	Responses []struct {
		Models struct {
			StudiableItem []studiableItem `json:"studiableItem"`
		} `json:"models"`
	} `json:"responses"`
}

type studiableItem struct {
	ID        int                     `json:"id"`
	CardSides []studiableItemCardSide `json:"cardSides"`
}

type studiableItemCardSide struct {
	Label cardSideLabel `json:"label"`
	Media []struct {
		PlainText string `json:"plainText"`
	}
}

type ModuleFetchingError struct {
	ID string
}

func (e *ModuleFetchingError) Error() string {
	return fmt.Sprintf("module \"%s\" fetching failed", e.ID)
}

type ModuleParsingError struct {
	ID string
}

func (e *ModuleParsingError) Error() string {
	return fmt.Sprintf("module \"%s\" parsing failed", e.ID)
}
