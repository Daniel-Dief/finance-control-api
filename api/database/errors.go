package database

import (
	"errors"

	"github.com/lib/pq"
)

// ErrNotFound is returned when a requested record does not exist in the
// database. Handlers can check for it with errors.Is to return a proper HTTP
// 404 response.
var ErrNotFound = errors.New("registro não encontrado")

// ErrDuplicateLogin is returned when trying to create a user whose login
// already exists.
var ErrDuplicateLogin = errors.New("Já existe um usuário com este login")

// Unique-violation sentinels, scoped per user.
var ErrDuplicateArea = errors.New("Já existe uma área com este nome")
var ErrDuplicateCategory = errors.New("Já existe uma categoria com este nome")
var ErrDuplicateBudget = errors.New("Já existe um orçamento para esta área neste período")

// isUniqueViolation reports whether err is a PostgreSQL unique-constraint
// violation (SQLSTATE 23505).
func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

var ErrGenericDatabase = errors.New("Falha ao consultar o banco de dados, em caso de persistencia contatar o suporte.")

var ErrProcessQuery = errors.New("Falha ao processar os resultados da consulta, em caso de persistencia contatar o suporte.")

var ErrBindQuery = errors.New("Falha indexar os resultados, em caso de persistencia contatar o suporte.")

func ErrRegisterObject(object string) error {
	return errors.New("Falha ao registrar um(a) " + object + " no banco de dados, em caso de persistencia contatar o suporte.")
}

func ErrUpdateObject(object string) error {
	return errors.New("Falha ao atualizar o(a) " + object + " no banco de dados, em caso de persistencia contatar o suporte.")
}

func ErrDeleteObject(object string) error {
	return errors.New("Falha ao deletar o(a) " + object + " no banco de dados, em caso de persistencia contatar o suporte.")
}
