package edgeDetection

import "math"

func (a Vector3) Add3(b Vector3) Vector3 {
	return Vector3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func (a Vector3) Sub3(b Vector3) Vector3 {
	return Vector3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func (a Vector3) MultiplyByScalar3(s float64) Vector3 {
	return Vector3{
		X: a.X * s,
		Y: a.Y * s,
		Z: a.Z * s,
	}
}

func (a Vector3) Dot3(b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector3) Length3() float64 {
	return math.Sqrt(a.Dot3(a))
}

func (a Vector3) Cross3(b Vector3) Vector3 {
	return Vector3{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (a Vector3) Normalize3() Vector3 {
	return a.MultiplyByScalar3(1. / a.Length3())
}

func (a Vector3) AngleTo3(b Vector3) float64 {

	theta := a.Dot3(b) / (a.Length3() * b.Length3())
	// clamp, to handle numerical problems
	return math.Acos(clamp(theta, -1, 1))
}

// Clamp clamps x to the provided closed interval [a, b]
func clamp(x, a, b float64) float64 {

	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

func (a Vector2) Add2(b Vector2) Vector2 {
	return Vector2{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func (a Vector2) Length2() float64 {
	return math.Sqrt(a.Dot2(a))
}

func (a Vector2) Dot2(b Vector2) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vector2) MultiplyByScalar2(s float64) Vector2 {
	return Vector2{
		X: a.X * s,
		Y: a.Y * s,
	}
}
