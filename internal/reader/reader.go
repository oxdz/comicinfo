package reader

import (
	"context"
	"errors"
	"time"
)

var (
	ErrorNotImplemented = errors.New("not implemented")
)

var (
	_ Reader = (*ComicReader)(nil)
)

type Reader interface{}

type ComicReader interface {
	ComicInfo(ctx context.Context) (*ComicInfo, error)
	EpisodeInfo(ctx context.Context) (*EpisodeInfo, error)
	DownloadPags(ctx context.Context, url string) (<-chan *PageData, error)
}

type ComicInfo struct {
	URL    string
	Title  string
	Author string

	Episodes []*EpisodeInfo
}

type EpisodeInfo struct {
	Title string
	Date  time.Time

	CoverEnc string
	Cover    []byte
}

type PageData struct {
	Directory string
	Filename  string

	ContentEnc string
	Content    []byte
}
