package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DeliveryRequest struct{ AssetID, SubscriberID, ContentID string }

func errorGroup(req DeliveryRequest, stage string) []string {
	return []string{"creator-commerce", stage, req.AssetID}
}

func processDelivery(req DeliveryRequest) (string, error) {
	if strings.TrimSpace(req.AssetID) == "" {
		return "", fmt.Errorf("asset_id is required")
	}
	if strings.TrimSpace(req.SubscriberID) == "" {
		return "", fmt.Errorf("subscriber_id is required")
	}
	if strings.TrimSpace(req.ContentID) == "" {
		return "", fmt.Errorf("content_id is required")
	}
	return "delivered", nil
}

func deliveryHandler(client *InfraiClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req DeliveryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		result, err := processDelivery(req)
		if err != nil {
			_ = client.Capture(map[string]any{"title": "delivery rejected", "message": err.Error(), "level": "error", "fingerprint": errorGroup(req, "delivery"), "exception": err.Error(), "context": map[string]any{"asset_id": req.AssetID, "subscriber_id": req.SubscriberID, "content_id": req.ContentID}})
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": result})
	}
}
