package params

type Orientation int

const (
	OrientationPortrait Orientation = iota
	OrientationLandscape
)

type RenderParam struct {
	WhiteBackground bool
	DeviceType      int
	Orientation     Orientation
}
