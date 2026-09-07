package auth

import (
	"github.com/khulnasoft/superkit/bootstrap/app/db"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/khulnasoft/superkit/kit"
	v "github.com/khulnasoft/superkit/validate"
	"golang.org/x/crypto/bcrypt"
)

const (
	userSessionName = "user-session"
)

var authSchema = v.Schema{
	"email":    v.Rules(v.Email),
	"password": v.Rules(v.Required),
}

func HandleLoginIndex(kit *kit.Kit) error {
	if kit.Auth().Check() {
		redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")
		return kit.Redirect(http.StatusSeeOther, redirectURL)
	}
	return kit.Render(LoginIndex(LoginIndexPageData{}))
}

func HandleLoginCreate(kit *kit.Kit) error {
	var values LoginFormValues
	errors, ok := v.Request(kit.Request, &values, authSchema)
	if !ok {
		return kit.Render(LoginForm(values, errors))
	}

	var user User
	err := db.GetSQL().QueryRow(
		"SELECT id, email, password_hash, email_verified_at FROM users WHERE email = ?",
		values.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerifiedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			errors.Add("credentials", "invalid credentials")
			return kit.Render(LoginForm(values, errors))
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(values.Password))
	if err != nil {
		errors.Add("credentials", "invalid credentials")
		return kit.Render(LoginForm(values, errors))
	}

	skipVerify := kit.Getenv("SUPERKIT_AUTH_SKIP_VERIFY", "false")
	if skipVerify != "true" {
		if !user.EmailVerifiedAt.Valid {
			errors.Add("verified", "please verify your email")
			return kit.Render(LoginForm(values, errors))
		}
	}

	sessionExpiryStr := kit.Getenv("SUPERKIT_AUTH_SESSION_EXPIRY_IN_HOURS", "48")
	sessionExpiry, err := strconv.Atoi(sessionExpiryStr)
	if err != nil {
		sessionExpiry = 48
	}
	session := Session{
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(sessionExpiry)),
	}
	if err = db.Get().Create(&session).Error; err != nil {
		return err
	}

	sess := kit.GetSession(userSessionName)
	sess.Values["sessionToken"] = session.Token
	sess.Save(kit.Request, kit.Response)
	redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")

	return kit.Redirect(http.StatusSeeOther, redirectURL)
}

func HandleLoginDelete(kit *kit.Kit) error {
	sess := kit.GetSession(userSessionName)
	defer func() {
		sess.Values = map[any]any{}
		sess.Save(kit.Request, kit.Response)
	}()
	token, _ := sess.Values["sessionToken"].(string)
	if token != "" {
		db.GetSQL().Exec("DELETE FROM sessions WHERE token = ?", token)
	}
	return kit.Redirect(http.StatusSeeOther, "/")
}

func HandleEmailVerify(kit *kit.Kit) error {
	tokenStr := kit.Request.URL.Query().Get("token")
	if len(tokenStr) == 0 {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	token, err := jwt.ParseWithClaims(
		tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SUPERKIT_SECRET")), nil
		}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}
	if !token.Valid {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	var user User
	err = db.Get().First(&user, userID).Error
	if err != nil {
		return err
	}

	if user.EmailVerifiedAt.Time.After(time.Time{}) {
		return kit.Render(EmailVerificationError("Email already verified"))
	}

	now := sql.NullTime{Time: time.Now(), Valid: true}
	user.EmailVerifiedAt = now
	err = db.Get().Save(&user).Error
	if err != nil {
		return err
	}

	return kit.Redirect(http.StatusSeeOther, "/login")
}

func AuthenticateUser(kit *kit.Kit) (kit.Auth, error) {
	return AuthenticateUserSQL(kit)
}

func AuthenticateUserSQL(kit *kit.Kit) (kit.Auth, error) {
	auth := Auth{}
	sess := kit.GetSession(userSessionName)
	token, ok := sess.Values["sessionToken"]
	if !ok {
		return auth, nil
	}

	sqlDB := db.GetSQL()
	var userID uint
	err := sqlDB.QueryRow(
		"SELECT user_id FROM sessions WHERE token = ? AND expires_at > ?",
		token, time.Now(),
	).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth, nil
		}
		return auth, err
	}
	if userID == 0 {
		return auth, nil
	}

	var email, firstName, lastName string
	err = sqlDB.QueryRow(
		"SELECT id, email, first_name, last_name FROM users WHERE id = ?",
		userID,
	).Scan(&userID, &email, &firstName, &lastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth, nil
		}
		return auth, err
	}

	return Auth{
		LoggedIn:  true,
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}
