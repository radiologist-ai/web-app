package rgen

import (
	"context"
	"errors"
	"fmt"
	"github.com/radiologist-ai/web-app/proto"
	"google.golang.org/grpc"
	"sync"
	"time"
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

func (c *Client) GenerateReportAsync(ctx context.Context, link2photo string, ch chan string, errCh chan error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	doneCh := make(chan struct{})
	go func() {
		defer func() { doneCh <- struct{}{} }()

		req := &proto.Request{
			PatientId:  "fake-id",
			LinkToXray: link2photo,
		}
		res, err := c.cli.GenerateReport(ctx, req)
		if err != nil {
			errCh <- err
			return
		}
		ch <- res.Report
		return
	}()
	select {
	case <-ctx.Done():
		return
	case <-doneCh:
		return
	case <-time.After(time.Second * 30):
		errCh <- fmt.Errorf("GenerateReportAsync timeout")
		return
	}
}
