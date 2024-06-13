package pulumi

/*
import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/codedeploy"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an S3 bucket to store the application deployment package
		bucket, err := s3.NewBucket(ctx, "appBucket", nil)
		if err != nil {
			return err
		}

		// Sync the local folder with the S3 bucket
		_, err = syncedfolder.NewS3BucketFolder(ctx, "appFolder", &syncedfolder.S3BucketFolderArgs{
			Acl:                pulumi.String("private"),
			Path:               pulumi.String("./app"), // Path to local application folder containing binary or source
			BucketName:         bucket.Bucket,
			ManagedObjects:     pulumi.Bool(true),
			IncludeHiddenFiles: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Create an IAM role for the EC2 instance
		role, err := iam.NewRole(ctx, "instanceRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Action": "sts:AssumeRole",
					"Principal": {
						"Service": "ec2.amazonaws.com"
					},
					"Effect": "Allow",
					"Sid": ""
				}]
			}`),
		})
		if err != nil {
			return err
		}

		// Attach the required policies to the role
		_, err = iam.NewRolePolicyAttachment(ctx, "instanceRolePolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      role.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"),
		})
		if err != nil {
			return err
		}

		// Look up the AMI
		ami, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
			Filters: []ec2.GetAmiFilter{
				{Name: "name", Values: []string{"amzn2-ami-hvm-*-x86_64-gp2"}},
			},
			Owners:     []string{"amazon"},
			MostRecent: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Create an EC2 instance
		instance, err := ec2.NewInstance(ctx, "appInstance", &ec2.InstanceArgs{
			InstanceType:       pulumi.String("t2.micro"),
			Ami:                pulumi.String(ami.Id),
			IamInstanceProfile: role.Name,
			UserData: pulumi.String(`#!/bin/bash
			cd /home/ec2-user
			aws s3 cp s3://" + bucket.Bucket + "/app.tar.gz ./app.tar.gz
			tar -xzf app.tar.gz
			cd app
			./app_binary  # Replace with commands to build and run if deploying source code
			`),
		})
		if err != nil {
			return err
		}

		// Create a CodeDeploy application
		app, err := codedeploy.NewApplication(ctx, "appDeploy", nil)
		if err != nil {
			return err
		}

		// Create a CodeDeploy deployment group
		_, err = codedeploy.NewDeploymentGroup(ctx, "appGroup", &codedeploy.DeploymentGroupArgs{
			AppName:              app.Name,
			ServiceRoleArn:       role.Arn,
			DeploymentGroupName:  pulumi.String("appDeployGroup"),
			DeploymentConfigName: pulumi.String("CodeDeployDefault.OneAtATime"),
			Ec2TagSets: []codedeploy.DeploymentGroupEc2TagSetArgs{
				{
					Ec2TagFilters: []codedeploy.DeploymentGroupEc2TagSetEc2TagFilterArgs{
						{
							Key:   pulumi.String("Name"),
							Value: pulumi.String("appInstance"),
							Type:  pulumi.String("KEY_AND_VALUE"),
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("bucketName", bucket.Bucket)
		ctx.Export("instanceId", instance.ID())
		return nil
	})
}
*/
