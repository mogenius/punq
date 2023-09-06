package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const SEC_KEY_PAIR = "keyPair"
const TOKEN_EXP_HOURS = 24 * 7 // 1 week

var KeyPairInstance *KeyPair

type KeyPair struct {
	PrivateKeyString string `json:"privateKey" validate:"required"`
	PublicKeyString  string `json:"publicKey" validate:"required"`

	PrivateKey *ecdsa.PrivateKey `json:"-"`
	PublicKey  any               `json:"-"`
}

type keyPairAlias KeyPair

type PunqClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

func (keyPair *KeyPair) UnmarshalJSON(data []byte) error {
	temp := keyPairAlias{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	*keyPair = KeyPair(temp)

	// extraction PEM block
	block, _ := pem.Decode([]byte(keyPair.PrivateKeyString))
	if block == nil {
		logger.Log.Error("failed extraction private key PEM block")
	}

	// recovering PEM block
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		logger.Log.Errorf("failed recovering private key PEM block %s", err)
	}
	keyPair.PrivateKey = privateKey

	// extraction PEM block
	block, _ = pem.Decode([]byte(keyPair.PublicKeyString))
	if block == nil {
		logger.Log.Error("failed extraction public key PEM block")
	}

	// recovering PEM block
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Log.Errorf("failed recovering public key PEM block %s", err)
	}
	keyPair.PublicKey = publicKey

	return nil
}

func generateAuthKeyPair() (*KeyPair, error) {
	keyPair := KeyPair{}

	// generate private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		logger.Log.Errorf("failed to generate private key %s", err)
		return nil, err
	}

	// generate private PEM
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		logger.Log.Errorf("failed to generate private key PEM %s", err)
		return nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	keyPair.PrivateKeyString = string(privateKeyPEM)

	// generate public key
	publicKey := &privateKey.PublicKey

	// generate public PEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		logger.Log.Errorf("failed to generate public key %s", err)
		return nil, err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	keyPair.PublicKeyString = string(publicKeyPEM)
	keyPair.PrivateKey = privateKey

	return &keyPair, nil
}

func InitAuthService() {
	CreateKeyPair()
}

func CreateKeyPair() *KeyPair {
	provider := kubernetes.NewKubeProvider()
	if provider == nil {
		logger.Log.Fatal("Failed to load kubeprovider.")
	}

	secretClient := provider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)
	existingSecret, getErr := secretClient.Get(context.TODO(), utils.JWTSECRET, metav1.GetOptions{})

	secret := utils.InitSecret()
	secret.ObjectMeta.Name = utils.JWTSECRET
	secret.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
	delete(secret.StringData, "exampleData") // delete example data

	keyPair, err := generateAuthKeyPair()
	if err != nil {
		logger.Log.Errorf("failed to generate key-pair %s", err)
	}
	rawKeyPair, err := json.Marshal(keyPair)
	if err != nil {
		logger.Log.Errorf("Error marshaling %s", err)
	}
	secret.StringData[SEC_KEY_PAIR] = string(rawKeyPair)

	// if not exist
	if existingSecret == nil || getErr != nil {
		fmt.Println("Creating new punq-auth secret ...")
		_, err := secretClient.Create(context.TODO(), &secret, kubernetes.MoCreateOptions())
		if err != nil {
			logger.Log.Error(err)
			return nil
		}
		fmt.Println("Created new punq-auth secret. ✅")

		return keyPair
	}

	// get from secret
	return GetKeyPair()
}

func RemoveKeyPair() {
	provider := kubernetes.NewKubeProvider()
	if provider == nil {
		logger.Log.Fatal("Failed to load kubeprovider.")
	}

	secretClient := provider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	fmt.Printf("Deleting %s/%s secret ...\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.JWTSECRET)
	deletePolicy := metav1.DeletePropagationForeground
	err := secretClient.Delete(context.TODO(), utils.JWTSECRET, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	fmt.Printf("Deleted %s/%s secret. ✅\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.JWTSECRET)
}

func GetKeyPair() *KeyPair {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.JWTSECRET)
	if secret == nil {
		logger.Log.Warningf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.JWTSECRET)
		return CreateKeyPair()
	}

	if secret.Data[SEC_KEY_PAIR] != nil {
		keyPair := KeyPair{}
		err := json.Unmarshal([]byte(secret.Data[SEC_KEY_PAIR]), &keyPair)
		if err != nil {
			logger.Log.Errorf("Failed to Unmarshal user '%s' %s.", SEC_KEY_PAIR, err)
		}
		return &keyPair
	}

	return nil
}

func GenerateToken(user *dtos.PunqUser) (*string, error) {
	if KeyPairInstance == nil {
		KeyPairInstance = GetKeyPair()
	}

	claims := jwt.MapClaims{}
	claims["accessLevel"] = user.AccessLevel
	claims["userId"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(TOKEN_EXP_HOURS)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)

	// sign JWT-Token with private key
	tokenString, err := token.SignedString(KeyPairInstance.PrivateKey)
	if err != nil {
		logger.Log.Errorf("sign JWT-Token with private key failed %s", err)
		return nil, err
	}
	return &tokenString, nil
}

func ValidationToken(tokenString string) *PunqClaims {
	if KeyPairInstance == nil {
		KeyPairInstance = GetKeyPair()
	}

	ecdsaPublicKey, ok := KeyPairInstance.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		logger.Log.Error("Invalid public key")
		return nil
	}

	// Validation
	token, err := jwt.ParseWithClaims(tokenString, &PunqClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return ecdsaPublicKey, nil
	})

	if claims, ok := token.Claims.(*PunqClaims); ok && token.Valid {
		return claims
	}
	logger.Log.Errorf("Invalid token %v", err)
	return nil
}