package ethRepo

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/config"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/pkg/eth"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IEthereumRepository interface {
	SetAllHandlers(payHandler eth.EventHandlerFunc[models.TrafficPaidEvent], cncPayHandler eth.EventHandlerFunc[models.CNCTokenDepositedEvent] /*, add other handlers */)
	BackfillAll(ctx context.Context)
	ListenToAllEvents(ctx context.Context)
}

type EthereumRepository struct {
	logger             log.ILogger
	config             config.Web3Config
	client             *ethclient.Client
	paymentListener    eth.EventListener[models.TrafficPaidEvent]
	cncPaymentListener eth.EventListener[models.CNCTokenDepositedEvent]
	// add other events listener here
}

func NewEthereumRepository(logger log.ILogger, client *ethclient.Client, conf config.Web3Config) IEthereumRepository {
	paymentListener := eth.EventListener[models.TrafficPaidEvent]{Enabled: false}
	if conf.TrafficPaidEnabled {
		paymentContractABI := eth.MustLoadABI(conf.PaymentContractABI)
		paymentContractAddr := common.HexToAddress(conf.PaymentContractAddr)
		paymentListener = eth.NewEventListener[models.TrafficPaidEvent](&paymentContractABI, paymentContractAddr, "TrafficPaid", false, conf.TrafficPaidEnabled)
	}

	cncPaymentListener := eth.EventListener[models.CNCTokenDepositedEvent]{Enabled: false}
	if conf.CNCPaymentEnabled {
		cncPaymentEventABI := eth.MustLoadABI(conf.CNCEventABI)
		cncPaymentContractAddr := common.HexToAddress(conf.CNCContractAddr)
		cncPaymentListener = eth.NewEventListener[models.CNCTokenDepositedEvent](&cncPaymentEventABI, cncPaymentContractAddr, "TokenDeposited", true, conf.CNCPaymentEnabled)
	}

	return &EthereumRepository{
		logger:             logger,
		config:             conf,
		client:             client,
		paymentListener:    paymentListener,
		cncPaymentListener: cncPaymentListener,
	}
}

func (r *EthereumRepository) BackfillAll(ctx context.Context) {
	if r.paymentListener.Enabled {
		start := eth.LoadLastBlock(r.paymentListener.EventName)
		eth.Backfill[models.TrafficPaidEvent](r.logger, ctx, r.client, r.paymentListener, start)
	}

	if r.cncPaymentListener.Enabled {
		start := eth.LoadLastBlock(r.cncPaymentListener.EventName)
		eth.Backfill[models.CNCTokenDepositedEvent](r.logger, ctx, r.client, r.cncPaymentListener, start)
	}
}

func (r *EthereumRepository) SetAllHandlers(payHandler eth.EventHandlerFunc[models.TrafficPaidEvent], cncPayHandler eth.EventHandlerFunc[models.CNCTokenDepositedEvent]) {
	r.paymentListener.SetHandler(payHandler)
	r.cncPaymentListener.SetHandler(cncPayHandler)
}

func (r *EthereumRepository) ListenToAllEvents(ctx context.Context) {
	var wg sync.WaitGroup

	if r.paymentListener.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.paymentListener.Start(ctx, r.logger, r.client)
		}()
	}

	if r.cncPaymentListener.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.cncPaymentListener.Start(ctx, r.logger, r.client)
		}()
	}

	wg.Wait()
}
