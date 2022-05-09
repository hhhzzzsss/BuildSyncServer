package util

import "math"

type Quaternion struct {
	W, X, Y, Z float64
}

func MakeQuaternion(w, x, y, z float64) Quaternion {
	return Quaternion{w, x, y, z}
}

func (q1 Quaternion) Add(q2 Quaternion) Quaternion {
	return Quaternion{
		q1.W + q2.W,
		q1.X + q2.X,
		q1.Y + q2.Y,
		q1.Z + q2.Z,
	}
}

func (q Quaternion) Scale(m float64) Quaternion {
	return Quaternion{
		q.W * m,
		q.X * m,
		q.Y * m,
		q.Z * m,
	}
}

func (q1 Quaternion) Multiply(q2 Quaternion) Quaternion {
	return Quaternion{
		q1.W*q2.W - q1.X*q2.X - q1.Y*q2.Y - q1.Z*q2.Z,
		q1.W*q2.X + q1.X*q2.W + q1.Y*q2.Z - q1.Z*q2.Y,
		q1.W*q2.Y - q1.X*q2.Z + q1.Y*q2.W + q1.Z*q2.X,
		q1.W*q2.Z + q1.X*q2.Y - q1.Y*q2.X + q1.Z*q2.W,
	}
}

func (q Quaternion) Square() Quaternion {
	return Quaternion{
		q.W*q.W - q.X*q.X - q.Y*q.Y - q.Z*q.Z,
		2 * q.W * q.X,
		2 * q.W * q.Y,
		2 * q.W * q.Z,
	}
}

func (q Quaternion) Cube() Quaternion {
	q2 := Quaternion{
		q.W * q.W,
		q.X * q.X,
		q.Y * q.Y,
		q.Z * q.Z,
	}
	rTemp := 3*q2.W - q2.X - q2.Y - q2.Z
	return Quaternion{
		q.W * (q2.W - 3*q2.X - 3*q2.Y - 3*q2.Z),
		q.X * rTemp,
		q.Y * rTemp,
		q.Z * rTemp,
	}
}

func (t Quaternion) LengthSquared() float64 {
	return t.W*t.W + t.X*t.X + t.Y*t.Y + t.Z*t.Z
}

func (t Quaternion) Length() float64 {
	return math.Sqrt(t.W*t.W + t.X*t.X + t.Y*t.Y + t.Z*t.Z)
}
