package integration

import (
	"context"
)

var _ = Describe("OrderService", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)
	})

	AfterEach(func() {
		// Чистим таблицы после теста
		err := testEnv.ClearOrderCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку таблиц orders и order_parts")

		cancel()
	})

	Describe("Order Management", func() {
		It("должен успешно подключиться к приложению", func() {
			// Простой тест проверки, что контейнеры запустились
			Expect(testEnv.App).ToNot(BeNil())
			Expect(testEnv.Postgres).ToNot(BeNil())
			Expect(testEnv.Postgres.Pool()).ToNot(BeNil())
		})

		It("должен успешно очистить таблицы", func() {
			// Проверяем, что очистка работает
			err := testEnv.ClearOrderCollection(ctx)
			Expect(err).ToNot(HaveOccurred())

			// Проверяем, что таблицы пустые
			pool := testEnv.Postgres.Pool()

			var count int
			err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders").Scan(&count)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(0), "таблица orders должна быть пустой")

			err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM order_parts").Scan(&count)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(0), "таблица order_parts должна быть пустой")
		})
	})
})
