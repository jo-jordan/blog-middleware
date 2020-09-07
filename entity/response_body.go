package entity

type ResponseBody struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}