package repositories

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type GroupRepository struct {
	core.Repository
	db *gorm.DB
}

func (gr *GroupRepository) CreateGroup(data interface {}) (interface {}, error) {
	group, ok := data.(*models.Group)

    if !ok {
        return nil, gr.GenError(gr.InvalidData, group)
    }

	err := group.BeforeCreateRecord()
	if err != nil {
		return nil, gr.GenError(gr.InvalidData, group)
	}

	err = gr.db.Table("groups").Create(group).Error
	if err != nil {
		return nil, gr.GenError(gr.CreateError, group)
	}

	return *group, nil
}

func (gr *GroupRepository) GetGroup(data interface {}) (interface {}, error) {
	id := data.(string)

	group := &models.Group{}

	err := gr.db.Table("groups").Where("id = ?", id).First(group).Error
	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, group)
	}

	return *group, nil
}

func (gr *GroupRepository) GetGroups(data interface {}) (interface {}, error) {
	groups, ok := data.(*[]models.Group)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, groups)
	}

	err := gr.db.Table("groups").Find(groups).Error
	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, groups)
	}

	return *groups, err
}

func (gr *GroupRepository) DeleteGroup(data interface {}) (interface {}, error) {
	id := data.(string)

	err := gr.db.Table("groups").Where("id = ?", id).Delete(&models.Group{}).Error
	if err != nil {
		return nil, gr.GenError(gr.DeleteError, nil)
	}

	return nil, err
}

func (gr *GroupRepository) AddGroupMember(data interface {}) (interface {}, error) {
	groupMember, ok := data.(*models.GroupMember)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, groupMember)
	}

	err := groupMember.BeforeCreateRecord()
	if err != nil {
		return nil, gr.GenError(gr.InvalidData, groupMember)
	}

	err = gr.db.Table("group_members").Create(groupMember).Error
	if err != nil {
		return nil, gr.GenError(gr.CreateError, groupMember)
	}

	return *groupMember, nil
}

func (gr *GroupRepository) GetGroupMembers(data interface {}) (interface {}, error) {
	groupID, ok := data.(string)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, groupID)
	}

	groupMembers := &[]models.GroupMember{}
	err := gr.db.Table("group_members").Where("group_id = ?", groupID).Find(groupMembers).Error
	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, groupMembers)
	}

	return *groupMembers, err
}

func (gr *GroupRepository) DeleteGroupMember(data interface {}) (interface {}, error) {
	id, ok := data.(string)
	
	if !ok {
		return nil, gr.GenError(gr.InvalidData, id)
	}

	err := gr.db.Table("group_members").Where("id = ?", id).Delete(&models.GroupMember{}).Error
	if err != nil {
		return nil, gr.GenError(gr.DeleteError, nil)
	}

	return nil, err
}

func (gr *GroupRepository) GetGroupMember(data interface {}) (interface {}, error) {
	groupMemberID, ok := data.(string)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, groupMemberID)
	}

	var groupMember models.GroupMember

	err := gr.db.Table("group_members").Where("id = ?", groupMemberID).First(&groupMember).Error
	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, groupMemberID)
	}

	return groupMember, nil
}

func (gr *GroupRepository) GetGroupsByProfile(data interface {}) (interface {}, error) {
	profileID, ok := data.(string)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, profileID)
	}

	groups := &[]models.Group{}

	err := gr.db.Table("groups").
		Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.profile_id = ?", profileID).
		Select("DISTINCT groups.*").
		Find(&groups).Error

	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, groups)
	}

	return *groups, err
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{
		Repository: *core.NewRepository(),
		db: db,
	}
}
