package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
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
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	// get
	g, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, g.Address, parcel.Address)
	require.Equal(t, g.Client, parcel.Client)
	require.Equal(t, g.CreatedAt, parcel.CreatedAt)
	require.Equal(t, g.Status, parcel.Status)

	// delete
	err = store.Delete(p)
	require.NoError(t, err)

	g, err = store.Get(p)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, g)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, number)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	// check
	g, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, g.Address, newAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, number)

	// set status
	newStatus := "new status"
	err = store.SetStatus(number, newStatus)
	require.NoError(t, err)

	// check
	g, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, g.Status, newStatus)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
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

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	// check
	for _, parcel := range storedParcels {
		par, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, par.Address, parcel.Address)
		require.Equal(t, par.Client, parcel.Client)
		require.Equal(t, par.CreatedAt, parcel.CreatedAt)
		require.Equal(t, par.Status, parcel.Status)
	}
}
