package controller

import (
	"context"

	"github.com/immune-gmbh/AttestationFailureAnalysisService/pkg/analysis"
	"github.com/immune-gmbh/AttestationFailureAnalysisService/pkg/types"
)

// FirmwareAccessor implements analysis.Blob given
type FirmwareAccessor struct {
	// == Non-serializable part ==

	downloadedContent []byte

	// == Serializable part ==

	ImageID types.ImageID
}

var _ analysis.Blob = (*FirmwareAccessor)(nil)

// Bytes implements analysis.Blob
func (fw *FirmwareAccessor) Bytes() []byte {
	if fw.downloadedContent == nil {
		panic("FirmwareAccessor is not initialized")
	}
	return fw.downloadedContent
}

// Init initializes the FirmwareAccessor after it was deserialized.
func (fw *FirmwareAccessor) Init(ctx context.Context, storage Storage) error {
	content, err := storage.GetFirmwareBytes(ctx, fw.ImageID)
	fw.downloadedContent = content
	return err
}
