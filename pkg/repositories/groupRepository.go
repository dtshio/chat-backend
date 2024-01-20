package repositories

import (
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"gorm.io/gorm"
)

type GroupRepository struct {
	core.Repository
}

func (gr *GroupRepository) CreateGroup(db *gorm.DB, data interface {}) (interface {}, error) {
	group, ok := data.(*models.Group)

    if !ok {
        return nil, gr.GenError(gr.InvalidData, group)
    }

	err := group.BeforeCreateRecord()
	if err != nil {
		return nil, gr.GenError(gr.InvalidData, group)
	}

	err = db.Table("groups").Create(group).Error
	if err != nil {
		return nil, gr.GenError(gr.CreateError, group)
	}

	return *group, nil
}

func (gr *GroupRepository) GetGroup(db *gorm.DB, data interface {}) (interface {}, error) {
	id := data.(string)

	group := &models.Group{}

	err := db.Table("groups").Where("id = ?", id).First(group).Error
	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, group)
	}

	return *group, nil
}

func (gr *GroupRepository) GetGroups(db *gorm.DB, data interface {}) (interface {}, error) {
	groups, ok := data.(*[]models.Group)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, groups)
	}

	err := db.Table("groups").Find(groups).Error
	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, groups)
	}

	return *groups, err
}

func (gr *GroupRepository) DeleteGroup(db *gorm.DB, data interface {}) (interface {}, error) {
	id := data.(string)

	err := db.Table("groups").Where("id = ?", id).Delete(&models.Group{}).Error
	if err != nil {
		return nil, gr.GenError(gr.DeleteError, nil)
	}

	return nil, err
}

func (gr *GroupRepository) AddGroupMember(db *gorm.DB, data interface {}) (interface {}, error) {
	groupMember, ok := data.(*models.GroupMember)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, groupMember)
	}

	err := groupMember.BeforeCreateRecord()
	if err != nil {
		return nil, gr.GenError(gr.InvalidData, groupMember)
	}

	err = db.Table("group_members").Create(groupMember).Error
	if err != nil {
		return nil, gr.GenError(gr.CreateError, groupMember)
	}

	return *groupMember, nil
}

func (gr *GroupRepository) GetGroupMembers(db *gorm.DB, data interface {}) (interface {}, error) {
	groupID, ok := data.(string)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, groupID)
	}

	groupMembers := &[]models.GroupMember{}
	err := db.Table("group_members").Where("group_id = ?", groupID).Find(groupMembers).Error
	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, groupMembers)
	}

	return *groupMembers, err
}

func (gr *GroupRepository) DeleteGroupMember(db *gorm.DB, data interface {}) (interface {}, error) {
	id, ok := data.(string)
	
	if !ok {
		return nil, gr.GenError(gr.InvalidData, id)
	}

	err := db.Table("group_members").Where("id = ?", id).Delete(&models.GroupMember{}).Error
	if err != nil {
		return nil, gr.GenError(gr.DeleteError, nil)
	}

	return nil, err
}

func (gr *GroupRepository) GetGroupMember(db *gorm.DB, data interface {}) (interface {}, error) {
	groupMemberID, ok := data.(string)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, groupMemberID)
	}

	var groupMember models.GroupMember

	err := db.Table("group_members").Where("id = ?", groupMemberID).First(&groupMember).Error
	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, groupMemberID)
	}

	return groupMember, nil
}

func (gr *GroupRepository) GetGroupsByProfile(db *gorm.DB, data interface {}) (interface {}, error) {
	profileID, ok := data.(string)

	if !ok {
		return nil, gr.GenError(gr.InvalidData, profileID)
	}

	groups := &[]models.Group{}

	db = db.Table("group_members").Joins("JOIN groups ON group_members.group_id = groups.id")
	db = db.Where("group_members.profile_id = ?", profileID)
	err := db.Select("DISTINCT groups.*").Find(&groups).Error
	if err != nil {
		return nil, gr.GenError(gr.NotFoundError, groups)
	}

	return *groups, err
}

func NewGroupRepository() *GroupRepository {
	repo := &GroupRepository{}

	return &GroupRepository{
		Repository: *core.NewRepository(
			&models.Group{},
			[]core.RepositoryMethod{
				repo.CreateGroup,
				repo.GetGroup,
				repo.GetGroups,
				repo.DeleteGroup,
				repo.AddGroupMember,
				repo.GetGroupMembers,
				repo.DeleteGroupMember,
				repo.GetGroupMember,
				repo.GetGroupsByProfile,
			},
		),
	}
}
