//go:generate error_helper -create-map=false -map-name=errMap
package errors

import (
	"github.com/semrush/zenrpc"
)

type ServiceError func(error, interface{}) *zenrpc.Error

var errMap map[int]string

const (
	Internal                                          = iota + 1 //eh:internal error
	UserNotFound                                                 //eh:user not found
	SpecifiedUserDoesNotHaveAvatar                               //eh:user does not have avatar
	EditOtherAvatar                                              //eh:u can't edit avatar of another user
	EnterNickname                                                //eh:enter nickname
	NicknameIsBusy                                               //eh:nickname is busy
	InvalidPhone                                                 //eh:invalid phone
	InvalidNickname                                              //eh:invalid nickname
	InvalidOrgTitle                                              //eh:invalid organization title
	InvalidOrgDescription                                        //eh:invalid organization description
	OrgNotFound                                                  //eh:organization not found
	OrgOperationNotPermitted                                     //eh:operation on organization is not permitted
	UserAlreadyInOrg                                             //eh:user already linked to org
	InvalidOrgNickname                                           //eh:invalid nickname of organization
	AssociateDirector                                            //eh:this user is director of company
	RemoveDirector                                               //eh:this user is director of company
	NotAssociatedUser                                            //eh:this user is not associated with this company
	OrganizationNotFound                                         //eh:organization not found
	NicknameAlreadyBound                                         //eh:nickname already bound
	InvalidSender                                                //eh:invalid sender
	SpecifiedOrgDoesNotHaveAvatar                                //eh:specified organization does not have avatar
	InvalidGroupTitle                                            //eh:invalid group title
	GroupNotFound                                                //eh:group not found
	NotMemberOfParentCreateGroup                                 //eh:create subgroup while not a member of parent, not director, not linked
	CreateGroupWhileNotInOrg                                     //eh:create group while not a member of organization
	CreateRootGroupWhileNotLinked                                //eh:create root group while not linked to organization
	CantDenyCreationRequest                                      //eh:this user can not deny group creation request
	CantAcceptCreationRequest                                    //eh:this user can not accept group creation request
	CreationRequestNotFound                                      //eh:creation request not found
	CreationRequestAlreadyClosed                                 //eh:creation request already closed
	CantWithdrawCreationRequest                                  //eh:this user can not withdraw this request
	FNSError                                                     //eh:fns error
	CantInviteToGroupPersonOutsideOrg                            //eh:cant invite to group person outside org
	CantViewThisGroup                                            //eh:cant view this group
	CantDissociateYourself                                       //eh:cant dissociate yourself
	TreeAccessDenied                                             //eh:access to tree of groups of organization is denied
	NotAssociated                                                //eh:u r not associated to organization
	AlreadyInGroup                                               //eh:target person already in group
	NotInOrg                                                     //eh:u r not in organization
	InvitePersonWhileNotAdminAndNotAssociated                    //eh:invite person while not admin and not associated
	GroupJoinRequestNotFound                                     //eh:group join request not found
	CantDenyGroupJoinRequest                                     //eh:only admin of group or org can deny group join request
	CantAcceptGroupJoinRequest                                   //eh:only admin of group or org can accept group join request
	CantWithdrawGroupJoinRequest                                 //eh:only initiator of group join request can withdraw
	GroupJoinRequestAlreadyClosed                                //eh:group join request already closed
	CantViewGroupJoinRequest                                     //eh:u have not access to this group join request
	IncorrectNotificationID                                      //eh:incorrect notification ID
	NotificationNotFound                                         //eh:notification not found
	CantViewRatingEventsOfThisGroup                              //eh:cant view rating events of this group
	NoRatingsInThisOrgForNow                                     //eh:no ratings in this org for now
	CantAccessRatingOfThisGroup                                  //eh:cant access rating of this group
	CantCreateSubgroup                                           //eh:cant create subgroup: you are not admin of parent group
	RatingEventNotFound                                          //eh:rating event not found
	ThisUserAlreadyEstimatedByU                                  //eh:this user already estimated by u
	RatingEventNotAvailable                                      //eh:rating event not available
	CantEstimateUserNotInMutualGroup                             //eh:cant estimate user not in mutual group
	InvalidEstimate                                              //eh:invalid estimate
	LeaveGroupWhileAdmin                                         //eh:tried to leave the group as its admin
	NotMemberOfGroup                                             //eh:user isn't member of this group
	EfficiencyScoreOutOfRange                                    //eh:efficiency score out of range
	LoyaltyScoreOutOfRange                                       //eh:loyalty score out of range
	TeamworkScoreOutOfRange                                      //eh:teamwork score out of range
	DisciplineScoreOutOfRange                                    //eh:discipline score out of range
	NotAllRequiredQuestionsAnswered                              //eh:not all required questions answered
	CantViewGroupsOfThisOrg                                      //eh:isn't admin of this org
	OrgJoinRequestNotFound                                       //eh:organization join request not found
	CantWithdrawOrgJoinRequest                                   //eh:this user can not withdraw this request
	OrgJoinRequestAlreadyClosed                                  //eh:organization join request already closed
	CantViewOrgJoinRequest                                       //eh:u have no access to this org join request
	GroupIsNotMemberOfOrg                                        //eh:group is not a member of this organization
	CantViewThisGroupCreationRequest                             //eh:cant view this group creation request
	ChatNotFound                                                 //eh:chat not found
	CantAccessChat                                               //eh:u cant access this chat
	OnlyAdminOfGroupCanSetNewAdmin                               //eh:only admin of group can set new admin
	UserNotInGroup                                               //eh:user not in group
	OnlyDirectorCanSetNewDirector                                //eh:only director can set new director
	UserNotAdminOfOrg                                            //eh:user not admin of org
	OnlyAdminOfOrgCanControlRating                               //eh:only admin of org can control rating
	RatingAlreadyEnabled                                         //eh:rating already enabled
	OnlyOrgAdminOrGroupAdminMayDeleteGroup                       //eh:only org admin or group admin may delete group
	OnlyOrgDirectorMayDeleteOrg                                  //eh:only org director may delete drg
	InvalidAuthCode                                              //eh:invalid auth code
	PhoneAlreadyRegistered                                       //eh:phone already registered
	RequestAlreadyOpened                                         //eh:request already opened
	UserDeleteWhileIsDirectorOfAnyOrg                            //eh:user delete while is director of any org
	UserIsNotInOrg                                               //eh:user is not in organization
	OnlyDirectorOrgAdminOrGroupAdminMaySetGroupAvatar            //eh:only org director, org admin or group admin may set group avatar
	UserDeleteWhileIsAdminOfAnyGroup                             //eh:user delete while is admin of any group
	InitiativeScoreOutOfRange                                    //eh:Initiative score out of range
	InvalidPassword                                              //eh:invalid password
)
