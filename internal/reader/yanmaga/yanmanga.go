package yanmaga

import (
	"context"
	"crypto/rand"
	"math"
	"math/big"

	"github.com/oxdz/comicinfo/internal/reader"
)

var _ reader.ComicReader = (*Viewer)(nil)

const (
	defaultRandomStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
)

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

func randomStr(n int, s string) string {
	if s == "" {
		s = defaultRandomStr
	}

	var r string
	for i := 0; i < int(n); i++ {
		v, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
		if v != nil {
			r += string(s[v.Int64()])
		} else {
			r += string(s[0])
		}
	}
	return r
}

func randomK(cid string) string {
	var li string
	for i := 0; i < int(math.Ceil(16.0/float64(len(cid)))); i++ {
		li += cid
	}

	subH := li[:16]

	var subT string
	if n := len(li) - 16; n < 10 {
		subT = li[:16]
	} else {
		subT = li[n:16]
	}

	var v [3]int
	str := randomStr(16, defaultRandomStr)

	var ret string
	for i, c := range []byte(str) {
		v[0] ^= int(c)
		v[1] ^= int(subH[i])
		v[2] ^= int(subT[i])
		ret += string(c) + string(defaultRandomStr[(v[0]+v[1]+v[2])&63])
	}

	return ret
}
