// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: scalekit/v1/directories/directories.proto

package directoriesconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	directories "github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/directories"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// DirectoryServiceName is the fully-qualified name of the DirectoryService service.
	DirectoryServiceName = "scalekit.v1.directories.DirectoryService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// DirectoryServiceCreateDirectoryProcedure is the fully-qualified name of the DirectoryService's
	// CreateDirectory RPC.
	DirectoryServiceCreateDirectoryProcedure = "/scalekit.v1.directories.DirectoryService/CreateDirectory"
	// DirectoryServiceUpdateDirectoryProcedure is the fully-qualified name of the DirectoryService's
	// UpdateDirectory RPC.
	DirectoryServiceUpdateDirectoryProcedure = "/scalekit.v1.directories.DirectoryService/UpdateDirectory"
	// DirectoryServiceAssignRolesProcedure is the fully-qualified name of the DirectoryService's
	// AssignRoles RPC.
	DirectoryServiceAssignRolesProcedure = "/scalekit.v1.directories.DirectoryService/AssignRoles"
	// DirectoryServiceUpdateAttributesProcedure is the fully-qualified name of the DirectoryService's
	// UpdateAttributes RPC.
	DirectoryServiceUpdateAttributesProcedure = "/scalekit.v1.directories.DirectoryService/UpdateAttributes"
	// DirectoryServiceGetDirectoryProcedure is the fully-qualified name of the DirectoryService's
	// GetDirectory RPC.
	DirectoryServiceGetDirectoryProcedure = "/scalekit.v1.directories.DirectoryService/GetDirectory"
	// DirectoryServiceListDirectoriesProcedure is the fully-qualified name of the DirectoryService's
	// ListDirectories RPC.
	DirectoryServiceListDirectoriesProcedure = "/scalekit.v1.directories.DirectoryService/ListDirectories"
	// DirectoryServiceEnableDirectoryProcedure is the fully-qualified name of the DirectoryService's
	// EnableDirectory RPC.
	DirectoryServiceEnableDirectoryProcedure = "/scalekit.v1.directories.DirectoryService/EnableDirectory"
	// DirectoryServiceDisableDirectoryProcedure is the fully-qualified name of the DirectoryService's
	// DisableDirectory RPC.
	DirectoryServiceDisableDirectoryProcedure = "/scalekit.v1.directories.DirectoryService/DisableDirectory"
	// DirectoryServiceListDirectoryUsersProcedure is the fully-qualified name of the DirectoryService's
	// ListDirectoryUsers RPC.
	DirectoryServiceListDirectoryUsersProcedure = "/scalekit.v1.directories.DirectoryService/ListDirectoryUsers"
	// DirectoryServiceListDirectoryGroupsProcedure is the fully-qualified name of the
	// DirectoryService's ListDirectoryGroups RPC.
	DirectoryServiceListDirectoryGroupsProcedure = "/scalekit.v1.directories.DirectoryService/ListDirectoryGroups"
	// DirectoryServiceCreateDirectorySecretProcedure is the fully-qualified name of the
	// DirectoryService's CreateDirectorySecret RPC.
	DirectoryServiceCreateDirectorySecretProcedure = "/scalekit.v1.directories.DirectoryService/CreateDirectorySecret"
	// DirectoryServiceRegenerateDirectorySecretProcedure is the fully-qualified name of the
	// DirectoryService's RegenerateDirectorySecret RPC.
	DirectoryServiceRegenerateDirectorySecretProcedure = "/scalekit.v1.directories.DirectoryService/RegenerateDirectorySecret"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	directoryServiceServiceDescriptor                         = directories.File_scalekit_v1_directories_directories_proto.Services().ByName("DirectoryService")
	directoryServiceCreateDirectoryMethodDescriptor           = directoryServiceServiceDescriptor.Methods().ByName("CreateDirectory")
	directoryServiceUpdateDirectoryMethodDescriptor           = directoryServiceServiceDescriptor.Methods().ByName("UpdateDirectory")
	directoryServiceAssignRolesMethodDescriptor               = directoryServiceServiceDescriptor.Methods().ByName("AssignRoles")
	directoryServiceUpdateAttributesMethodDescriptor          = directoryServiceServiceDescriptor.Methods().ByName("UpdateAttributes")
	directoryServiceGetDirectoryMethodDescriptor              = directoryServiceServiceDescriptor.Methods().ByName("GetDirectory")
	directoryServiceListDirectoriesMethodDescriptor           = directoryServiceServiceDescriptor.Methods().ByName("ListDirectories")
	directoryServiceEnableDirectoryMethodDescriptor           = directoryServiceServiceDescriptor.Methods().ByName("EnableDirectory")
	directoryServiceDisableDirectoryMethodDescriptor          = directoryServiceServiceDescriptor.Methods().ByName("DisableDirectory")
	directoryServiceListDirectoryUsersMethodDescriptor        = directoryServiceServiceDescriptor.Methods().ByName("ListDirectoryUsers")
	directoryServiceListDirectoryGroupsMethodDescriptor       = directoryServiceServiceDescriptor.Methods().ByName("ListDirectoryGroups")
	directoryServiceCreateDirectorySecretMethodDescriptor     = directoryServiceServiceDescriptor.Methods().ByName("CreateDirectorySecret")
	directoryServiceRegenerateDirectorySecretMethodDescriptor = directoryServiceServiceDescriptor.Methods().ByName("RegenerateDirectorySecret")
)

