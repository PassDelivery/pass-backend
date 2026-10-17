package core_files

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

const MaxPhotoSize = 5 << 20

const (
	PhotoKindRestaurants = "restaurants"
	PhotoKindDishes      = "dishes"
)

var photoExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

var photoName = regexp.MustCompile(`^[0-9a-f-]{36}\.(jpg|png|webp)$`)

func uploadsDir() string {
	if dir := os.Getenv("UPLOADS_PATH"); dir != "" {
		return dir
	}
	return "out/uploads"
}

func isPhotoKind(kind string) bool {
	return kind == PhotoKindRestaurants || kind == PhotoKindDishes
}

func SavePhoto(kind string, data []byte) (string, error) {
	if !isPhotoKind(kind) {
		return "", fmt.Errorf("unknown photo kind %s: %w", kind, core_errors.ErrInvalidArgument)
	}
	if len(data) == 0 || len(data) > MaxPhotoSize {
		return "", fmt.Errorf("photo size must be 1 byte..5 MB: %w", core_errors.ErrInvalidArgument)
	}

	ext, ok := photoExtensions[http.DetectContentType(data)]
	if !ok {
		return "", fmt.Errorf("photo must be jpeg, png or webp: %w", core_errors.ErrInvalidArgument)
	}

	folder := filepath.Join(uploadsDir(), kind)
	if err := os.MkdirAll(folder, 0o755); err != nil {
		return "", fmt.Errorf("create uploads folder: %w", err)
	}

	name := uuid.NewString() + ext
	if err := os.WriteFile(filepath.Join(folder, name), data, 0o644); err != nil {
		return "", fmt.Errorf("write photo: %w", err)
	}

	return "/uploads/" + kind + "/" + name, nil
}

func PhotoPath(kind string, name string) (string, error) {
	if !isPhotoKind(kind) || !photoName.MatchString(name) {
		return "", fmt.Errorf("photo not found: %w", core_errors.ErrNotFound)
	}

	return filepath.Join(uploadsDir(), kind, name), nil
}

func RemovePhoto(url *string) {
	if url == nil {
		return
	}

	parts := strings.Split(strings.TrimPrefix(*url, "/uploads/"), "/")
	if len(parts) != 2 {
		return
	}

	path, err := PhotoPath(parts[0], parts[1])
	if err != nil {
		return
	}

	_ = os.Remove(path)
}
