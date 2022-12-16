package models

type Performance struct {
	PerformanceId  int
	SongId         int
	ClearRate      float64
	FullComboRate  float64
	AllPerfectRate float64
	GeneralRate    float64
	Try            int
}
