package repository

import (
	"main/infrastructure/database"
	rAppModule "main/repository/app-module"
	rApplication "main/repository/application"
	rApprovalRequestApprover "main/repository/approval-request-approver"
	rInvolvement "main/repository/involvement"
	rIpdrInvolvement "main/repository/ipdr-involvement"
	rIPDRequest "main/repository/ipdrequest"
	rItem "main/repository/item"
	rItemActivity "main/repository/item-activity"
	rLegalConsultation "main/repository/legal-consultation"
)

type Repository struct {
	Application             rApplication.ApplicationRepository
	ApplicationModule       rAppModule.ApplicationModuleRepository
	ApprovalRequestApprover rApprovalRequestApprover.ApprovalRequestApproverRepository
	Involvement             rInvolvement.InvolvementRepository
	IPDRequest              rIPDRequest.IpdRequestRepository
	IpdrInvolvement         rIpdrInvolvement.IpdrInvolvementRepository
	Item                    rItem.ItemRepository
	ItemActivity            rItemActivity.ItemActivityRepository
	LegalConsultation       rLegalConsultation.LegalConsultationRepository
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

func NewInvolvement(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.Involvement = rInvolvement.NewInvolvementRepository(db)
	}
}

func NewIPDRequest(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.IPDRequest = rIPDRequest.NewIpdRequestRepository(db)
	}
}

func NewIpdrInvolvement(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.IpdrInvolvement = rIpdrInvolvement.NewIpdrInvolvementRepository(db)
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

func NewItemActivity(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ItemActivity = rItemActivity.NewItemActivityRepository(db)
	}
}

func NewLegalConsultation(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.LegalConsultation = rLegalConsultation.NewLegalConsultationRepository(db)
	}
}
