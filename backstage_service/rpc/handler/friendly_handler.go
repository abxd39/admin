package handler

import (
	"admin/backstage_service/log"
	"admin/backstage_service/module"
	def "admin/proto/common"
	proto "admin/proto/rpc"
	"context"
	"fmt"
)

type RPCServer struct{}

func (s *RPCServer) Test(ctx context.Context, req *proto.TestRequest, rsp *proto.TestResponse) error {
	rsp.Code = def.ERRCODE_SUCCESS
	rsp.Msg = "Hello BACKSTAGE SERVICE ^v^ !!"
	return nil
}

func (s *RPCServer) AddFriendlyLink(ctx context.Context, req *proto.AddFriendlyLinkRequest, rsp *proto.AddFriendlyLinkResponse) error {
	f := module.FriendlyLink{}
	err := f.Add(req, rsp)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}

func (s *RPCServer) GetFriendlyLink(ctx context.Context, req *proto.FriendlyLinkRequest, rsp *proto.FriendlyLinkResponse) error {
	fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
	f := module.FriendlyLink{}
	err := f.GetFriendlyLinkList(req, rsp)
	if err != nil {
		log.Log.Errorf(err.Error())
	}
	return nil
}
