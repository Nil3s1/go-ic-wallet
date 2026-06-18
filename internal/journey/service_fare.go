package journey

const BaseFare = 200

type CalculatedFare struct {
	fare     int
	distance int
}

func (v CalculatedFare) Fare() int     { return v.fare }
func (v CalculatedFare) Distance() int { return v.distance }

type FareCalculator struct {
	provider StationProvider
}

func (fc *FareCalculator) CalculateFare(start string, end string) (cf CalculatedFare, err error) {
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

	fare = BaseFare + (farePerDistance * distance)

	calculatedFare := CalculatedFare{
		fare:     fare,
		distance: distance,
	}

	return calculatedFare, nil
}
