package src

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/transport"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHandler returns a handler for the catalog service.
func MakeHandler(svc GolfService, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	createGolfHandler := kithttp.NewServer(
		makeCreateGolfEndpoint(svc),
		decodeGolfReqCreate,
		encodeResponse,
		opts...,
	)

	getGolfHandler := kithttp.NewServer(
		makeGetGolfEndpoint(svc),
		decodeGolfReqGet,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/catalog", createGolfHandler).Methods("POST")
	r.Handle("/catalog/{id}", getGolfHandler).Methods("GET")

	return r
}

// Endcode/Decode
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Println("encodeResponse Called")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeGolfReqCreate(_ context.Context, r *http.Request) (interface{}, error) {
	fmt.Println("decode requests called")
	var req CreateGolfRequest
	body, _ := ioutil.ReadAll(r.Body)
	err := req.UnmarshalJSON(body)
	if err != nil {
		fmt.Printf("Error %+v", err)
		return nil, err
	}
	return req, nil
}

func decodeGolfReqGet(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetGolfRequest
	vars := mux.Vars(r)

	req = GetGolfRequest{
		id: vars["id"],
	}
	return req, nil
}

func (req *CreateGolfRequest) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["id"], &req.id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["name"], &req.name)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["location"], &req.location)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["amenities"], &req.amenities)
	if err != nil {
		return err
	}

	return nil
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	//case argo.ErrUnknown:
	//	w.WriteHeader(http.StatusNotFound)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
