package mw

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
)

func Logging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	raw, _ := protojson.Marshal((req).(proto.Message))
	log.Printf("Request: Method: %v, REQ: %v\n", info.FullMethod, string(raw))

	if resp, err = handler(ctx, req); err != nil {
		log.Printf("Response: Method: %v, err: %v\n", info.FullMethod, err)
		return
	}

	rawResp, _ := protojson.Marshal((resp).(proto.Message))
	log.Printf("Response: Method: %v, RESP: %v\n", info.FullMethod, string(rawResp))

	return
}

func HTTPLogging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: Method=%s, URL=%s, Headers=%v", r.Method, r.URL.String(), r.Header)
		next.ServeHTTP(w, r)
		log.Printf("Request processed: Method=%s, URL=%s", r.Method, r.URL.String())
	}

	return http.HandlerFunc(fn)
}
