// Copyright 2023 Meta Platforms, Inc. and affiliates.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package firmwarewand

import (
	"context"
	"reflect"

	"github.com/immune-gmbh/AttestationFailureAnalysisService/if/generated/afas"

	"github.com/facebookincubator/go-belt/tool/logger"
)

// Search asks the firmware analysis service to provide information on already
// collected images, which satisfies selected filters.
func (fwwand *FirmwareWand) Search(
	ctx context.Context,
	filters afas.SearchFirmwareFilters,
	shouldFetchContent bool,
) (*afas.SearchFirmwareResult_, error) {
	l := logger.FromCtx(ctx)

	if reflect.ValueOf(filters).IsZero() {
		return nil, ErrInvalidInput{Desc: "filters cannot be completely empty"}
	}

	request := afas.SearchFirmwareRequest{
		FetchContent: shouldFetchContent,
	}
	request.OrFilters = append(request.OrFilters, &filters)

	l.Debugf("sending the request to firmware analyzer service...")
	result, err := fwwand.afasClient.SearchFirmware(ctx, &request)
	l.Debugf("received a response from the firmware analyzer service; err == %v, result == %+v", err, result)

	return result, err
}
