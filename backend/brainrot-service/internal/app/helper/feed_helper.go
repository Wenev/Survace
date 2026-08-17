package helper

import (
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"math/rand"
	"time"
)

func WeightedShuffle(a, b []*domain.Video, weightA float64, limit int) []*domain.Video {
	rand.Seed(time.Now().UnixNano())
	result := make([]*domain.Video, 0, limit)
	ai, bi := 0, 0
	for len(result) < limit && (ai < len(a) || bi < len(b)) {
		pickA := rand.Float64() < weightA
		if pickA && ai < len(a) {
			result = append(result, a[ai])
			ai++
		} else if bi < len(b) {
			result = append(result, b[bi])
			bi++
		} else if ai < len(a) {
			result = append(result, a[ai])
			ai++
		}
	}
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
	return result
}

func FilterUnwatched(videos []*domain.Video, watchedIDs map[int32]struct{}) []*domain.Video {
	var out []*domain.Video
	for _, v := range videos {
		if _, watched := watchedIDs[v.ID]; !watched {
			out = append(out, v)
		}
	}
	return out
}
