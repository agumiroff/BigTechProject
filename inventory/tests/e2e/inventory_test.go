package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	invV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		invClient invV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		invClient = invV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearInventoryCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("CreatePart", func() {
		It("должен успешно создавать новую запчасть", func() {
			part := env.GetTestPartData()

			resp, err := invClient.CreatePart(ctx, &invV1.CreatePartRequest{
				Part: part,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetUuid()).ToNot(BeEmpty())
			Expect(resp.GetUuid()).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))
		})
	})

	Describe("GetPart", func() {
		var partUUID string

		BeforeEach(func() {
			// Вставляем тестовую запчасть
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой запчасти в MongoDB")
		})

		It("должен успешно возвращать запчасть по UUID", func() {
			resp, err := invClient.GetPart(ctx, &invV1.GetPartRequest{
				Uuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().Uuid).To(Equal(partUUID))
			Expect(resp.GetPart().Name).ToNot(BeEmpty())
			Expect(resp.GetPart().Description).ToNot(BeEmpty())
			Expect(resp.GetPart().Price).To(BeNumerically(">", 0))
			Expect(resp.GetPart().CreatedAt).ToNot(BeZero())
		})

		It("должен возвращать ошибку для несуществующей запчасти", func() {
			resp, err := invClient.GetPart(ctx, &invV1.GetPartRequest{
				Uuid: "non-existent-uuid",
			})

			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})

	Describe("ListParts", func() {
		BeforeEach(func() {
			// Вставляем несколько тестовых запчастей
			for i := 0; i < 3; i++ {
				_, err := env.InsertTestPart(ctx)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("должен возвращать список запчастей", func() {
			resp, err := invClient.ListParts(ctx, &invV1.ListPartsRequest{
				Filter: &invV1.PartsFilter{},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeEmpty())
			Expect(len(resp.GetParts())).To(BeNumerically(">=", 3))
		})

		It("должен фильтровать запчасти по категории", func() {
			resp, err := invClient.ListParts(ctx, &invV1.ListPartsRequest{
				Filter: &invV1.PartsFilter{
					Categories: []invV1.Category{invV1.Category_CATEGORY_ENGINE},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeEmpty())
			for _, part := range resp.GetParts() {
				Expect(part.Category).To(Equal(invV1.Category_CATEGORY_ENGINE))
			}
		})
	})

	Describe("Полный жизненный цикл", func() {
		It("должен поддерживать создание и получение запчасти", func() {
			// 1. Создаем запчасть
			partData := env.CreateTestPart()
			createResp, err := invClient.CreatePart(ctx, &invV1.CreatePartRequest{
				Part: partData,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(createResp.GetUuid()).ToNot(BeEmpty())
			uuid := createResp.GetUuid()

			// 2. Получаем созданную запчасть
			getResp, err := invClient.GetPart(ctx, &invV1.GetPartRequest{
				Uuid: uuid,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(getResp.GetPart().Uuid).To(Equal(uuid))
			Expect(getResp.GetPart().Name).To(Equal(partData.Name))
			Expect(getResp.GetPart().Description).To(Equal(partData.Description))
			Expect(getResp.GetPart().Price).To(Equal(partData.Price))
			Expect(getResp.GetPart().Category).To(Equal(partData.Category))

			// 3. Проверяем, что запчасть есть в списке
			listResp, err := invClient.ListParts(ctx, &invV1.ListPartsRequest{
				Filter: &invV1.PartsFilter{
					Uuids: []string{uuid},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(listResp.GetParts()).To(HaveLen(1))
			Expect(listResp.GetParts()[0].Uuid).To(Equal(uuid))
		})
	})
})
