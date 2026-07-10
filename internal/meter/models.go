package meter

import "time"

type Type string

const (
	TypeSinglePhase           Type = "SINGLE_PHASE"
	TypeThreePhaseDirect      Type = "THREE_PHASE_DIRECT"
	TypeThreePhaseTransformer Type = "THREE_PHASE_TRANSFORMER"
)

func (t Type) Valid() bool {
	switch t {
	case TypeSinglePhase, TypeThreePhaseDirect, TypeThreePhaseTransformer:
		return true
	default:
		return false
	}
}

type SealState string

const (
	SealIntact  SealState = "INTACT"
	SealBroken  SealState = "BROKEN"
	SealMissing SealState = "MISSING"
)

func (s SealState) Valid() bool {
	switch s {
	case "", SealIntact, SealBroken, SealMissing:
		return true
	default:
		return false
	}
}

type Meter struct {
	ID                  int64
	Type                Type
	SerialNumber        string
	ManufactureYear     *int32
	VerificationDate    *time.Time
	SealState           SealState
	TransformationRatio *int32
	InspectionActID     int64
	CreatedAt           time.Time
}
