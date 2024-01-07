package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauthv1beta1 "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

var (
	gcpScopes = []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
	}
)

func Gcp(ctx context.Context) error {
	cred, err := google.FindDefaultCredentials(ctx, gcpScopes...)
	if err != nil {
		return err
	}
	if cred == nil {
		return errors.New("failed finding default credentials, cred is nil")
	}
	token, err := cred.TokenSource.Token()
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("failed retrieving token from credentials")
	}
	ec := newExecCredential(token.AccessToken, token.Expiry)
	credString := formatJSON(ec)
	_, _ = fmt.Fprint(os.Stdout, credString)
	return nil
}

func formatJSON(ec *clientauthv1beta1.ExecCredential) string {
	enc, _ := json.MarshalIndent(ec, "", "  ")
	return string(enc)
}

func newExecCredential(token string, expiration time.Time) *clientauthv1beta1.ExecCredential {
	expirationTimestamp := metav1.NewTime(expiration)
	return &clientauthv1beta1.ExecCredential{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "client.authentication.k8s.io/v1beta1",
			Kind:       "ExecCredential",
		},
		Status: &clientauthv1beta1.ExecCredentialStatus{
			ExpirationTimestamp: &expirationTimestamp,
			Token:               token,
		},
	}
}
