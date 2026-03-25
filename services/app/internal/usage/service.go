package usage

import "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"

type UsageService struct {
	filesService *files.FilesService
}

type CalculateForLocationResult struct {
	Location string
	Size     uint64
}

func NewService(filesService *files.FilesService) UsageService {
	return UsageService{
		filesService: filesService,
	}
}

func (s *UsageService) CalculateForLocation(location string) (CalculateForLocationResult, error) {
	files, err := s.filesService.FindByLocation(location)

	if err != nil {
		return CalculateForLocationResult{}, err
	}

	size := uint64(0)

	for _, file := range files {
		size += file.Size
	}

	return CalculateForLocationResult{
		Location: location,
		Size:     size,
	}, nil
}