// DirectoryServiceClient is a client for the scalekit.v1.directories.DirectoryService service.
type DirectoryServiceClient interface {
	CreateDirectory(context.Context, *connect.Request[directories.CreateDirectoryRequest]) (*connect.Response[directories.CreateDirectoryResponse], error)
	UpdateDirectory(context.Context, *connect.Request[directories.UpdateDirectoryRequest]) (*connect.Response[directories.UpdateDirectoryResponse], error)
	AssignRoles(context.Context, *connect.Request[directories.AssignRolesRequest]) (*connect.Response[directories.AssignRolesResponse], error)
	UpdateAttributes(context.Context, *connect.Request[directories.UpdateAttributesRequest]) (*connect.Response[directories.UpdateAttributesResponse], error)
	GetDirectory(context.Context, *connect.Request[directories.GetDirectoryRequest]) (*connect.Response[directories.GetDirectoryResponse], error)
	ListDirectories(context.Context, *connect.Request[directories.ListDirectoriesRequest]) (*connect.Response[directories.ListDirectoriesResponse], error)
	EnableDirectory(context.Context, *connect.Request[directories.ToggleDirectoryRequest]) (*connect.Response[directories.ToggleDirectoryResponse], error)
	DisableDirectory(context.Context, *connect.Request[directories.ToggleDirectoryRequest]) (*connect.Response[directories.ToggleDirectoryResponse], error)
	ListDirectoryUsers(context.Context, *connect.Request[directories.ListDirectoryUsersRequest]) (*connect.Response[directories.ListDirectoryUsersResponse], error)
	ListDirectoryGroups(context.Context, *connect.Request[directories.ListDirectoryGroupsRequest]) (*connect.Response[directories.ListDirectoryGroupsResponse], error)
	CreateDirectorySecret(context.Context, *connect.Request[directories.CreateDirectorySecretRequest]) (*connect.Response[directories.CreateDirectorySecretResponse], error)
	RegenerateDirectorySecret(context.Context, *connect.Request[directories.RegenerateDirectorySecretRequest]) (*connect.Response[directories.RegenerateDirectorySecretResponse], error)
}

