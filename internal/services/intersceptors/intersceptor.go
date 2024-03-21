package intersceptors

// import (
// 	"context"
// 	"fmt"

// 	"google.golang.org/grpc"
// )

// func UnaryInterscaptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
// 	fmt.Println(info.FullMethod, " ---this method was called")
// 	compID := ctx.Value("companyID")
// 	if compID != nil {
// 		fmt.Println(compID.(string), " --- companyID")
// 	} else {
// 		fmt.Println("compID was nil")
// 	}
// 	return handler(ctx, req)
// }
