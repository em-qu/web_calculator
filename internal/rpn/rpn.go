package rpn

import (
    "fmt"
    "strings"
    "strconv"
)

// разбиваем строку на токены
func split_to_tokens(str string) []string {
    var token strings.Builder
    var tokens []string
    for _, ch := range str {
        switch ch {
        case ' ':
            continue // пробелы пропускаем
        case '(', ')', '+', '-', '/', '*': //токен закончился
            if token.Len() > 0 { 
                tokens = append(tokens, token.String())
                token.Reset()
            }
            tokens = append(tokens, string(ch))
        default: // цифра
            token.WriteRune(ch)
        }
    }
    if token.Len() > 0 {  // добавляем последний токен если есть
        tokens = append(tokens, token.String())
    }
    return tokens
}

// может ли токен быть конвертирован к float64
func is_number(token string) bool {
    _, err := strconv.ParseFloat(token, 64)
    if err != nil {
        return false
    }
    return true
}

// являетс ли токен операцией
func is_operation(token string) bool {
    if token == "+" || token == "-" || token == "/" || token == "*" { 
        return true
    }
    return false 
}

// возвращает приоритет операции (больше приоритет - выполняется первым)
func priority(operation string) int {
    switch operation {
    case "+", "-":
        return 1
    case "/", "*":
        return 2
    default:
        return 0
    }
}


// Преобразует слайс(массив) токенов в обратную польскую запись (Reverse Polish notation)
func rPN(tokens []string) ([]string, error) {
    var rpn []string
    var stack []string

    for _, token := range tokens {
        if is_number(token) {  
            rpn = append(rpn, token)
        
        } else if token == "(" {
            stack = append(stack, token)

        } else if token == ")" {
            for len(stack) > 0 && stack[len(stack)-1] != "(" {
                rpn = append(rpn, stack[len(stack)-1])
                stack = stack[:len(stack)-1]
            }
            if len(stack) == 0  {
                return nil, fmt.Errorf("Parenthesis mismatch")
            }
            stack = stack[:len(stack)-1]

        } else if is_operation(token) {
            for len(stack) > 0 && priority(stack[len(stack)-1]) >= priority(token) {
                rpn = append(rpn, stack[len(stack)-1])
                stack = stack[:len(stack)-1]
            }
            stack = append(stack, token)

        } else {
            return nil, fmt.Errorf("Invalid token %s", token)
        }
    }

    for len(stack) > 0 {
        if stack[len(stack)-1] == "(" {
            return nil, fmt.Errorf("Parenthesis mismatch")
        }
        rpn = append(rpn, stack[len(stack)-1])
        stack = stack[:len(stack)-1]
    }

    return rpn, nil
}

// Вычисляет значение выражения по обратной польской записи
func calc_rpn(rpn []string) (float64, error) {

    var stack []float64

    for _, token := range rpn {
        if is_number(token) {
            f, _ := strconv.ParseFloat(token, 64)
            stack = append(stack, f)
            
        } else if is_operation(token) {
            if len(stack) < 2 {
                return 0, fmt.Errorf("Incorrect expression near operation %s", token)
            }
            y := stack[len(stack)-1]
            x := stack[len(stack)-2]
            stack = stack[:len(stack)-2]
            switch token {
            case "+":
                stack = append(stack, x + y)
            case "-":
                stack = append(stack, x - y)
            case "*":
                stack = append(stack, x * y)
            case "/":
                if y == 0 {
                    return 0, fmt.Errorf("Division by zero")
                }
                stack = append(stack, x / y)
            }
            
        } else  {
            return 0, fmt.Errorf("Incorrect token: %s", token)
        }
    }

    if len(stack) != 1 {
        return 0, fmt.Errorf("Incorrect expression")
    }

    return stack[0], nil
}


// Вычисляет значение выражения
func Calc(str string) (float64, error) {
    tokens := split_to_tokens(str)
    rpn, err := rPN(tokens)
    if err != nil {
        return 0, err
    }
    return calc_rpn(rpn)
}