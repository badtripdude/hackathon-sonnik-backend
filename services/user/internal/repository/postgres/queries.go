package postgres

const (
	// ---- Users ----
	saveUserQuery = `
	INSERT INTO users (email, password_hash, username, birth_date)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	getUserByEmailQuery = `
	SELECT id, email, password_hash, created_at
	FROM users
	WHERE email = $1;
	`

	getUserByIDQuery = `
	SELECT id, email, password_hash, created_at
	FROM users
	WHERE id = $1;
	`

	deleteUserByIDQuery = `
	DELETE FROM users
	WHERE id = $1;
	`

	// ---- Refresh Tokens ----
	saveRefreshTokenQuery = `
	INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
	VALUES (:id, :user_id, :token_hash, :expires_at, :created_at);
	`

	getRefreshTokenByTokenHashQuery = `
	SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
	FROM refresh_tokens
	WHERE token_hash = $1;
	`

	revokeRefreshTokenByTokenHashQuery = `
	UPDATE refresh_tokens
	SET revoked_at = $1
	WHERE token_hash = $2;
	`

	revokeAllRefreshTokenByUserIdQuery = `
	UPDATE refresh_tokens
	SET revoked_at = $1
	WHERE user_id = $2;
	`
)
