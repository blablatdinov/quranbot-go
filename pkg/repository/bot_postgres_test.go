package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateConditionForDeactivatingSubscribers(t *testing.T) {
	result := GenerateConditionForDeactivatingSubscribers([]int64{1, 2, 3})
	expected := "where tg_chat_id=1 or tg_chat_id=2 or tg_chat_id=3"
	assert.Equalf(t, expected, result, "")
}