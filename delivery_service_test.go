package main

import "testing"

func TestProcessDeliveryDecision(t *testing.T) {
	tests := []struct {
		name    string
		req     DeliveryRequest
		want    string
		wantErr bool
	}{
		{"complete delivery", DeliveryRequest{"asset-7", "sub-2", "post-9"}, "delivered", false},
		{"missing content is rejected", DeliveryRequest{"asset-7", "sub-2", ""}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processDelivery(tt.req)
			if (err != nil) != tt.wantErr || got != tt.want {
				t.Fatalf("got %q, err %v", got, err)
			}
		})
	}
}

func TestErrorGroupKeepsStageAndAsset(t *testing.T) {
	got := errorGroup(DeliveryRequest{AssetID: "asset-7"}, "content-processing")
	want := []string{"creator-commerce", "content-processing", "asset-7"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("group[%d] = %q", i, got[i])
		}
	}
}
