package entities

import "github.com/nalej/grpc-authx-go"

// RoleData is the structure that is stored in the provider.
type RoleData struct {
	OrganizationID string
	RoleID         string
	Name           string
	Internal       bool
	Primitives     []string
}

// NewRoleData create a new instance of the structure.
func NewRoleData(organizationID string, roleID string, name string, internal bool, primitives []string) *RoleData {
	return &RoleData{
		OrganizationID: organizationID,
		RoleID:         roleID,
		Name:           name,
		Internal:       internal,
		Primitives:     primitives,
	}
}

func PrimitiveToGRPC(name string) grpc_authx_go.AccessPrimitive {
	switch name {
	case grpc_authx_go.AccessPrimitive_ORG.String():
		return grpc_authx_go.AccessPrimitive_ORG
	case grpc_authx_go.AccessPrimitive_APPS.String():
		return grpc_authx_go.AccessPrimitive_APPS
	case grpc_authx_go.AccessPrimitive_RESOURCES.String():
		return grpc_authx_go.AccessPrimitive_RESOURCES
	case grpc_authx_go.AccessPrimitive_PROFILE.String():
		return grpc_authx_go.AccessPrimitive_PROFILE
	case grpc_authx_go.AccessPrimitive_APPCLUSTEROPS.String():
		return grpc_authx_go.AccessPrimitive_APPCLUSTEROPS
	}
	panic("access primitive not found")
}

func (r *RoleData) ToGRPC() *grpc_authx_go.Role {
	primitives := make([]grpc_authx_go.AccessPrimitive, 0)
	for _, p := range r.Primitives {
		primitives = append(primitives, PrimitiveToGRPC(p))
	}
	return &grpc_authx_go.Role{
		OrganizationId: r.OrganizationID,
		RoleId:         r.RoleID,
		Name:           r.Name,
		Internal:       r.Internal,
		Primitives:     primitives,
	}
}

// EditRoleData is the structure that is used to edit the data in the provider.
type EditRoleData struct {
	Name       *string
	Primitives *[]string
}

//WithName update the name of the role.
func (d *EditRoleData) WithName(name string) *EditRoleData {
	d.Name = &name
	return d
}

//WithPrimitives update the primitives.
func (d *EditRoleData) WithPrimitives(primitives []string) *EditRoleData {
	d.Primitives = &primitives
	return d
}

//NewEditRoleData create a new instance of the structure.
func NewEditRoleData() *EditRoleData {
	return &EditRoleData{}
}
