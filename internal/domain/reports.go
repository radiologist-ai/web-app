package domain

import "context"

type RGen interface {
	GenerateReport(ctx context.Context, link2photo string) (string, error)
}
