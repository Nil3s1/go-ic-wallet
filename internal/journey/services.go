package journey

type CalculatedFare struct {
	fare     int
	distance int
}

func (v CalculatedFare) Fare() int     { return v.fare }
func (v CalculatedFare) Distance() int { return v.distance }

type StationInfo struct {
	StationId     string
	SequenceOrder int
	TrackPosition int
}

type StationProvider interface {
	GetStationInfo(stationId string) (StationInfo, error)
}

type FareCalculator struct {
	provider StationProvider
}

func (fc *FareCalculator) CalculateFare(start string, end string) (cf CalculatedFare, err error) {
	baseFare := 200
	farePerDistance := 10
	fare := 0
	startInfo, err := fc.provider.GetStationInfo(start)
	if err != nil {
		return CalculatedFare{}, err
	}

	endInfo, err := fc.provider.GetStationInfo(end)
	if err != nil {
		return CalculatedFare{}, err
	}

	distance := endInfo.TrackPosition - startInfo.TrackPosition

	if distance < 0 {
		distance = -distance
	}

	fare = baseFare + (farePerDistance * distance)

	calculatedFare := CalculatedFare{
		fare:     fare,
		distance: distance,
	}

	return calculatedFare, nil
}
