package journey

type StationInfo struct {
	StationId     string
	SequenceOrder int
	TrackPosition int
}

type StationProvider interface {
	GetStationInfo(stationId string) (StationInfo, error)
}
