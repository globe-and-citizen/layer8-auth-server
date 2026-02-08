package ethRepo

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/config"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/pkg/eth"
	"globe-and-citizen/layer8/auth-server/pkg/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IEthereumRepository interface {
	SetAllHandlers(payHandler eth.EventHandlerFunc[models.TrafficPaidEvent] /*, add other handlers */)
	BackfillAll(ctx context.Context)
	ListenToAllEvents(ctx context.Context)
}

type EthereumRepository struct {
	logger          log.ILogger
	config          config.Web3Config
	client          *ethclient.Client
	paymentListener eth.EventListener[models.TrafficPaidEvent]
	// add other events listener here
}

func NewEthereumRepository(logger log.ILogger, client *ethclient.Client, conf config.Web3Config) IEthereumRepository {
	paymentContractABI := eth.MustLoadABI(conf.PaymentContractABI)
	paymentContractAddr := common.HexToAddress(conf.PaymentContractAddr)

	paymentListener := eth.NewEventListener[models.TrafficPaidEvent](&paymentContractABI, paymentContractAddr, "TrafficPaid")

	return &EthereumRepository{
		logger:          logger,
		config:          conf,
		client:          client,
		paymentListener: paymentListener,
	}
}

func (r *EthereumRepository) BackfillAll(ctx context.Context) {
	start := eth.LoadLastBlock()
	eth.Backfill[models.TrafficPaidEvent](r.logger, ctx, r.client, r.paymentListener, start)
}

func (r *EthereumRepository) SetAllHandlers(payHandler eth.EventHandlerFunc[models.TrafficPaidEvent]) {
	r.paymentListener.SetHandler(payHandler)
}

func (r *EthereumRepository) ListenToAllEvents(ctx context.Context) {
	r.paymentListener.Start(ctx, r.logger, r.client)
}
