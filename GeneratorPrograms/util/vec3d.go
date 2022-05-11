package util

import "math"

type Vec3d struct {
	X, Y, Z float64
}

func MakeVec3d(x, y, z float64) Vec3d {
	return Vec3d{x, y, z}
}

func (v1 Vec3d) Add(v2 Vec3d) Vec3d {
	return Vec3d{
		v1.X + v2.X,
		v1.Y + v2.Y,
		v1.Z + v2.Z,
	}
}

func (v1 Vec3d) Sub(v2 Vec3d) Vec3d {
	return Vec3d{
		v1.X - v2.X,
		v1.Y - v2.Y,
		v1.Z - v2.Z,
	}
}

func (v Vec3d) Scale(factor float64) Vec3d {
	return Vec3d{
		v.X * factor,
		v.Y * factor,
		v.Z * factor,
	}
}

func (v Vec3d) InvScale(factor float64) Vec3d {
	return Vec3d{
		v.X / factor,
		v.Y / factor,
		v.Z / factor,
	}
}

func (v Vec3d) Normalize() Vec3d {
	l2 := v.LengthSquared()
	if l2 == 0 {
		return Vec3d{0, 0, 0}
	} else {
		return v.InvScale(v.Length())
	}
}

func (v Vec3d) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3d) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}
