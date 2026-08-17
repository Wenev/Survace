package helper

import "github.com/Wenev/Survace/brainrot-service/internal/app/domain"

func DeduplicateVideos(videos []*domain.Video) []*domain.Video {
	seen := make(map[int32]struct{})
	result := make([]*domain.Video, 0, len(videos))
	for _, v := range videos {
		if _, ok := seen[v.ID]; !ok {
			seen[v.ID] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
