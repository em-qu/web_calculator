package application

import (
   "testing"
   "net/http"
   "net/http/httptest"
   "encoding/json"
   "strconv"
   "bytes"
   "log/slog"
   "io"
)

func TestHandlerCalculate(t *testing.T) {
   cases := []struct {
      name           string
      method         string
      expr           string
      expectedStatus   int
      expectedResult float64
   }{
      {"Valid case", http.MethodPost, "(1+3) / ((2+0)*2)", http.StatusOK, 1 },
   // {"Wrong answer simulation", http.MethodPost, "2-4", http.StatusOK, 2 },
      {"Wrong expression", http.MethodPost, "(1+3 / (2a*2)", http.StatusUnprocessableEntity, 0 },
      {"Division by zero", http.MethodPost, "3/(2-2)", http.StatusUnprocessableEntity, 0 },
      {"Internal server error", http.MethodDelete, "(1+3) / ((2+0)*2)", http.StatusInternalServerError, 0 },
   }

   app := New()
   app.log = slog.New(slog.NewTextHandler(io.Discard, nil)) //easy way to suppress logger without implementing mock
   handlerFunc := app.handlerCalculate()
   for _, tc := range cases {
      t.Run(tc.name, func(t *testing.T) {
         jsn, _ := json.Marshal(Expression{tc.expr})
         req, err := http.NewRequest(tc.method, endpoint_calculate, bytes.NewReader(jsn))
         if err != nil  { t.Fatal(err) }
         rr := httptest.NewRecorder()
         handlerFunc(rr, req)
         if rr.Code != tc.expectedStatus {
            t.Errorf("Wrong status code: expected %v , got %v", tc.expectedStatus, rr.Code)
         }
         var resp Response
         if err := json.NewDecoder(rr.Body).Decode(&resp);  err != nil {  // Incorrect json
            t.Errorf("Error decoding output json - %v. Response body:\n%v", err, rr.Body.String())
         }
         if rr.Code == http.StatusOK { // Result returned
            if resp.Error != "" {
               t.Errorf("Expected valid result %v , but returned error %v", tc.expectedResult, resp.Error)   
            }
            result, err := strconv.ParseFloat(resp.Result, 64)
            if err != nil {
               t.Errorf("Error parsing result %v - %v", resp.Result, err)   
            } else if result != tc.expectedResult {
               t.Errorf("Wrong result: expected %v , got %v. Expression %v", tc.expectedResult, result, tc.expr)
            }
         } else { // Error returned
            if resp.Result !="" || resp.Error == ""{
               t.Errorf("Expected error, but returned something other. Response body:\n%v", rr.Body.String())
            }
         }
      })
   }
}
