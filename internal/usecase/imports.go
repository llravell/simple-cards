package usecase

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"github.com/llravell/simple-cards/internal/entity"
	"github.com/rs/zerolog"
)

const csvRecordMinLength = 2

type QuizletImportWork struct {
	repo                ModulesRepository
	quizletModuleParser QuizletModuleParser
	log                 *zerolog.Logger
	module              *entity.Module
	quizletModuleID     string
}

func (w *QuizletImportWork) Do(ctx context.Context) {
	quizletCards, err := w.quizletModuleParser.Parse(ctx, w.quizletModuleID)
	if err != nil {
		w.log.Error().Err(err).Msg("quizlet module parsing failed")

		return
	}

	if len(quizletCards) == 0 {
		return
	}

	w.log.Info().Msgf("quizlet module \"%s\" parsed", w.quizletModuleID)

	moduleCards := make([]*entity.Card, 0, len(quizletCards))

	for _, quizletCard := range quizletCards {
		card := &entity.Card{
			Term:    quizletCard.Front,
			Meaning: quizletCard.Back,
		}

		moduleCards = append(moduleCards, card)
	}

	err = w.repo.CreateNewModuleWithCards(
		ctx,
		&entity.ModuleWithCards{
			Module: *w.module,
			Cards:  moduleCards,
		},
	)
	if err != nil {
		w.log.Error().Err(err).Msg("module from quizlet storing failed")
	} else {
		w.log.Info().Msgf("quizlet module \"%s\" imported", w.quizletModuleID)
	}
}

type CSVImportWork struct {
	repo   ModulesRepository
	log    *zerolog.Logger
	module *entity.Module
	reader io.ReadCloser
}

func (w *CSVImportWork) Do(ctx context.Context) {
	defer w.reader.Close()

	moduleCards := make([]*entity.Card, 0)
	csvReader := csv.NewReader(w.reader)

	for {
		var (
			record []string
			err    error
		)

		select {
		case <-ctx.Done():
			w.log.Error().Msg("import work has been interrupted")

			return
		default:
			record, err = csvReader.Read()
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			w.log.Error().Err(err).Msg("csv reading error")

			return
		}

		if len(record) < csvRecordMinLength {
			continue
		}

		card := &entity.Card{
			Term:       strings.TrimSpace(record[0]),
			Meaning:    strings.TrimSpace(record[1]),
			ModuleUUID: w.module.UUID,
		}

		if card.Term != "" && card.Meaning != "" {
			moduleCards = append(moduleCards, card)
		}
	}

	err := w.repo.CreateNewModuleWithCards(
		ctx,
		&entity.ModuleWithCards{
			Module: *w.module,
			Cards:  moduleCards,
		},
	)
	if err != nil {
		w.log.Error().Err(err).Msg("module from csv storing failed")
	} else {
		w.log.Info().Msg("csv module imported")
	}
}
