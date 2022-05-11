package treegen

import (
	"math"
)

// This specifically implements 3D trees because that's all I need

type KDTree struct {
	points   []Point
	treeSize int
}

func MakeKDTree() KDTree {
	return KDTree{
		make([]Point, 0),
		0,
	}
}

// Adds a node outside of the tree
func (kdtree *KDTree) Add(p Point) {
	kdtree.points = append(kdtree.points, p)
}

func (kdtree *KDTree) NearestNeighbor(target Point) Point {
	var closest Point = nil
	var closestDist float64 = math.Inf(1)
	nearestHelper(kdtree.points, 0, 0, kdtree.treeSize-1, target, &closest, &closestDist)
	for i := kdtree.treeSize; i < len(kdtree.points); i++ {
		p := kdtree.points[i]
		pDist := PointDistSq(p, target)
		if pDist < closestDist {
			closest = p
			closestDist = pDist
		}
	}
	return closest
}

func nearestHelper(
	arr []Point, dim int, left, right int,
	target Point,
	closest *Point, closestDist *float64,
) {
	if right < left {
		return
	}

	middleIdx := (left + right) / 2
	currentPoint := arr[middleIdx]
	currentDist := PointDistSq(currentPoint, target)
	if currentDist < *closestDist {
		*closest = currentPoint
		*closestDist = currentDist
	}

	// region a is the child containing target, region b is the other child
	var aLeft, aRight, bLeft, bRight int
	if target.GetDim(dim) < currentPoint.GetDim(dim) {
		aLeft = left
		aRight = middleIdx - 1
		bLeft = middleIdx + 1
		bRight = right
	} else {
		bLeft = left
		bRight = middleIdx - 1
		aLeft = middleIdx + 1
		aRight = right
	}
	nextDim := (dim + 1) % 3
	nearestHelper(arr, nextDim, aLeft, aRight, target, closest, closestDist)
	dimDist := target.GetDim(dim) - currentPoint.GetDim(dim)
	if dimDist*dimDist < *closestDist {
		nearestHelper(arr, nextDim, bLeft, bRight, target, closest, closestDist)
	}
}

// Balances the tree if the number of nodes stored outside the tree is greater than the threshold
func (kdtree *KDTree) BalanceWithThreshold(threshold int) {
	if len(kdtree.points)-kdtree.treeSize > threshold {
		kdtree.Balance()
	}
}

// Adds the all the nodes that were originally stored outside the tree and rebalances
func (kdtree *KDTree) Balance() {
	kdtree.treeSize = len(kdtree.points)
	balanceHelper(kdtree.points, 0, 0, kdtree.treeSize-1)
}

func balanceHelper(arr []Point, dim int, left, right int) {
	if right < left {
		return
	}

	median := (left + right) / 2
	quickSelect(arr, dim, left, right, median)

	nextDim := (dim + 1) % 3
	balanceHelper(arr, nextDim, left, median-1)
	balanceHelper(arr, nextDim, median+1, right)
}

// Does a quick select for the kth element
func quickSelect(arr []Point, dim int, left, right, targetIdx int) {
	for left != right {
		pivotIdx := partition(arr, dim, left, right)
		if targetIdx == pivotIdx {
			return
		} else if targetIdx < pivotIdx {
			right = pivotIdx - 1
		} else {
			left = pivotIdx + 1
		}
	}
}

// Does a partition and returns the final index of the pivot
func partition(arr []Point, dim int, left, right int) int {
	pivotIdx := (left + right) / 2
	pivotValue := arr[pivotIdx].GetDim(dim)
	arr[pivotIdx], arr[right] = arr[right], arr[pivotIdx]
	storeIdx := left
	for i := left; i < right; i++ {
		if arr[i].GetDim(dim) < pivotValue {
			arr[storeIdx], arr[i] = arr[i], arr[storeIdx]
			storeIdx++
		}
	}
	arr[storeIdx], arr[right] = arr[right], arr[storeIdx]
	return storeIdx
}

func (kdtree *KDTree) Size() int {
	return len(kdtree.points)
}
