package util

import (
	"strings"

	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"
)

type Hasher interface {
	ComputePHash(filePath string) (string, error)
	ComparePHash(h1, h2 string) (bool, error)
}

type ImagePHasher struct {
	threshold int
}

func NewImagePHasher() Hasher {
	return ImagePHasher{
		threshold: 5,
	}
}

// ComputePHash computes the perceptual hash of an image
func (i ImagePHasher) ComputePHash(filePath string) (string, error) {
	img, err := imaging.Open(filePath)
	if err != nil {
		return "", err
	}

	phash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return "", err
	}

	return phash.ToString(), nil
}

// ComparePHash returns true if the images are similar (threshold = 5 bits)
func (i ImagePHasher) ComparePHash(h1, h2 string) (bool, error) {
	hash1, err := goimagehash.LoadImageHash(strings.NewReader(h1))
	if err != nil {
		return false, err
	}

	hash2, err := goimagehash.LoadImageHash(strings.NewReader(h2))
	if err != nil {
		return false, err
	}

	distance, err := hash1.Distance(hash2)
	if err != nil {
		return false, err
	}

	return distance <= i.threshold, nil
}
