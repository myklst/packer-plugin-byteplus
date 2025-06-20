//go:generate packer-sdc mapstructure-to-hcl2 -type DataSourceImageConfig,DataSourceImageOutput,Image,Tag,DetectionResults,DetectionItem,TagFilters
package image

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	pc "github.com/volcengine/packer-plugin-volcengine/builder/ecs"
	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/service/vpc"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
	"github.com/zclconf/go-cty/cty"
)

type Datasource struct {
	config DataSourceImageConfig
}

type DataSourceImageConfig struct {
	pc.VolcengineClientConfig `mapstructure:",squash"`
	ImageId                   string       `mapstructure:"image_id"`
	ImageName                 string       `mapstructure:"image_name"`
	Platform                  string       `mapstructure:"platform"`
	Status                    []string     `mapstructure:"status"`
	InstanceTypeId            string       `mapstructure:"instance_type_id"`
	Visibility                string       `mapstructure:"visibility"`
	IsSupportCloudInit        bool         `mapstructure:"is_support_cloud_init"`
	IsLTS                     bool         `mapstructure:"is_lts"`
	ProjectName               string       `mapstructure:"project_name"`
	OsType                    string       `mapstructure:"os_type"`
	TagFilters                []TagFilters `mapstructure:"tag_filters"`
}

type DataSourceImageOutput struct {
	Images []Image `mapstructure:"images"`
}

type Image struct {
	ImageId                  string           `mapstructure:"image_id"`
	ImageName                string           `mapstructure:"image_name"`
	Description              string           `mapstructure:"description"`
	Platform                 string           `mapstructure:"platform"`
	PlatformVersion          string           `mapstructure:"platform_version"`
	Visibility               string           `mapstructure:"visibility"`
	IsSupportCloudInit       bool             `mapstructure:"is_support_cloud_init"`
	OsType                   string           `mapstructure:"os_type"`
	Status                   string           `mapstructure:"status"`
	Architecture             string           `mapstructure:"architecture"`
	OsName                   string           `mapstructure:"os_name"`
	ShareStatus              string           `mapstructure:"share_status"`
	Size                     int32            `mapstructure:"size"`
	BootMode                 string           `mapstructure:"boot_mode"`
	CreatedAt                string           `mapstructure:"created_at"`
	UpdatedAt                string           `mapstructure:"updated_at"`
	LicenseType              string           `mapstructure:"license_type"`
	IsLTS                    bool             `mapstructure:"is_lts"`
	ImageOwnerId             string           `mapstructure:"image_owner_id"`
	Tags                     []Tag            `mapstructure:"tags"`
	Kernel                   string           `mapstructure:"kernel"`
	IsInstallRunCommandAgent bool             `mapstructure:"is_install_run_command_agent"`
	ProjectName              string           `mapstructure:"project_name"`
	DetectionResults         DetectionResults `mapstructure:"detection_results"`
}

type DetectionItem struct {
	Name      string `mapstructure:"name"`
	Result    string `mapstructure:"result"`
	RiskCode  string `mapstructure:"risk_code"`
	RiskLevel string `mapstructure:"risk_level"`
}

type DetectionResults struct {
	DetectionStatus string          `mapstructure:"detection_status"`
	Items           []DetectionItem `mapstructure:"items"`
}

func (i *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return i.config.FlatMapstructure().HCL2Spec()
}

func (i *Datasource) Configure(raws ...interface{}) error {
	if err := config.Decode(&i.config, nil, raws...); err != nil {
		return fmt.Errorf("error decoding config: %s", err)
	}

	var err *packer.MultiError
	if i.config.VolcengineAccessKey == "" {
		err = packer.MultiErrorAppend(err, fmt.Errorf("access_key is required"))
	}
	if i.config.VolcengineSecretKey == "" {
		err = packer.MultiErrorAppend(err, fmt.Errorf("secret_key is required"))
	}
	if i.config.VolcengineRegion == "" {
		err = packer.MultiErrorAppend(err, fmt.Errorf("region is required"))
	}

	if err != nil && len(err.Errors) > 0 {
		return err
	}
	return nil
}

func (i *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DataSourceImageOutput{}).FlatMapstructure().HCL2Spec()
}

func (i *Datasource) getClient() *pc.VolcengineClientWrapper {
	c := volcengine.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(
			i.config.VolcengineAccessKey,
			i.config.VolcengineSecretKey,
			i.config.VolcengineSessionKey),
		).
		//WithDisableSSL(*i.config.VolcengineDisableSSL).
		WithRegion(i.config.VolcengineRegion)

	if i.config.VolcengineEndpoint != "" {
		c.WithEndpoint(i.config.VolcengineEndpoint)
	}

	sess, _ := session.NewSession(c)

	return &pc.VolcengineClientWrapper{
		EcsClient: ecs.New(sess),
		VpcClient: vpc.New(sess),
	}
}

