package skill

import (
	"embed"

	"github.com/xhd2015/skills/skillcmd"
)

//go:embed SKILL.md
var skillContent string

// Nested topics use path/TOPIC.md (not nested SKILL.md).
//
//go:embed android-images
var skillTreeFS embed.FS

const SkillName = "analyse-my-disk-space"

const skillHelp = `Usage: disk-usage-analyser skill --show [--header] [<topic-path>]
       disk-usage-analyser skill <topic-path> --show [--header]
       disk-usage-analyser skill --install [OPTIONS] [<dir>]
       disk-usage-analyser skill --list

Show the root SKILL.md index or a nested topic (path/TOPIC.md).
Install copies SKILL.md and nested TOPIC.md topics into agent skill directories.
List prints the skill name and every available topic path.
--help also lists available topics (see below).

Examples:
  disk-usage-analyser skill --show
  disk-usage-analyser skill --show android-images
  disk-usage-analyser skill android-images --show
  disk-usage-analyser skill --list
  disk-usage-analyser skill --install --dry-run
  disk-usage-analyser skill --install --help

Options:
  --show [--header] [path]   Print skill or topic content (header-only with --header)
  --install [OPTIONS] [dir]  Install skill files (see --install --help)
  --list                     Print skill name and all topic paths
  -h, --help                 Show this help and available topics
`

// Single returns the embedded Shape-3 single-skill host.
func Single() *skillcmd.SingleSkill {
	return &skillcmd.SingleSkill{
		Name:        SkillName,
		RootContent: skillContent,
		TreeFS:      skillTreeFS,
		Usage:       "disk-usage-analyser skill --install",
		Help:        skillHelp,
	}
}

// Handle dispatches skill --show / --install / --list / --help.
func Handle(args []string) error {
	return Single().Handle(args)
}

// HandleInstall is the top-level install alias of skill --install.
func HandleInstall(args []string) error {
	return Single().Handle(append([]string{"--install"}, args...))
}
