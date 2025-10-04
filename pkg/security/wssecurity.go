package security

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"time"
)

// WSSecurity represents WS-Security configuration
type WSSecurity struct {
	Username string
	Password string
	UseDigest bool
}

// SecurityHeader represents the WS-Security header
type SecurityHeader struct {
	XMLName   xml.Name        `xml:"wsse:Security"`
	WSSE      string          `xml:"xmlns:wsse,attr"`
	WSU       string          `xml:"xmlns:wsu,attr"`
	Timestamp *Timestamp      `xml:"wsu:Timestamp,omitempty"`
	UsernameToken *UsernameToken `xml:"wsse:UsernameToken,omitempty"`
}

// Timestamp represents WS-Security timestamp
type Timestamp struct {
	XMLName xml.Name `xml:"wsu:Timestamp"`
	Created string   `xml:"wsu:Created"`
	Expires string   `xml:"wsu:Expires"`
}

// UsernameToken represents WS-Security username token
type UsernameToken struct {
	XMLName  xml.Name  `xml:"wsse:UsernameToken"`
	Username string    `xml:"wsse:Username"`
	Password *Password `xml:"wsse:Password"`
	Nonce    *Nonce    `xml:"wsse:Nonce,omitempty"`
	Created  string    `xml:"wsu:Created,omitempty"`
}

// Password represents the password element
type Password struct {
	XMLName xml.Name `xml:"wsse:Password"`
	Type    string   `xml:"Type,attr"`
	Value   string   `xml:",chardata"`
}

// Nonce represents a nonce value
type Nonce struct {
	XMLName      xml.Name `xml:"wsse:Nonce"`
	EncodingType string   `xml:"EncodingType,attr"`
	Value        string   `xml:",chardata"`
}

// NewSecurityHeader creates a new WS-Security header
func NewSecurityHeader(ws *WSSecurity) *SecurityHeader {
	if ws == nil {
		return nil
	}

	header := &SecurityHeader{
		WSSE: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
		WSU:  "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
	}

	// Add timestamp
	header.Timestamp = createTimestamp()

	// Add username token
	if ws.Username != "" {
		header.UsernameToken = createUsernameToken(ws)
	}

	return header
}

// createTimestamp creates a timestamp element
func createTimestamp() *Timestamp {
	now := time.Now().UTC()
	expires := now.Add(5 * time.Minute)

	return &Timestamp{
		Created: now.Format(time.RFC3339),
		Expires: expires.Format(time.RFC3339),
	}
}

// createUsernameToken creates a username token element
func createUsernameToken(ws *WSSecurity) *UsernameToken {
	token := &UsernameToken{
		Username: ws.Username,
	}

	if ws.UseDigest {
		// Create nonce
		nonce := generateNonce()
		created := time.Now().UTC().Format(time.RFC3339)

		// Create password digest
		digest := createPasswordDigest(nonce, created, ws.Password)

		token.Password = &Password{
			Type:  "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordDigest",
			Value: digest,
		}
		token.Nonce = &Nonce{
			EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary",
			Value:        base64.StdEncoding.EncodeToString(nonce),
		}
		token.Created = created
	} else {
		// Plain text password
		token.Password = &Password{
			Type:  "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText",
			Value: ws.Password,
		}
	}

	return token
}

// generateNonce generates a random nonce
func generateNonce() []byte {
	nonce := make([]byte, 16)
	rand.Read(nonce)
	return nonce
}

// createPasswordDigest creates a password digest
// Digest = Base64(SHA1(nonce + created + password))
func createPasswordDigest(nonce []byte, created, password string) string {
	h := sha1.New()
	h.Write(nonce)
	h.Write([]byte(created))
	h.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
