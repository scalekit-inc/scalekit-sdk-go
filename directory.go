package scalekit

import (
	"context"
	"time"

	directoriesv1 "github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/directories"
	"github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/directories/directoriesconnect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListDirectoriesResponse = directoriesv1.ListDirectoriesResponse
type ListDirectoryUsersResponse = directoriesv1.ListDirectoryUsersResponse
type ListDirectoryGroupsResponse = directoriesv1.ListDirectoryGroupsResponse

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
	ListDirectories(ctx context.Context, organizationId string) (*ListDirectoriesResponse, error)
	ListDirectoryUsers(ctx context.Context, organizationId string, directoryId string, options *ListDirectoryUsersOptions) (*ListDirectoryUsersResponse, error)
	ListDirectoryGroups(ctx context.Context, organizationId string, directoryId string, options *ListDirectoryGroupsOptions) (*ListDirectoryGroupsResponse, error)
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

func (d *directory) ListDirectories(ctx context.Context, organizationId string) (*ListDirectoriesResponse, error) {
	return newConnectExecuter(
		d.coreClient,
		d.client.ListDirectory,
		&directoriesv1.ListDirectoriesRequest{
			OrganizationId: organizationId,
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