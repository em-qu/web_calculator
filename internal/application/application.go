package application

import (
	"context"
   "log/slog"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"
	"encoding/json"
   "fmt"
   
	"github.com/em-qu/web_calculator/internal/config"
	"github.com/em-qu/web_calculator/internal/rpn"
)

// Data common to the entire application
type Application struct {
   log *slog.Logger  
}

// Constructor
func New() *Application {
	return &Application{}
}

// Application entry point
func (a *Application) Run() error {
	cfg := config.Load()
   a.log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	a.log.Info("starting web_calculator app")
	handler := http.NewServeMux()
	handler.HandleFunc("/api/v1/calculate", a.handlerCalculate())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      handler, 
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			a.log.Error("failed to start server")
		}
	}()

	a.log.Info(fmt.Sprintf("server started on %s", cfg.Address) )
	<-done
	a.log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.log.Error("failed to stop server")
		return err
	}
	a.log.Info("server stopped")
   return nil
}


// input expression
type Expression struct {
	Expr string `json:"expression"`
}
// output response
type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}
// Handler of endpoint "calculate"
func (a *Application) handlerCalculate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
      a.log.Info(fmt.Sprintf("connection from %v, method %s", r.RemoteAddr, r.Method))
      w.Header().Set("Content-Type", "application/json")
      var expr Expression 
      var resp Response   
      
      if r.Method == http.MethodPost || r.Method == http.MethodGet  { // Allowed http methods
         decoder := json.NewDecoder(r.Body)
         if err := decoder.Decode(&expr);  err != nil {  // Incorrect input json
            w.WriteHeader(http.StatusBadRequest)
            resp = Response{Result: "", Error: fmt.Sprintf("Bad request. %v", err)}
         } else { // Correct input json
            res, err := rpn.Calc(expr.Expr)
            if err != nil {  // Error when calculating expression
               resp = Response{Result: "", Error: err.Error()}
               w.WriteHeader(http.StatusUnprocessableEntity)
            } else { // Expression calculated successfully
               resp = Response{Result: fmt.Sprintf("%f", res), Error: ""}
               w.WriteHeader(http.StatusOK)
            }
         }
      } else { // Other http methods
         w.WriteHeader(http.StatusInternalServerError)
         resp = Response{Result: "", Error: "Internal server error"}
      }
      jsn, _ := json.Marshal(resp)
      w.Write(jsn)
   }  
}

