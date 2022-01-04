package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateConditionForDeactivatingSubscribers(t *testing.T) {
	result := GenerateConditionForUpdatingSubscribers([]int64{1, 2, 3})
	expected := "WHERE tg_chat_id=1 OR tg_chat_id=2 OR tg_chat_id=3"
	assert.Equalf(t, expected, result, "")
}
