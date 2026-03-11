// ทำ response format กลางๆ สำหรับ API โดยใช้ struct APIResponse ที่มีฟิลด์ Success, Message และ Data เพื่อให้สามารถส่งข้อมูลกลับไปยังผู้ใช้ได้อย่างมีประสิทธิภาพและเป็นมาตรฐาน.
package dto

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
