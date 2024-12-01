package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// JSONElement は動的なJSON要素を表現する構造体
type JSONElement struct {
	mu       sync.RWMutex
	data     map[string]interface{}
	children map[string]*JSONElement
}

// NewJSONElement は新しいJSONElementを作成
func NewJSONElement() *JSONElement {
	return &JSONElement{
		data:     make(map[string]interface{}),
		children: make(map[string]*JSONElement),
	}
}

// SetValue は指定されたキーに値を設定
func (je *JSONElement) SetValue(key string, value interface{}) {
	je.mu.Lock()
	defer je.mu.Unlock()
	je.data[key] = value
}

// AddChild は新しい子要素を追加
func (je *JSONElement) AddChild(key string) *JSONElement {
	je.mu.Lock()
	defer je.mu.Unlock()

	child := NewJSONElement()
	je.children[key] = child
	return child
}

// ToJSON は現在の要素をJSON形式に変換
func (je *JSONElement) ToJSON() ([]byte, error) {
	je.mu.RLock()
	defer je.mu.RUnlock()

	// データと子要素をマージ
	fullData := make(map[string]interface{})

	// 通常のデータをコピー
	for k, v := range je.data {
		fullData[k] = v
	}

	// 子要素をマージ
	for k, child := range je.children {
		childJSON, err := child.ToJSON()
		if err != nil {
			return nil, err
		}

		var childData map[string]interface{}
		if err := json.Unmarshal(childJSON, &childData); err != nil {
			return nil, err
		}

		fullData[k] = childData
	}

	return json.MarshalIndent(fullData, "", "  ")
}

// デモンストレーション用のmain関数
func edit() {
	// 使用例
	root := NewJSONElement()

	// ルート要素にデータを追加
	root.SetValue("projectName", "Dynamic JSON Builder")
	root.SetValue("version", "1.0.0")

	// 子要素の追加
	person := root.AddChild("person")
	person.SetValue("name", "John Doe")
	person.SetValue("age", 30)

	// 子要素のネスト
	address := person.AddChild("address")
	address.SetValue("street", "123 Main St")
	address.SetValue("city", "Anytown")
	address.SetValue("country", "USA")

	// JSONに変換して出力
	jsonData, err := root.ToJSON()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))
}
