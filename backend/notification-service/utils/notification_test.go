//package app
//
//import (
//	"context"
//	"testing"
//	"time"
//
//	app "github.com/Wenev/Survace/notification-service/internal/app"
//	"github.com/Wenev/Survace/notification-service/internal/app/domain"
//)
//
//type mockRepo struct {
//	FindByUserIDFunc      func(ctx context.Context, userID int32) ([]*domain.Notification, error)
//	FindByIDFunc          func(ctx context.Context, id int32) (*domain.Notification, error)
//	UpdateFunc            func(ctx context.Context, notif *domain.Notification) error
//	DeleteFunc            func(ctx context.Context, id int32) error
//	DeleteAllByUserIDFunc func(ctx context.Context, userID int32) error
//	CreateFunc            func(ctx context.Context, notif *domain.Notification) error
//}
//
//func (m *mockRepo) FindByUserID(ctx context.Context, userID int32) ([]*domain.Notification, error) {
//	return m.FindByUserIDFunc(ctx, userID)
//}
//func (m *mockRepo) FindByID(ctx context.Context, id int32) (*domain.Notification, error) {
//	return m.FindByIDFunc(ctx, id)
//}
//func (m *mockRepo) Update(ctx context.Context, notif *domain.Notification) error {
//	return m.UpdateFunc(ctx, notif)
//}
//func (m *mockRepo) Delete(ctx context.Context, id int32) error {
//	return m.DeleteFunc(ctx, id)
//}
//func (m *mockRepo) DeleteAllByUserID(ctx context.Context, userID int32) error {
//	return m.DeleteAllByUserIDFunc(ctx, userID)
//}
//func (m *mockRepo) Create(ctx context.Context, notif *domain.Notification) error {
//	return m.CreateFunc(ctx, notif)
//}
//
//type mockCache struct {
//	GetFunc    func(key string, dest interface{}) error
//	SetFunc    func(key string, value interface{}, expiration time.Duration) error
//	DeleteFunc func(key string) error
//}
//
//func (m *mockCache) Get(key string, dest interface{}) error {
//	return m.GetFunc(key, dest)
//}
//func (m *mockCache) Set(key string, value interface{}, expiration time.Duration) error {
//	return m.SetFunc(key, value, expiration)
//}
//func (m *mockCache) Delete(key string) error {
//	return m.DeleteFunc(key)
//}
//
//func TestNotificationService_AllEndpoints(t *testing.T) {
//	repo := &mockRepo{}
//	cache := &mockCache{}
//	svc := &app.NotificationServiceImpl{repo: repo, cache: cache}
//	ctx := context.Background()
//
//	repo.FindByUserIDFunc = func(ctx context.Context, userID int32) ([]*domain.Notification, error) {
//		return []*domain.Notification{{ID: 1, UserID: userID, Type: "message", Content: "hi", Read: false, CreatedAt: time.Now()}}, nil
//	}
//	cache.GetFunc = func(key string, dest interface{}) error { return assertError }
//	cache.SetFunc = func(key string, value interface{}, expiration time.Duration) error { return nil }
//	notifsCode, _, notifs, err := svc.GetNotificationsByUserID(ctx, 42)
//	if err != nil {
//		t.Errorf("unexpected error: %v", err)
//	}
//	if notifsCode != 0 {
//		t.Errorf("expected code 0, got %v", notifsCode)
//	}
//	if len(notifs) != 1 {
//		t.Errorf("expected 1 notification, got %v", len(notifs))
//	}
//
//	// Test MarkAsRead
//	repo.FindByIDFunc = func(ctx context.Context, id int32) (*domain.Notification, error) {
//		return &domain.Notification{ID: id, UserID: 42, Read: false}, nil
//	}
//	repo.UpdateFunc = func(ctx context.Context, notif *domain.Notification) error {
//		if !notif.Read {
//			t.Errorf("expected notification to be marked as read")
//		}
//		return nil
//	}
//	cache.DeleteFunc = func(key string) error { return nil }
//	code, msg, err := svc.MarkAsRead(ctx, 1)
//	if err != nil {
//		t.Errorf("unexpected error: %v", err)
//	}
//	if code != 0 {
//		t.Errorf("expected code 0, got %v", code)
//	}
//	if msg == "" || msg == "ok" {
//		t.Errorf("expected message to contain 'read', got %v", msg)
//	}
//
//	// Test MarkAsUnread
//	repo.FindByIDFunc = func(ctx context.Context, id int32) (*domain.Notification, error) {
//		return &domain.Notification{ID: id, UserID: 42, Read: true}, nil
//	}
//	repo.UpdateFunc = func(ctx context.Context, notif *domain.Notification) error {
//		if notif.Read {
//			t.Errorf("expected notification to be marked as unread")
//		}
//		return nil
//	}
//	code, msg, err = svc.MarkAsUnread(ctx, 1)
//	if err != nil {
//		t.Errorf("unexpected error: %v", err)
//	}
//	if code != 0 {
//		t.Errorf("expected code 0, got %v", code)
//	}
//	if msg == "" || msg == "ok" {
//		t.Errorf("expected message to contain 'unread', got %v", msg)
//	}
//
//	// Test DeleteNotification
//	repo.FindByIDFunc = func(ctx context.Context, id int32) (*domain.Notification, error) {
//		return &domain.Notification{ID: id, UserID: 42}, nil
//	}
//	repo.DeleteFunc = func(ctx context.Context, id int32) error { return nil }
//	cache.DeleteFunc = func(key string) error { return nil }
//	code, msg, err = svc.DeleteNotification(ctx, 1)
//	if err != nil {
//		t.Errorf("unexpected error: %v", err)
//	}
//	if code != 0 {
//		t.Errorf("expected code 0, got %v", code)
//	}
//	if msg == "" || msg == "ok" {
//		t.Errorf("expected message to contain 'deleted', got %v", msg)
//	}
//
//	// Test DeleteAllNotificationsByUserID
//	repo.DeleteAllByUserIDFunc = func(ctx context.Context, userID int32) error { return nil }
//	cache.DeleteFunc = func(key string) error { return nil }
//	code, msg, err = svc.DeleteAllNotificationsByUserID(ctx, 42)
//	if err != nil {
//		t.Errorf("unexpected error: %v", err)
//	}
//	if code != 0 {
//		t.Errorf("expected code 0, got %v", code)
//	}
//	if msg == "" || msg == "ok" {
//		t.Errorf("expected message to contain 'All notifications deleted', got %v", msg)
//	}
//
//	repo.CreateFunc = func(ctx context.Context, notif *domain.Notification) error {
//		if notif.Type != "message" {
//			t.Errorf("expected type 'message', got %v", notif.Type)
//		}
//		return nil
//	}
//	cache.DeleteFunc = func(key string) error { return nil }
//	code, msg, err = svc.CreateNotification(ctx, &domain.Notification{UserID: 42, Type: "message", Content: "hi"})
//	if err != nil {
//		t.Errorf("unexpected error: %v", err)
//	}
//	if code != 0 {
//		t.Errorf("expected code 0, got %v", code)
//	}
//	if msg == "" || msg == "ok" {
//		t.Errorf("expected message to contain 'created', got %v", msg)
//	}
//
//	code, msg, err = svc.CreateNotification(ctx, &domain.Notification{UserID: 42, Type: "invalid", Content: "hi"})
//	if err == nil {
//		t.Errorf("expected error for invalid type")
//	}
//	if code == 0 {
//		t.Errorf("expected non-zero code for invalid type")
//	}
//	if msg == "" || msg == "ok" {
//		t.Errorf("expected message to contain 'invalid notification type', got %v", msg)
//	}
//}
//
//var assertError = &testError{"cache miss"}
//
//type testError struct{ msg string }
//
//func (e *testError) Error() string { return e.msg }
//
