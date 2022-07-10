package converter

import (
	"context"
	"errors"
	"os"
	"time"

	pb "github.com/amdf/conv-make-img/svc"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const imgDir = "/var/tengwar/img"

type TengwarConverter struct {
	ClientGRPC pb.TengwarConverterClient
	ClientConn *grpc.ClientConn
}

func NewTengwarConverter(tracer opentracing.Tracer, svcAddr string) (s *TengwarConverter, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	var conn *grpc.ClientConn

	conn, err = grpc.DialContext(ctx, svcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		return
	}

	c := pb.NewTengwarConverterClient(conn)
	s = &TengwarConverter{ClientGRPC: c, ClientConn: conn}
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

func SaveImage(ctx context.Context, id string, bytes []byte) error {
	if "" == id {
		return errors.New("empty id!")
	}
	parent := opentracing.SpanFromContext(ctx)
	sp := opentracing.StartSpan("SaveImage", opentracing.ChildOf(parent.Context()))
	defer sp.Finish()

	filename := imgDir + "/" + id + ".png"

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)

	return err
}
