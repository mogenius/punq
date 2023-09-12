package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SecKeyPair    = "keyPair"
	TokenExpHours = 24 * 7 // 1 week
)

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

func CreateKeyPair() (*KeyPair, error) {
	provider := kubernetes.NewKubeProvider(nil)
	if provider == nil {
		msg := fmt.Sprintf("Failed to load kubeprovider.")
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	secretClient := provider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)
	existingSecret, getErr := secretClient.Get(context.TODO(), utils.JWTSECRET, metav1.GetOptions{})

	secret := utils.InitSecret()
	secret.ObjectMeta.Name = utils.JWTSECRET
	secret.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
	delete(secret.StringData, "exampleData") // delete example data

	keyPair, err := generateAuthKeyPair()
	if err != nil {
		msg := fmt.Sprintf("failed to generate key-pair %v", err)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}
	rawKeyPair, err := json.Marshal(keyPair)
	if err != nil {
		msg := fmt.Sprintf("failed marshaling %v", err)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}
	secret.StringData[SecKeyPair] = string(rawKeyPair)

	// if not exist
	if existingSecret == nil || getErr != nil {
		fmt.Println("Creating new punq-auth secret ...")
		_, err := secretClient.Create(context.TODO(), &secret, kubernetes.MoCreateOptions())
		if err != nil {
			return nil, err
		}
		fmt.Println("Created new punq-auth secret. ✅")

		return keyPair, nil
	}

	// get from secret
	return GetKeyPair()
}

func RemoveKeyPair() {
	provider := kubernetes.NewKubeProvider(nil)
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

func GetKeyPair() (*KeyPair, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.JWTSECRET, nil)
	if secret == nil {
		logger.Log.Warningf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.JWTSECRET)
		return CreateKeyPair()
	}

	if secret.Data[SecKeyPair] != nil {
		keyPair := KeyPair{}
		err := json.Unmarshal([]byte(secret.Data[SecKeyPair]), &keyPair)
		if err != nil {
			msg := fmt.Sprintf("failed to Unmarshal user '%s' %v.", SecKeyPair, err)
			logger.Log.Error(msg)
			return nil, err
		}
		return &keyPair, nil
	}

	return CreateKeyPair()
}

func GenerateToken(user *dtos.PunqUser) (*dtos.PunqToken, error) {
	if KeyPairInstance == nil {
		keyPair, err := GetKeyPair()
		if err != nil {
			return nil, err
		}
		KeyPairInstance = keyPair
	}

	claims := jwt.MapClaims{}
	claims["accessLevel"] = user.AccessLevel
	claims["userId"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(TokenExpHours)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)

	// sign JWT-Token with private key
	tokenString, err := token.SignedString(KeyPairInstance.PrivateKey)
	if err != nil {
		logger.Log.Errorf("sign JWT-Token with private key failed %s", err)
		return nil, err
	}
	return dtos.CreateToken(tokenString), nil
}

func ValidationToken(tokenString string) (*PunqClaims, error) {
	if KeyPairInstance == nil {
		keyPair, err := GetKeyPair()
		if err != nil {
			return nil, err
		}
		KeyPairInstance = keyPair
	}

	ecdsaPublicKey, ok := KeyPairInstance.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		msg := fmt.Sprintf("Invalid public key")
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	// Validation
	token, err := jwt.ParseWithClaims(tokenString, &PunqClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return ecdsaPublicKey, nil
	})
	claims, ok := token.Claims.(*PunqClaims)
	if !ok {
		msg := fmt.Sprintf("Type Assertion failed. Expected PunqClaims but received something different.")
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}
	if !token.Valid {
		msg := fmt.Sprintf("Invalid token %v", err)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}
	return claims, nil
}
