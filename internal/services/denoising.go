package services

import (
	"log"
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

type DenoisingAlgorithm interface {
	Process(matrix *mat.Dense, nComponents int) *mat.Dense
}

type PCADenoising struct {
	maxComponents int
}

func NewPCADenoising(maxComps int) *PCADenoising {
	return &PCADenoising{maxComponents: maxComps}
}

func (p *PCADenoising) Process(matrix *mat.Dense, nComponents int) *mat.Dense {
	rows, cols := matrix.Dims()

	if rows == 0 || cols == 0 {
		log.Println("PCA: empty input matrix")
		return matrix
	}

	if nComponents <= 0 {
		nComponents = cols
	}
	if nComponents > cols {
		nComponents = cols
	}
	if nComponents > p.maxComponents {
		nComponents = p.maxComponents
	}

	centered := mat.NewDense(rows, cols, nil)
	means := make([]float64, cols)

	for j := 0; j < cols; j++ {
		col := matrix.ColView(j)
		means[j] = mat.Sum(col) / float64(rows)
		for i := 0; i < rows; i++ {
			centered.Set(i, j, matrix.At(i, j)-means[j])
		}
	}

	var cov mat.SymDense
	cov.SymOuterK(1.0/float64(rows-1), centered)

	var eig mat.EigenSym
	if ok := eig.Factorize(&cov, true); !ok {
		log.Println("PCA: failed to factorize covariance matrix")
		return matrix
	}

	var vecs mat.Dense
	eig.VectorsTo(&vecs)

	vecRows, vecCols := vecs.Dims()
	if vecRows != cols || vecCols != cols {
		log.Printf("PCA: unexpected eigenvectors dimensions: got %dx%d, expected %dx%d",
			vecRows, vecCols, cols, cols)
		return matrix
	}

	var components mat.Dense
	if nComponents > vecCols {
		nComponents = vecCols
	}
	components = *vecs.Slice(0, vecRows, 0, nComponents).(*mat.Dense)

	var projected, reconstructed mat.Dense
	projected.Mul(centered, &components)
	reconstructed.Mul(&projected, components.T())

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := reconstructed.At(i, j) + means[j]
			val = math.Max(0, math.Min(255, val))
			reconstructed.Set(i, j, math.Round(val))
		}
	}

	return &reconstructed
}

type NMFDenoising struct {
	maxIterations int
}

func NewNMFDenoising(maxIters int) *NMFDenoising {
	return &NMFDenoising{maxIterations: maxIters}
}

func (n *NMFDenoising) Process(matrix *mat.Dense, nComponents int) *mat.Dense {
	rows, cols := matrix.Dims()

	if rows == 0 || cols == 0 {
		log.Println("NMF: empty input matrix")
		return matrix
	}

	if nComponents <= 0 {
		nComponents = 1
	}
	if nComponents > cols {
		nComponents = cols
	}

	W := mat.NewDense(rows, nComponents, nil)
	H := mat.NewDense(nComponents, cols, nil)

	for i := 0; i < rows; i++ {
		for j := 0; j < nComponents; j++ {
			W.Set(i, j, 0.5+rand.Float64())
		}
	}
	for i := 0; i < nComponents; i++ {
		for j := 0; j < cols; j++ {
			H.Set(i, j, 0.5+rand.Float64())
		}
	}

	for iter := 0; iter < n.maxIterations; iter++ {
		var WH mat.Dense
		WH.Mul(W, H)

		var WT mat.Dense
		WT.CloneFrom(W.T())

		var numeratorH, denominatorH mat.Dense
		numeratorH.Mul(&WT, matrix)
		denominatorH.Mul(&WT, &WH)

		for i := 0; i < nComponents; i++ {
			for j := 0; j < cols; j++ {
				newVal := H.At(i, j) * numeratorH.At(i, j) / (denominatorH.At(i, j) + 1e-10)
				H.Set(i, j, math.Max(0.1, math.Min(1.0, newVal)))
			}
		}

		var HT mat.Dense
		HT.CloneFrom(H.T())

		var numeratorW, denominatorW mat.Dense
		numeratorW.Mul(matrix, &HT)
		denominatorW.Mul(&WH, &HT)

		for i := 0; i < rows; i++ {
			for j := 0; j < nComponents; j++ {
				newVal := W.At(i, j) * numeratorW.At(i, j) / (denominatorW.At(i, j) + 1e-10)
				W.Set(i, j, math.Max(0.1, math.Min(1.0, newVal)))
			}
		}
	}

	var reconstructed mat.Dense
	reconstructed.Mul(W, H)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := reconstructed.At(i, j) * 255
			val = math.Max(0, math.Min(255, val))
			reconstructed.Set(i, j, math.Round(val))
		}
	}

	return &reconstructed
}
