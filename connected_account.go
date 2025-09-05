package scalekit

import (
	"context"
	"errors"

	connectedAccountsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/connected_accounts"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/connected_accounts/connected_accountsconnect"
)

type ListConnectedAccountsResponse = connectedAccountsv1.ListConnectedAccountsResponse
type CreateConnectedAccountResponse = connectedAccountsv1.CreateConnectedAccountResponse
type UpdateConnectedAccountResponse = connectedAccountsv1.UpdateConnectedAccountResponse
type DeleteConnectedAccountResponse = connectedAccountsv1.DeleteConnectedAccountResponse
type GetMagicLinkForConnectedAccountResponse = connectedAccountsv1.GetMagicLinkForConnectedAccountResponse
type GetConnectedAccountByIdentifierResponse = connectedAccountsv1.GetConnectedAccountByIdentifierResponse

type ListConnectedAccountsOptions struct {
	OrganizationId *string
	UserId         *string
	Connector      *string
	Identifier     *string
	Provider       *string
	PageSize       uint32
	PageToken      string
	Query          string
}

type CreateConnectedAccountOptions struct {
	OrganizationId   *string
	UserId           *string
	Connector        *string
	Identifier       *string
	ConnectedAccount *connectedAccountsv1.CreateConnectedAccount
}

type UpdateConnectedAccountOptions struct {
	OrganizationId   *string
	UserId           *string
	Connector        *string
	Identifier       *string
	Id               *string
	ConnectedAccount *connectedAccountsv1.UpdateConnectedAccount
}

type DeleteConnectedAccountOptions struct {
	OrganizationId *string
	UserId         *string
	Connector      *string
	Identifier     *string
	Id             *string
}

type MagicLinkOptions struct {
	OrganizationId *string
	UserId         *string
	Connector      *string
	Identifier     *string
	Id             *string
}

type GetConnectedAccountAuthOptions struct {
	OrganizationId *string
	UserId         *string
	Connector      *string
	Identifier     *string
	Id             *string
}

type ConnectedAccount interface {
	ListConnectedAccounts(ctx context.Context, options *ListConnectedAccountsOptions) (*ListConnectedAccountsResponse, error)
	CreateConnectedAccount(ctx context.Context, options *CreateConnectedAccountOptions) (*CreateConnectedAccountResponse, error)
	UpdateConnectedAccount(ctx context.Context, options *UpdateConnectedAccountOptions) (*UpdateConnectedAccountResponse, error)
	DeleteConnectedAccount(ctx context.Context, options *DeleteConnectedAccountOptions) (*DeleteConnectedAccountResponse, error)
	GetMagicLinkForConnectedAccount(ctx context.Context, options *MagicLinkOptions) (*GetMagicLinkForConnectedAccountResponse, error)
	GetConnectedAccountAuth(ctx context.Context, options *GetConnectedAccountAuthOptions) (*GetConnectedAccountByIdentifierResponse, error)
}

type connectedAccount struct {
	coreClient *coreClient
	client     connected_accountsconnect.ConnectedAccountServiceClient
}

func newConnectedAccountClient(coreClient *coreClient) ConnectedAccount {
	return &connectedAccount{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, connected_accountsconnect.NewConnectedAccountServiceClient),
	}
}

func (ca *connectedAccount) ListConnectedAccounts(ctx context.Context, options *ListConnectedAccountsOptions) (*ListConnectedAccountsResponse, error) {
	requestData := &connectedAccountsv1.ListConnectedAccountsRequest{}

	if options != nil {
		if options.OrganizationId != nil {
			requestData.OrganizationId = options.OrganizationId
		}
		if options.UserId != nil {
			requestData.UserId = options.UserId
		}
		if options.Connector != nil {
			requestData.Connector = options.Connector
		}
		if options.Identifier != nil {
			requestData.Identifier = options.Identifier
		}
		if options.Provider != nil {
			requestData.Provider = *options.Provider
		}
		requestData.PageSize = options.PageSize
		requestData.PageToken = options.PageToken
		requestData.Query = options.Query
	}

	return newConnectExecuter(
		ca.coreClient,
		ca.client.ListConnectedAccounts,
		requestData,
	).exec(ctx)
}

