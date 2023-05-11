package helpers

import (
	"context"

	"github.com/9elements/converged-security-suite/v2/pkg/uefi"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/logger"
	fianoUEFI "github.com/linuxboot/fiano/pkg/uefi"
)

// GetAlignedImage returns a part of fullImage which is aligned
// (and ready to be compared) with partialImage.
//
// Some images might contain only the BIOS region instead of the whole image,
// and to properly compare two BIOS regions we return the same regios of
// originalFirmware as we have in receivedImage.
func GetAlignedImage(
	ctx context.Context,
	fullImage *uefi.UEFI,
	partialImage []byte,
) (firmware *uefi.UEFI, offset uint64, err error) {
	if fullImage == nil {
		return nil, 0, ErrNoOrigImageToCompareWith{}
	}
	log := logger.FromCtx(ctx).WithFields(
		field.Map[any]{
			"fullSize":    len(fullImage.Buf()),
			"partialSize": len(partialImage),
		},
	)

	defer func() {
		log.WithField("offset", offset).WithField("err", err).Debugf("")
	}()

	switch {
	case len(fullImage.Buf()) == len(partialImage):
		return fullImage, 0, nil
	case len(fullImage.Buf()) < len(partialImage):
		return nil, 0, ErrImageLengthDoesNotMatch{ExpectedLength: uint(len(fullImage.Buf())), ReceivedLength: uint(len(partialImage))}
	}

	// There are different ways to dump a firmware image. And sometimes we
	// we get only the BIOS region instead of the whole UEFI image.
	// Here we adjust to the case when the original image is the whole
	// UEFI image, while the received image is only a BIOS region.
	//
	// Since we have only a part of the image and we want to make it work,
	// we reduce the original image to BIOS region as well, to make offsets
	// and sizes the same as in the received image.
	log.Debugf("original image is larger than received image (%d > %d), assuming the received image is only the BIOS region",
		len(fullImage.Buf()), len(partialImage))
	biosRegions, err := fullImage.GetByRegionType(fianoUEFI.RegionTypeBIOS)
	if err != nil {
		return nil, 0, ErrUnableToFindBIOSRegion{Err: err}
	}
	if len(biosRegions) != 1 {
		return nil, 0, ErrUnexpectedAmountOfBIOSRegions{FoundCount: uint(len(biosRegions))}
	}
	biosRegion := biosRegions[0]

	if len(biosRegion.Buf()) != len(partialImage) {
		// If we got here, something is wrong. We expect the sizes be the same
		// due to line "fullImage = biosRegions[0]" above.
		return nil, 0, ErrImageLengthDoesNotMatch{ExpectedLength: uint(len(biosRegion.Buf())), ReceivedLength: uint(len(partialImage))}
	}

	// Here we force to compare a BIOS region with a BIOS region (instead
	// of the whole image):
	return &uefi.UEFI{Node: *biosRegion}, biosRegion.Offset, nil
}
