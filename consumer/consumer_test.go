package consumer

import (
	"testing"
)

func TestProcessMessageWithTag(t *testing.T) {
	msg := "0b137205-3291-5f5b-5832-ab2458b9936a@nginx-1 This is test logmessage"
	uuid, tag, msg, _ := processMessage(msg)

	if uuid != "0b137205-3291-5f5b-5832-ab2458b9936a" {
		t.Error("Wrong uuid:", uuid)
	}
	if tag != "nginx-1" {
		t.Error("Wrong tag:", tag)
	}
	if msg != "This is test logmessage" {
		t.Error("Wrong message", msg)
	}
}

func TestProcessMessageWithoutTag(t *testing.T) {
	msg := "0b137205-3291-5f5b-5832-ab2458b9936a This is test logmessage"
	uuid, tag, msg, _ := processMessage(msg)

	if uuid != "0b137205-3291-5f5b-5832-ab2458b9936a" {
		t.Error("Wrong uuid:", uuid)
	}
	if tag != "default" {
		t.Error("Wrong tag:", tag)
	}
	if msg != "This is test logmessage" {
		t.Error("Wrong message", msg)
	}
}

func TestProcessMessageError(t *testing.T) {
	msg := "This is test logmessage"
	_, _, _, err := processMessage(msg)

	if err == nil {
		t.Error("No error")
	}
}
func BenchmarkProcessMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processMessage("0b137205-3291-5f5b-5832-ab2458b9936a This is test logmessage")
	}
}
