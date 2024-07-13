package auction

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"labs-auction/configuration/logger"
	"labs-auction/internal/entity/auction_entity"
	"labs-auction/internal/internal_error"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func GetAuctionDuration() (time.Duration, error) {
	durationStr := os.Getenv("AUCTION_INTERVAL")
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		return 0, err
	}
	return time.Duration(duration) * time.Minute, nil
}

func StartAuctionExpirationRoutine(db *mongo.Database) {
	go func() {
		auctionInterval, err := time.ParseDuration(os.Getenv("AUCTION_INTERVAL"))
		if err != nil {
			log.Fatalf("Error parsing AUCTION_INTERVAL: %v", err)
			return
		}
		closeExpiredAuctions(db)
		for {
			time.Sleep(auctionInterval)
			closeExpiredAuctions(db)
		}
	}()
}

func closeExpiredAuctions(db *mongo.Database) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.Background()
	//defer cancel()

	//duration, err := GetAuctionDuration()
	duration, err := time.ParseDuration(os.Getenv("AUCTION_EXPIRED"))
	if err != nil {
		log.Printf("Error getting auction duration: %v", err)
		return
	}

	expirationTime := time.Now().Add(-duration)
	filter := bson.M{"timestamp": bson.M{"$lt": expirationTime.Unix()}, "status": auction_entity.Active}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	_, err = db.Collection("auctions").UpdateMany(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating expired auctions: %v", err)
	}
}
