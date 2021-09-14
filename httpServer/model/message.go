/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-14 10:37:28
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 11:02:32
@FilePath: /goAdapter-Raw/httpServer/model/message.go
*/
package model

type Response struct {
	Code    string      `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}
