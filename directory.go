package scalekit

import (
	"context"
	"errors"
	"time"

	directoriesv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/directories"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/directories/directoriesconnect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListDirectoriesResponse = directoriesv1.ListDirectoriesResponse
type GetDirectoryResponse = directoriesv1.GetDirectoryResponse
type ListDirectoryUsersResponse = directoriesv1.ListDirectoryUsersResponse
type ListDirectoryGroupsResponse = directoriesv1.ListDirectoryGroupsResponse
type ToggleDirectoryResponse = directoriesv1.ToggleDirectoryResponse
type CreateDirectoryResponse = directoriesv1.CreateDirectoryResponse

type ListDirectoryUsersOptions struct {
	PageSize         uint32
	PageToken        string
	IncludeDetail    *bool
	DirectoryGroupId *string
	UpdatedAfter     *time.Time
}

type ListDirectoryGroupsOptions struct {
	PageSize      uint32
	PageToken     string
	IncludeDetail *bool
	UpdatedAfter  *time.Time
}

type Directory interface {
	CreateDirectory(ctx context.Context, organizationId string, directory *directoriesv1.CreateDirectory) (*CreateDirectoryResponse, error)
	ListDirectories(ctx context.Context, organizationId string) (*ListDirectoriesResponse, error)
	ListDirectoryUsers(ctx context.Context, organizationId string, directoryId string, options *ListDirectoryUsersOptions) (*ListDirectoryUsersResponse, error)
	ListDirectoryGroups(ctx context.Context, organizationId string, directoryId string, options *ListDirectoryGroupsOptions) (*ListDirectoryGroupsResponse, error)
	GetPrimaryDirectoryByOrganizationId(ctx context.Context, organizationId string) (*GetDirectoryResponse, error)
	EnableDirectory(ctx context.Context, organizationId string, directoryId string) (*ToggleDirectoryResponse, error)
	DisableDirectory(ctx context.Context, organizationId string, directoryId string) (*ToggleDirectoryResponse, error)
	GetDirectory(ctx context.Context, organizationId string, directoryId string) (*GetDirectoryResponse, error)
	DeleteDirectory(ctx context.Context, organizationId string, directoryId string) error
}

type directory struct {
	coreClient *coreClient
	client     directoriesconnect.DirectoryServiceClient
}

func newDirectoryClient(coreClient *coreClient) Directory {
	return &directory{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, directoriesconnect.NewDirectoryServiceClient),
	}
}

func (d *directory) CreateDirectory(ctx context.Context, organizationId string, directory *directoriesv1.CreateDirectory) (*CreateDirectoryResponse, error) {
	return newConnectExecuter(
		d.coreClient,
		d.client.CreateDirectory,
		&directoriesv1.CreateDirectoryRequest{
			OrganizationId: organizationId,
			Directory:      directory,
		},
	).exec(ctx)
}

func (d *directory) ListDirectories(ctx context.Context, organizationId string) (*ListDirectoriesResponse, error) {
	return newConnectExecuter(
		d.coreClient,
		d.client.ListDirectories,
		&directoriesv1.ListDirectoriesRequest{
			OrganizationId: organizationId,
		},
	).exec(ctx)
}

func (d *directory) GetDirectory(ctx context.Context, organizationId string, directoryId string) (*GetDirectoryResponse, error) {
	return newConnectExecuter(
		d.coreClient,
		d.client.GetDirectory,
		&directoriesv1.GetDirectoryRequest{
			OrganizationId: organizationId,
			Id:             directoryId,
		},
	).exec(ctx)
}

func (d *directory) ListDirectoryUsers(ctx context.Context, organizationId string, directoryId string, options *ListDirectoryUsersOptions) (*ListDirectoryUsersResponse, error) {
	requestData := &directoriesv1.ListDirectoryUsersRequest{
		OrganizationId: organizationId,
		DirectoryId:    directoryId,
	}
	if options != nil {
		requestData.PageSize = options.PageSize
		requestData.PageToken = options.PageToken
		requestData.IncludeDetail = options.IncludeDetail
		requestData.DirectoryGroupId = options.DirectoryGroupId
		if options.UpdatedAfter != nil {
			requestData.UpdatedAfter = timestamppb.New(*options.UpdatedAfter)
		}
	}

	return newConnectExecuter(
		d.coreClient,
		d.client.ListDirectoryUsers,
		requestData,
	).exec(ctx)
}

func (d *directory) ListDirectoryGroups(ctx context.Context, organizationId string, directoryId string, options *ListDirectoryGroupsOptions) (*ListDirectoryGroupsResponse, error) {
	requestData := &directoriesv1.ListDirectoryGroupsRequest{
		OrganizationId: organizationId,
		DirectoryId:    directoryId,
	}
	if options != nil {
		requestData.PageSize = options.PageSize
		requestData.PageToken = options.PageToken
		requestData.IncludeDetail = options.IncludeDetail
		if options.UpdatedAfter != nil {
			requestData.UpdatedAfter = timestamppb.New(*options.UpdatedAfter)
		}
	}

	return newConnectExecuter(
		d.coreClient,
		d.client.ListDirectoryGroups,
		requestData,
	).exec(ctx)
}

func (d *directory) EnableDirectory(ctx context.Context, organizationId string, directoryId string) (*ToggleDirectoryResponse, error) {
	return newConnectExecuter(
		d.coreClient,
		d.client.EnableDirectory,
		&directoriesv1.ToggleDirectoryRequest{
			OrganizationId: organizationId,
			Id:             directoryId,
		},
	).exec(ctx)
}

func (d *directory) DisableDirectory(ctx context.Context, organizationId string, directoryId string) (*ToggleDirectoryResponse, error) {
	return newConnectExecuter(
		d.coreClient,
		d.client.DisableDirectory,
		&directoriesv1.ToggleDirectoryRequest{
			OrganizationId: organizationId,
			Id:             directoryId,
		},
	).exec(ctx)
}

func (d *directory) GetPrimaryDirectoryByOrganizationId(ctx context.Context, organizationId string) (*GetDirectoryResponse, error) {
	listDirectories, err := d.ListDirectories(ctx, organizationId)
	if err != nil {
		return nil, err
	}
	if len(listDirectories.GetDirectories()) == 0 {
		return nil, errors.New("directory does not exist for organization")
	}
	response := &GetDirectoryResponse{
		Directory: listDirectories.GetDirectories()[0],
	}
	return response, nil
}

func (d *directory) DeleteDirectory(ctx context.Context, organizationId string, directoryId string) error {
	_, err := newConnectExecuter(
		d.coreClient,
		d.client.DeleteDirectory,
		&directoriesv1.DeleteDirectoryRequest{
			OrganizationId: organizationId,
			Id:             directoryId,
		},
	).exec(ctx)
	return err
}
