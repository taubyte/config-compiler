package compile

import (
	"fmt"

	projectSchema "github.com/taubyte/go-project-schema/project"
)

func messaging(name string, application string, project projectSchema.Project) (_id string, returnMap map[string]interface{}, err error) {
	iFace, err := project.Messaging(name, application)
	if err != nil {
		return "", nil, err
	}

	getter := iFace.Get()
	_id = getter.Id()

	returnMap = map[string]interface{}{
		"name":        getter.Name(),
		"description": getter.Description(),
		"local":       getter.Local(),
		"regex":       getter.Regex(),
		"match":       getter.ChannelMatch(),
		"mqtt":        getter.MQTT(),
		"webSocket":   getter.WebSocket(),
	}

	_tags := getter.Tags()
	if len(_tags) > 0 {
		returnMap["tags"] = _tags
	}

	err = attachSmartOpsFromTags(returnMap, _tags, application, project, "")
	if err != nil {
		return "", nil, fmt.Errorf("Messaging( %s/`%s` ): Getting smartOps failed with: %v", application, name, err)
	}

	return _id, returnMap, nil
}
