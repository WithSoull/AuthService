package user

import (
	"github.com/WithSoull/AuthService/internal/model"
	pb "github.com/WithSoull/AuthService/pkg/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromModelToProtoUserInfo(model model.UserInfo) *pb.UserInfo {
	return &pb.UserInfo{
		Name: model.Name,
		Email: model.Email,
		Role:  pb.Role(model.Role),
	}
}

func FromProtoToModelUserInfo(proto *pb.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name:  proto.GetName(),
		Email: proto.GetEmail(),
		Role:  model.Role(proto.GetRole()),
	}
}

func FromModelToProtoUser(model model.User) *pb.User {
	return &pb.User{
		Id: model.Id,
		UserInfo: FromModelToProtoUserInfo(model.UserInfo),
		CreatedAt: timestamppb.New(model.CreatedAt),
		UpdatedAt: timestamppb.New(model.UpdatedAt),
	}
}

func FromProtoToModelUser(proto *pb.User) model.User {
	return model.User{
		Id: proto.GetId(),
		UserInfo: FromProtoToModelUserInfo(proto.UserInfo),
		CreatedAt: proto.GetCreatedAt().AsTime(),
		UpdatedAt: proto.GetUpdatedAt().AsTime(),
	}
}
