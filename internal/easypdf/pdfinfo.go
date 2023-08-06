package easypdf

import (
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFInfo struct {
	PageCount int
}

func ExtractPDFInfo(pdfFile *os.File) (*PDFInfo, error) {
	ctx, err := api.ReadContext(pdfFile, model.NewDefaultConfiguration())
	if err != nil {
		return nil, err
	}
	if err := api.ValidateContext(ctx); err != nil {
		return nil, err
	}
	return &PDFInfo{PageCount: ctx.PageCount}, nil
}
