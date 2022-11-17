package sdk

import (
	cacheSchema "github.com/everFinance/go-everpay/cache/schema"
	"github.com/everFinance/go-everpay/sdk/schema"
	"sync"
	"time"
)

type SubscribeTx struct {
	client *Client

	ch          chan cacheSchema.TxResponse
	filterQuery schema.FilterQuery
	quit        chan struct{}
	quitOnce    sync.Once
}

func newSubscribeTx(c *Client, fq schema.FilterQuery) *SubscribeTx {
	return &SubscribeTx{
		client:      c,
		ch:          make(chan cacheSchema.TxResponse),
		filterQuery: fq,
		quit:        make(chan struct{}),
	}
}

func (s *SubscribeTx) run() {
	isAccTxs := s.filterQuery.Address != ""
	interval := 1 * time.Second
	t1 := time.NewTimer(interval)
	cursorId := s.filterQuery.StartCursor
	for {
		var txs schema.Txs
		var err error
		t1.Reset(interval)
		select {
		case <-t1.C:
			if isAccTxs {
				txs, err = s.client.CursorTxsByAcc(s.filterQuery.Address, cursorId, s.filterQuery.TokenTag, s.filterQuery.Action, s.filterQuery.WithoutAction)
			} else {
				txs, err = s.client.CursorTxs(cursorId, s.filterQuery.TokenTag, s.filterQuery.Action, s.filterQuery.WithoutAction)
			}

			if err != nil {
				interval = 10 * time.Second
				continue
			}

			for _, tx := range txs.Txs {
				s.ch <- *tx
			}

			num := len(txs.Txs)
			if num > 0 {
				cursorId = txs.Txs[num-1].RawId
				interval = 1 * time.Second
			} else {
				interval = 5 * time.Second
			}
		case <-s.quit:
			log.Debug("Unsubscribe txs")
			return
		}
	}
}

func (s *SubscribeTx) Subscribe() <-chan cacheSchema.TxResponse {
	return s.ch
}

func (s *SubscribeTx) Unsubscribe() {
	s.quitOnce.Do(func() {
		close(s.quit)
	})
}
