package model

type RecipeStep struct {
	StepComplete    string `json:"rcp_stp_cmp"`
	StepCount       string `json:"rcp_stp"`
	RealTime        string `json:"rtm"`
	RealTemperature string `json:"rtp"`
	SetTime         string `json:"stm"`
}
