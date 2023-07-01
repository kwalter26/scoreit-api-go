package util

// BaseballPosition is a custom type to represent baseball positions
type BaseballPosition string

// Constants representing baseball positions
const (
	Pitcher          BaseballPosition = "PITCHER"
	Catcher          BaseballPosition = "CATCHER"
	FirstBase        BaseballPosition = "FIRST_BASE"
	SecondBase       BaseballPosition = "SECOND_BASE"
	ThirdBase        BaseballPosition = "THIRD_BASE"
	ShortStop        BaseballPosition = "SHORT_STOP"
	LeftField        BaseballPosition = "LEFT_FIELD"
	CenterField      BaseballPosition = "CENTER_FIELD"
	RightField       BaseballPosition = "RIGHT_FIELD"
	DesignatedHitter BaseballPosition = "DESIGNATED_HITTER"
	RightCenterField BaseballPosition = "RIGHT_CENTER_FIELD"
	LeftCenterField  BaseballPosition = "LEFT_CENTER_FIELD"
)
