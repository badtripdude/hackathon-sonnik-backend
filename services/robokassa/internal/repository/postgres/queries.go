package postgres

const (
	savePaymentQuery = `
	INSERT INTO payments (id, user_id, amount, currency, status)
	VALUES (:id, :user_id, :amount, :currency, :status)
	RETURNING created_at, updated_at;
	`

	getPaymentByIDQuery = `
	SELECT id, user_id, amount, currency, status, created_at, updated_at
	FROM payments
	WHERE id = $1;
	`

	getPaymentsByUserIDQuery = `
	SELECT id, user_id, amount, currency, status, created_at, updated_at
	FROM payments
	WHERE user_id = $1
	ORDER BY created_at DESC;
	`

	updatePaymentStatusQuery = `
	UPDATE payments
	SET status = $1, updated_at = NOW()
	WHERE id = $2;
	`
)
