package util

import "math"

type Triplex struct {
	X, Y, Z float64
}

func MakeTriplex(x, y, z float64) Triplex {
	return Triplex{x, y, z}
}

func (t Triplex) Add(o Triplex) Triplex {
	return Triplex{t.X + o.X, t.Y + o.Y, t.Z + o.Z}
}

func (t Triplex) Multiply(m float64) Triplex {
	return Triplex{t.X * m, t.Y * m, t.Z * m}
}

func (t Triplex) Pow(exp float64) Triplex {
	r := t.Length()
	theta := math.Acos(t.Z / r)
	phi := math.Atan2(t.Y, t.X)
	if math.IsNaN(theta) {
		theta = 0
	}
	nr := math.Pow(r, exp)
	ntheta := exp * theta
	nphi := exp * phi
	return Triplex{
		nr * math.Sin(ntheta) * math.Cos(nphi),
		nr * math.Sin(ntheta) * math.Sin(nphi),
		nr * math.Cos(ntheta),
	}
}

func (t Triplex) LengthSquared() float64 {
	return t.X*t.X + t.Y*t.Y + t.Z*t.Z
}

func (t Triplex) Length() float64 {
	return math.Sqrt(t.X*t.X + t.Y*t.Y + t.Z*t.Z)
}
