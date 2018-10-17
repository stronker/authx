/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package interceptor

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/nalej/authx/internal/app/authx/handler"
	"github.com/nalej/authx/internal/app/authx/manager"
	"github.com/nalej/authx/pkg/token"
	pbAuthx "github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"time"
)

var _ = ginkgo.Describe("Authorize method", func() {
	ginkgo.Context("empty authorization list with AllowsAll", func() {

		duration, _ := time.ParseDuration("1d")
		claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{}, "o1"),
			"i1", time.Now(), duration)
		cfg := NewConfig(&AuthorizationConfig{Permissions: map[string]Permission{}},
			"myLittleSecret", "auth", true)

		ginkgo.It("allows any method", func() {
			err := authorize("service1", claim, cfg)
			gomega.Expect(err).To(gomega.Succeed())
		})

	})

	ginkgo.Context("empty authorization list without AllowsAll", func() {

		duration, _ := time.ParseDuration("1d")
		claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{}, "o1"),
			"i1", time.Now(), duration)
		cfg := NewConfig(&AuthorizationConfig{Permissions: map[string]Permission{}},
			"myLittleSecret", "auth", false)

		ginkgo.It("allows any method", func() {
			err := authorize("service1", claim, cfg)
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})

	ginkgo.Context("with authorization list with AllowsAll", func() {

		duration, _ := time.ParseDuration("1d")
		unknownMethod := "unknownService"
		method1 := "method1"
		method2 := "method2"

		primitive1 := "primitive1"
		primitive2 := "primitive2"

		cfg := NewConfig(&AuthorizationConfig{Permissions: map[string]Permission{
			method1: {Must: [] string{primitive1}},
			method2: {Must: [] string{primitive2}},
		}},
			"myLittleSecret", "auth", true)

		ginkgo.Context("without primitives", func() {
			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{}, "o1"),
				"i1", time.Now(), duration)
			ginkgo.It("allows unknown method", func() {
				err := authorize(unknownMethod, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})

			ginkgo.It("doesn't allow method1", func() {
				err := authorize(method1, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
			ginkgo.It("doesn't allow method2", func() {
				err := authorize(method2, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("with primitive1", func() {
			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{primitive1}, "o1"),
				"i1", time.Now(), duration)
			ginkgo.It("allows unknown method", func() {
				err := authorize(unknownMethod, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})

			ginkgo.It("allows method1", func() {
				err := authorize(method1, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})
			ginkgo.It("doesn't allow method2", func() {
				err := authorize(method2, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
		ginkgo.Context("with primitive2", func() {
			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{primitive2}, "o1"),
				"i1", time.Now(), duration)
			ginkgo.It("allows unknown method", func() {
				err := authorize(unknownMethod, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})

			ginkgo.It("doesn't allow method1", func() {
				err := authorize(method1, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
			ginkgo.It("allows method2", func() {
				err := authorize(method2, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})
		})

		ginkgo.Context("with all primitives", func() {
			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{primitive1, primitive2}, "o1"),
				"i1", time.Now(), duration)
			ginkgo.It("allows unknown method", func() {
				err := authorize(unknownMethod, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})

			ginkgo.It("allows method1", func() {
				err := authorize(method1, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})
			ginkgo.It("allows method2", func() {
				err := authorize(method2, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})
		})

	})

	ginkgo.Context("with authorization list without AllowsAll", func() {

		duration, _ := time.ParseDuration("1d")
		unknownMethod := "unknownService"
		method1 := "method1"
		method2 := "method2"

		primitive1 := "primitive1"
		primitive2 := "primitive2"

		cfg := NewConfig(&AuthorizationConfig{Permissions: map[string]Permission{
			method1: {Must: [] string{primitive1}},
			method2: {Must: [] string{primitive2}},
		}},
			"myLittleSecret", "auth", false)

		ginkgo.Context("without primitives", func() {
			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{}, "o1"),
				"i1", time.Now(), duration)
			ginkgo.It("doesn't allow unknown method", func() {
				err := authorize(unknownMethod, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})

			ginkgo.It("doesn't allow method1", func() {
				err := authorize(method1, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
			ginkgo.It("doesn't allow method2", func() {
				err := authorize(method2, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})

		ginkgo.Context("with primitive1", func() {
			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{primitive1}, "o1"),
				"i1", time.Now(), duration)
			ginkgo.It("doesn't allow unknown method", func() {
				err := authorize(unknownMethod, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})

			ginkgo.It("allows method1", func() {
				err := authorize(method1, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})
			ginkgo.It("doesn't allow method2", func() {
				err := authorize(method2, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
		})
		ginkgo.Context("with primitive2", func() {
			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{primitive2}, "o1"),
				"i1", time.Now(), duration)
			ginkgo.It("doesn't allow unknown method", func() {
				err := authorize(unknownMethod, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})

			ginkgo.It("doesn't allow method1", func() {
				err := authorize(method1, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
			ginkgo.It("allows method2", func() {
				err := authorize(method2, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})
		})

		ginkgo.Context("with all primitives", func() {
			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{primitive1, primitive2}, "o1"),
				"i1", time.Now(), duration)
			ginkgo.It("doesn't allow unknown method", func() {
				err := authorize(unknownMethod, claim, cfg)
				gomega.Expect(err).To(gomega.HaveOccurred())
			})

			ginkgo.It("allows method1", func() {
				err := authorize(method1, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})
			ginkgo.It("allows method2", func() {
				err := authorize(method2, claim, cfg)
				gomega.Expect(err).To(gomega.Succeed())
			})
		})

	})

})

var _ = ginkgo.Describe("checkJWT method", func() {
	ginkgo.Context("with valid JWT", func() {
		duration, _ := time.ParseDuration("1d")
		secret := "myLittleSecret"
		header := "auth"
		claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{}, "o1"),
			"i1", time.Now(), duration)
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, _ := t.SignedString([]byte(secret))
		cfg := NewConfig(&AuthorizationConfig{Permissions: map[string]Permission{}},
			secret, header, true)
		md := metadata.New(map[string]string{header: tokenString})

		ctx := metadata.NewIncomingContext(context.Background(), md)

		ginkgo.It("works", func() {
			claim, err := checkJWT(ctx, cfg)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(claim).NotTo(gomega.BeNil())
		})

	})
	ginkgo.Context("with invalid JWT", func() {
		duration, _ := time.ParseDuration("1d")
		secret := "myLittleSecret"
		wrongSecret := "myWrongSecret"
		header := "auth"

		cfg := NewConfig(&AuthorizationConfig{Permissions: map[string]Permission{}},
			secret, header, true)

		claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{}, "o1"),
			"i1", time.Now(), duration)
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, _ := t.SignedString([]byte(wrongSecret))

		md := metadata.New(map[string]string{header: tokenString})

		ctx := metadata.NewIncomingContext(context.Background(), md)

		ginkgo.It("doesn't work", func() {
			claim, err := checkJWT(ctx, cfg)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(claim).To(gomega.BeNil())
		})
	})
	ginkgo.Context("with wrong header", func() {
		duration, _ := time.ParseDuration("1d")
		secret := "myLittleSecret"
		header := "auth"
		wrongHeader := "authx"
		claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{}, "o1"),
			"i1", time.Now(), duration)
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, _ := t.SignedString([]byte(secret))
		cfg := NewConfig(&AuthorizationConfig{Permissions: map[string]Permission{}},
			secret, header, true)

		md := metadata.New(map[string]string{wrongHeader: tokenString})
		parentCtx, _ := context.WithTimeout(context.TODO(), duration)
		ctx := metadata.NewIncomingContext(parentCtx, md)

		ginkgo.It("doesn't works", func() {
			claim, err := checkJWT(ctx, cfg)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(claim).To(gomega.BeNil())
		})

	})
	ginkgo.Context("with wrong MD", func() {
		duration, _ := time.ParseDuration("1d")
		secret := "myLittleSecret"
		header := "auth"
		claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1", [] string{}, "o1"),
			"i1", time.Now(), duration)
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, _ := t.SignedString([]byte(secret))
		cfg := NewConfig(&AuthorizationConfig{Permissions: map[string]Permission{}},
			secret, header, true)

		md := metadata.New(map[string]string{header: tokenString})
		parentCtx, _ := context.WithTimeout(context.TODO(), duration)
		ctx := metadata.NewOutgoingContext(parentCtx, md)

		ginkgo.It("doesn't works", func() {
			claim, err := checkJWT(ctx, cfg)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(claim).To(gomega.BeNil())
		})

	})

})

var _ = ginkgo.Describe("GRP interceptor method ", func() {

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client pbAuthx.AuthxClient

	var mgr *manager.Authx

	duration, _ := time.ParseDuration("1d")

	method1 := "/authx.Authx/AddBasicCredentials"
	method2 := "/authx.Authx/DeleteCredentials"

	primitive1 := "primitive1"
	primitive2 := "primitive2"


	ginkgo.Context("with AllowsAll", func() {
		cfg := NewConfig(&AuthorizationConfig{Permissions: map[string]Permission{
			method1: {Must: [] string{primitive1}},
			method2: {Must: [] string{primitive2}},
		}}, "myLittleSecret", "auth", true)

		ginkgo.BeforeSuite(func() {
			listener = test.GetDefaultListener()
			server = grpc.NewServer(WithServerUnaryInterceptor(cfg))

			mgr = manager.NewAuthxMockup()
			handler := handler.NewAuthx(mgr)

			pbAuthx.RegisterAuthxServer(server, handler)

			test.LaunchServer(server, listener)

			conn, err := test.GetConn(*listener)
			gomega.Expect(err).Should(gomega.Succeed())
			client = pbAuthx.NewAuthxClient(conn)
		})

		userName := "u1"
		organizationID := "o1"
		roleID := "r1"
		pass := "MyLittlePassword"

		ginkgo.BeforeEach(func() {
			role := &pbAuthx.Role{
				OrganizationId: organizationID,
				RoleId:         roleID,
				Name:           "rName1",
				Primitives:     []pbAuthx.AccessPrimitive{pbAuthx.AccessPrimitive_ORG},
			}
			success, err := client.AddRole(context.Background(), role)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})

		ginkgo.It("add basic credentials with correct roleID and correct JWT", func() {

			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1",
				[] string{primitive1, primitive2}, "o1"),
				"i1", time.Now(), duration)
			t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
			tokenString, _ := t.SignedString([]byte(cfg.Secret))

			md := metadata.New(map[string]string{cfg.Header: tokenString})

			ctx := metadata.NewOutgoingContext(context.Background(), md)
			success, err := client.AddBasicCredentials(ctx,
				&pbAuthx.AddBasicCredentialRequest{OrganizationId: organizationID,
					RoleId:   roleID,
					Username: userName,
					Password: pass,
				})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})

		ginkgo.It("add basic credentials with correct roleID and incorrect JWT", func() {

			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1",
				[] string{primitive1, primitive2}, "o1"),
				"i1", time.Now(), duration)
			t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
			tokenString, _ := t.SignedString([]byte("wrongSecret"))

			md := metadata.New(map[string]string{cfg.Header: tokenString})

			ctx := metadata.NewOutgoingContext(context.Background(), md)
			success, err := client.AddBasicCredentials(ctx,
				&pbAuthx.AddBasicCredentialRequest{OrganizationId: organizationID,
					RoleId:   roleID,
					Username: userName,
					Password: pass,
				})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(success).To(gomega.BeNil())
		})

		ginkgo.It("add basic credentials with correct roleID and correct JWT", func() {

			claim := token.NewClaim(*token.NewPersonalClaim("u1", "r1",
				[] string{primitive2}, "o1"),
				"i1", time.Now(), duration)
			t := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
			tokenString, _ := t.SignedString([]byte(cfg.Secret))

			md := metadata.New(map[string]string{cfg.Header: tokenString})

			ctx := metadata.NewOutgoingContext(context.Background(), md)
			success, err := client.AddBasicCredentials(ctx,
				&pbAuthx.AddBasicCredentialRequest{OrganizationId: organizationID,
					RoleId:   roleID,
					Username: userName,
					Password: pass,
				})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(success).To(gomega.BeNil())
		})

		ginkgo.AfterEach(func() {
			err := mgr.Clean()
			gomega.Expect(err).To(gomega.Succeed())
		})

		ginkgo.AfterSuite(func() {
			server.Stop()
			listener.Close()
		})
	})

})
