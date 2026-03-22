package repositoryrpc

import (
	"context"

	"google.golang.org/grpc"
)

const (
	ServiceName           = "collector.RepositoryService"
	GetRepositoryFullName = "/" + ServiceName + "/GetRepository"
)

type GetRepositoryRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type RepositoryResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	Forks       int64  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type RepositoryServiceServer interface {
	GetRepository(context.Context, *GetRepositoryRequest) (*RepositoryResponse, error)
}

type RepositoryServiceClient interface {
	GetRepository(ctx context.Context, req *GetRepositoryRequest, opts ...grpc.CallOption) (*RepositoryResponse, error)
}

type repositoryServiceClient struct {
	conn grpc.ClientConnInterface
}

func NewRepositoryServiceClient(conn grpc.ClientConnInterface) RepositoryServiceClient {
	return &repositoryServiceClient{conn: conn}
}

func (c *repositoryServiceClient) GetRepository(ctx context.Context, req *GetRepositoryRequest, opts ...grpc.CallOption) (*RepositoryResponse, error) {
	resp := new(RepositoryResponse)
	if err := c.conn.Invoke(ctx, GetRepositoryFullName, req, resp, opts...); err != nil {
		return nil, err
	}

	return resp, nil
}

func RegisterRepositoryServiceServer(server grpc.ServiceRegistrar, svc RepositoryServiceServer) {
	server.RegisterService(&grpc.ServiceDesc{
		ServiceName: ServiceName,
		HandlerType: (*RepositoryServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetRepository",
				Handler:    getRepositoryHandler(svc),
			},
		},
	}, svc)
}

func getRepositoryHandler(svc RepositoryServiceServer) grpc.MethodHandler {
	return func(server any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
		req := new(GetRepositoryRequest)
		if err := dec(req); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return svc.GetRepository(ctx, req)
		}

		info := &grpc.UnaryServerInfo{
			Server:     server,
			FullMethod: GetRepositoryFullName,
		}

		handler := func(ctx context.Context, request any) (any, error) {
			return svc.GetRepository(ctx, request.(*GetRepositoryRequest))
		}

		return interceptor(ctx, req, info, handler)
	}
}
