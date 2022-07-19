////////////////////////////////////////////////////////////////////////////////
// wraith - the wraith game engine and Server
// Copyright (c) 2022 Michael D. Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
////////////////////////////////////////////////////////////////////////////////

package cheese

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mdhender/wraith/internal/orders"
	"github.com/mdhender/wraith/internal/osk"
	"github.com/mdhender/wraith/models"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func (s *Server) clusterGetHandler(templates string) http.HandlerFunc {
	t := osk.New(templates, "cluster_list.html")
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		if _, ok := s.claims[strings.ToLower(claims["user_id"].(string))]; !ok {
			log.Printf("%s: %s: fetchClusterListByGame: not ok\n", r.Method, r.URL.Path)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// todo: use claim to fetch game

		pGameName := chi.URLParam(r, "game")
		game, err := s.store.LookupGameByName(pGameName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		x, _ := strconv.Atoi(chi.URLParam(r, "x"))
		y, _ := strconv.Atoi(chi.URLParam(r, "y"))
		z, _ := strconv.Atoi(chi.URLParam(r, "z"))

		systems, err := s.store.FetchClusterListByGame(game.Id)
		if err != nil {
			log.Printf("%s: %s: game %d: %v\n", r.Method, r.URL.Path, game.Id, err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		scans := make(ClusterList, len(systems), len(systems))
		for i, sys := range systems {
			dx, dy, dz := sys.X-x, sys.Y-y, sys.Z-z
			scans[i] = &ClusterListItem{
				Distance: math.Sqrt(float64(dx*dx + dy*dy + dz*dz)),
				X:        sys.X,
				Y:        sys.Y,
				Z:        sys.Z,
				QtyStars: sys.QtyStars,
			}
		}
		sort.Sort(scans)

		type row struct {
			X, Y, Z, QtyStars int
			Distance          string
			URL               string
		}
		clusterList := struct {
			X, Y, Z int // origin
			Systems []row
		}{
			X:       x,
			Y:       y,
			Z:       z,
			Systems: make([]row, len(scans), len(scans)),
		}
		for i, scan := range scans {
			clusterList.Systems[i] = row{
				X:        scan.X,
				Y:        scan.Y,
				Z:        scan.Z,
				QtyStars: scan.QtyStars,
				Distance: fmt.Sprintf("%.3f ly", scan.Distance),
				URL:      fmt.Sprintf("/ui/games/PT-1/cluster/%d/%d/%d", scan.X, scan.Y, scan.Z),
			}
		}

		t.Handle(w, r, clusterList)
	}
}

func (s *Server) homeGetHandler(templates string) http.HandlerFunc {
	t := osk.New(templates, "home.html")
	return func(w http.ResponseWriter, r *http.Request) {
		var claim *models.Claim
		_, claims, _ := jwtauth.FromContext(r.Context())
		if userId, ok := claims["user_id"].(string); !ok {
			log.Printf("%s: %s: claims[%q]: not a string\n", r.Method, r.URL.Path, "user_id")
		} else if claim, ok = s.claims[strings.ToLower(userId)]; !ok {
			log.Printf("%s: %s: claims[%q]: not ok\n", r.Method, r.URL.Path, strings.ToLower(userId))
		}

		t.Handle(w, r, claim)
	}
}

func (s *Server) loginGetHandler(templates, cookieName, token string) http.HandlerFunc {
	t := osk.New(templates, "login.html")
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Server: %s: %s\n", r.Method, r.URL.Path)
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Path:     "/",
			MaxAge:   -1, // delete any existing cookie
			HttpOnly: true,
		})
		t.ServeHTTP(w, r)
	}
}

