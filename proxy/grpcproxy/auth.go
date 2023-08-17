package grpcproxy

import (
	"context"

	"oldnicke/etcd/clientv3"
	pb "oldnicke/etcd/etcdserver/etcdserverpb"
)

type AuthProxy struct {
	client *clientv3.Client
}

func NewAuthProxy(c *clientv3.Client) pb.AuthServer {
	return &AuthProxy{client: c}
}

func (ap *AuthProxy) AuthEnable(ctx context.Context, r *pb.AuthEnableRequest) (*pb.AuthEnableResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).AuthEnable(ctx, r)
}

func (ap *AuthProxy) AuthDisable(ctx context.Context, r *pb.AuthDisableRequest) (*pb.AuthDisableResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).AuthDisable(ctx, r)
}

func (ap *AuthProxy) Authenticate(ctx context.Context, r *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).Authenticate(ctx, r)
}

func (ap *AuthProxy) RoleAdd(ctx context.Context, r *pb.AuthRoleAddRequest) (*pb.AuthRoleAddResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).RoleAdd(ctx, r)
}

func (ap *AuthProxy) RoleDelete(ctx context.Context, r *pb.AuthRoleDeleteRequest) (*pb.AuthRoleDeleteResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).RoleDelete(ctx, r)
}

func (ap *AuthProxy) RoleGet(ctx context.Context, r *pb.AuthRoleGetRequest) (*pb.AuthRoleGetResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).RoleGet(ctx, r)
}

func (ap *AuthProxy) RoleList(ctx context.Context, r *pb.AuthRoleListRequest) (*pb.AuthRoleListResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).RoleList(ctx, r)
}

func (ap *AuthProxy) RoleRevokePermission(ctx context.Context, r *pb.AuthRoleRevokePermissionRequest) (*pb.AuthRoleRevokePermissionResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).RoleRevokePermission(ctx, r)
}

func (ap *AuthProxy) RoleGrantPermission(ctx context.Context, r *pb.AuthRoleGrantPermissionRequest) (*pb.AuthRoleGrantPermissionResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).RoleGrantPermission(ctx, r)
}

func (ap *AuthProxy) UserAdd(ctx context.Context, r *pb.AuthUserAddRequest) (*pb.AuthUserAddResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).UserAdd(ctx, r)
}

func (ap *AuthProxy) UserDelete(ctx context.Context, r *pb.AuthUserDeleteRequest) (*pb.AuthUserDeleteResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).UserDelete(ctx, r)
}

func (ap *AuthProxy) UserGet(ctx context.Context, r *pb.AuthUserGetRequest) (*pb.AuthUserGetResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).UserGet(ctx, r)
}

func (ap *AuthProxy) UserList(ctx context.Context, r *pb.AuthUserListRequest) (*pb.AuthUserListResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).UserList(ctx, r)
}

func (ap *AuthProxy) UserGrantRole(ctx context.Context, r *pb.AuthUserGrantRoleRequest) (*pb.AuthUserGrantRoleResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).UserGrantRole(ctx, r)
}

func (ap *AuthProxy) UserRevokeRole(ctx context.Context, r *pb.AuthUserRevokeRoleRequest) (*pb.AuthUserRevokeRoleResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).UserRevokeRole(ctx, r)
}

func (ap *AuthProxy) UserChangePassword(ctx context.Context, r *pb.AuthUserChangePasswordRequest) (*pb.AuthUserChangePasswordResponse, error) {
	conn := ap.client.ActiveConnection()
	return pb.NewAuthClient(conn).UserChangePassword(ctx, r)
}
