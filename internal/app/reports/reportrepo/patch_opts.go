package reportrepo

import "github.com/radiologist-ai/web-app/internal/domain"

func WithReportText(text string) domain.PatchOpt {
	return func(p *domain.PatchConf) {
		p.ReportText = &text
	}
}

func WithApproved(approved bool) domain.PatchOpt {
	return func(p *domain.PatchConf) {
		p.Approved = &approved
	}
}
