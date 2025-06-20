## Installation et Lancement avec Docker Compose

1. **Installer Docker** : 
    https://docs.docker.com/get-docker/
    
2.  **Cloner le dépôt** :
    ```sh
    git clone https://github.com/Baglyy/impression_tracking.git
    cd impression_tracking
    ```

3.  **Construire et lancer les conteneurs** :
    ```sh
    docker-compose build
    docker-compose up -d
    ```

## Utilisation de l'Environnement Devbox (pour le développement)

1. **Installer DevBox**
    https://www.jetpack.io/devbox/docs/installing-devbox/

2.  **Activer le shell Devbox** (à la racine du projet)
    ```sh
    devbox shell
    ```

3.  **Lancer les tests unitaires** :
    ```sh
    go test ./...
    ```

## Comment Tester l'Endpoint gRPC

1. **Installer grpcurl** :
    ```sh
    wget https://github.com/fullstorydev/grpcurl/releases/download/v1.9.1/grpcurl_1.9.1_linux_x86_64.tar.gz
    tar -xzf grpcurl_1.9.1_linux_x86_64.tar.gz
    sudo mv grpcurl /usr/local/bin/
    ```

1.  **Lister les services disponibles** :
    ```sh
    grpcurl -plaintext localhost:50051 list
    ```

2.  **Envoyer une impression** :
    ```sh
    grpcurl -plaintext -d '{"ad_id": "ad_test_99"}' localhost:50051 arago.ImpressionTrackerService/TrackImpression
    ```
    **Réponse attendue :**
    ```json
    {
      "ad_id": "ad_test_99",
      "total_impressions": "1"
    }
    ```
    **Tester plusieurs fois pour voir le compteur s'incrémenter**

3.  **Tester la gestion d'erreur** (ID vide) :
    ```sh
    grpcurl -plaintext -d '{"ad_id": ""}' localhost:50051 arago.ImpressionTrackerService/TrackImpression
    ```
    **Réponse attendue :**
    ERROR:
      Code: InvalidArgument
      Message: ad_id cannot be empty