func (s *Server) loginGetHandleSecretHandler(cookieName string, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("Server: %s: %s: %q %q\n", r.Method, r.URL.Path, chi.URLParam(r, "handle"), chi.URLParam(r, "secret"))
		u, err := s.store.FetchUserByCredentials(chi.URLParam(r, "handle"), chi.URLParam(r, "secret"))
		if err != nil {
			log.Printf("Server: %s: %s: fetchUsersByCredentials: %v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		log.Printf("Server: %s: %s: fetchUsersByCredentials: %q\n", r.Method, r.URL.Path, u.Handle)

		maxAge := 14 * 24 * 60 * 60
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Path:     "/",
			Value:    token,
			MaxAge:   maxAge,
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: handle %q handle %q email %q", r.Method, r.URL.Path, u.Handle, u.Profiles[0].Handle, u.Profiles[0].Email)))
		//_, _ = w.Write([]byte(fmt.Sprintf("%s: %s: claims[user_id] %q", r.Method, r.URL.Path, claims["user_id"])))
	}
}

func (s *Server) loginPostHandler(cookieName string, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("Server: %s: %s\n", r.Method, r.URL.Path)

		var input struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		}

		contentType := r.Header.Get("Content-type")
		switch contentType {
		case "application/json":
			r.Body = http.MaxBytesReader(w, r.Body, 1024) // enforce a maximum read of 1kb from the response body
			dec := json.NewDecoder(r.Body)                // create a json decoder that will accept only our specific fields
			dec.DisallowUnknownFields()
			if err := dec.Decode(&input); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			// call decode again to confirm that the request contained only a single JSON object
			if err := dec.Decode(&struct{}{}); err != io.EOF {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			//log.Printf("Server: %s %q: json: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
		case "application/x-www-form-urlencoded":
			if err := r.ParseForm(); err != nil {
				log.Printf("Server: %s %q: form: %+v\n", r.Method, r.URL.Path, err)
				http.SetCookie(w, &http.Cookie{Name: cookieName, Path: "/", MaxAge: -1, HttpOnly: true})
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			//log.Printf("Server: %s %q: form: %v\n", r.Method, r.URL.Path, r.PostForm)
			for k, v := range r.Form {
				switch k {
				case "username":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Username = v[0]
				case "password":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Password = v[0]
				}
			}
			//log.Printf("Server: %s %q: form: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
		case "text/html":
			if err := r.ParseForm(); err != nil {
				log.Printf("Server: %s %q: html: %+v\n", r.Method, r.URL.Path, err)
				http.SetCookie(w, &http.Cookie{Name: cookieName, Path: "/", MaxAge: -1, HttpOnly: true})
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			//log.Printf("Server: %s %q: html: %v\n", r.Method, r.URL.Path, r.PostForm)
			for k, v := range r.Form {
				switch k {
				case "username":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Username = v[0]
				case "password":
					if len(v) != 1 || !utf8.ValidString(v[0]) {
						http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
						return
					}
					input.Password = v[0]
				}
			}
			//log.Printf("Server: %s %q: html: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
		default:
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		}

		if input.Username == "" || input.Password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		//log.Printf("Server: %s %q: %v\n", r.Method, r.URL.Path, input)

		u, err := s.store.FetchUserByCredentials(input.Username, input.Password)
		if err != nil {
			log.Printf("Server: %s: %s: fetchUsersByCredentials: %v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		log.Printf("Server: %s: %s: fetchUsersByCredentials: %q\n", r.Method, r.URL.Path, u.Handle)

		claims := map[string]interface{}{"user_id": strings.ToLower(u.Handle)}

		jwtauth.SetExpiryIn(claims, time.Second*7*24*60*60)
		_, tokenString, _ := jwtauth.New("HS256", s.key, nil).Encode(claims)

		switch contentType {
		case "application/json":
			//log.Printf("Server: %s %q: json: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
			var response struct {
				Links struct {
					Self string `json:"self"`
				} `json:"links"`
				Data struct {
					Token string `json:"token"`
				} `json:"data,omitempty"`
			}
			response.Links.Self = r.URL.Path
			response.Data.Token = tokenString
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
			return
		case "application/x-www-form-urlencoded":
			//log.Printf("Server: %s %q: form: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Path:     "/",
				Value:    tokenString,
				MaxAge:   14 * 24 * 60 * 60,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/ui", http.StatusSeeOther)
			return
		case "text/html":
			//log.Printf("Server: %s %q: html: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Path:     "/",
				Value:    tokenString,
				MaxAge:   14 * 24 * 60 * 60,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/ui", http.StatusSeeOther)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) logoutHandler(cookieName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Path:     "/",
			MaxAge:   -1, // delete cookie
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) ordersGetHandler(templates string) http.HandlerFunc {
	t := osk.New(templates, "order_entry.html")

	type orderEntry struct {
		Game       string
		Year       string
		Quarter    string
		NationNo   int
		Rows, Cols int
		Orders     string
		Validate   string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s: ordersPath %q\n", r.Method, r.URL.Path, s.store.OrdersPath())

		var claim *models.Claim
		_, claims, _ := jwtauth.FromContext(r.Context())
		if userId, ok := claims["user_id"].(string); !ok {
			log.Printf("%s: %s: claims[%q]: not a string\n", r.Method, r.URL.Path, "user_id")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else if claim, ok = s.claims[strings.ToLower(userId)]; !ok {
			log.Printf("%s: %s: claims[%q]: not ok\n", r.Method, r.URL.Path, strings.ToLower(userId))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		pGameName := chi.URLParam(r, "game")
		game, err := s.store.LookupGameByName(pGameName)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		oe := orderEntry{
			Game:     game.ShortName,
			Year:     chi.URLParam(r, "year"),
			Quarter:  chi.URLParam(r, "quarter"),
			NationNo: claim.NationNo,
			Rows:     18,
			Cols:     80,
		}

		if oe.Year == "" || oe.Quarter == "" {
			oe.Year = "0000" //fmt.Sprintf("%04d", game.CurrentTurn.Year)
			oe.Quarter = "0" //fmt.Sprintf("%d", game.CurrentTurn.Quarter)
			http.Redirect(w, r, "/ui/games/PT-1/orders/0000/0", http.StatusTemporaryRedirect)
		} else if oe.Quarter == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ordersFile := filepath.Join(s.store.OrdersPath(), fmt.Sprintf("%s.%s.%s.%d.txt", game.ShortName, chi.URLParam(r, "year"), chi.URLParam(r, "quarter"), claim.NationNo))
		b, err := os.ReadFile(ordersFile)
		if err == nil {
			oe.Orders = string(b)
		}

		// we accept a boolean query parameter to validate the orders file
		if r.URL.Query().Get("validate") == "true" {
			if p, err := orders.Parse(b); err != nil {
				oe.Validate = fmt.Sprintf(";; sorry, but there was an error validating\n;; %+v\n", err)
			} else {
				bb := &bytes.Buffer{}
				for _, order := range p {
					bb.WriteString(order.String())
					bb.Write([]byte{'\n'})
				}
				oe.Validate = string(bb.Bytes())
			}
		}

		t.Handle(w, r, oe)
	}
}

func (s *Server) ordersGetRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		userId, ok := claims["user_id"].(string)
		if !ok {
			log.Printf("%s: %s: claims[%q]: not a string\n", r.Method, r.URL.Path, "user_id")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else if _, ok = s.claims[strings.ToLower(userId)]; !ok {
			log.Printf("%s: %s: claims[%q]: not ok\n", r.Method, r.URL.Path, strings.ToLower(userId))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		pGameName := chi.URLParam(r, "game")
		t, err := s.store.FetchCurrentTurn(userId, pGameName)
		if err != nil {
			log.Printf("%s: %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/ui/games/%s/orders/%04d/%d", pGameName, t.Year, t.Quarter), http.StatusTemporaryRedirect)
	}
}

func (s *Server) ordersPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s: entered\n", r.Method, r.URL.Path)

		_, claims, _ := jwtauth.FromContext(r.Context())
		userId, ok := claims["user_id"].(string)
		if !ok {
			log.Printf("%s: %s: claims[%q]: not a string\n", r.Method, r.URL.Path, "user_id")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		claim, ok := s.claims[strings.ToLower(userId)]
		if !ok {
			log.Printf("%s: %s: claims[%q]: not ok\n", r.Method, r.URL.Path, strings.ToLower(userId))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		pGameName, pYear, pQuarter := chi.URLParam(r, "game"), chi.URLParam(r, "year"), chi.URLParam(r, "quarter")
		t, err := s.store.FetchCurrentTurn(userId, pGameName)
		if err != nil {
			log.Printf("%s: %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		currentTurn := pYear + "/" + pQuarter
		if currentTurn != t.String() {
			log.Printf("%s: %s: not current turn: %q %q\n", r.Method, r.URL.Path, currentTurn, t.String())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("%s: %s: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		//log.Printf("Server: %s %q: %v\n", r.Method, r.URL.Path, r.PostForm)
		var input struct {
			orders   string
			validate bool
		}
		for k, v := range r.Form {
			switch k {
			case "orders":
				if len(v) != 1 {
					log.Printf("%s: %s: too many forms.orders: %d\n", r.Method, r.URL.Path, len(v))
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				} else if len(v[0]) < 1 || len(v[0]) > 64*1024 {
					log.Printf("%s: %s: invalid orders length: %d\n", r.Method, r.URL.Path, len(input.orders))
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				} else if !utf8.ValidString(v[0]) || len(v[0]) < 1 || len(v[0]) > 64*1024 {
					log.Printf("%s: %s: invalid utf-8 string\n", r.Method, r.URL.Path)
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				input.orders = v[0]
			case "validate":
				input.validate = true
			}
		}

		if len(input.orders) == 0 {
			log.Printf("%s: %s: invalid orders length: %d\n", r.Method, r.URL.Path, len(input.orders))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// try to replace characters we know the parser doesn't like
		o := input.orders
		for _, pair := range [][]string{
			{"\r\n", "\n"},
			{"\t", " "},
			{"\u2013", `-`},
			{"\u2014", `-`},
			{"\u2015", `-`},
			{"\u2017", `_`},
			{"\u2018", `'`},
			{"\u2019", `'`},
			{"\u201a", `,`},
			{"\u201b", `'`},
			{"\u201c", `"`},
			{"\u201d", `"`},
			{"\u201e", `"`},
			{"\u201f", `"`},
			{"\u2026", `...`},
			{"\u2032", `'`},
			{"\u2033", `"`},
		} {
			o = strings.ReplaceAll(o, pair[0], pair[1])
		}

		//log.Printf("%s: %s: ordersPath %q\n", r.Method, r.URL.Path, s.store.OrdersPath())

		ordersFile := filepath.Join(s.store.OrdersPath(), fmt.Sprintf("%s.%s.%s.%d.txt", pGameName, chi.URLParam(r, "year"), chi.URLParam(r, "quarter"), claim.NationNo))
		date := time.Now().UTC().Format(time.RFC3339)
		o = fmt.Sprintf(";; %s %d %s %s\n\n", pGameName, claim.NationNo, currentTurn, date) + o + "\n"
		if err := os.WriteFile(ordersFile, []byte(o), 0644); err != nil {
			log.Printf("%s: %s: writeFile %q: %v\n", r.Method, r.URL.Path, ordersFile, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if input.validate {
			http.Redirect(w, r, r.URL.Path+"?validate=true", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
		return
	}
}

func (s *Server) unitsGetHandler(templates string) http.HandlerFunc {
	units := s.store.FetchUnits()
	t := osk.New(templates, "units.html")
	return func(w http.ResponseWriter, r *http.Request) {
		t.Handle(w, r, units)
	}
}
