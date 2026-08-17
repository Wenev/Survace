package helper

import (
	"fmt"
	"math"

	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
)

func CalculateJaroDistance(s1, s2 string) float64 {
	if len(s1) == 0 && len(s2) == 0 {
		return 1.0
	}

	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	matchDistance := int(math.Max(float64(len(s1)), float64(len(s2)))/2.0) - 1
	if matchDistance < 0 {
		matchDistance = 0
	}

	s1Matches := make([]bool, len(s1))
	s2Matches := make([]bool, len(s2))
	matchCount := 0

	for i := 0; i < len(s1); i++ {
		start := int(math.Max(0, float64(i-matchDistance)))
		end := int(math.Min(float64(len(s2)-1), float64(i+matchDistance)))

		for j := start; j <= end; j++ {
			if !s2Matches[j] && s1[i] == s2[j] {
				s1Matches[i] = true
				s2Matches[j] = true
				matchCount++
				break
			}
		}
	}

	if matchCount == 0 {
		return 0.0
	}

	transpositions := 0
	k := 0

	for i := 0; i < len(s1); i++ {
		if s1Matches[i] {
			for !s2Matches[k] {
				k++
			}

			if s1[i] != s2[k] {
				transpositions++
			}
			k++
		}
	}

	transpositions = transpositions / 2

	m := float64(matchCount)
	jaroDistance := (m/float64(len(s1)) + m/float64(len(s2)) + (m-float64(transpositions))/m) / 3.0

	return jaroDistance
}

func FilterVideosByJaroDistance(videos []*domain.Video, query string, threshold float64) []*domain.Video {
	var result []*domain.Video

	for _, video := range videos {
		jaroDistance := CalculateJaroDistance(query, video.Title)
		fmt.Print(video.Title)
		fmt.Print("---")
		if jaroDistance >= threshold {
			result = append(result, video)
		}
	}

	return result
}
