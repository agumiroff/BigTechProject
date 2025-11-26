package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

var _ = Describe("PaymentService", func() {
	var (
		ctx           context.Context
		cancel        context.CancelFunc
		paymentClient paymentv1.PaymentServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		paymentClient = paymentv1.NewPaymentServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPaymentsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции payments")

		cancel()
	})

	Describe("PayOrder", func() {
		It("должен успешно обрабатывать платеж", func() {
			payment := env.GetTestPaymentData()

			resp, err := paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
				Payment: payment,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetTransactionUuid()).ToNot(BeEmpty())
			Expect(resp.GetTransactionUuid()).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))
		})

		It("должен успешно обрабатывать платеж картой", func() {
			payment := &paymentv1.Payment{
				OrderUuid:     gofakeit.UUID(),
				UserUuid:      gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			}

			resp, err := paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
				Payment: payment,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetTransactionUuid()).ToNot(BeEmpty())
		})

		It("должен успешно обрабатывать платеж через СБП", func() {
			payment := &paymentv1.Payment{
				OrderUuid:     gofakeit.UUID(),
				UserUuid:      gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
			}

			resp, err := paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
				Payment: payment,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetTransactionUuid()).ToNot(BeEmpty())
		})

		It("должен возвращать ошибку при невалидном методе оплаты", func() {
			payment := &paymentv1.Payment{
				OrderUuid:     gofakeit.UUID(),
				UserUuid:      gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			}

			resp, err := paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
				Payment: payment,
			})

			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})

		It("должен возвращать ошибку при пустом order_uuid", func() {
			payment := &paymentv1.Payment{
				OrderUuid:     "",
				UserUuid:      gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			}

			resp, err := paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
				Payment: payment,
			})

			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})

		It("должен возвращать ошибку при пустом user_uuid", func() {
			payment := &paymentv1.Payment{
				OrderUuid:     gofakeit.UUID(),
				UserUuid:      "",
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			}

			resp, err := paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
				Payment: payment,
			})

			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})

	Describe("Полный жизненный цикл платежа", func() {
		It("должен успешно создавать платеж и возвращать transaction_uuid", func() {
			// 1. Создаем платеж
			payment := &paymentv1.Payment{
				OrderUuid:     gofakeit.UUID(),
				UserUuid:      gofakeit.UUID(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
			}

			resp, err := paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
				Payment: payment,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetTransactionUuid()).ToNot(BeEmpty())

			transactionUUID := resp.GetTransactionUuid()
			Expect(transactionUUID).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))
		})

		It("должен запрещать дублирующие платежи для одного заказа", func() {
			orderUUID := gofakeit.UUID()
			userUUID := gofakeit.UUID()

			// 1. Первый платеж должен пройти успешно
			payment := &paymentv1.Payment{
				OrderUuid:     orderUUID,
				UserUuid:      userUUID,
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			}

			resp, err := paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
				Payment: payment,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetTransactionUuid()).ToNot(BeEmpty())

			// 2. Повторный платеж для того же заказа должен вернуть ошибку
			resp2, err := paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
				Payment: payment,
			})

			Expect(err).To(HaveOccurred())
			Expect(resp2).To(BeNil())
		})
	})
})
