package database

import (
	"context"
	"database/sql"
	"errors"
	"finance-control-api/models"
	"log"
)

// ListCategories retrieves a list of categories from the database, optionally filtered by name.
// If name is empty, all categories are returned.
func ListCategories(name *string) ([]models.Category, error) {
	query := `
		SELECT "Id", "Name"
		FROM "Categories"
		WHERE ($1::text IS NULL OR "Name" ILIKE '%' || $1 || '%')
	`

	rows, err := Pool.QueryContext(context.Background(), query, name)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, errors.New("Falha ao consultar o banco de dados, em caso de persistencia contatar o suporte.")
	} else if rows.Err() != nil {
		log.Println("Error with rows:", rows.Err())
		return nil, errors.New("Falha ao processar os resultados da consulta, em caso de persistencia contatar o suporte.")
	}
	defer rows.Close()

	var result []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			log.Println("Error scanning row:", err)
			return nil, errors.New("Falha indexar os resultados, em caso de persistencia contatar o suporte.")
		}
		result = append(result, c)
	}

	return result, nil
}

// GetCategoryByID retrieves a category from the database by its ID.
func GetCategoryByID(id int) (models.Category, error) {
	query := `
		SELECT "Id", "Name"
		FROM "Categories"
		WHERE "Id" = $1
	`

	var c models.Category
	err := Pool.QueryRowContext(context.Background(), query, id).Scan(&c.ID, &c.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Category{}, nil
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Category{}, errors.New("Falha ao consultar o banco de dados, em caso de persistencia contatar o suporte.")
	}

	return c, nil
}

// CreateCategory inserts a new category into the database and returns the created category.
func CreateCategory(name string) (models.Category, error) {
	query := `
		INSERT INTO "Categories" ("Name")
		VALUES ($1)
		RETURNING "Id", "Name"
	`

	var c models.Category
	err := Pool.QueryRowContext(context.Background(), query, name).Scan(&c.ID, &c.Name)
	if err != nil {
		log.Println("Error executing query:", err)
		return models.Category{}, errors.New("Falha ao registrar a categoria no banco de dados, em caso de persistencia contatar o suporte.")
	}

	return c, nil
}

// UpdateCategory updates the name of an existing category in the database and returns the updated category.
func UpdateCategory(id int, name string) (models.Category, error) {
	query := `
		UPDATE "Categories"
		SET "Name" = $2
		WHERE "Id" = $1
		RETURNING "Id", "Name"
	`

	var c models.Category
	err := Pool.QueryRowContext(context.Background(), query, id, name).Scan(&c.ID, &c.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Category{}, errors.New("Falha ao atualizar a categoria: categoria não encontrada")
	} else if err != nil {
		log.Println("Error executing query:", err)
		return models.Category{}, errors.New("Falha ao atualizar a categoria no banco de dados, em caso de persistencia contatar o suporte.")
	}

	return c, nil
}

// DeleteCategory removes a category from the database by its ID.
func DeleteCategory(id int) error {
	query := `
		DELETE FROM "Categories"
		WHERE "Id" = $1
	`

	res, err := Pool.ExecContext(context.Background(), query, id)
	if err != nil {
		log.Println("Error executing query:", err)
		return errors.New("Falha ao deletar a categoria no banco de dados, em caso de persistencia contatar o suporte.")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error executing query:", err)
		return errors.New("Falha ao deletar a categoria no banco de dados, em caso de persistencia contatar o suporte.")
	}

	if rowsAffected == 0 {
		return errors.New("Categoria não encontrada")
	}

	return nil
}
