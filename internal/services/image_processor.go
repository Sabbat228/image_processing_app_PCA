package services

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"image"
	"image-processing-app/internal/utils"
)

type ImageProcessor struct {
	pcaAlgo *PCADenoising
	nmfAlgo *NMFDenoising
}

func NewImageProcessor(pcaMaxComps, nmfMaxIters int) *ImageProcessor {
	return &ImageProcessor{
		pcaAlgo: NewPCADenoising(pcaMaxComps),
		nmfAlgo: NewNMFDenoising(nmfMaxIters),
	}
}

func (ip *ImageProcessor) ProcessImage(method string, img image.Image, nFactors int) (image.Image, error) {
	if img == nil {
		return nil, fmt.Errorf("input image is nil")
	}
	matrix := utils.ConvertImageToMatrix(img)

	var result *mat.Dense
	switch method {
	case "pca":
		result = ip.pcaAlgo.Process(matrix, nFactors)
	case "nmf":
		result = ip.nmfAlgo.Process(matrix, nFactors)
	default:
		return nil, utils.ErrInvalidMethod
	}

	return utils.ConvertMatrixToImage(result), nil
}

func (ip *ImageProcessor) ApplyPCA(matrix *mat.Dense, nComponents int) *mat.Dense {
	return ip.pcaAlgo.Process(matrix, nComponents)
}

func (ip *ImageProcessor) ApplyNMF(matrix *mat.Dense, nComponents int) *mat.Dense {
	return ip.nmfAlgo.Process(matrix, nComponents)
}
