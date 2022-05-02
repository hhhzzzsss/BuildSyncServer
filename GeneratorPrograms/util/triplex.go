package util

import "math"

type Triplex struct {
	x, y, z float64
}

func (t Triplex) Add(o Triplex) Triplex {
	return Triplex{t.x + o.x, t.y + o.y, t.z + o.z}
}

func (t Triplex) Multiply(m float64) Triplex {
	return Triplex{t.x * m, t.y * m, t.z * m}
}

func (t Triplex) Pow(exp float64) Triplex {
	r := t.Length()
	theta := math.Atan2(t.y, t.x)
	phi := math.Acos(t.z / r)
	if math.IsNaN(phi) {
		phi = 0
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
	return t.x*t.x + t.y*t.y + t.z*t.z
}

func (t Triplex) Length() float64 {
	return math.Sqrt(t.x*t.x + t.y*t.y + t.z*t.z)
}
