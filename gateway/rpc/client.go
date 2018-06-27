package rpc

import "admin/gateway/rpc/client"

var InnerService *RPCClient

type RPCClient struct {
	BackstageSevice *client.ContentRPCCli
	// CurrencyService *client.CurrencyRPCCli
	// TokenService    *client.TokenRPCCli
	// WallService     *client.WalletRPCCli
	// PublicService   *client.PublciRPCCli
	// WalletSevice    *client.WalletRPCCli
	// KineService     *client.KlineRPCCli
}

func NewRPCClient() (c *RPCClient) {
	c = &RPCClient{

		BackstageSevice: client.NewBackstageRPCCli(),
		// CurrencyService: client.NewCurrencyRPCCli(),
		// TokenService:    client.NewTokenRPCCli(),
		// WallService:     client.NewWalletRPCCli(),
		// PublicService:   client.NewPublciRPCCli(),
		// WalletSevice:    client.NewWalletRPCCli(),
		// KineService:     client.NewKlineRPCCli(),
	}
	return c
}

func InitInnerService() {
	InnerService = NewRPCClient()
}
