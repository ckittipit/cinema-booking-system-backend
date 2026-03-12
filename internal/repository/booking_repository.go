package repository

import (
	"cinema-booking/backend/internal/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepository struct {
	collection *mongo.Collection
}

func NewBookingRepository(db *mongo.Database) *BookingRepository {
	return &BookingRepository{
		collection: db.Collection("bookings"),
	}
}

func (r *BookingRepository) Create(ctx context.Context, booking *model.Booking) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, booking)
	if err != nil {
		return err
	}

	booking.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *BookingRepository) FindConfirmSeatIDsByShowtimeID(
	ctx context.Context, showtimeID primitive.ObjectID,
) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//defer เอาไว้ทำท้ายสุด ในกรณีนี้คือทำสักแย่างกับ ctx(context)เสร็จ แล้วให้ cancel() ต่อทันที //กันmemoryบวม
	defer cancel()

	filter := bson.M{
		"showtime_id": showtimeID,
		"status":      model.BookingStatusConformed,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []model.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	seatIDs := make([]string, 0, len(bookings))
	//ใส่ _ เพราะเวลาfor เป็น index, value แต่เราไม่ใช้ index เลยต้องใส่ _ ให้ทราบว่าไม่มีค่า
	for _, booking := range bookings {
		seatIDs = append(seatIDs, booking.SeatID)
	}

	return seatIDs, nil
}

func (r *BookingRepository) ExistConfirmedBookingByShowtimeAndSeat(
	ctx context.Context,
	showtimeID primitive.ObjectID,
	seatID string,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{
		"showtime_id": showtimeID,
		"seat_id":     seatID,
		"status":      model.BookingStatusConformed,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *BookingRepository) FindByID(
	ctx context.Context,
	bookingID primitive.ObjectID,
) (*model.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var booking model.Booking
	if err := r.collection.FindOne(ctx, bson.M{"_id": bookingID}).Decode(&booking); err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) Update(ctx context.Context, booking *model.Booking) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	booking.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": booking.ID},
		bson.M{
			"$set": bson.M{
				"status":     booking.Status,
				"locked_at":  booking.LockedAt,
				"expires_at": booking.ExpiresAt,
				"booked_at":  booking.BookedAt,
				"price":      booking.Price,
				"updated_at": booking.UpdatedAt,
			},
		},
	)

	return err
}

func (r *BookingRepository) ExistsActiveLockedBookingByShowTimeAndSeat(
	ctx context.Context,
	showtimeID primitive.ObjectID,
	seatID string,
	now time.Time,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{
		"showtime_id": showtimeID,
		"seat_id":     seatID,
		"status":      model.BookingStatusLocked,
		"expires_at": bson.M{
			"$gt": now,
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *BookingRepository) FindActiveLockedBookingsByShowtimeID(
	ctx context.Context,
	showtimeID primitive.ObjectID,
	now time.Time,
) ([]model.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{
		"showtim_id": showtimeID,
		"status":     model.BookingStatusLocked,
		"expires_at": bson.M{
			"$gt": now,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []model.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}
