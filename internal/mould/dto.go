package mould

type CreateMoldRequest struct {
	Type        string `json:"type"`
	MoldName    string `json:"mold_name"`
	MoldNo      string `json:"mold_no"`
	Description string `json:"description"`
	Cavities    int    `json:"cavities"`
	TargetShots int64  `json:"target_shots"`
}

type UpdateMoldRequest struct {
	MoldNo      string `json:"mold_no"`
	Description string `json:"description"`
	Cavities    int    `json:"cavities"`
	TargetShots int64  `json:"target_shots"`
}
