package entity

type Markdown struct {
	Content string `json:"content"`
}

type RobotSendReq struct {
	Msgtype  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}
