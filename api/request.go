package api

type SetStatusReq struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ProgramsName string `json:"programs_name"`
	Status       bool   `json:"status"`
}
type CheckStatusReq struct {
	Name string `json:"name"`
}
