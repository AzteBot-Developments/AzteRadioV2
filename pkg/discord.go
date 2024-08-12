package pkg

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func MemberIsAdmin(guildId string, s *discordgo.Session, i discordgo.Interaction, m discordgo.Member) bool {

	hasAdminPermissions := false

	for _, roleID := range m.Roles {
		role, err := s.State.Role(guildId, roleID)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving role from Discord: %v\n", err)
			continue
		}

		if role.Permissions&discordgo.PermissionAdministrator != 0 {
			hasAdminPermissions = true
			break
		}
	}

	// Check if the member has the "Administrator" permission directly (e.g., server owner or other)
	if i.Member.Permissions&discordgo.PermissionAdministrator != 0 {
		hasAdminPermissions = true
	}

	if hasAdminPermissions {
		return true
	} else {
		return false
	}
}
