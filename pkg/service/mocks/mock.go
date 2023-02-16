// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	models "cmd/pkg/repository/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuthorization) CreateUser(user models.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthorizationMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthorization)(nil).CreateUser), user)
}

// GenerateToken mocks base method.
func (m *MockAuthorization) GenerateToken(username, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", username, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthorizationMockRecorder) GenerateToken(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuthorization)(nil).GenerateToken), username, password)
}

// GetByName mocks base method.
func (m *MockAuthorization) GetByName(name string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", name)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockAuthorizationMockRecorder) GetByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockAuthorization)(nil).GetByName), name)
}

// GetUserById mocks base method.
func (m *MockAuthorization) GetUserById(userId int) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", userId)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockAuthorizationMockRecorder) GetUserById(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockAuthorization)(nil).GetUserById), userId)
}

// ParseToken mocks base method.
func (m *MockAuthorization) ParseToken(token string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", token)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockAuthorizationMockRecorder) ParseToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockAuthorization)(nil).ParseToken), token)
}

// UpdateData mocks base method.
func (m *MockAuthorization) UpdateData(user models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateData", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateData indicates an expected call of UpdateData.
func (mr *MockAuthorizationMockRecorder) UpdateData(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateData", reflect.TypeOf((*MockAuthorization)(nil).UpdateData), user)
}

// UpdatePassword mocks base method.
func (m *MockAuthorization) UpdatePassword(user models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockAuthorizationMockRecorder) UpdatePassword(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockAuthorization)(nil).UpdatePassword), user)
}

// MockChat is a mock of Chat interface.
type MockChat struct {
	ctrl     *gomock.Controller
	recorder *MockChatMockRecorder
}

// MockChatMockRecorder is the mock recorder for MockChat.
type MockChatMockRecorder struct {
	mock *MockChat
}

// NewMockChat creates a new mock instance.
func NewMockChat(ctrl *gomock.Controller) *MockChat {
	mock := &MockChat{ctrl: ctrl}
	mock.recorder = &MockChatMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChat) EXPECT() *MockChatMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockChat) AddUser(users models.ChatUsers) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", users)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockChatMockRecorder) AddUser(users interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockChat)(nil).AddUser), users)
}

// Create mocks base method.
func (m *MockChat) Create(chat models.Chat) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", chat)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockChatMockRecorder) Create(chat interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockChat)(nil).Create), chat)
}

// Delete mocks base method.
func (m *MockChat) Delete(chatId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", chatId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockChatMockRecorder) Delete(chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockChat)(nil).Delete), chatId)
}

// DeleteAllMessages mocks base method.
func (m *MockChat) DeleteAllMessages(chatId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllMessages", chatId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllMessages indicates an expected call of DeleteAllMessages.
func (mr *MockChatMockRecorder) DeleteAllMessages(chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllMessages", reflect.TypeOf((*MockChat)(nil).DeleteAllMessages), chatId)
}

// DeleteUser mocks base method.
func (m *MockChat) DeleteUser(userId, chatId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", userId, chatId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockChatMockRecorder) DeleteUser(userId, chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockChat)(nil).DeleteUser), userId, chatId)
}

// Get mocks base method.
func (m *MockChat) Get(chatId int) (models.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", chatId)
	ret0, _ := ret[0].(models.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockChatMockRecorder) Get(chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockChat)(nil).Get), chatId)
}

