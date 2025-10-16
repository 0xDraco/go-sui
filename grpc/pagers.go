package grpc

import (
	"context"
	"errors"

	v2 "github.com/0xdraco/go-sui/proto/sui/rpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type OwnedObjectsPager struct {
	iter *pageIterator[*v2.Object]
}

func (c *GRPCClient) OwnedObjectsPager(req *v2.ListOwnedObjectsRequest, opts ...grpc.CallOption) (*OwnedObjectsPager, error) {
	if c == nil {
		return nil, errors.New("nil client")
	}
	if req == nil {
		return nil, errors.New("nil request")
	}

	base := proto.Clone(req).(*v2.ListOwnedObjectsRequest)
	fetch := func(ctx context.Context, token []byte) ([]*v2.Object, []byte, error) {
		if len(token) == 0 {
			base.PageToken = nil
		} else {
			base.PageToken = cloneBytes(token)
		}
		resp, err := c.StateClient().ListOwnedObjects(ctx, base, opts...)
		if err != nil {
			return nil, nil, err
		}
		next := cloneBytes(resp.GetNextPageToken())
		base.PageToken = next
		return resp.GetObjects(), next, nil
	}

	iter, err := newPageIterator(base.GetPageToken(), fetch)
	if err != nil {
		return nil, err
	}

	return &OwnedObjectsPager{iter: iter}, nil
}

func (p *OwnedObjectsPager) Next(ctx context.Context) ([]*v2.Object, error) {
	if p == nil {
		return nil, errors.New("nil pager")
	}
	return p.iter.Next(ctx)
}

func (p *OwnedObjectsPager) ForEach(ctx context.Context, fn func(*v2.Object) error) error {
	if p == nil {
		return errors.New("nil pager")
	}
	return p.iter.ForEach(ctx, fn)
}

func (p *OwnedObjectsPager) Collect(ctx context.Context) ([]*v2.Object, error) {
	if p == nil {
		return nil, errors.New("nil pager")
	}
	return p.iter.Collect(ctx)
}

type DynamicFieldsPager struct {
	iter *pageIterator[*v2.DynamicField]
}

func (c *GRPCClient) DynamicFieldsPager(req *v2.ListDynamicFieldsRequest, opts ...grpc.CallOption) (*DynamicFieldsPager, error) {
	if c == nil {
		return nil, errors.New("nil client")
	}
	if req == nil {
		return nil, errors.New("nil request")
	}

	base := proto.Clone(req).(*v2.ListDynamicFieldsRequest)
	fetch := func(ctx context.Context, token []byte) ([]*v2.DynamicField, []byte, error) {
		if len(token) == 0 {
			base.PageToken = nil
		} else {
			base.PageToken = cloneBytes(token)
		}
		resp, err := c.StateClient().ListDynamicFields(ctx, base, opts...)
		if err != nil {
			return nil, nil, err
		}

		next := cloneBytes(resp.GetNextPageToken())
		base.PageToken = next
		return resp.GetDynamicFields(), next, nil
	}

	iter, err := newPageIterator(base.GetPageToken(), fetch)
	if err != nil {
		return nil, err
	}

	return &DynamicFieldsPager{iter: iter}, nil
}

func (p *DynamicFieldsPager) Next(ctx context.Context) ([]*v2.DynamicField, error) {
	if p == nil {
		return nil, errors.New("nil pager")
	}
	return p.iter.Next(ctx)
}

func (p *DynamicFieldsPager) ForEach(ctx context.Context, fn func(*v2.DynamicField) error) error {
	if p == nil {
		return errors.New("nil pager")
	}
	return p.iter.ForEach(ctx, fn)
}

func (p *DynamicFieldsPager) Collect(ctx context.Context) ([]*v2.DynamicField, error) {
	if p == nil {
		return nil, errors.New("nil pager")
	}
	return p.iter.Collect(ctx)
}

type BalancesPager struct {
	iter *pageIterator[*v2.Balance]
}

func (c *GRPCClient) BalancesPager(req *v2.ListBalancesRequest, opts ...grpc.CallOption) (*BalancesPager, error) {
	if c == nil {
		return nil, errors.New("nil client")
	}
	if req == nil {
		return nil, errors.New("nil request")
	}

	base := proto.Clone(req).(*v2.ListBalancesRequest)
	fetch := func(ctx context.Context, token []byte) ([]*v2.Balance, []byte, error) {
		if len(token) == 0 {
			base.PageToken = nil
		} else {
			base.PageToken = cloneBytes(token)
		}
		resp, err := c.StateClient().ListBalances(ctx, base, opts...)
		if err != nil {
			return nil, nil, err
		}
		next := cloneBytes(resp.GetNextPageToken())
		base.PageToken = next
		return resp.GetBalances(), next, nil
	}

	iter, err := newPageIterator(base.GetPageToken(), fetch)
	if err != nil {
		return nil, err
	}

	return &BalancesPager{iter: iter}, nil
}

func (p *BalancesPager) Next(ctx context.Context) ([]*v2.Balance, error) {
	if p == nil {
		return nil, errors.New("nil pager")
	}
	return p.iter.Next(ctx)
}

func (p *BalancesPager) ForEach(ctx context.Context, fn func(*v2.Balance) error) error {
	if p == nil {
		return errors.New("nil pager")
	}
	return p.iter.ForEach(ctx, fn)
}

func (p *BalancesPager) Collect(ctx context.Context) ([]*v2.Balance, error) {
	if p == nil {
		return nil, errors.New("nil pager")
	}
	return p.iter.Collect(ctx)
}

type PackageVersionsPager struct {
	iter *pageIterator[*v2.PackageVersion]
}

func (c *GRPCClient) PackageVersionsPager(req *v2.ListPackageVersionsRequest, opts ...grpc.CallOption) (*PackageVersionsPager, error) {
	if c == nil {
		return nil, errors.New("nil client")
	}
	if req == nil {
		return nil, errors.New("nil request")
	}

	base := proto.Clone(req).(*v2.ListPackageVersionsRequest)
	fetch := func(ctx context.Context, token []byte) ([]*v2.PackageVersion, []byte, error) {
		if len(token) == 0 {
			base.PageToken = nil
		} else {
			base.PageToken = cloneBytes(token)
		}
		resp, err := c.MovePackageClient().ListPackageVersions(ctx, base, opts...)
		if err != nil {
			return nil, nil, err
		}
		next := cloneBytes(resp.GetNextPageToken())
		base.PageToken = next
		return resp.GetVersions(), next, nil
	}

	iter, err := newPageIterator(base.GetPageToken(), fetch)
	if err != nil {
		return nil, err
	}

	return &PackageVersionsPager{iter: iter}, nil
}

func (p *PackageVersionsPager) Next(ctx context.Context) ([]*v2.PackageVersion, error) {
	if p == nil {
		return nil, errors.New("nil pager")
	}
	return p.iter.Next(ctx)
}

func (p *PackageVersionsPager) ForEach(ctx context.Context, fn func(*v2.PackageVersion) error) error {
	if p == nil {
		return errors.New("nil pager")
	}
	return p.iter.ForEach(ctx, fn)
}

func (p *PackageVersionsPager) Collect(ctx context.Context) ([]*v2.PackageVersion, error) {
	if p == nil {
		return nil, errors.New("nil pager")
	}
	return p.iter.Collect(ctx)
}
