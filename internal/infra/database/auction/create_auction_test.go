package auction

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"labs-auction/configuration/database/mongodb"
	"labs-auction/internal/entity/auction_entity"
	"log"
	"os"
	"testing"
	"time"
)

func TestShouldNotBeExpired(t *testing.T) {
	ctx := context.Background()

	os.Setenv("BATCH_INSERT_INTERVAL", "20s")
	os.Setenv("MAX_BATCH_SIZE", "4")
	os.Setenv("AUCTION_INTERVAL", "20s")
	os.Setenv("MONGO_INITDB_ROOT_USERNAME", "admin")
	os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "admin")
	os.Setenv("MONGODB_URL", "mongodb://admin:admin@mongodb:27017/auctions?authSource=admin")
	os.Setenv("MONGODB_DB", "auctions")
	os.Setenv("AUCTION_EXPIRED", "20s")
	db, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Create a test auction that should be closed
	testAuction1 := &AuctionEntityMongo{
		Id:          uuid.New().String(),
		ProductName: "Test Product 1",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Unix(),
	}

	_, err = db.Collection("auctions").InsertOne(context.Background(), testAuction1)
	if err != nil {
		t.Fatalf("Failed to insert auction: %v", err)
	}

	closeExpiredAuctions(db)

	// Check if the first auction is closed
	var fetchedAuction AuctionEntityMongo
	err = db.Collection("auctions").FindOne(context.Background(), bson.M{"_id": testAuction1.Id}).Decode(&fetchedAuction)
	if err != nil {
		t.Fatalf("Failed to find auction: %v", err)
	}

	if fetchedAuction.Status != auction_entity.Active {
		t.Errorf("Expected auction 1 to be closed")
	}

}

func TestShouldBeExpired(t *testing.T) {
	ctx := context.Background()

	os.Setenv("BATCH_INSERT_INTERVAL", "20s")
	os.Setenv("MAX_BATCH_SIZE", "4")
	os.Setenv("AUCTION_INTERVAL", "20s")
	os.Setenv("MONGO_INITDB_ROOT_USERNAME", "admin")
	os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "admin")
	os.Setenv("MONGODB_URL", "mongodb://admin:admin@mongodb:27017/auctions?authSource=admin")
	os.Setenv("MONGODB_DB", "auctions")
	os.Setenv("AUCTION_EXPIRED", "0.5s")

	db, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Create a test auction that should be closed
	testAuction1 := &AuctionEntityMongo{
		Id:          uuid.New().String(),
		ProductName: "Test Product 1",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Unix(),
	}

	_, err = db.Collection("auctions").InsertOne(context.Background(), testAuction1)
	if err != nil {
		t.Fatalf("Failed to insert auction: %v", err)
	}

	time.Sleep(1 * time.Second)

	closeExpiredAuctions(db)

	// Check if the first auction is closed
	var fetchedAuction AuctionEntityMongo
	err = db.Collection("auctions").FindOne(context.Background(), bson.M{"_id": testAuction1.Id}).Decode(&fetchedAuction)
	if err != nil {
		t.Fatalf("Failed to find auction: %v", err)
	}

	if fetchedAuction.Status != auction_entity.Completed {
		t.Errorf("Expected auction 1 to be closed")
	}
}
