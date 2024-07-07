package icons

type IconType byte

const (
	INCREASE IconType = iota
	DECREASE
	PEN
	ERASER
	FILL
	SAVE
)

type Icon struct {
	Data   *[]byte
	IcType IconType
}
