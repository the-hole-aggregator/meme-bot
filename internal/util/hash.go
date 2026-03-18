package util

import (
	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"
)

// ComputePHash вычисляет перцептивный хэш изображения
func ComputePHash(filePath string) (*goimagehash.ImageHash, error) {
	img, err := imaging.Open(filePath)
	if err != nil {
		return nil, err
	}

	phash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, err
	}

	return phash, nil
}

// ComparePHash возвращает true, если изображения похожи (threshold = 5 бит)
func ComparePHash(h1, h2 *goimagehash.ImageHash) (bool, error) {
	distance, err := h1.Distance(h2)
	if err != nil {
		return false, err
	}

	const threshold = 5
	return distance <= threshold, nil
}
