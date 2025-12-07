package model

// ConsumResponse 核销响应
type ConsumResponse struct {
	Status      int
	Message     string
	Name        string
	CardNo      string
	Money       int
	Subsidy     float64
	Times       int
	Integral    float64
	InTime      string
	OutTime     string
	Cumulative  string
	Amount      string
	VoiceID     string
	Text        string
}