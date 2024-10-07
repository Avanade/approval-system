package repository

import (
	"main/infrastructure/database"
	rAppModule "main/repository/app-module"
	rApplication "main/repository/application"
	rApprovalRequestApprover "main/repository/approval-request-approver"
	rItem "main/repository/item"
)

type Repository struct {
	Application             rApplication.ApplicationRepository
	ApplicationModule       rAppModule.ApplicationModuleRepository
	ApprovalRequestApprover rApprovalRequestApprover.ApprovalRequestApproverRepository
	Item                    rItem.ItemRepository
}

type RepositoryOptionFunc func(*Repository)

func NewRepository(opts ...RepositoryOptionFunc) *Repository {
	repo := &Repository{}

	for _, opt := range opts {
		opt(repo)
	}

	return repo
}

func NewApplication(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.Application = rApplication.NewApplicationRepository(db)
	}
}

func NewApplicationModule(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ApplicationModule = rAppModule.NewApplicationModuleRepository(db)
	}
}

func NewItem(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.Item = rItem.NewItemRepository(db)
	}
}

func NewApprovalRequestApprover(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ApprovalRequestApprover = rApprovalRequestApprover.NewApprovalRequestApproverRepository(db)
	}
}
