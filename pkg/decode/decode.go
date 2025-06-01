package decode

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
)

type EpisodeInfo struct {
	Status            string         `json:"status"`
	TitleID           int            `json:"title_id"`
	EpisodeID         int            `json:"episode_id"`
	PageStartPosition int            `json:"page_start_position"`
	ScrambleSeed      int            `json:"scramble_seed"`
	PageList          []string       `json:"page_list"`
	PreviousEpisode   map[string]int `json:"previous_episode"`
	NextEpisode       map[string]int `json:"next_episode"`
}

func MangeHash(episode_id int, platform int) string {
	e := strconv.FormatInt(int64(episode_id), 10)
	p := strconv.FormatInt(int64(platform), 10)
	v1 := Sum256("platform")
	v2 := Sum512(p)
	v3 := Sum256("episode_id")
	v4 := Sum512(e)

	v6 := Sum256("")
	v7 := Sum512("")

	v5 := Sum256(
		fmt.Sprintf("%s_%s,%s_%s", v3, v4, v1, v2),
	)
	return Sum512(v5 + v6 + "_" + v7)
}

func Sum256(v string) string {
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:])
}

func Sum512(v string) string {
	h := sha512.Sum512([]byte(v))
	return hex.EncodeToString(h[:])
}

type WordArray struct {
	Words    []uint32
	SigBytes int
}

// Clamp 方法实现
func (wa *WordArray) Clamp() {
	p := wa.SigBytes
	if len(wa.Words) > 0 && p > 0 {
		// 清除超出有效字节的位
		wa.Words[p>>2] &= 0xFFFFFFFF << (32 - (p%4)*8)
	}
	// 调整words数组长度为ceil(p/4)
	newLength := int(math.Ceil(float64(p) / 4))
	if newLength < len(wa.Words) {
		wa.Words = wa.Words[:newLength]
	}
}

// Concat 方法实现
func (wa *WordArray) Concat(g *WordArray) *WordArray {
	p := wa.Words
	wordsG := g.Words
	w := wa.SigBytes
	P := g.SigBytes

	// 先执行clamp操作
	wa.Clamp()

	if w%4 != 0 {
		// 情况1：当前字节数不是4的倍数，按字节拼接
		for C := 0; C < P; C++ {
			// 从g中提取一个字节
			S := (wordsG[C>>2] >> (24 - (C%4)*8)) & 255
			// 将该字节放入p的适当位置
			index := (w + C) >> 2
			shift := 24 - ((w+C)%4)*8
			if index >= len(p) {
				// 需要扩展p数组
				newLen := index + 1
				for len(p) < newLen {
					p = append(p, 0)
				}
				wa.Words = p
			}
			p[index] |= S << shift
		}
	} else {
		// 情况2：当前字节数是4的倍数，按字（32位）拼接
		for x := 0; x < P; x += 4 {
			index := (w + x) >> 2
			srcIndex := x >> 2
			if srcIndex < len(wordsG) {
				if index >= len(p) {
					// 需要扩展p数组
					newLen := index + 1
					for len(p) < newLen {
						p = append(p, 0)
					}
					wa.Words = p
				}
				p[index] = wordsG[srcIndex]
			}
		}
	}

	// 更新有效字节数
	wa.SigBytes += P
	return wa
}

func (wa *WordArray) Init(g string) {
	p := len(g)
	arr := make([]uint32, (p+3)/4) // 等价于JavaScript中的 _[w >>> 2]

	for w := 0; w < p; w++ {
		index := w >> 2 // 等价于JavaScript中的 w >>> 2
		shift := 24 - (w%4)*8
		charCode := uint32(g[w]) & 255 // 等价于JavaScript中的 g.charCodeAt(w) & 255
		arr[index] |= charCode << shift
	}
	wa.Words = arr
	wa.SigBytes = p
}

func ParseW(g string) WordArray {
	p := len(g)
	arr := make([]uint32, (p+3)/4) // 等价于JavaScript中的 _[w >>> 2]

	for w := 0; w < p; w++ {
		index := w >> 2 // 等价于JavaScript中的 w >>> 2
		shift := 24 - (w%4)*8
		charCode := uint32(g[w]) & 255 // 等价于JavaScript中的 g.charCodeAt(w) & 255
		arr[index] |= charCode << shift
	}

	return WordArray{
		arr, p,
	}
}

func (wa *WordArray) PreP(nbytes int) {
	m := wa
	h := m.Words
	b := nbytes * 8
	v := m.SigBytes * 8

	// 确保h数组足够大，能够容纳v>>>5位置的元素
	requiredIndex := v >> 5
	if requiredIndex >= len(h) {
		newLen := requiredIndex + 1
		for i := len(h); i < newLen; i++ {
			h = append(h, 0)
		}
		m.Words = h // 更新_words
	}

	// h[v >>> 5] |= 128 << 24 - v % 32
	h[v>>5] |= 128 << (24 - v%32)

	// 计算需要的数组长度
	index14 := ((v + 64) >> 9 << 4) + 14
	index15 := ((v + 64) >> 9 << 4) + 15
	maxIndex := index14
	if index15 > maxIndex {
		maxIndex = index15
	}

	// 确保h数组足够大
	if maxIndex >= len(h) {
		newLen := maxIndex + 1
		for i := len(h); i < newLen; i++ {
			h = append(h, 0)
		}
		m.Words = h // 更新_words
	}

	// h[(v + 64 >>> 9 << 4) + 14] = n.floor(b / 4294967296)
	h[index14] = uint32(math.Floor(float64(b) / 4294967296))

	// h[(v + 64 >>> 9 << 4) + 15] = b
	h[index15] = uint32(b)
}
