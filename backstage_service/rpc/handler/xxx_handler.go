package handler

import (
	def "admin/proto/common"
	proto "admin/proto/rpc"
	"context"
)

type RPCServer struct{}

func (s *RPCServer) Test(ctx context.Context, req *proto.TestRequest, rsp *proto.TestResponse) error {
	rsp.Code = def.ERRCODE_SUCCESS
	rsp.Msg = "Hello BACKSTAGE SERVICE ^v^ !!"
	return nil
}