func (ca *connectedAccount) CreateConnectedAccount(ctx context.Context, options *CreateConnectedAccountOptions) (*CreateConnectedAccountResponse, error) {
	if options == nil || options.ConnectedAccount == nil {
		return nil, errors.New("connected account data is required")
	}

	requestData := &connectedAccountsv1.CreateConnectedAccountRequest{
		ConnectedAccount: options.ConnectedAccount,
	}

	if options.OrganizationId != nil {
		requestData.OrganizationId = options.OrganizationId
	}
	if options.UserId != nil {
		requestData.UserId = options.UserId
	}
	if options.Connector != nil {
		requestData.Connector = options.Connector
	}
	if options.Identifier != nil {
		requestData.Identifier = options.Identifier
	}

	return newConnectExecuter(
		ca.coreClient,
		ca.client.CreateConnectedAccount,
		requestData,
	).exec(ctx)
}

func (ca *connectedAccount) UpdateConnectedAccount(ctx context.Context, options *UpdateConnectedAccountOptions) (*UpdateConnectedAccountResponse, error) {
	if options == nil || options.ConnectedAccount == nil {
		return nil, errors.New("connected account data is required")
	}

	requestData := &connectedAccountsv1.UpdateConnectedAccountRequest{
		ConnectedAccount: options.ConnectedAccount,
	}

	if options.OrganizationId != nil {
		requestData.OrganizationId = options.OrganizationId
	}
	if options.UserId != nil {
		requestData.UserId = options.UserId
	}
	if options.Connector != nil {
		requestData.Connector = options.Connector
	}
	if options.Identifier != nil {
		requestData.Identifier = options.Identifier
	}
	if options.Id != nil {
		requestData.Id = options.Id
	}

	return newConnectExecuter(
		ca.coreClient,
		ca.client.UpdateConnectedAccount,
		requestData,
	).exec(ctx)
}

func (ca *connectedAccount) DeleteConnectedAccount(ctx context.Context, options *DeleteConnectedAccountOptions) (*DeleteConnectedAccountResponse, error) {
	requestData := &connectedAccountsv1.DeleteConnectedAccountRequest{}

	if options != nil {
		if options.OrganizationId != nil {
			requestData.OrganizationId = options.OrganizationId
		}
		if options.UserId != nil {
			requestData.UserId = options.UserId
		}
		if options.Connector != nil {
			requestData.Connector = options.Connector
		}
		if options.Identifier != nil {
			requestData.Identifier = options.Identifier
		}
		if options.Id != nil {
			requestData.Id = options.Id
		}
	}

	return newConnectExecuter(
		ca.coreClient,
		ca.client.DeleteConnectedAccount,
		requestData,
	).exec(ctx)
}

func (ca *connectedAccount) GetMagicLinkForConnectedAccount(ctx context.Context, options *MagicLinkOptions) (*GetMagicLinkForConnectedAccountResponse, error) {
	requestData := &connectedAccountsv1.GetMagicLinkForConnectedAccountRequest{}

	if options != nil {
		if options.OrganizationId != nil {
			requestData.OrganizationId = options.OrganizationId
		}
		if options.UserId != nil {
			requestData.UserId = options.UserId
		}
		if options.Connector != nil {
			requestData.Connector = options.Connector
		}
		if options.Identifier != nil {
			requestData.Identifier = options.Identifier
		}
		if options.Id != nil {
			requestData.Id = options.Id
		}
	}

	return newConnectExecuter(
		ca.coreClient,
		ca.client.GetMagicLinkForConnectedAccount,
		requestData,
	).exec(ctx)
}

func (ca *connectedAccount) GetConnectedAccountAuth(ctx context.Context, options *GetConnectedAccountAuthOptions) (*GetConnectedAccountByIdentifierResponse, error) {
	requestData := &connectedAccountsv1.GetConnectedAccountByIdentifierRequest{}

	if options != nil {
		if options.OrganizationId != nil {
			requestData.OrganizationId = options.OrganizationId
		}
		if options.UserId != nil {
			requestData.UserId = options.UserId
		}
		if options.Connector != nil {
			requestData.Connector = options.Connector
		}
		if options.Identifier != nil {
			requestData.Identifier = options.Identifier
		}
		if options.Id != nil {
			requestData.Id = options.Id
		}
	}

	return newConnectExecuter(
		ca.coreClient,
		ca.client.GetConnectedAccountAuth,
		requestData,
	).exec(ctx)
}
