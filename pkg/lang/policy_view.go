package lang

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/object"
)

// PolicyView an struct which allows to view/manage policy on behalf on a certain user
// It will enforce all ACLs, allowing the user to only perform actions which he is allowed to perform
type PolicyView struct {
	Policy *Policy
	User   *User
}

// NewPolicyView creates a new PolicyView
func NewPolicyView(policy *Policy, user *User) *PolicyView {
	return &PolicyView{
		Policy: policy,
		User:   user,
	}
}

// AddObject adds an object into the policy, putting it into the corresponding namespace
// If an ACL-related error occurs or user doesn't have permissions to perform an operation, then ACL error will be returned
func (policyView *PolicyView) AddObject(obj object.Base) error {
	privilege, err := policyView.Policy.aclResolver.GetUserPrivileges(policyView.User, obj)
	if err != nil {
		return err
	}
	if !privilege.Manage {
		return fmt.Errorf("user '%s' doesn't have permissions to manage object '%s/%s/%s'", policyView.User.ID, obj.GetNamespace(), obj.GetKind(), obj.GetName())
	}
	policyView.Policy.AddObject(obj)
	return nil
}

// ViewObject looks up and returns an object from the policy, given its kind, locator ([namespace/]name), and namespace relative to which the call is being made
// If policy lookup error occurs or user doesn't have permissions to view an object, then ACL error will be returned
func (policyView *PolicyView) ViewObject(kind string, locator string, currentNs string) (object.Base, error) {
	obj, err := policyView.Policy.GetObject(kind, locator, currentNs)
	if err != nil {
		return nil, err
	}
	privilege, err := policyView.Policy.aclResolver.GetUserPrivileges(policyView.User, obj)
	if err != nil {
		return nil, err
	}
	if !privilege.View {
		return nil, fmt.Errorf("user '%s' doesn't have permissions to view object '%s/%s/%s'", policyView.User.ID, obj.GetNamespace(), obj.GetKind(), obj.GetName())
	}
	return obj, nil
}

// ManageObject looks up and returns an object from the policy, given its kind, locator ([namespace/]name), and namespace relative to which the call is being made
// If policy lookup error occurs or user doesn't have permissions to manage an object, then ACL error will be returned
func (policyView *PolicyView) ManageObject(kind string, locator string, currentNs string) (object.Base, error) {
	obj, err := policyView.Policy.GetObject(kind, locator, currentNs)
	if err != nil {
		return nil, err
	}
	privilege, err := policyView.Policy.aclResolver.GetUserPrivileges(policyView.User, obj)
	if err != nil {
		return nil, err
	}
	if !privilege.Manage {
		return nil, fmt.Errorf("user '%s' doesn't have permissions to manage object '%s/%s/%s'", policyView.User.ID, obj.GetNamespace(), obj.GetKind(), obj.GetName())
	}
	return obj, nil
}