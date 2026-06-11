package httputil_test

import (
	"encoding/json"
	"testing"

	"github.com/Jarukit-PM/TicketBookingSystem/api/pkg/httputil"
)

func TestJSONSliceNilEncodesEmptyArray(t *testing.T) {
	var items []string
	encoded, err := json.Marshal(httputil.JSONSlice(items))
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != "[]" {
		t.Fatalf("encoded = %s, want []", encoded)
	}
}
