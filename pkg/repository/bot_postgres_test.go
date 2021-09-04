package repository

import (
	"testing"
)

func TestGenerateConditionForDeactivatingSubscribers(t *testing.T) {
	result := GenerateConditionForDeactivatingSubscribers([]int64{1, 2, 3})
	expected := "where tg_chat_id=1 or tg_chat_id=2 or tg_chat_id=3"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}