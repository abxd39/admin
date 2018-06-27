package client

import (
	cf "admin/gateway/conf"
	proto "admin/proto/rpc"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
)

type ContentRPCCli struct {
	conn proto.BackstageRPCService
}

func NewCurrencyRPCCli() (u *ContentRPCCli) {
	consul_addr := cf.Cfg.MustValue("consul", "addr")
	r := consul.NewRegistry(registry.Addrs(consul_addr))
	service := micro.NewService(
		micro.Name("currency.client"),
		micro.Registry(r),
	)
	service.Init()

	service_name := cf.Cfg.MustValue("base", "service_client_currency")
	greeter := proto.NewBackstageRPCService(service_name, service.Client())
	u = &ContentRPCCli{
		conn: greeter,
	}
	return
}

// func (s *contentRPCCli) CallAdmin(name string) (rsp *proto.AdminResponse, err error) {
// 	rsp, err = s.conn.AdminCmd(context.TODO(), &proto.AdminRequest{})
// 	if err != nil {
// 		Log.Errorln(err.Error())
// 		return
// 	}
// 	return
// }
