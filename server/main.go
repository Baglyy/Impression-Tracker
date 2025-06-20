package main

import (
	"errors"
	"context" // Requis par gRPC, on ne s'en sert pas directement ici.
	"log"     // Pour afficher des messages dans la console (ex: "Le serveur démarre...").
	"net"     // Pour tout ce qui est réseau, comme ouvrir un port.

	// Import du code que buf a généré
	pb "github.com/Baglyy/impression_tracking/proto" // Code généré à partir de proto
	"github.com/redis/go-redis/v9" // Lib pour Dragonfly/Redis
	"google.golang.org/grpc" // lib grpc
	"google.golang.org/grpc/reflection" // lib test grpcurl 

)

// Référence au client Dragonfly.
type server struct {
	pb.UnimplementedImpressionTrackerServiceServer
	dragonflyClient *redis.Client
}

// NewServer crée une nouvelle instance de notre serveur.
func NewServer(dflyClient *redis.Client) *server {
	return &server{dragonflyClient: dflyClient}
}

// Implémentation de la méthode RPC.
func (s *server) TrackImpression(ctx context.Context, req *pb.TrackImpressionRequest) (*pb.TrackImpressionResponse, error) {
	
	// 1. Récupérer l'ID de la publicité depuis la requête.
	adID := req.GetAdId()

    // 2. Valider l'input (Point Bonus).
	if adID == "" {
		return nil, errors.New("L'ad_id est vide")
	}

    // 3. Construire la clé à utiliser dans Dragonfly.
	key := adID

    // 4. Utiliser la commande atomique INCR de Redis.
	newCount, err := s.dragonflyClient.Incr(ctx, key).Result() 	// s.dragonflyClient donne accès au champ qu'on a définie dans struct.

	// Test d'erreur
	if err != nil {
		log.Println("Echec de l'incrémentation de ad_id : ", adID, err)
		return nil, err
	}

    // 5. Si tout s'est bien passé, construire la réponse.
	reponse := &pb.TrackImpressionResponse{
		AdId: adID,
		TotalImpressions: newCount,
	}

	return reponse, nil
}

func main() {

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("Echec de l'ouverture du port 50051 : ", err)
	}

	dbClient := redis.NewClient(&redis.Options{
		Addr: "dragonfly:6379",
	})

	// Création et enregistrement du serveur gRPC
	grpcServer := grpc.NewServer()

	// Passer au serveur la connexion à la BDD.
	myServer := &server{
		dragonflyClient: dbClient,
	}

	// Dire au serveur gRPC d'utiliser notre logique.
	pb.RegisterImpressionTrackerServiceServer(grpcServer, myServer)

	// Activez la réflexion sur le serveur gRPC. C'est ce qui permet à grpcurl de fonctionner.
	reflection.Register(grpcServer)

	log.Println("Le serveur démarre sur le port 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalln("Le serveur a planté : ", err)
	}
}