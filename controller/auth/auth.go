package auth

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "time"
    "strings"

    "product/database"
    "product/model/user"

    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type Credentials struct {
    Name      string    `json:"name"`
    NPM      int    `json:"npm"`
    Password string `json:"password"`
    Role string `json:"role"`
}

type Claims struct {
    NPM int `json:"npm"`
    jwt.StandardClaims
}

func Register(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Periksa apakah NPM sudah ada di database
    var existingUser user.User
    err = database.DB.QueryRow("SELECT id FROM users WHERE npm=$1", creds.NPM).Scan(&existingUser.ID)
    if err != nil && err != sql.ErrNoRows {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    if existingUser.ID != 0 {
        http.Error(w, "NPM already exists", http.StatusBadRequest)
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Internal server error Password", http.StatusInternalServerError)
        return
    }

    // Simpan pengguna baru ke database
    _, err = database.DB.Exec("INSERT INTO users (name,npm, password,role,created_at) VALUES ($1, $2,$3,$4, NOW())", creds.Name,creds.NPM ,hashedPassword,creds.Role)
    if err != nil {
        http.Error(w, "Internal server error Insert", http.StatusInternalServerError)
        return
    }

    // Berikan respon sukses
    w.Header().Set("Content-Type", "application/json")
    response := map[string]interface{}{
        "message": "User registered successfully",
    }
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}

func Login(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var user user.User
    err = database.DB.QueryRow("SELECT id, npm, password FROM users WHERE npm=$1", creds.NPM).Scan(&user.ID, &user.NPM, &user.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
    if err != nil {
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(60 * time.Minute)
    claims := &Claims{
        NPM: creds.NPM,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    response := map[string]interface{}{
        "message": "Login successful",
        "token":   tokenString,
    }
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}

    func ValidateToken(tokenString string) (bool, error) {
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil {
            return false, err
        }

        return token.Valid, nil
    }

    func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            bearerToken := r.Header.Get("Authorization")
            sttArr := strings.Split(bearerToken, " ")
            if len(sttArr) == 2 {
                isValid, _ := ValidateToken(sttArr[1])
                if isValid {
                    next.ServeHTTP(w, r)
                } else {
                    http.Error(w, "Unauthorized", http.StatusUnauthorized)
                }
            } else {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
            }
        })
    }
