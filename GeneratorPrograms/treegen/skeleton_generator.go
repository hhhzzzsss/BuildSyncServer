package treegen

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hhhzzzsss/build-sync-generator/util"
)

type SkeletonNode struct {
	pos       [3]float64
	children  []*SkeletonNode
	parent    *SkeletonNode
	thickness float64
}

type Attractor struct {
	pos [3]float64
}

func (node *SkeletonNode) GetDim(dim int) float64 {
	return node.pos[dim]
}

func (node *Attractor) GetDim(dim int) float64 {
	return node.pos[dim]
}

func NewSkeletonNode(x, y, z float64) *SkeletonNode {
	return &SkeletonNode{
		[3]float64{x, y, z},
		make([]*SkeletonNode, 0),
		nil,
		0,
	}
}

func NewAttractor(x, y, z float64) *Attractor {
	return &Attractor{[3]float64{x, y, z}}
}

type Skeleton struct {
	roots     []*SkeletonNode
	nodeCache []*SkeletonNode
}

type GeneratorSettings struct {
	StepSize         float64
	KillDistance     float64
	AttractionRadius float64

	BranchPower    float64
	BranchDecay    float64
	BaseThickness  float64
	BaseBulge      float64
	TaperThreshold float64
	TaperRate      float64 // should be less than StepSize to have an effect

	BalancingThreshold int
}

func GetDefaultSettings() GeneratorSettings {
	return GeneratorSettings{

		StepSize:         1,
		KillDistance:     5,
		AttractionRadius: math.Inf(1),

		BranchPower:    2,
		BranchDecay:    0,
		BaseThickness:  15,
		BaseBulge:      0,
		TaperThreshold: math.Inf(1),
		TaperRate:      0.7,

		BalancingThreshold: 20,
	}
}

func GenerateSkeleton(roots []*SkeletonNode, attractors []*Attractor, settings GeneratorSettings) Skeleton {
	kdtree := MakeKDTree()
	for _, root := range roots {
		kdtree.Add(root)
	}
	iterations := 0
	fmt.Print("Generating skeleton (0 iterations)...")
	for skeletonIteration(&kdtree, &attractors, settings) {
		iterations++
		fmt.Printf("\rGenerating skeleton... (%d iterations)", iterations)
	}
	fmt.Print(" Done.\n")
	fmt.Print("Calculating node diameters...")
	for _, root := range roots {
		calculateThickness(root, &settings)
		scaleThickness(root, settings.BaseThickness/root.thickness)
		root.thickness += settings.BaseBulge
		taperThickness(root, &settings)
	}
	fmt.Print(" Done.\n")

	fmt.Print("Creating node cache for final tree...")
	skeleton := Skeleton{roots, nil}
	for _, root := range roots {
		skeleton.buildNodeCache(root)
	}
	fmt.Print(" Done.\n")

	return skeleton
}

// Returns true if there are no more attractors in reach (i.e. finished), false otherwise
func skeletonIteration(kdtree *KDTree, attractors *[]*Attractor, settings GeneratorSettings) bool {
	kdtree.BalanceWithThreshold(settings.BalancingThreshold)

	attractionMap := make(map[*SkeletonNode][]util.Vec3d)
	for _, attractor := range *attractors {
		nearestNode := kdtree.NearestNeighbor(attractor).(*SkeletonNode)
		nnPos := PointToVec3d(nearestNode)
		attrPos := PointToVec3d(attractor)
		posDif := attrPos.Sub(nnPos)
		if posDif.LengthSquared() > settings.AttractionRadius*settings.AttractionRadius {
			continue
		}
		attractionMap[nearestNode] = append(attractionMap[nearestNode], posDif.Normalize())
	}

	if len(attractionMap) == 0 {
		return false
	}

	addedNodes := 0
	for node, attractDirs := range attractionMap {
		// Add up attraction directions
		attraction := util.Vec3d{}
		for _, dir := range attractDirs {
			attraction = attraction.Add(dir)
		}

		// If attraction is very small, try removing a random attractor to see if it breaks the symmetry
		if len(attractDirs) >= 2 && attraction.Length() < 0.5 {
			exclusion := rand.Intn(len(attractDirs))
			attraction = util.Vec3d{}
			for i, dir := range attractDirs {
				if i != exclusion {
					attraction = attraction.Add(dir)
				}
			}
		}

		attraction = attraction.Normalize().Scale(settings.StepSize)
		newNodePos := PointToVec3d(node).Add(attraction)
		newNode := NewSkeletonNode(newNodePos.X, newNodePos.Y, newNodePos.Z)
		newNode.parent = node
		if PointDistSq(newNode, newNode.parent) < 0.1*settings.StepSize {
			continue
		}
		for _, child := range node.children {
			if PointDistSq(newNode, child) < 0.01*settings.StepSize {
				goto endloop
			}
		}
		node.children = append(node.children, newNode)
		kdtree.Add(newNode)
		addedNodes++
		*attractors = removeNearbyAttractors(*attractors, newNode, settings.KillDistance)
	endloop:
	}
	if addedNodes == 0 {
		return false
	}

	return true
}

func removeNearbyAttractors(attractors []*Attractor, point Point, killDist float64) []*Attractor {
	for i := 0; i < len(attractors); i++ {
		for PointDistSq(attractors[i], point) < killDist*killDist {
			attractors[i] = attractors[len(attractors)-1]
			attractors = attractors[:len(attractors)-1]
			if i >= len(attractors) {
				return attractors
			}
		}
	}
	return attractors
}

func calculateThickness(node *SkeletonNode, settings *GeneratorSettings) {
	if len(node.children) == 0 {
		node.thickness = 1
		return
	} else {
		node.thickness = 0
		for _, child := range node.children {
			calculateThickness(child, settings)
			node.thickness += math.Pow(child.thickness, settings.BranchPower)
		}
		node.thickness = math.Pow(node.thickness, 1/settings.BranchPower)
		node.thickness /= (1.0 - settings.BranchDecay)
	}
}

func scaleThickness(node *SkeletonNode, factor float64) {
	node.thickness *= factor
	for _, child := range node.children {
		scaleThickness(child, factor)
	}
}

func taperThickness(node *SkeletonNode, settings *GeneratorSettings) {
	if node.parent != nil && node.thickness >= settings.TaperThreshold {
		node.thickness = math.Max(node.thickness, node.parent.thickness-settings.TaperRate/2)
	}
	for _, child := range node.children {
		taperThickness(child, settings)
	}
}

func (s *Skeleton) buildNodeCache(node *SkeletonNode) {
	s.nodeCache = append(s.nodeCache, node)
	for _, child := range node.children {
		s.buildNodeCache(child)
	}
}

func (s Skeleton) ForEachNode(f func(node *SkeletonNode)) {
	var bar util.ProgressBar
	bar.Initialize(len(s.nodeCache))
	for i, node := range s.nodeCache {
		f(node)
		bar.Play(i + 1)
	}
	bar.Finish()
}

func (n *SkeletonNode) GetThickness() float64 {
	return n.thickness
}

func (n *SkeletonNode) IsRoot() bool {
	return n.parent == nil
}
