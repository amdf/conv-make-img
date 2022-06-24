package converter

import (
	pb "github.com/amdf/conv-make-img/svc"
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
