package main

import (
	"testing"
)

// Stress test con valori casuali per controllare panic/crash (da fare sulle route "complete" per verificare che i dati non creino problemi da qualche parte)

func FuzzRegisterUserRoute(f *testing.F) {

	f.Add([]byte(`{"first_name":"Mario","last_name":"Rossi","birthday":"31-12","email":"mario@example.com","password":"Password%123!","password_confirmation":"Password%123!"}`))

	f.Fuzz(func(t *testing.T, firstName, lastName, birthday, email, password, passwordConf string) {
		reqStruct := makeRegisterUserReq(firstName, lastName, birthday, email, password, passwordConf)

		_ = makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", reqStruct)
	})
}

func FuzzLoginUserRoute(f *testing.F) {

	f.Add([]byte(`{"email":"mario@example.com","password":"Password%123!"}`))

	f.Fuzz(func(t *testing.T, email, password string) {
		reqStruct := makeLoginUserReq(email, password)

		_ = makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", reqStruct)
	})
}
