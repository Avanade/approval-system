package repository

import (
	"main/infrastructure/database"
	rAppModule "main/repository/app-module"
	rItem "main/repository/item"
)

type Repository struct {
	ApplicationModule rAppModule.ApplicationModuleRepository
	Item              rItem.ItemRepository
}

type RepositoryOptionFunc func(*Repository)

func NewRepository(opts ...RepositoryOptionFunc) *Repository {
	repo := &Repository{}

	for _, opt := range opts {
		opt(repo)
	}

	return repo
}

func NewApplicationModule(db database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ApplicationModule = rAppModule.NewApplicationModuleRepository(db)
	}
}

func NewItem(db database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.Item = rItem.NewItemRepository(db)
	}
}
