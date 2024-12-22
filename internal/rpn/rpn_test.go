package rpn

import (
   "testing"
)

func TestCalc(t *testing.T) {
   cases := []struct {
      expr string
      expected   float64
      expectedError  bool
   }{
      // valid
      {"2+4.7", 6.7, false},
      {"2-7", -5, false},
      {"0*7", 0, false},
      {"7/2", 3.5, false},
      {"1+2 * 5", 11, false},
      {"(1+2) * 5", 15, false},
      {"(1-5) / ((2+(0))*2)", -1, false},
      // invalid
      {"", 0, true}, 
      {"9/(5*0)", 0, true}, 
      {"(1-5) / ((2+(0)*2)", 0, true},
      {"3+O", 0, true},
      {`1\2`, 0, true},
   }

   for _, tc := range cases {
      result, err := Calc(tc.expr)
      if tc.expectedError != (err != nil) {
         t.Errorf("For expression %q expected error presence: %v, but got: %v", tc.expr, tc.expectedError, err)
      }
      if !tc.expectedError && (tc.expected != result){
         t.Errorf("For expression %q expected result %v, but got: %v", tc.expr, tc.expected, result)
      }
   }
}