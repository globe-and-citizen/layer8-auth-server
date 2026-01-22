package ethRepo

import (
	"context"
	"globe-and-citizen/layer8/auth-server/internal/config"
	"globe-and-citizen/layer8/auth-server/internal/models"
	eth2 "globe-and-citizen/layer8/auth-server/pkg/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IEthereumRepository interface {
	SetAllHandlers(payHandler eth2.EventHandlerFunc[models.TrafficPaidEvent] /*, add other handlers */)
	BackfillAll(ctx context.Context)
	ListenToAllEvents(ctx context.Context)
}

type EthereumRepository struct {
	config          config.Web3Config
	client          *ethclient.Client
	paymentListener eth2.EventListener[models.TrafficPaidEvent]
	// add other events listener here
}

func NewEthereumRepository(client *ethclient.Client, conf config.Web3Config) IEthereumRepository {
	paymentContractABI := eth2.MustLoadABI(conf.PaymentContractABI)
	paymentContractAddr := common.HexToAddress(conf.PaymentContractAddr)

	paymentListener := eth2.NewEventListener[models.TrafficPaidEvent](&paymentContractABI, paymentContractAddr, "TrafficPaid")

	return &EthereumRepository{
		config:          conf,
		client:          client,
		paymentListener: paymentListener,
	}
}

func (r *EthereumRepository) BackfillAll(ctx context.Context) {
	start := eth2.LoadLastBlock()
	eth2.Backfill[models.TrafficPaidEvent](ctx, r.client, r.paymentListener, start)
}

func (r *EthereumRepository) SetAllHandlers(payHandler eth2.EventHandlerFunc[models.TrafficPaidEvent]) {
	r.paymentListener.SetHandler(payHandler)
}

func (r *EthereumRepository) ListenToAllEvents(ctx context.Context) {
	r.paymentListener.Start(ctx, r.client)
}
