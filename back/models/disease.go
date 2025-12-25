package models

type Disease struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`  // 内科、外科、儿科等
	Symptoms    string `json:"symptoms"`  // 症状描述
	Treatment   string `json:"treatment"` // 治疗方法
}
