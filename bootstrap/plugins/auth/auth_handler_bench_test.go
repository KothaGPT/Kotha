package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/khulnasoft/superkit/bootstrap/app/conf"
	"github.com/khulnasoft/superkit/bootstrap/app/db"
	"github.com/khulnasoft/superkit/kit"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupBenchmarkDB(b *testing.B) *gorm.DB {
	os.Setenv("SUPERKIT_SECRET", "12345678901234567890123456789012")
	os.Setenv("DB_DRIVER", "sqlite3")
	os.Setenv("DB_NAME", ":memory:")
	os.Setenv("SUPERKIT_AUTH_SKIP_VERIFY", "true")

	conf.Load()
	kit.Setup()

	if err := db.Initialize(); err != nil {
		b.Fatalf("failed to initialize db: %v", err)
	}

	gormDB := db.Get()
	if err := gormDB.AutoMigrate(&User{}, &Session{}); err != nil {
		b.Fatalf("failed to migrate: %v", err)
	}

	return gormDB
}

func BenchmarkAuthenticateUserWithSession(b *testing.B) {
	gormDB := setupBenchmarkDB(b)
	b.Cleanup(func() { db.Close() })

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := User{
		Email:        "bench@example.com",
		FirstName:    "Bench",
		LastName:     "User",
		PasswordHash: string(hash),
	}
	gormDB.Create(&user)

	session := Session{
		UserID:    user.ID,
		Token:     "bench-session-token",
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	gormDB.Create(&session)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	k := &kit.Kit{
		Response: w,
		Request:  req,
	}
	sess := k.GetSession(userSessionName)
	sess.Values["sessionToken"] = "bench-session-token"
	sess.Save(req, w)
	for _, cookie := range w.Result().Cookies() {
		req.AddCookie(cookie)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := &kit.Kit{
			Response: w,
			Request:  req,
		}
		auth, err := AuthenticateUser(k)
		assert.NoError(b, err)
		assert.True(b, auth.Check())
	}
}

func BenchmarkAuthenticateUserSQLWithSession(b *testing.B) {
	gormDB := setupBenchmarkDB(b)
	b.Cleanup(func() { db.Close() })

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := User{
		Email:        "bench@example.com",
		FirstName:    "Bench",
		LastName:     "User",
		PasswordHash: string(hash),
	}
	gormDB.Create(&user)

	session := Session{
		UserID:    user.ID,
		Token:     "bench-session-token",
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	gormDB.Create(&session)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	k := &kit.Kit{
		Response: w,
		Request:  req,
	}
	sess := k.GetSession(userSessionName)
	sess.Values["sessionToken"] = "bench-session-token"
	sess.Save(req, w)
	for _, cookie := range w.Result().Cookies() {
		req.AddCookie(cookie)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := &kit.Kit{
			Response: w,
			Request:  req,
		}
		auth, err := AuthenticateUserSQL(k)
		assert.NoError(b, err)
		assert.True(b, auth.Check())
	}
}
