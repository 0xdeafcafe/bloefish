package app

import (
	"github.com/0xdeafcafe/bloefish/services/skillset/internal/domain/ports"
)

type App struct {
	SkillSetRepository ports.SkillSetRepository
}
