package main

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	pb "github.com/Baglyy/impression_tracking/proto"
)

func TestTrackImpression(t *testing.T) {
	ctx := context.Background()

	// Sous-test 1: Le cas où tout se passe bien
	t.Run("Cas de succès", func(t *testing.T) {

		// Créer une fausse base de données pour ce test
		db, mock := redismock.NewClientMock()

		// Créer le serveur avec la fausse base de données
		server := &server{dragonflyClient: db}

		adID := "ad_test_123"

		key := int64(5)

		// Appel à Incr avec adID. Si l'appel est réalisé, réponse = 5
		mock.ExpectIncr(adID).SetVal(key)

		// Simuler une requête gRPC
		req := &pb.TrackImpressionRequest{AdId: adID}

		// Appel de TrackImpression
		res, err := server.TrackImpression(ctx, req)

		// Vérifier les erreurs
		if err != nil {
			t.Fatal("Echec du test")
		}

		// Vérifier les données
		if res.GetTotalImpressions() != key {
			t.Fatal("Echec du test")
		}

	})

	t.Run("Cas avec ad_id vide", func(t *testing.T) {

		db, _ := redismock.NewClientMock()
		s := &server{dragonflyClient: db}

		// Requête avec ad_id vide
		req := &pb.TrackImpressionRequest{AdId: ""}

		// On appelle notre fonction à tester. On ignore la réponse (`_`)
		_, err := s.TrackImpression(ctx, req)

		// Test erreur
		if err == nil {
			t.Fatal("Echec du test")
		} 

	})
}