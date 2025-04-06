package controller

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"

	storage "github.com/barcek2281/MyEcho/internal/app/store"
	"github.com/barcek2281/MyEcho/pkg/utils"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type ControllerMs struct {
	storage *storage.Storage
	session sessions.Store
	log     *logrus.Logger
}

func NewControllerMs(store *storage.Storage, session sessions.Store, log *logrus.Logger) *ControllerMs {
	return &ControllerMs{
		storage: store,
		session: session,
		log:     log,
	}
}

func (c *ControllerMs) PaymentPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/payment.html")
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			c.log.Errorf("Error with parse file: %v", err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			c.log.Errorf("Error with execute file: %v", err)
			return
		}
		c.log.Infof("handle payment/ GET")
	}
}

func (c *ControllerMs) PaymentPost() http.HandlerFunc {
	type Request struct {
		Id    string  `json:"id"`
		Plan  string  `json:"plan"`
		Price float64 `json:"price"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			c.log.Error("Error:", err)
			return
		}

		session, err := c.session.Get(r, sessionName)
		if err != nil {
			utils.Error(w, r, 403, err)
			c.log.Error("Error:", err)
			return
		}
		user_id, ok := session.Values["user_id"].(int)
		if !ok {
			utils.Error(w, r, http.StatusForbidden, err)
			c.log.Error("Error:", err)
			return
		}
		m, err := c.storage.User().FindById(user_id)
		if err != nil {
			utils.Error(w, r, http.StatusForbidden, err)
			c.log.Error("Error:", err)
			return
		}

		type MicroRequest struct {
			RRequest struct {
				Id    string  `json:"id"`
				Plan  string  `json:"plan"`
				Price float64 `json:"price"`
			} `json:"CartItem"`
			Customer struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"customer"`
		}

		mresquest := MicroRequest{RRequest: req, Customer: struct {
			ID    string "json:\"id\""
			Name  string "json:\"name\""
			Email string "json:\"email\""
		}{strconv.Itoa(m.ID), m.Login, m.Email}}

		c.log.Infof("request to microservice: %+v", mresquest)

		data, err := json.Marshal(mresquest)
		if err != nil {
			c.log.Error("Error:", err)
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		resp, err := http.Post("http://localhost:8081/payment", "application/json", bytes.NewBuffer(data))

		if resp.StatusCode == 200 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				c.log.Error("Error decoding JSON:", err)
				utils.Error(w, r, http.StatusBadRequest, nil)
				return
			}
			c.log.Info("Raw response body:", string(body)) // Посмотреть ответ
			var response struct {
				Status         string `json:"status"`
				Transaction_id int    `json:"transaction_id"`
			}
			if err := json.Unmarshal(body, &response); err != nil {
				c.log.Error("Error decoding JSON:", err)
				return
			}

			utils.Response(w, r, 200, response)
		} else {
			utils.Error(w, r, http.StatusBadRequest, nil)
			return
		}

		c.log.Infof("Handle /payment POST")
	}
}

func (c *ControllerMs) ProcessPaymentPost() http.HandlerFunc {
	type request struct {
		ID             string `json:"id"`
		CardNumber     string `json:"cardNumber"`
		ExpirationDate string `json:"expirationDate"`
		CVV            string `json:"cvv"`
		Name           string `json:"name"`
		Address        string `json:"address"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			c.log.Errorf("error, bad request: %v", err)
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}
		type MicroRequest struct {
			CardNumber     string `json:"cardNumber"`
			ExpirationDate string `json:"expirationDate"`
			CVV            string `json:"cvv"`
			Name           string `json:"name"`
			Address        string `json:"address"`
		}

		mreq := MicroRequest{
			CardNumber:     req.CardNumber,
			ExpirationDate: req.ExpirationDate,
			CVV:            req.CVV,
			Name:           req.Name,
			Address:        req.Address,
		}

		s, err := c.session.Get(r, sessionName)
		if err != nil {
			utils.Error(w, r, http.StatusForbidden, err)
			return
		}

		user_id, ok := s.Values["user_id"].(int)
		if !ok {
			utils.Error(w, r, http.StatusForbidden, err)
			c.log.Error("Error:", err)
			return
		}

		data, err := json.Marshal(mreq)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		resp, err := http.Post("http://localhost:8081/process-payment/"+req.ID, "application/json", bytes.NewBuffer(data))
		if err != nil {
			utils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		if resp.StatusCode == http.StatusOK {
			c.log.Info(req.ID, req.Address)
			err = c.storage.User().Prime(user_id)
			if err != nil {
				utils.Error(w, r, http.StatusNotAcceptable, err)
				return
			}
			utils.Response(w, r, 200, resp.Body)
		} else {
			c.log.Info(req.ID, req.Address, "cannot make")
			utils.Response(w, r, http.StatusNotAcceptable, resp.Body)
		}
		c.log.Info("Handle /process-payment POST", req.ID)
	}
}

func (c *ControllerMs) RemovePrime() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := c.session.Get(r, sessionName)
		if err != nil {
			utils.Error(w, r, http.StatusForbidden, err)
			return
		}

		user_id, ok := s.Values["user_id"].(int)
		if !ok {
			utils.Error(w, r, http.StatusForbidden, err)
			c.log.Error("Error:", err)
			return
		}

		err = c.storage.User().DeactivatePrime(user_id)
		if err != nil {
			utils.Error(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusAccepted, nil)
		c.log.Infof("hanle /removePrime GET")
	}
}
