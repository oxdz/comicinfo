package yanmaga

import (
	"context"

	"github.com/oxdz/comicinfo/internal/reader"
)

var _ reader.ComicReader = (*Viewer)(nil)

type Viewer struct {
}

func (v *Viewer) ComicInfo(ctx context.Context) (*reader.ComicInfo, error) {
	return nil, reader.ErrorNotImplemented
}

func (v *Viewer) EpisodeInfo(ctx context.Context) (*reader.EpisodeInfo, error) {
	return nil, reader.ErrorNotImplemented
}

func (v *Viewer) DownloadPags(ctx context.Context, url string) (<-chan *reader.PageData, error) {
	return nil, nil
}
