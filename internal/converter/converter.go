package converter

import (
	"context"

	pb "github.com/amdf/conv-make-img/svc"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const svcAddr = "[::]:50051"

type TengwarConverter struct {
	ClientGRPC pb.TengwarConverterClient
}

func NewTengwarConverter() (s *TengwarConverter, err error) {

	conn, err := grpc.Dial(svcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return
	}

	//TODO: close conn somewhere?

	c := pb.NewTengwarConverterClient(conn)
	s = &TengwarConverter{ClientGRPC: c}
	return
}

func (client TengwarConverter) MakeImage(ctx context.Context, rq ConvertRequest) (bytes []byte, err error) {

	var body *httpbody.HttpBody
	body, err = client.ClientGRPC.MakeImage(ctx, &pb.ConvertRequest{
		InputText: rq.InputText,
		FontSize:  rq.FontSize,
		FontFile:  rq.FontFile,
		FontStyle: pb.ConvertRequest_FontStyles(rq.FontSize),
	})

	if err != nil {
		return
	}

	bytes = body.Data

	return
}
