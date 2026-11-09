package day3

// JWT in Go
// A commonly used package is:

// golang-jwt/jwt
// Generating a JWT
// Simplified example:


package main

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

func generateToken(
    userID string,
    secret []byte,
) (string, error) {

    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(
            15 * time.Minute,
        ).Unix(),
    }

    token := jwt.NewWithClaims(
        jwt.SigningMethodHS256,
        claims,
    )

    return token.SignedString(secret)
}

func main() {

    token, err := generateToken(
        "user-123",
        []byte("replace-with-a-secret"),
    )

    if err != nil {
        panic(err)
    }

    fmt.Println(token)
}


// This is an educational example. In an actual application, the secret must come from secure configuration, not source code.
// 5. Verifying a JWT

func verifyToken(
    tokenString string,
    secret []byte,
) (string, error) {

    token, err := jwt.Parse(
        tokenString,
        func(token *jwt.Token) (any, error) {

            if token.Method != jwt.SigningMethodHS256 {
                return nil, fmt.Errorf("unexpected signing method")
            }

            return secret, nil
        },
    )

    if err != nil {
        return "", err
    }

    claims, ok := token.Claims.(jwt.MapClaims)

    if !ok || !token.Valid {
        return "", fmt.Errorf("invalid token")
    }

    userID, ok := claims["sub"].(string)

    if !ok {
        return "", fmt.Errorf("missing subject")
    }

    return userID, nil
}


// For production, also enforce issuer, audience, expiry, appropriate signing algorithms and key rotation policies.

// Password hashing
// Never store passwords as plaintext.

// Use a password hashing algorithm designed for password storage, such as Argon2id or bcrypt.

// A popular Go package is:
// go get golang.org/x/crypto/bcrypt


Example:
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password),
    bcrypt.DefaultCost,
)


Verify:
err := bcrypt.CompareHashAndPassword(
    hashedPassword,
    []byte(password),
)

if err != nil {
    // Invalid password
}