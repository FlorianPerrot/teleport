/*
Copyright 2023 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"github.com/gravitational/trace"
	"google.golang.org/protobuf/types/known/timestamppb"

	accesslistv1 "github.com/gravitational/teleport/api/gen/proto/go/teleport/accesslist/v1"
	"github.com/gravitational/teleport/api/types/accesslist"
	headerv1 "github.com/gravitational/teleport/api/types/header/convert/v1"
)

// FromMemberProto converts a v1 access list member into an internal access list member object.
func FromMemberProto(msg *accesslistv1.Member) (*accesslist.AccessListMember, error) {
	if msg == nil {
		return nil, trace.BadParameter("access list message is nil")
	}

	if msg.Spec == nil {
		return nil, trace.BadParameter("spec is missing")
	}

	member, err := accesslist.NewAccessListMember(headerv1.FromMetadataProto(msg.Header.Metadata), accesslist.AccessListMemberSpec{
		AccessList: msg.Spec.AccessList,
		Name:       msg.Spec.Name,
		Joined:     msg.Spec.Joined.AsTime(),
		Expires:    msg.Spec.Expires.AsTime(),
		Reason:     msg.Spec.Reason,
		AddedBy:    msg.Spec.AddedBy,
	})

	return member, trace.Wrap(err)
}

// FromMembersProto converts a list of v1 access list members into a list of internal access list members.
func FromMembersProto(msgs []*accesslistv1.Member) ([]*accesslist.AccessListMember, error) {
	members := make([]*accesslist.AccessListMember, len(msgs))
	for i, msg := range msgs {
		var err error
		members[i], err = FromMemberProto(msg)
		if err != nil {
			return nil, trace.Wrap(err)
		}
	}
	return members, nil
}

// ToMemberProto converts an internal access list member into a v1 access list member object.
func ToMemberProto(member *accesslist.AccessListMember) *accesslistv1.Member {
	return &accesslistv1.Member{
		Header: headerv1.ToResourceHeaderProto(member.ResourceHeader),
		Spec: &accesslistv1.MemberSpec{
			AccessList: member.Spec.AccessList,
			Name:       member.Spec.Name,
			Joined:     timestamppb.New(member.Spec.Joined),
			Expires:    timestamppb.New(member.Spec.Expires),
			Reason:     member.Spec.Reason,
			AddedBy:    member.Spec.AddedBy,
		},
	}
}

// ToMembersProto converts a list of internal access list members into a list of v1 access list members.
func ToMembersProto(members []*accesslist.AccessListMember) []*accesslistv1.Member {
	out := make([]*accesslistv1.Member, len(members))
	for i, member := range members {
		out[i] = ToMemberProto(member)
	}
	return out
}
