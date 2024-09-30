package helper

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func IntParam(r *http.Request, key string) (int, error) {
	return strconv.Atoi(chi.URLParam(r, key))
}

func Float64Param(r *http.Request, key string) (float64, error) {
	v, err := strconv.ParseFloat(chi.URLParam(r, key), 64)
	return float64(v), err
}

func BoolParam(r *http.Request, key string) (bool, error) {
	v, err := strconv.ParseBool(chi.URLParam(r, key))
	return bool(v), err
}

func QueryStringArray(r *http.Request, key string) ([]string, error) {
	values, ok := r.URL.Query()[key]
	if !ok {
		return nil, fmt.Errorf("invalid params")
	}
	return values, nil
}

func QueryIntArray(r *http.Request, key string) ([]int, error) {
	values, err := QueryStringArray(r, key)
	if err != nil {
		return nil, err
	}

	res := make([]int, len(values))
	for i, value := range values {
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func QueryFloat64Array(r *http.Request, key string) ([]float64, error) {
	values, err := QueryStringArray(r, key)
	if err != nil {
		return nil, err
	}

	res := make([]float64, len(values))
	for i, value := range values {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		res[i] = float64(v)
	} 
	return res, nil
}

func VadlidateID(id string) bool {
	v, err := strconv.Atoi(id)
	if err != nil || id == "" {
		return false
	}
	return v > 0
}