// GetPrivateChats mocks base method.
func (m *MockChat) GetPrivateChats(userId int) ([]models.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrivateChats", userId)
	ret0, _ := ret[0].([]models.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPrivateChats indicates an expected call of GetPrivateChats.
func (mr *MockChatMockRecorder) GetPrivateChats(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrivateChats", reflect.TypeOf((*MockChat)(nil).GetPrivateChats), userId)
}

// GetPrivates mocks base method.
func (m *MockChat) GetPrivates(firstUser, secondUser int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrivates", firstUser, secondUser)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPrivates indicates an expected call of GetPrivates.
func (mr *MockChatMockRecorder) GetPrivates(firstUser, secondUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrivates", reflect.TypeOf((*MockChat)(nil).GetPrivates), firstUser, secondUser)
}

// GetPublicChats mocks base method.
func (m *MockChat) GetPublicChats(userId int) ([]models.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicChats", userId)
	ret0, _ := ret[0].([]models.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPublicChats indicates an expected call of GetPublicChats.
func (mr *MockChatMockRecorder) GetPublicChats(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicChats", reflect.TypeOf((*MockChat)(nil).GetPublicChats), userId)
}

// GetUserById mocks base method.
func (m *MockChat) GetUserById(userId int) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", userId)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockChatMockRecorder) GetUserById(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockChat)(nil).GetUserById), userId)
}

// GetUsers mocks base method.
func (m *MockChat) GetUsers(chatId int) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", chatId)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockChatMockRecorder) GetUsers(chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockChat)(nil).GetUsers), chatId)
}

// SearchChat mocks base method.
func (m *MockChat) SearchChat(name string) ([]models.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchChat", name)
	ret0, _ := ret[0].([]models.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchChat indicates an expected call of SearchChat.
func (mr *MockChatMockRecorder) SearchChat(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchChat", reflect.TypeOf((*MockChat)(nil).SearchChat), name)
}

// Update mocks base method.
func (m *MockChat) Update(chat models.Chat) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", chat)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockChatMockRecorder) Update(chat interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockChat)(nil).Update), chat)
}

// MockStatus is a mock of Status interface.
type MockStatus struct {
	ctrl     *gomock.Controller
	recorder *MockStatusMockRecorder
}

// MockStatusMockRecorder is the mock recorder for MockStatus.
type MockStatusMockRecorder struct {
	mock *MockStatus
}

// NewMockStatus creates a new mock instance.
func NewMockStatus(ctrl *gomock.Controller) *MockStatus {
	mock := &MockStatus{ctrl: ctrl}
	mock.recorder = &MockStatusMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatus) EXPECT() *MockStatusMockRecorder {
	return m.recorder
}

// AddStatus mocks base method.
func (m *MockStatus) AddStatus(status models.Status) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddStatus", status)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddStatus indicates an expected call of AddStatus.
func (mr *MockStatusMockRecorder) AddStatus(status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddStatus", reflect.TypeOf((*MockStatus)(nil).AddStatus), status)
}

// DeleteStatus mocks base method.
func (m *MockStatus) DeleteStatus(status models.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStatus", status)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStatus indicates an expected call of DeleteStatus.
func (mr *MockStatusMockRecorder) DeleteStatus(status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStatus", reflect.TypeOf((*MockStatus)(nil).DeleteStatus), status)
}

// GetBlackList mocks base method.
func (m *MockStatus) GetBlackList(userId int) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlackList", userId)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlackList indicates an expected call of GetBlackList.
func (mr *MockStatusMockRecorder) GetBlackList(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlackList", reflect.TypeOf((*MockStatus)(nil).GetBlackList), userId)
}

// GetBlackListToUser mocks base method.
func (m *MockStatus) GetBlackListToUser(userId int) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlackListToUser", userId)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlackListToUser indicates an expected call of GetBlackListToUser.
func (mr *MockStatusMockRecorder) GetBlackListToUser(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlackListToUser", reflect.TypeOf((*MockStatus)(nil).GetBlackListToUser), userId)
}

// GetFriends mocks base method.
func (m *MockStatus) GetFriends(userId int) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriends", userId)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriends indicates an expected call of GetFriends.
func (mr *MockStatusMockRecorder) GetFriends(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriends", reflect.TypeOf((*MockStatus)(nil).GetFriends), userId)
}

// GetInvites mocks base method.
func (m *MockStatus) GetInvites(userId int) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvites", userId)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInvites indicates an expected call of GetInvites.
func (mr *MockStatusMockRecorder) GetInvites(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvites", reflect.TypeOf((*MockStatus)(nil).GetInvites), userId)
}

// GetSentInvites mocks base method.
func (m *MockStatus) GetSentInvites(userId int) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSentInvites", userId)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSentInvites indicates an expected call of GetSentInvites.
func (mr *MockStatusMockRecorder) GetSentInvites(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSentInvites", reflect.TypeOf((*MockStatus)(nil).GetSentInvites), userId)
}

// GetStatuses mocks base method.
func (m *MockStatus) GetStatuses(senderId, recipientId int) ([]models.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatuses", senderId, recipientId)
	ret0, _ := ret[0].([]models.Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatuses indicates an expected call of GetStatuses.
func (mr *MockStatusMockRecorder) GetStatuses(senderId, recipientId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatuses", reflect.TypeOf((*MockStatus)(nil).GetStatuses), senderId, recipientId)
}

// GetUserById mocks base method.
func (m *MockStatus) GetUserById(userId int) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", userId)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockStatusMockRecorder) GetUserById(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockStatus)(nil).GetUserById), userId)
}

// SearchUser mocks base method.
func (m *MockStatus) SearchUser(username string) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUser", username)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchUser indicates an expected call of SearchUser.
func (mr *MockStatusMockRecorder) SearchUser(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUser", reflect.TypeOf((*MockStatus)(nil).SearchUser), username)
}

// UpdateStatus mocks base method.
func (m *MockStatus) UpdateStatus(status models.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockStatusMockRecorder) UpdateStatus(status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockStatus)(nil).UpdateStatus), status)
}

// MockMessage is a mock of Message interface.
type MockMessage struct {
	ctrl     *gomock.Controller
	recorder *MockMessageMockRecorder
}

// MockMessageMockRecorder is the mock recorder for MockMessage.
type MockMessageMockRecorder struct {
	mock *MockMessage
}

// NewMockMessage creates a new mock instance.
func NewMockMessage(ctrl *gomock.Controller) *MockMessage {
	mock := &MockMessage{ctrl: ctrl}
	mock.recorder = &MockMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessage) EXPECT() *MockMessageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockMessage) Create(msg models.Message) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", msg)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockMessageMockRecorder) Create(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMessage)(nil).Create), msg)
}

// DeleteAll mocks base method.
func (m *MockMessage) DeleteAll(chatId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", chatId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockMessageMockRecorder) DeleteAll(chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockMessage)(nil).DeleteAll), chatId)
}

// Get mocks base method.
func (m *MockMessage) Get(msgId int) (models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", msgId)
	ret0, _ := ret[0].(models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMessageMockRecorder) Get(msgId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMessage)(nil).Get), msgId)
}

// GetAll mocks base method.
func (m *MockMessage) GetAll(chatId int) ([]models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", chatId)
	ret0, _ := ret[0].([]models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockMessageMockRecorder) GetAll(chatId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMessage)(nil).GetAll), chatId)
}

// GetLimit mocks base method.
func (m *MockMessage) GetLimit(chatId, limit int) ([]models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLimit", chatId, limit)
	ret0, _ := ret[0].([]models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLimit indicates an expected call of GetLimit.
func (mr *MockMessageMockRecorder) GetLimit(chatId, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLimit", reflect.TypeOf((*MockMessage)(nil).GetLimit), chatId, limit)
}
