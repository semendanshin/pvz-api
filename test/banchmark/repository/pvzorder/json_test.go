package pvzorder

import (
	"github.com/stretchr/testify/assert"
	"homework/internal/domain"
	"homework/internal/infrastructure/repositories/pvzorder"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func BenchmarkJsonPVZOrderRepository_CreateOrder(b *testing.B) {
	repo := pvzorder.NewJSONRepository("test.json")
	defer os.Remove("test.json")

	order := domain.NewPVZOrder(
		"100",
		"1",
		"1",
		1000,
		1000,
		24*time.Hour,
		domain.PackagingTypeBox,
		false,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := repo.CreateOrder(order)
		assert.NoError(b, err)
	}

	b.StopTimer()
}

func BenchmarkJsonPVZOrderRepository_GetOrder(b *testing.B) {
	const maxID = 1000

	repo := pvzorder.NewJSONRepository("test.json")
	defer os.Remove("test.json")

	for i := 0; i < maxID; i++ {
		order := domain.NewPVZOrder(
			strconv.Itoa(i),
			"1",
			"1",
			1000,
			1000,
			24*time.Hour,
			domain.PackagingTypeBox,
			false,
		)
		err := repo.CreateOrder(order)
		assert.NoError(b, err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id := rand.Intn(maxID)
		_, err := repo.GetOrder(strconv.Itoa(id))
		assert.NoError(b, err)
	}

	b.StopTimer()
}
