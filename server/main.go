package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/redis/go-redis/v9"      // Lib pour Dragonfly/Redis
	"google.golang.org/grpc"            // lib grpc
	"google.golang.org/grpc/reflection" // lib test grpcurl

	pb "github.com/Baglyy/impression_tracking/proto" // Code généré à partir de proto
)

// Référence au client Dragonfly.
type server struct {
	pb.UnimplementedImpressionTrackerServiceServer
	dragonflyClient *redis.Client
}

// NOTE: Not used
// Créer une nouvelle instance du serveur
func NewServer(dflyClient *redis.Client) *server {
	return &server{dragonflyClient: dflyClient}
}

// Implémentation de la méthode RPC.
func (s *server) TrackImpression(ctx context.Context, req *pb.TrackImpressionRequest) (*pb.TrackImpressionResponse, error) {
	// 1. Récupérer l'ID de la publicité depuis la requête.
	adID := req.GetAdId()

	// 2. Valider l'input (Point Bonus).
	if adID == "" {
		// NOTE: Error in French ?
		return nil, errors.New("L'ad_id est vide")
	}

	// 3. Construire la clé à utiliser dans Dragonfly.
	key := adID

	// NOTE: Can use adID as key, no need to reassign a var
	// 4. Utiliser la commande atomique INCR de Redis.
	newCount, err := s.dragonflyClient.Incr(ctx, key).Result() // s.dragonflyClient donne accès au champ qu'on a définie dans struct.
	// Test d'erreur
	if err != nil {
		log.Println("Echec de l'incrémentation de ad_id : ", adID, err)
		return nil, err
	}

	// NOTE: Can return directly the reponse
	// 5. Si tout s'est bien passé, construire la réponse.
	reponse := &pb.TrackImpressionResponse{
		AdId:             adID,
		TotalImpressions: newCount,
	}

	return reponse, nil
}

// NOTE: Don't do logs in french
func main() {
	// NOTE: This should be near the Serve below. Group things that work together.
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("Echec de l'ouverture du port 50051 : ", err)
	}

	dbClient := redis.NewClient(&redis.Options{
		Addr: "dragonfly:6379",
	})

	// Création et enregistrement du serveur gRPC
	grpcServer := grpc.NewServer()

	// NOTE: You should use your NewServer function to initialize the server.
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
