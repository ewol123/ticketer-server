package json

import (
	"encoding/json"
	u "github.com/ewol123/ticketer-server/user-service/user"
	"github.com/pkg/errors"
)

// User : define our user struct so we can implement user serializer interface on it
type User struct{}

type TypeName string

const (
	USER TypeName = "User"
	ROLE TypeName = "Role"
	GET_USER_REQUEST_MODEL TypeName = "GetUserRequestModel"
	GET_USER_RESPONSE_MODEL TypeName = "GetUserResponseModel"
	GET_ALL_USER_REQUEST_MODEL TypeName = "GetAllUserRequestModel"
	GET_ALL_USER_RESPONSE_MODEL TypeName = "GetAllUserResponseModel"
	UPDATE_USER_REQUEST_MODEL TypeName = "UpdateUserRequestModel"
	REGISTER_REQUEST_MODEL TypeName = "RegisterRequestModel"
	REGISTER_RETURN_MODEL TypeName = "RegisterReturnModel"
	DELETE_USER_REQUEST_MODEL TypeName = "DeleteUserRequestModel"
	CONFIRM_REGISTRATION_REQUEST_MODEL TypeName = "ConfirmRegistrationRequestModel"
	LOGIN_REQUEST_MODEL TypeName = "LoginRequestModel"
	LOGIN_RESPONSE_MODEL TypeName = "LoginResponseModel"
	RESET_PASSWORD_REQUEST_MODEL TypeName = "ResetPasswordRequestModel"
	SEND_PASSWD_RESET_REQUEST_MODEL TypeName = "SendPasswdResetRequestModel"
)

func TypeFactory(typeName TypeName) (interface{}, error) {
	switch typeName {
	case USER: return &u.User{}, nil
	case ROLE: return &u.Role{}, nil
	case GET_USER_REQUEST_MODEL: return &u.GetUserRequestModel{}, nil
	case GET_USER_RESPONSE_MODEL: return &u.GetUserResponseModel{}, nil
	case GET_ALL_USER_REQUEST_MODEL: return &u.GetAllUserRequestModel{}, nil
	case GET_ALL_USER_RESPONSE_MODEL: return &u.GetAllUserResponseModel{}, nil
	case UPDATE_USER_REQUEST_MODEL: return &u.UpdateUserRequestModel{},nil
	case REGISTER_REQUEST_MODEL: return &u.RegisterRequestModel{}, nil
	case REGISTER_RETURN_MODEL: return &u.RegisterReturnModel{}, nil
	case DELETE_USER_REQUEST_MODEL: return &u.DeleteUserRequestModel{}, nil
	case CONFIRM_REGISTRATION_REQUEST_MODEL: return &u.ConfirmRegistrationRequestModel{}, nil
	case LOGIN_REQUEST_MODEL: return &u.LoginRequestModel{}, nil
	case LOGIN_RESPONSE_MODEL: return &u.LoginResponseModel{}, nil
	case RESET_PASSWORD_REQUEST_MODEL: return &u.ResetPasswordRequestModel{}, nil
	case SEND_PASSWD_RESET_REQUEST_MODEL: return &u.SendPasswdResetRequestModel{}, nil
	}

	return nil, errors.New("oops, no such typeName")
}

// Decode : decode our input to a struct
func (r *User) Decode(input []byte, typeName TypeName) (interface{}, error) {
	newType, err := TypeFactory(typeName)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(input, newType); err != nil {
		return nil, errors.Wrap(err, "serializer.User.Decode")
	}
	return &newType, nil
}


// Encode : encode any interface to bytes
func (r *User) Encode(input interface{}) ([]byte, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.User.EncodeAll")
	}

	return raw, nil
}
