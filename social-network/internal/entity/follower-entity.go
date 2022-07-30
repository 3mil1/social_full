package entity

type Follower struct {
	SourceId string `json:"source_id,omitempty"`
	TargetId string `json:"target_id,omitempty"`
	Status   int    `json:"status"`
}
