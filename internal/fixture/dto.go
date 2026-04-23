package fixture

type CreateFixtureRequest struct {
	Type      string `json:"type"`
	FixtureNo string `json:"fixture_no"`

	FixtureName *string `json:"fixture_name"`
	Description *string `json:"description"`

	Cavities  *int   `json:"cavities"`
	LifeShots *int64 `json:"life_shots"`

	FixtureType  *string                `json:"fixture_type"`
	Material     *string                `json:"material"`
	SpecialNotes map[string]interface{} `json:"special_notes"`
}

type UpdateFixtureRequest struct {
	FixtureNo *string `json:"fixture_no"`

	FixtureName *string `json:"fixture_name"`
	Description *string `json:"description"`

	Cavities  *int   `json:"cavities"`
	LifeShots *int64 `json:"life_shots"`

	FixtureType *string `json:"fixture_type"`
	Material    *string `json:"material"`
}
