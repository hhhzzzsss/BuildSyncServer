package treegen

import (
	"math"
	"math/rand"
	"testing"
)

type testPoint struct {
	pos [3]float64
}

func (p *testPoint) GetDim(dim int) float64 {
	return p.pos[dim]
}

func makeTestPoint(x, y, z float64) testPoint {
	return testPoint{[3]float64{x, y, z}}
}

func TestKDTreeNoBalance(t *testing.T) {

	numPoints := 10000
	numTests := 100
	for testNum := 0; testNum < numTests; testNum++ {
		var kdtree KDTree
		points := make([]testPoint, numPoints)
		for i := 0; i < numPoints; i++ {
			points[i] = makeTestPoint(rand.Float64(), rand.Float64(), rand.Float64())
			kdtree.Add(&points[i])
		}

		targetPoint := makeTestPoint(rand.Float64(), rand.Float64(), rand.Float64())

		closestPoint := makeTestPoint(10, 10, 10)
		for i := 0; i < numPoints; i++ {
			if PointDistSq(&points[i], &targetPoint) < PointDistSq(&closestPoint, &targetPoint) {
				closestPoint = points[i]
			}
		}

		kdClosestPoint := kdtree.NearestNeighbor(&targetPoint)

		actualDist := PointDistSq(&closestPoint, &targetPoint)
		kdDist := PointDistSq(kdClosestPoint, &targetPoint)

		if math.Abs(kdDist-actualDist) > 1e-10 {
			t.Errorf("Closest distance given by KDTree does not match actual (%f != %f)", kdDist, actualDist)
			return
		}
	}
}

func TestKDTreeBalanced(t *testing.T) {
	numPoints := 10000
	numTests := 100
	for testNum := 0; testNum < numTests; testNum++ {
		var kdtree KDTree
		points := make([]testPoint, numPoints)
		for i := 0; i < numPoints; i++ {
			points[i] = makeTestPoint(rand.Float64(), rand.Float64(), rand.Float64())
			kdtree.Add(&points[i])
		}
		kdtree.Balance()

		targetPoint := makeTestPoint(rand.Float64(), rand.Float64(), rand.Float64())

		closestPoint := makeTestPoint(10, 10, 10)
		for i := 0; i < numPoints; i++ {
			if PointDistSq(&points[i], &targetPoint) < PointDistSq(&closestPoint, &targetPoint) {
				closestPoint = points[i]
			}
		}

		kdClosestPoint := kdtree.NearestNeighbor(&targetPoint)

		actualDist := PointDistSq(&closestPoint, &targetPoint)
		kdDist := PointDistSq(kdClosestPoint, &targetPoint)

		if math.Abs(kdDist-actualDist) > 1e-10 {
			t.Errorf("Closest distance given by KDTree does not match actual (%f != %f)", kdDist, actualDist)
			return
		}
	}
}

func TestKDTreePartiallyBalanced(t *testing.T) {
	numBalancedPoints := 10000
	numUnbalancedPoints := 10000
	numTests := 100
	for testNum := 0; testNum < numTests; testNum++ {
		var kdtree KDTree
		points := make([]testPoint, numBalancedPoints+numUnbalancedPoints)
		for i := 0; i < numBalancedPoints; i++ {
			points[i] = makeTestPoint(rand.Float64(), rand.Float64(), rand.Float64())
			kdtree.Add(&points[i])
		}
		kdtree.Balance()
		for i := numBalancedPoints; i < numBalancedPoints+numUnbalancedPoints; i++ {
			points[i] = makeTestPoint(rand.Float64(), rand.Float64(), rand.Float64())
			kdtree.Add(&points[i])
		}

		targetPoint := makeTestPoint(rand.Float64(), rand.Float64(), rand.Float64())

		closestPoint := makeTestPoint(10, 10, 10)
		for i := 0; i < numBalancedPoints+numUnbalancedPoints; i++ {
			if PointDistSq(&points[i], &targetPoint) < PointDistSq(&closestPoint, &targetPoint) {
				closestPoint = points[i]
			}
		}

		kdClosestPoint := kdtree.NearestNeighbor(&targetPoint)

		actualDist := PointDistSq(&closestPoint, &targetPoint)
		kdDist := PointDistSq(kdClosestPoint, &targetPoint)

		if math.Abs(kdDist-actualDist) > 1e-10 {
			t.Errorf("Closest distance given by KDTree does not match actual (%f != %f)", kdDist, actualDist)
			return
		}
	}
}
