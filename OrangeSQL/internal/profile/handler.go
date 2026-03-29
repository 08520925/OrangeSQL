package profile

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ListHandler は GET /api/profiles のハンドラを返す。
func ListHandler(cm *ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		store, err := cm.Store().Load()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(ProfilesResponse{
			Profiles: store.Profiles,
			ActiveID: cm.ActiveID(),
		})
	}
}

// CreateHandler は POST /api/profiles のハンドラを返す。
func CreateHandler(cm *ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON"})
			return
		}

		if req.Driver == "" {
			req.Driver = "sqlite"
		}

		// バリデーション用に仮 Profile を構築して検証
		tmp := Profile{Name: strings.TrimSpace(req.Name), Driver: req.Driver, Path: strings.TrimSpace(req.Path), Host: strings.TrimSpace(req.Host), DBName: strings.TrimSpace(req.DBName)}
		if err := tmp.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		store, err := cm.Store().Load()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		p := NewProfile(NextID(store), req)
		store.Profiles = append(store.Profiles, p)

		if err := cm.Store().Save(store); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

// UpdateHandler は PUT /api/profiles/{id} のハンドラを返す。
func UpdateHandler(cm *ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := r.PathValue("id")

		var req UpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON"})
			return
		}

		store, err := cm.Store().Load()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		found := false
		var updated Profile
		for i, p := range store.Profiles {
			if p.ID == id {
				if req.Name != "" {
					store.Profiles[i].Name = req.Name
				}
				if req.Path != "" {
					store.Profiles[i].Path = req.Path
				}
				if req.Host != "" {
					store.Profiles[i].Host = req.Host
				}
				if req.Port > 0 {
					store.Profiles[i].Port = req.Port
				}
				if req.User != "" {
					store.Profiles[i].User = req.User
				}
				if req.Password != "" {
					store.Profiles[i].Password = req.Password
				}
				if req.DBName != "" {
					store.Profiles[i].DBName = req.DBName
				}
				if req.SSLMode != "" {
					store.Profiles[i].SSLMode = req.SSLMode
				}
				updated = store.Profiles[i]
				found = true
				break
			}
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "profile not found"})
			return
		}

		if err := cm.Store().Save(store); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(updated)
	}
}

// DeleteHandler は DELETE /api/profiles/{id} のハンドラを返す。
func DeleteHandler(cm *ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := r.PathValue("id")

		if id == cm.ActiveID() {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "cannot delete active profile"})
			return
		}

		store, err := cm.Store().Load()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		found := false
		filtered := make([]Profile, 0, len(store.Profiles))
		for _, p := range store.Profiles {
			if p.ID == id {
				found = true
				continue
			}
			filtered = append(filtered, p)
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "profile not found"})
			return
		}

		store.Profiles = filtered
		if err := cm.Store().Save(store); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
	}
}

// ConnectHandler は POST /api/profiles/{id}/connect のハンドラを返す。
func ConnectHandler(cm *ConnectionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := r.PathValue("id")

		store, err := cm.Store().Load()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		var target *Profile
		for _, p := range store.Profiles {
			if p.ID == id {
				target = &p
				break
			}
		}

		if target == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "profile not found"})
			return
		}

		if err := cm.SwitchTo(*target); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"connected": true,
			"database":  cm.DB().Name(),
		})
	}
}
