package service

import (
	"context"

	"rest-api-database-pricelists/internal/dto"
	"rest-api-database-pricelists/internal/repository"

	"go.uber.org/zap"
)

type SearchService struct {
	repo   *repository.ProductRepository
	logger *zap.Logger
}

func NewSearchService(repo *repository.ProductRepository, logger *zap.Logger) *SearchService {
	return &SearchService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SearchService) Search(ctx context.Context, mpn string, qty int) ([]dto.SearchItem, error) {

	if s.logger != nil {
		s.logger.Info("service search started")
		zap.String("mpn:", mpn)
		zap.Int("qty", qty)
	}

	rows, err := s.repo.SearchByName(ctx, mpn)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []dto.SearchItem{
			{
				RequestedMPN: mpn,
				RequestedQty: qty,
				Status:       "Не найдено",
			},
		}, nil
	}

	grouped := make(map[string][]repository.ProductRow)

	for _, r := range rows {
		grouped[r.Supplier] = append(grouped[r.Supplier], r)
	}

	var result []dto.SearchItem

	for supplier, items := range grouped {

		first := items[0]

		var breaks []dto.PriceBreak
		var basePrice float64
		var currency string

		for _, item := range items {
			if qty >= item.Moq {
				breaks = append(breaks, dto.PriceBreak{
					Quantity: item.Quant,
					Price:    item.Price,
					Currency: item.Currency,
				})
			}

			if item.Quant <= qty {
				basePrice = item.Price
				currency = item.Currency
			}
		}

		result = append(result, dto.SearchItem{
			MPN:          first.Code,
			RequestedMPN: mpn,
			RequestedQty: qty,
			Manufacturer: first.Producer,
			Stock:        first.QtyAvailable,
			Status:       "Найдено",
			Price:        basePrice,
			Currency:     currency,
			PriceBreaks:  breaks,
			Supplier:     supplier,
		})
	}

	s.logger.Info("service search completed",
		zap.Int("result_count", len(result)),
	)

	return result, nil
}
