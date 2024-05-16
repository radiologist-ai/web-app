package rgen

import (
	"context"
	"errors"
	"github.com/radiologist-ai/web-app/proto"
	"google.golang.org/grpc"
	"sync"
)

type Client struct {
	cli proto.RGenClient
	mu  sync.Mutex
}

func NewClient(cc grpc.ClientConnInterface) (*Client, error) {
	if cc == nil {
		return nil, errors.New("grpc conn can't be nil")
	}
	cli := proto.NewRGenClient(cc)
	return &Client{cli: cli}, nil
}

func (c *Client) GenerateReport(ctx context.Context, link2photo string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	req := &proto.Request{
		PatientId:  "fake-id",
		LinkToXray: link2photo,
	}
	res, err := c.cli.GenerateReport(ctx, req)
	if err != nil {
		return "", err
	}
	return res.Report, nil
}
