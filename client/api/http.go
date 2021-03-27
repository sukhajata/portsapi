package api

import (
	"encoding/json"
	"fmt"
	"github.com/sukhajata/portsapi/client/internal/portdomain"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/rs/cors"
	pb "github.com/sukhajata/portsapi/service/pkg/proto"
)

type HTTPServer struct {
	Ready bool
	Live  bool
	portDomainClient *portdomain.PortDomainClient
}

// readinessHandler reports ready status
func (s *HTTPServer) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if !s.Ready {
		http.Error(w, "Not ready", http.StatusInternalServerError)
		return
	}
	_, err := fmt.Fprintf(w, "Ready")
	if err != nil {
		fmt.Println(err.Error())
	}
}

// livenessHandler reports live status
func (s *HTTPServer) livenessHandler(w http.ResponseWriter, r *http.Request) {
	if !s.Live {
		http.Error(w, "Not live", http.StatusInternalServerError)
		return
	}
	_, err := fmt.Fprintf(w, "Live")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (s *HTTPServer) getPortHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "missing connection id", http.StatusBadRequest)
		return
	}

	port, err := s.portDomainClient.GetPort(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *HTTPServer) postPortHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var content pb.UpsertPortRequest
	err := decoder.Decode(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	port, err := s.portDomainClient.UpsertPort(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		fmt.Println(err)
	}

}

func NewHTTPServer(portDomainClient *portdomain.PortDomainClient) *HTTPServer {
	server := &HTTPServer{
		portDomainClient: portDomainClient,
		Live:        true,
		Ready:       true,
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		//Debug: true,
	})

	router := mux.NewRouter()

	router.HandleFunc("/health/ready", server.readinessHandler).Methods("GET")
	router.HandleFunc("/health/live", server.livenessHandler).Methods("GET")
	router.HandleFunc("/ports/{id}", server.getPortHandler).Methods("GET")
	router.HandleFunc("/ports", server.postPortHandler).Methods("PUT", "POST")


	n := negroni.New()
	n.Use(negroni.NewRecovery())

	n.UseHandler(router)
	n.Use(c)

	// start on new goroutine
	go func() {
		n.Run(":80")
	}()

	return server
}
