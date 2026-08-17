package domain

const BaseFare = 200

type CalculatedFare struct {
	fare     uint
	distance uint
}

func (v CalculatedFare) Fare() uint     { return v.fare }
func (v CalculatedFare) Distance() uint { return v.distance }

type FareCalculator struct {
	provider StationProvider
}

func NewFareCalculator(provider StationProvider) *FareCalculator {
	return &FareCalculator{
		provider: provider,
	}
}

func (fc *FareCalculator) CalculateFare(start string, end string) (cf CalculatedFare, err error) {
	farePerDistance := uint(10)
	fare := uint(0)
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

	fare = BaseFare + (farePerDistance * uint(distance))

	calculatedFare := CalculatedFare{
		fare:     fare,
		distance: uint(distance),
	}

	return calculatedFare, nil
}