// NewDirectoryServiceClient constructs a client for the scalekit.v1.directories.DirectoryService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewDirectoryServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) DirectoryServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &directoryServiceClient{
		createDirectory: connect.NewClient[directories.CreateDirectoryRequest, directories.CreateDirectoryResponse](
			httpClient,
			baseURL+DirectoryServiceCreateDirectoryProcedure,
			connect.WithSchema(directoryServiceCreateDirectoryMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateDirectory: connect.NewClient[directories.UpdateDirectoryRequest, directories.UpdateDirectoryResponse](
			httpClient,
			baseURL+DirectoryServiceUpdateDirectoryProcedure,
			connect.WithSchema(directoryServiceUpdateDirectoryMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		assignRoles: connect.NewClient[directories.AssignRolesRequest, directories.AssignRolesResponse](
			httpClient,
			baseURL+DirectoryServiceAssignRolesProcedure,
			connect.WithSchema(directoryServiceAssignRolesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateAttributes: connect.NewClient[directories.UpdateAttributesRequest, directories.UpdateAttributesResponse](
			httpClient,
			baseURL+DirectoryServiceUpdateAttributesProcedure,
			connect.WithSchema(directoryServiceUpdateAttributesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getDirectory: connect.NewClient[directories.GetDirectoryRequest, directories.GetDirectoryResponse](
			httpClient,
			baseURL+DirectoryServiceGetDirectoryProcedure,
			connect.WithSchema(directoryServiceGetDirectoryMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listDirectories: connect.NewClient[directories.ListDirectoriesRequest, directories.ListDirectoriesResponse](
			httpClient,
			baseURL+DirectoryServiceListDirectoriesProcedure,
			connect.WithSchema(directoryServiceListDirectoriesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		enableDirectory: connect.NewClient[directories.ToggleDirectoryRequest, directories.ToggleDirectoryResponse](
			httpClient,
			baseURL+DirectoryServiceEnableDirectoryProcedure,
			connect.WithSchema(directoryServiceEnableDirectoryMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		disableDirectory: connect.NewClient[directories.ToggleDirectoryRequest, directories.ToggleDirectoryResponse](
			httpClient,
			baseURL+DirectoryServiceDisableDirectoryProcedure,
			connect.WithSchema(directoryServiceDisableDirectoryMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listDirectoryUsers: connect.NewClient[directories.ListDirectoryUsersRequest, directories.ListDirectoryUsersResponse](
			httpClient,
			baseURL+DirectoryServiceListDirectoryUsersProcedure,
			connect.WithSchema(directoryServiceListDirectoryUsersMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listDirectoryGroups: connect.NewClient[directories.ListDirectoryGroupsRequest, directories.ListDirectoryGroupsResponse](
			httpClient,
			baseURL+DirectoryServiceListDirectoryGroupsProcedure,
			connect.WithSchema(directoryServiceListDirectoryGroupsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createDirectorySecret: connect.NewClient[directories.CreateDirectorySecretRequest, directories.CreateDirectorySecretResponse](
			httpClient,
			baseURL+DirectoryServiceCreateDirectorySecretProcedure,
			connect.WithSchema(directoryServiceCreateDirectorySecretMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		regenerateDirectorySecret: connect.NewClient[directories.RegenerateDirectorySecretRequest, directories.RegenerateDirectorySecretResponse](
			httpClient,
			baseURL+DirectoryServiceRegenerateDirectorySecretProcedure,
			connect.WithSchema(directoryServiceRegenerateDirectorySecretMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// directoryServiceClient implements DirectoryServiceClient.
type directoryServiceClient struct {
	createDirectory           *connect.Client[directories.CreateDirectoryRequest, directories.CreateDirectoryResponse]
	updateDirectory           *connect.Client[directories.UpdateDirectoryRequest, directories.UpdateDirectoryResponse]
	assignRoles               *connect.Client[directories.AssignRolesRequest, directories.AssignRolesResponse]
	updateAttributes          *connect.Client[directories.UpdateAttributesRequest, directories.UpdateAttributesResponse]
	getDirectory              *connect.Client[directories.GetDirectoryRequest, directories.GetDirectoryResponse]
	listDirectories           *connect.Client[directories.ListDirectoriesRequest, directories.ListDirectoriesResponse]
	enableDirectory           *connect.Client[directories.ToggleDirectoryRequest, directories.ToggleDirectoryResponse]
	disableDirectory          *connect.Client[directories.ToggleDirectoryRequest, directories.ToggleDirectoryResponse]
	listDirectoryUsers        *connect.Client[directories.ListDirectoryUsersRequest, directories.ListDirectoryUsersResponse]
	listDirectoryGroups       *connect.Client[directories.ListDirectoryGroupsRequest, directories.ListDirectoryGroupsResponse]
	createDirectorySecret     *connect.Client[directories.CreateDirectorySecretRequest, directories.CreateDirectorySecretResponse]
	regenerateDirectorySecret *connect.Client[directories.RegenerateDirectorySecretRequest, directories.RegenerateDirectorySecretResponse]
}

// CreateDirectory calls scalekit.v1.directories.DirectoryService.CreateDirectory.
func (c *directoryServiceClient) CreateDirectory(ctx context.Context, req *connect.Request[directories.CreateDirectoryRequest]) (*connect.Response[directories.CreateDirectoryResponse], error) {
	return c.createDirectory.CallUnary(ctx, req)
}

// UpdateDirectory calls scalekit.v1.directories.DirectoryService.UpdateDirectory.
func (c *directoryServiceClient) UpdateDirectory(ctx context.Context, req *connect.Request[directories.UpdateDirectoryRequest]) (*connect.Response[directories.UpdateDirectoryResponse], error) {
	return c.updateDirectory.CallUnary(ctx, req)
}

// AssignRoles calls scalekit.v1.directories.DirectoryService.AssignRoles.
func (c *directoryServiceClient) AssignRoles(ctx context.Context, req *connect.Request[directories.AssignRolesRequest]) (*connect.Response[directories.AssignRolesResponse], error) {
	return c.assignRoles.CallUnary(ctx, req)
}

// UpdateAttributes calls scalekit.v1.directories.DirectoryService.UpdateAttributes.
func (c *directoryServiceClient) UpdateAttributes(ctx context.Context, req *connect.Request[directories.UpdateAttributesRequest]) (*connect.Response[directories.UpdateAttributesResponse], error) {
	return c.updateAttributes.CallUnary(ctx, req)
}

// GetDirectory calls scalekit.v1.directories.DirectoryService.GetDirectory.
func (c *directoryServiceClient) GetDirectory(ctx context.Context, req *connect.Request[directories.GetDirectoryRequest]) (*connect.Response[directories.GetDirectoryResponse], error) {
	return c.getDirectory.CallUnary(ctx, req)
}

// ListDirectories calls scalekit.v1.directories.DirectoryService.ListDirectories.
func (c *directoryServiceClient) ListDirectories(ctx context.Context, req *connect.Request[directories.ListDirectoriesRequest]) (*connect.Response[directories.ListDirectoriesResponse], error) {
	return c.listDirectories.CallUnary(ctx, req)
}

// EnableDirectory calls scalekit.v1.directories.DirectoryService.EnableDirectory.
func (c *directoryServiceClient) EnableDirectory(ctx context.Context, req *connect.Request[directories.ToggleDirectoryRequest]) (*connect.Response[directories.ToggleDirectoryResponse], error) {
	return c.enableDirectory.CallUnary(ctx, req)
}

// DisableDirectory calls scalekit.v1.directories.DirectoryService.DisableDirectory.
func (c *directoryServiceClient) DisableDirectory(ctx context.Context, req *connect.Request[directories.ToggleDirectoryRequest]) (*connect.Response[directories.ToggleDirectoryResponse], error) {
	return c.disableDirectory.CallUnary(ctx, req)
}

// ListDirectoryUsers calls scalekit.v1.directories.DirectoryService.ListDirectoryUsers.
func (c *directoryServiceClient) ListDirectoryUsers(ctx context.Context, req *connect.Request[directories.ListDirectoryUsersRequest]) (*connect.Response[directories.ListDirectoryUsersResponse], error) {
	return c.listDirectoryUsers.CallUnary(ctx, req)
}

// ListDirectoryGroups calls scalekit.v1.directories.DirectoryService.ListDirectoryGroups.
func (c *directoryServiceClient) ListDirectoryGroups(ctx context.Context, req *connect.Request[directories.ListDirectoryGroupsRequest]) (*connect.Response[directories.ListDirectoryGroupsResponse], error) {
	return c.listDirectoryGroups.CallUnary(ctx, req)
}

// CreateDirectorySecret calls scalekit.v1.directories.DirectoryService.CreateDirectorySecret.
func (c *directoryServiceClient) CreateDirectorySecret(ctx context.Context, req *connect.Request[directories.CreateDirectorySecretRequest]) (*connect.Response[directories.CreateDirectorySecretResponse], error) {
	return c.createDirectorySecret.CallUnary(ctx, req)
}

// RegenerateDirectorySecret calls
// scalekit.v1.directories.DirectoryService.RegenerateDirectorySecret.
func (c *directoryServiceClient) RegenerateDirectorySecret(ctx context.Context, req *connect.Request[directories.RegenerateDirectorySecretRequest]) (*connect.Response[directories.RegenerateDirectorySecretResponse], error) {
	return c.regenerateDirectorySecret.CallUnary(ctx, req)
}

// DirectoryServiceHandler is an implementation of the scalekit.v1.directories.DirectoryService
// service.
type DirectoryServiceHandler interface {
	CreateDirectory(context.Context, *connect.Request[directories.CreateDirectoryRequest]) (*connect.Response[directories.CreateDirectoryResponse], error)
	UpdateDirectory(context.Context, *connect.Request[directories.UpdateDirectoryRequest]) (*connect.Response[directories.UpdateDirectoryResponse], error)
	AssignRoles(context.Context, *connect.Request[directories.AssignRolesRequest]) (*connect.Response[directories.AssignRolesResponse], error)
	UpdateAttributes(context.Context, *connect.Request[directories.UpdateAttributesRequest]) (*connect.Response[directories.UpdateAttributesResponse], error)
	GetDirectory(context.Context, *connect.Request[directories.GetDirectoryRequest]) (*connect.Response[directories.GetDirectoryResponse], error)
	ListDirectories(context.Context, *connect.Request[directories.ListDirectoriesRequest]) (*connect.Response[directories.ListDirectoriesResponse], error)
	EnableDirectory(context.Context, *connect.Request[directories.ToggleDirectoryRequest]) (*connect.Response[directories.ToggleDirectoryResponse], error)
	DisableDirectory(context.Context, *connect.Request[directories.ToggleDirectoryRequest]) (*connect.Response[directories.ToggleDirectoryResponse], error)
	ListDirectoryUsers(context.Context, *connect.Request[directories.ListDirectoryUsersRequest]) (*connect.Response[directories.ListDirectoryUsersResponse], error)
	ListDirectoryGroups(context.Context, *connect.Request[directories.ListDirectoryGroupsRequest]) (*connect.Response[directories.ListDirectoryGroupsResponse], error)
	CreateDirectorySecret(context.Context, *connect.Request[directories.CreateDirectorySecretRequest]) (*connect.Response[directories.CreateDirectorySecretResponse], error)
	RegenerateDirectorySecret(context.Context, *connect.Request[directories.RegenerateDirectorySecretRequest]) (*connect.Response[directories.RegenerateDirectorySecretResponse], error)
}

// NewDirectoryServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewDirectoryServiceHandler(svc DirectoryServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	directoryServiceCreateDirectoryHandler := connect.NewUnaryHandler(
		DirectoryServiceCreateDirectoryProcedure,
		svc.CreateDirectory,
		connect.WithSchema(directoryServiceCreateDirectoryMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceUpdateDirectoryHandler := connect.NewUnaryHandler(
		DirectoryServiceUpdateDirectoryProcedure,
		svc.UpdateDirectory,
		connect.WithSchema(directoryServiceUpdateDirectoryMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceAssignRolesHandler := connect.NewUnaryHandler(
		DirectoryServiceAssignRolesProcedure,
		svc.AssignRoles,
		connect.WithSchema(directoryServiceAssignRolesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceUpdateAttributesHandler := connect.NewUnaryHandler(
		DirectoryServiceUpdateAttributesProcedure,
		svc.UpdateAttributes,
		connect.WithSchema(directoryServiceUpdateAttributesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceGetDirectoryHandler := connect.NewUnaryHandler(
		DirectoryServiceGetDirectoryProcedure,
		svc.GetDirectory,
		connect.WithSchema(directoryServiceGetDirectoryMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceListDirectoriesHandler := connect.NewUnaryHandler(
		DirectoryServiceListDirectoriesProcedure,
		svc.ListDirectories,
		connect.WithSchema(directoryServiceListDirectoriesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceEnableDirectoryHandler := connect.NewUnaryHandler(
		DirectoryServiceEnableDirectoryProcedure,
		svc.EnableDirectory,
		connect.WithSchema(directoryServiceEnableDirectoryMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceDisableDirectoryHandler := connect.NewUnaryHandler(
		DirectoryServiceDisableDirectoryProcedure,
		svc.DisableDirectory,
		connect.WithSchema(directoryServiceDisableDirectoryMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceListDirectoryUsersHandler := connect.NewUnaryHandler(
		DirectoryServiceListDirectoryUsersProcedure,
		svc.ListDirectoryUsers,
		connect.WithSchema(directoryServiceListDirectoryUsersMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceListDirectoryGroupsHandler := connect.NewUnaryHandler(
		DirectoryServiceListDirectoryGroupsProcedure,
		svc.ListDirectoryGroups,
		connect.WithSchema(directoryServiceListDirectoryGroupsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceCreateDirectorySecretHandler := connect.NewUnaryHandler(
		DirectoryServiceCreateDirectorySecretProcedure,
		svc.CreateDirectorySecret,
		connect.WithSchema(directoryServiceCreateDirectorySecretMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	directoryServiceRegenerateDirectorySecretHandler := connect.NewUnaryHandler(
		DirectoryServiceRegenerateDirectorySecretProcedure,
		svc.RegenerateDirectorySecret,
		connect.WithSchema(directoryServiceRegenerateDirectorySecretMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/scalekit.v1.directories.DirectoryService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case DirectoryServiceCreateDirectoryProcedure:
			directoryServiceCreateDirectoryHandler.ServeHTTP(w, r)
		case DirectoryServiceUpdateDirectoryProcedure:
			directoryServiceUpdateDirectoryHandler.ServeHTTP(w, r)
		case DirectoryServiceAssignRolesProcedure:
			directoryServiceAssignRolesHandler.ServeHTTP(w, r)
		case DirectoryServiceUpdateAttributesProcedure:
			directoryServiceUpdateAttributesHandler.ServeHTTP(w, r)
		case DirectoryServiceGetDirectoryProcedure:
			directoryServiceGetDirectoryHandler.ServeHTTP(w, r)
		case DirectoryServiceListDirectoriesProcedure:
			directoryServiceListDirectoriesHandler.ServeHTTP(w, r)
		case DirectoryServiceEnableDirectoryProcedure:
			directoryServiceEnableDirectoryHandler.ServeHTTP(w, r)
		case DirectoryServiceDisableDirectoryProcedure:
			directoryServiceDisableDirectoryHandler.ServeHTTP(w, r)
		case DirectoryServiceListDirectoryUsersProcedure:
			directoryServiceListDirectoryUsersHandler.ServeHTTP(w, r)
		case DirectoryServiceListDirectoryGroupsProcedure:
			directoryServiceListDirectoryGroupsHandler.ServeHTTP(w, r)
		case DirectoryServiceCreateDirectorySecretProcedure:
			directoryServiceCreateDirectorySecretHandler.ServeHTTP(w, r)
		case DirectoryServiceRegenerateDirectorySecretProcedure:
			directoryServiceRegenerateDirectorySecretHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedDirectoryServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedDirectoryServiceHandler struct{}

func (UnimplementedDirectoryServiceHandler) CreateDirectory(context.Context, *connect.Request[directories.CreateDirectoryRequest]) (*connect.Response[directories.CreateDirectoryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.CreateDirectory is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) UpdateDirectory(context.Context, *connect.Request[directories.UpdateDirectoryRequest]) (*connect.Response[directories.UpdateDirectoryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.UpdateDirectory is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) AssignRoles(context.Context, *connect.Request[directories.AssignRolesRequest]) (*connect.Response[directories.AssignRolesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.AssignRoles is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) UpdateAttributes(context.Context, *connect.Request[directories.UpdateAttributesRequest]) (*connect.Response[directories.UpdateAttributesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.UpdateAttributes is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) GetDirectory(context.Context, *connect.Request[directories.GetDirectoryRequest]) (*connect.Response[directories.GetDirectoryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.GetDirectory is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) ListDirectories(context.Context, *connect.Request[directories.ListDirectoriesRequest]) (*connect.Response[directories.ListDirectoriesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.ListDirectories is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) EnableDirectory(context.Context, *connect.Request[directories.ToggleDirectoryRequest]) (*connect.Response[directories.ToggleDirectoryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.EnableDirectory is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) DisableDirectory(context.Context, *connect.Request[directories.ToggleDirectoryRequest]) (*connect.Response[directories.ToggleDirectoryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.DisableDirectory is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) ListDirectoryUsers(context.Context, *connect.Request[directories.ListDirectoryUsersRequest]) (*connect.Response[directories.ListDirectoryUsersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.ListDirectoryUsers is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) ListDirectoryGroups(context.Context, *connect.Request[directories.ListDirectoryGroupsRequest]) (*connect.Response[directories.ListDirectoryGroupsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.ListDirectoryGroups is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) CreateDirectorySecret(context.Context, *connect.Request[directories.CreateDirectorySecretRequest]) (*connect.Response[directories.CreateDirectorySecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.CreateDirectorySecret is not implemented"))
}

func (UnimplementedDirectoryServiceHandler) RegenerateDirectorySecret(context.Context, *connect.Request[directories.RegenerateDirectorySecretRequest]) (*connect.Response[directories.RegenerateDirectorySecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("scalekit.v1.directories.DirectoryService.RegenerateDirectorySecret is not implemented"))
}
