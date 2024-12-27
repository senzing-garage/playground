package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/senzing-garage/sz-sdk-go-grpc/szabstractfactory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ctx         = context.TODO()
	grpcAddress = "localhost:8261"
)

func testErr(err error) {
	if err != nil {
		panic(err)
	}
}

func asPrettyJSON(str string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return str
	}
	return prettyJSON.String()
}

func main() {
	grpcConnection, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	testErr(err)
	szAbstractFactory := &szabstractfactory.Szabstractfactory{
		GrpcConnection: grpcConnection,
	}

	szProduct, err := szAbstractFactory.CreateProduct(ctx)
	testErr(err)

	version, err := szProduct.GetVersion(ctx)
	testErr(err)

	fmt.Println(asPrettyJSON(version))

	fmt.Println("Hello, World!")
}
