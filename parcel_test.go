package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	store := NewParcelStore("sqlite", "tracker.db")
	defer store.db.Close()

	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, number)

	addedParcel, err := store.Get(number)
	require.NoError(t, err)

	parcel.Number = number
	assert.Equal(t, addedParcel, parcel)

	err = store.Delete(number)
	require.NoError(t, err)

	_, err = store.Get(number)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	store := NewParcelStore("sqlite", "tracker.db")
	defer store.db.Close()
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	addedParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, addedParcel.Address, newAddress)

}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	store := NewParcelStore("sqlite", "tracker.db")
	defer store.db.Close()

	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, number)

	err = store.SetStatus(number, ParcelStatusDelivered)
	require.NoError(t, err)

	addedParcel, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, addedParcel.Status, ParcelStatusDelivered)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	store := NewParcelStore("sqlite", "tracker.db")
	defer store.db.Close()

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		number, err := store.Add(parcels[i])
		require.NoError(t, err)
		assert.NotEmpty(t, number)

		parcels[i].Number = number
		parcelMap[number] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)

	require.Equal(t, len(parcels), len(storedParcels))
	for _, parcel := range storedParcels {
		assert.Equal(t, parcel, parcelMap[parcel.Number])
	}
}
