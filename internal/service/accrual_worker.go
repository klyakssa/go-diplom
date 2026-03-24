package service

import (
	"context"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	ordersPkg "github.com/klyakssa/go-diplom.git/internal/domain/orders"
	"github.com/klyakssa/go-diplom.git/internal/logger"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// AccrualWorker
type AccrualWorker struct {
	client    *resty.Client
	url       string
	orderRepo ordersPkg.Repository
	log       *logger.Logger
}

// NewAccrualWorker creates new instance of AccrualWorker
func NewAccrualWorker(log *logger.Logger, orderRepo ordersPkg.Repository, url string) *AccrualWorker {
	return &AccrualWorker{
		client: resty.New().
			SetTimeout(2 * time.Second),
		url:       url,
		orderRepo: orderRepo,
		log:       log,
	}
}

// Start starts the AccrualWorker
func (w *AccrualWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}

// AccrualResponse
type AccrualResponse struct {
	Order   string          `json:"order"` // order number
	Status  string          `json:"status"` // order status
	Accrual decimal.Decimal `json:"accrual"` // order accrual
}

// process processes pending orders
func (w *AccrualWorker) process(ctx context.Context) {
	orders, err := w.orderRepo.GetPendingOrders(ctx)
	if err != nil {
		w.log.Error("failed to get pending orders", zap.Error(err))
		return
	}

	for _, order := range orders {

		w.log.Debug("get order",
			zap.String("number", order.Number),
			zap.String("user_id", order.UserID),
			zap.String("status", order.Status),
			zap.Any("accrual", order.Accrual),
		)

		var resp AccrualResponse

		res, err := w.client.R().
			SetContext(ctx).
			SetResult(&resp).
			Get(w.url + "/api/orders/" + order.Number)

		if err != nil {
			w.log.Error("failed to get order", zap.Error(err))
			continue
		}

		if res.IsError() {
			w.log.Error("failed to get order", zap.Int("code", res.StatusCode()))
			continue
		}

		if res.StatusCode() == http.StatusTooManyRequests {
			w.log.Warn("too many requests", zap.Int("code", res.StatusCode()))
			time.Sleep(60 * time.Second)
			continue
		}

		w.log.Debug("accrual response received",
			zap.String("order", resp.Order),
			zap.String("status", resp.Status),
			zap.Any("accrual", resp.Accrual),
			zap.Int("code", res.StatusCode()),
		)

		err = w.orderRepo.UpdateOrder(ctx, resp.Order, resp.Status, resp.Accrual)
		if err != nil {
			w.log.Error("failed to update order", zap.Error(err))
			continue
		}

		if resp.Status == "PROCESSED" && resp.Accrual.GreaterThan(decimal.Zero) {
			err = w.orderRepo.ApplyAccrual(ctx, &ordersPkg.Order{
				Number:  resp.Order,
				UserID:  order.UserID,
				Accrual: resp.Accrual,
			})
			if err != nil {
				w.log.Error("failed to add balance or apply accrual", zap.Error(err))
				continue
			}
		}
	}
}
