package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/kevinmarcellius/go-simple-auth/internal/model"
	"github.com/kevinmarcellius/go-simple-auth/internal/repository/mocks"
	"github.com/kevinmarcellius/go-simple-auth/internal/utils"
)

func TestUserService_Login(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")
	testCases := []struct {
		name        string
		req         model.LoginRequest
		mockRepo    func(mock *mocks.MockUserRepository)
		expectedMsg string
		expectedErr error
	}{
		{
			name: "Success",
			req: model.LoginRequest{
				Email: "test@mail.id",
				Password: "password123",
			},
			mockRepo: func(mock *mocks.MockUserRepository) {
				mock.EXPECT().GetUserByEmail("test@mail.id").Return(model.User{
					ID:           uuid.New(),
					Email:        "test@mail.id",
					PasswordHash: hashedPassword,
					IsAdmin: false,
					Username: "testuser",
				}, nil)
			},
			expectedMsg: "Login successful",
			expectedErr: nil,
	},}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			
			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			tc.mockRepo(mockUserRepo)
			
			userService := NewUserService(mockUserRepo, "test-secret-key")
			_, err := userService.Login(context.Background(), tc.req)
			
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				
			} else {
				assert.NoError(t, err)
			
			}
		})
	}
		
}

func TestUserService_CreateUser(t *testing.T) {
	testCases := []struct {
		name        string
		req         model.UserRequest
		mockRepo    func(mock *mocks.MockUserRepository)
		expectedMsg string
		expectedErr error
	}{
		{
			name: "Success",
			req: model.UserRequest{
				Username: "newuser",
				Email:    "new@example.com",
				Password: "password123",
			},
			mockRepo: func(mock *mocks.MockUserRepository) {
				// We expect CreateUser to be called with any user object, since the ID and hashed password are created inside the service
				mock.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			expectedMsg: "User created successfully",
			expectedErr: nil,
		},
		{
			name: "Database Error",
			req: model.UserRequest{
				Username: "newuser",
				Email:    "new@example.com",
				Password: "password123",
			},
			mockRepo: func(mock *mocks.MockUserRepository) {
				mock.EXPECT().CreateUser(gomock.Any()).Return(assert.AnError)
			},
			expectedMsg: "Failed to create user",
			expectedErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			tc.mockRepo(mockUserRepo)

			userService := NewUserService(mockUserRepo, "test-secret-key")

			res, err := userService.CreateUser(context.Background(), tc.req)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				assert.Equal(t, tc.expectedMsg, res.Message)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMsg, res.Message)
			}
		})
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	password := "oldpassword"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	testUser := model.User{
		ID:           uuid.New(),
		PasswordHash: hashedPassword,
	}

	testCases := []struct {
		name        string
		userID      string
		req         model.UpdatePasswordRequest
		mockRepo    func(mock *mocks.MockUserRepository)
		expectedErr error
	}{
		{
			name:   "Success",
			userID: testUser.ID.String(),
			req: model.UpdatePasswordRequest{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			},
			mockRepo: func(mock *mocks.MockUserRepository) {
				mock.EXPECT().FindUserByID(testUser.ID).Return(testUser, nil)
				// We expect the update to be called with a user model where the password hash has changed
				mock.EXPECT().UpdateUserById(testUser.ID, gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "User Not Found",
			userID: uuid.New().String(),
			req:    model.UpdatePasswordRequest{},
			mockRepo: func(mock *mocks.MockUserRepository) {
				mock.EXPECT().FindUserByID(gomock.Any()).Return(model.User{}, assert.AnError)
			},
			expectedErr: assert.AnError,
		},
		{
			name:   "Invalid Old Password",
			userID: testUser.ID.String(),
			req: model.UpdatePasswordRequest{
				OldPassword: "wrongpassword",
				NewPassword: "newpassword",
			},
			mockRepo: func(mock *mocks.MockUserRepository) {
				mock.EXPECT().FindUserByID(testUser.ID).Return(testUser, nil)
			},
			expectedErr: utils.ErrInvalidPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepository(ctrl)
			tc.mockRepo(mockUserRepo)

			userService := NewUserService(mockUserRepo, "test-secret-key")

			err := userService.UpdatePassword(context.Background(), tc.userID, tc.req)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_Refresh(t *testing.T) {
    jwtKey := "test-secret-key"
    testUser := model.User{
        ID:       uuid.New(),
        Username: "testuser",
        Email:    "test@example.com",
    }

    // Generate a valid refresh token for the test user
    _, refreshToken, err := utils.GenerateJWT(testUser, jwtKey)
    assert.NoError(t, err)

    testCases := []struct {
        name          string
        req           model.RefreshTokenRequest
        mockRepo      func(mock *mocks.MockUserRepository)
        expectedToken bool
        expectedErr   error
    }{
        {
            name: "Success",
            req:  model.RefreshTokenRequest{RefreshToken: refreshToken},
            mockRepo: func(mock *mocks.MockUserRepository) {
                // The user ID inside the token should be used to find the user
                mock.EXPECT().FindUserByID(testUser.ID).Return(testUser, nil)
            },
            expectedToken: true,
            expectedErr:   nil,
        },
        {
            name: "Invalid Token",
            req:  model.RefreshTokenRequest{RefreshToken: "invalid-token"},
            mockRepo: func(mock *mocks.MockUserRepository) {
                // No repo calls should be made if the token is invalid
            },
            expectedToken: false,
            expectedErr:   assert.AnError, // Expect a generic error from the JWT library
        },
        {
            name: "User Not Found From Token",
            req:  model.RefreshTokenRequest{RefreshToken: refreshToken},
            mockRepo: func(mock *mocks.MockUserRepository) {
                mock.EXPECT().FindUserByID(testUser.ID).Return(model.User{}, assert.AnError)
            },
            expectedToken: false,
            expectedErr:   assert.AnError,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockUserRepo := mocks.NewMockUserRepository(ctrl)
            tc.mockRepo(mockUserRepo)

            userService := NewUserService(mockUserRepo, jwtKey)

            res, err := userService.Refresh(context.Background(), tc.req)

            if tc.expectedErr != nil {
                assert.Error(t, err)
                // For generic errors, we just check that an error occurred
                if tc.expectedErr != assert.AnError {
                    assert.Equal(t, tc.expectedErr, err)
                }
                assert.Empty(t, res.AccessToken)
            } else {
                assert.NoError(t, err)
                if tc.expectedToken {
                    assert.NotEmpty(t, res.AccessToken)
                }
            }
        })
    }
}
