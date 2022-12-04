package resizeproxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

var ErrNotAcceptableResponseStatus = fmt.Errorf("not acceptable server response status")

type ImageClient struct{}

func NewImageClient() ImageGetter {
	return &ImageClient{}
}

func (ic ImageClient) GetImage(ctx context.Context, imgRequest ImageRequest) (*ImageResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imgRequest.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = imgRequest.Headers

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrNotAcceptableResponseStatus
	}

	responseContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &ImageResponse{responseContent, resp.Header}, nil
}
