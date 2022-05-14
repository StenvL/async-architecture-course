package repository

import (
	"github.com/StenvL/async-architecture-course/billing/app/model"
)

func (r Repository) MakePayments() ([]model.AccountSmall, error) {
	accounts, err := r.getPositiveBalancedAccounts()
	if err != nil {
		return nil, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	for _, acc := range accounts {
		if err = r.AddToAuditLogTx(tx, model.AuditLogRecord{
			AccountID:     acc.ID,
			Type:          model.PaymentMadeRecordType,
			BalanceChange: acc.Balance.Neg(),
		}); err != nil {
			_ = tx.Rollback()
			return nil, err
		}

		if err = r.ChangeAccountBalanceTx(tx, acc.ID, acc.Balance.Neg()); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	return accounts, tx.Commit()
}
