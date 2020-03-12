package stores

import "github.com/lib/pq"

func IsNeedVacuumFullErr(err error) bool {
	if pqError, ok := err.(*pq.Error); ok {
		return pqError.Code == "42710" || pqError.Code == "22000"
	}

	return false
}

func IsDuplicate(err error) bool {
	if pqError, ok := err.(*pq.Error); ok {
		return pqError.Code == "23505"
	}

	return false
}
