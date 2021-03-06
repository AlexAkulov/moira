package controller

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/moira-alert/moira"
	"github.com/moira-alert/moira/api"
	"github.com/moira-alert/moira/api/dto"
	"github.com/moira-alert/moira/mock/moira-alert"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetUserSubscriptions(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	database := mock_moira_alert.NewMockDatabase(mockCtrl)
	login := "user"

	Convey("Two subscriptions", t, func() {
		subscriptionIDs := []string{uuid.NewV4().String(), uuid.NewV4().String()}
		subscriptions := []*moira.SubscriptionData{{ID: subscriptionIDs[0]}, {ID: subscriptionIDs[1]}}
		database.EXPECT().GetUserSubscriptionIDs(login).Return(subscriptionIDs, nil)
		database.EXPECT().GetSubscriptions(subscriptionIDs).Return(subscriptions, nil)
		list, err := GetUserSubscriptions(database, login)
		So(err, ShouldBeNil)
		So(list, ShouldResemble, &dto.SubscriptionList{List: []moira.SubscriptionData{*subscriptions[0], *subscriptions[1]}})
	})

	Convey("Two ids, one subscription", t, func() {
		subscriptionIDs := []string{uuid.NewV4().String(), uuid.NewV4().String()}
		subscriptions := []*moira.SubscriptionData{{ID: subscriptionIDs[1]}}
		database.EXPECT().GetUserSubscriptionIDs(login).Return(subscriptionIDs, nil)
		database.EXPECT().GetSubscriptions(subscriptionIDs).Return(subscriptions, nil)
		list, err := GetUserSubscriptions(database, login)
		So(err, ShouldBeNil)
		So(list, ShouldResemble, &dto.SubscriptionList{List: []moira.SubscriptionData{*subscriptions[0]}})
	})

	Convey("Errors", t, func() {
		Convey("GetUserSubscriptionIDs", func() {
			expected := fmt.Errorf("Oh no!!!11 Cant get subscription ids")
			database.EXPECT().GetUserSubscriptionIDs(login).Return(nil, expected)
			list, err := GetUserSubscriptions(database, login)
			So(err, ShouldResemble, api.ErrorInternalServer(expected))
			So(list, ShouldBeNil)
		})

		Convey("GetSubscriptions", func() {
			expected := fmt.Errorf("Oh no!!!11 Cant get subscriptions")
			subscriptionIDs := []string{uuid.NewV4().String(), uuid.NewV4().String()}
			database.EXPECT().GetUserSubscriptionIDs(login).Return(subscriptionIDs, nil)
			database.EXPECT().GetSubscriptions(subscriptionIDs).Return(nil, expected)
			list, err := GetUserSubscriptions(database, login)
			So(err, ShouldResemble, api.ErrorInternalServer(expected))
			So(list, ShouldBeNil)
		})
	})
}

func TestRemoveSubscription(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	database := mock_moira_alert.NewMockDatabase(mockCtrl)
	login := "user"
	id := uuid.NewV4().String()

	Convey("Success", t, func() {
		database.EXPECT().RemoveSubscription(id).Return(nil)
		err := RemoveSubscription(database, id, login)
		So(err, ShouldBeNil)
	})

	Convey("Error", t, func() {
		expected := fmt.Errorf("Oooops! Can not remove subscription")
		database.EXPECT().RemoveSubscription(id).Return(expected)
		err := RemoveSubscription(database, id, login)
		So(err, ShouldResemble, api.ErrorInternalServer(expected))
	})
}

func TestSendTestNotification(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	database := mock_moira_alert.NewMockDatabase(mockCtrl)
	id := uuid.NewV4().String()

	Convey("Success", t, func() {
		database.EXPECT().PushNotificationEvent(gomock.Any(), false).Return(nil)
		err := SendTestNotification(database, id)
		So(err, ShouldBeNil)
	})

	Convey("Error", t, func() {
		expected := fmt.Errorf("Oooops! Can not push event")
		database.EXPECT().PushNotificationEvent(gomock.Any(), false).Return(expected)
		err := SendTestNotification(database, id)
		So(err, ShouldResemble, api.ErrorInternalServer(expected))
	})
}

func TestWriteSubscription(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	database := mock_moira_alert.NewMockDatabase(mockCtrl)
	login := "user"

	subscription := dto.Subscription{ID: ""}
	Convey("Success", t, func() {
		database.EXPECT().SaveSubscription(gomock.Any()).Return(nil)
		err := WriteSubscription(database, login, &subscription)
		So(err, ShouldBeNil)
	})

	Convey("Error", t, func() {
		expected := fmt.Errorf("Oooops! Can not create subscription")
		database.EXPECT().SaveSubscription(gomock.Any()).Return(expected)
		err := WriteSubscription(database, login, &subscription)
		So(err, ShouldResemble, api.ErrorInternalServer(expected))
	})
}
