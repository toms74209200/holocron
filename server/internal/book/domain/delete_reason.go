package domain

import "errors"

var ErrInvalidDeleteReason = errors.New("delete reason must be one of: transfer, disposal, lost, other")

type DeleteReason string

const (
	DeleteReasonTransfer DeleteReason = "transfer"
	DeleteReasonDisposal DeleteReason = "disposal"
	DeleteReasonLost     DeleteReason = "lost"
	DeleteReasonOther    DeleteReason = "other"
)

func ParseDeleteReason(s string) (DeleteReason, error) {
	switch s {
	case string(DeleteReasonTransfer):
		return DeleteReasonTransfer, nil
	case string(DeleteReasonDisposal):
		return DeleteReasonDisposal, nil
	case string(DeleteReasonLost):
		return DeleteReasonLost, nil
	case string(DeleteReasonOther):
		return DeleteReasonOther, nil
	default:
		return "", ErrInvalidDeleteReason
	}
}
