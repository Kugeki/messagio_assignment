package dto

import "messagioassignment/internal/domain/message"

type GetStatsResp struct {
	All       int `json:"all"`
	Processed int `json:"processed"`
}

func (r *GetStatsResp) FromDomain(stats *message.Stats) {
	r.All = stats.All
	r.Processed = stats.Processed
}
