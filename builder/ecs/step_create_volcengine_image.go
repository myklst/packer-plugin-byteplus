package ecs

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type stepCreateVolcengineImage struct {
	VolcengineEcsConfig *VolcengineEcsConfig
}

func (s *stepCreateVolcengineImage) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packer.Ui)
	client := stateBag.Get("client").(*VolcengineClientWrapper)
	instanceId := stateBag.Get("instanceId").(string)
	tags := s.VolcengineEcsConfig.Tags
	var volcTags []*ecs.TagForCreateImageInput
	if len(tags) > 0 {
		volcTags = make([]*ecs.TagForCreateImageInput, 0)
		for _, tag := range tags {
			volcTags = append(volcTags, &ecs.TagForCreateImageInput{
				Key:   volcengine.String(tag.Key),
				Value: volcengine.String(tag.Value),
			})
		}
	}
	ui.Say(fmt.Sprintf("create new image "))
	input := ecs.CreateImageInput{
		InstanceId: volcengine.String(instanceId),
		ImageName:  volcengine.String(s.VolcengineEcsConfig.TargetImageName),
	}
	if len(volcTags) > 0 {
		input.Tags = volcTags
	}
	output, err := client.EcsClient.CreateImageWithContext(ctx, &input)
	if err != nil {
		return Halt(stateBag, err, "Error create image")
	}
	_, err = WaitImageStatus(stateBag, *output.ImageId, "available")
	if err != nil {
		return Halt(stateBag, err, "Error stop ecs instance")
	}
	stateBag.Put("TargetImageId", *output.ImageId)
	return multistep.ActionContinue
}

func (stepCreateVolcengineImage) Cleanup(bag multistep.StateBag) {

}