func (i *Datasource) Execute() (cty.Value, error) {
	client := i.getClient()

	describeImageReq := &ecs.DescribeImagesInput{}

	if i.config.ImageName != "" {
		describeImageReq.ImageName = volcengine.String(i.config.ImageName)
	}

	if len(i.config.Status) > 0 {
		describeImageReq.Status = volcengine.StringSlice(i.config.Status)
	}
	if i.config.InstanceTypeId != "" {
		describeImageReq.InstanceTypeId = volcengine.String(i.config.InstanceTypeId)
	}
	if i.config.Visibility != "" {
		describeImageReq.Visibility = volcengine.String(i.config.Visibility)
	}
	if i.config.IsSupportCloudInit {
		describeImageReq.IsSupportCloudInit = volcengine.Bool(i.config.IsSupportCloudInit)
	}
	if i.config.IsLTS {
		describeImageReq.IsLTS = volcengine.Bool(i.config.IsLTS)
	}
	if i.config.ProjectName != "" {
		describeImageReq.ProjectName = volcengine.String(i.config.ProjectName)
	}
	if i.config.OsType != "" {
		describeImageReq.OsType = volcengine.String(i.config.OsType)
	}
	if i.config.Platform != "" {
		describeImageReq.Platform = volcengine.String(i.config.Platform)
	}
	if len(i.config.TagFilters) > 0 {
		var tagFilters []*ecs.TagFilterForDescribeImagesInput
		for _, tagFilter := range i.config.TagFilters {
			tagFilters = append(tagFilters, &ecs.TagFilterForDescribeImagesInput{
				Key:    volcengine.String(tagFilter.Key),
				Values: volcengine.StringSlice(tagFilter.Values),
			})
		}
		describeImageReq.TagFilters = tagFilters
	}

	var dataSourceOutput DataSourceImageOutput
	resp, err := client.EcsClient.DescribeImages(describeImageReq)
	if err != nil {
		return cty.NilVal, err
	}
	var images []Image
	for _, image := range resp.Images {
		var tags []Tag
		for _, tag := range image.Tags {
			tags = append(tags, Tag{
				Key:   volcengine.StringValue(tag.Key),
				Value: volcengine.StringValue(tag.Value),
			})
		}
		var detectionResults DetectionResults
		if image.DetectionResults != nil {
			detectionResults = DetectionResults{
				DetectionStatus: volcengine.StringValue(image.DetectionResults.DetectionStatus),
				Items:           nil,
			}
			for _, item := range image.DetectionResults.Items {
				detectionResults.Items = append(detectionResults.Items, struct {
					Name      string `mapstructure:"name"`
					Result    string `mapstructure:"result"`
					RiskCode  string `mapstructure:"risk_code"`
					RiskLevel string `mapstructure:"risk_level"`
				}{
					Name:      volcengine.StringValue(item.Name),
					Result:    volcengine.StringValue(item.Result),
					RiskCode:  volcengine.StringValue(item.RiskCode),
					RiskLevel: volcengine.StringValue(item.RiskLevel),
				})
			}
		}
		images = append(images, Image{
			ImageId:                  volcengine.StringValue(image.ImageId),
			ImageName:                volcengine.StringValue(image.ImageName),
			Description:              volcengine.StringValue(image.Description),
			Platform:                 volcengine.StringValue(image.Platform),
			PlatformVersion:          volcengine.StringValue(image.PlatformVersion),
			Visibility:               volcengine.StringValue(image.Visibility),
			IsSupportCloudInit:       volcengine.BoolValue(image.IsSupportCloudInit),
			OsType:                   volcengine.StringValue(image.OsType),
			Status:                   volcengine.StringValue(image.Status),
			Architecture:             volcengine.StringValue(image.Architecture),
			OsName:                   volcengine.StringValue(image.OsName),
			ShareStatus:              volcengine.StringValue(image.ShareStatus),
			Size:                     volcengine.Int32Value(image.Size),
			BootMode:                 volcengine.StringValue(image.BootMode),
			CreatedAt:                volcengine.StringValue(image.CreatedAt),
			UpdatedAt:                volcengine.StringValue(image.UpdatedAt),
			LicenseType:              volcengine.StringValue(image.LicenseType),
			IsLTS:                    volcengine.BoolValue(image.IsLTS),
			ImageOwnerId:             volcengine.StringValue(image.ImageOwnerId),
			Tags:                     tags,
			Kernel:                   volcengine.StringValue(image.Kernel),
			IsInstallRunCommandAgent: volcengine.BoolValue(image.IsInstallRunCommandAgent),
			ProjectName:              volcengine.StringValue(image.ProjectName),
			DetectionResults:         detectionResults,
		})
	}
	dataSourceOutput.Images = images

	return hcl2helper.HCL2ValueFromConfig(dataSourceOutput, i.OutputSpec()), nil
}

type TagFilters struct {
	Key    string   `mapstructure:"key"`
	Values []string `mapstructure:"values"`
}

type Tag struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}
