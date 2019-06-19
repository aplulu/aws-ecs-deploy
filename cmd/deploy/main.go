package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"

	"github.com/aplulu/aws-ecs-deploy/pkg/config"
)

func main() {
	config.ParseFlags()
	if err := config.LoadConf(); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}
	conf := config.GetConf()

	sess, err := session.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %v\n", err)
		os.Exit(1)
	}

	svc := ecs.New(sess)

	td, err := describeTaskDefinition(svc, conf.TaskDefinition)
	if err != nil {
		fmt.Printf("Failed to retrieve TaskDefinition: %v\n", err)
		os.Exit(1)
	}

	cds := td.ContainerDefinitions
	for i := 0; i < len(cds); i++ {
		if cds[i].Name != nil && *cds[i].Name == conf.Container {
			cds[i].Image = &conf.Image
		}
	}

	ntd, err := registerTaskDefinition(svc, conf.TaskDefinition, cds, td)
	if err != nil {
		fmt.Printf("Failed to register TaskDefinition: %v\n", err)
		os.Exit(1)
	}

	ntdf := fmt.Sprintf("%s:%d", *ntd.Family, *ntd.Revision)
	fmt.Printf("Registered New TaskDefinition: %s\n", *ntd.TaskDefinitionArn)

	if conf.Service == "" {
		fmt.Printf("Deployed.\n")
		os.Exit(0)
	}

	if err := updateService(svc, conf.Cluster, conf.Service, ntdf); err != nil {
		fmt.Printf("Failed to update Service: %v\n", err)
		os.Exit(1)
	}

	if conf.SkipVerify {
		fmt.Printf("Deployed.\n")
		os.Exit(0)
	}

	for attempt := 0; attempt <= conf.WaitCount; attempt++ {
		s, err := describeService(svc, conf.Cluster, conf.Service)
		if err != nil {
			fmt.Printf("Failed to Retrieve Service: %v\n", err)
			os.Exit(1)
		}

		log.Println(*s.TaskDefinition)

		if *s.TaskDefinition == *ntd.TaskDefinitionArn {
			fmt.Printf("Deployed.\n")
			os.Exit(0)
		}

		fmt.Printf("Waiting deployment...\n")
		time.Sleep(time.Duration(conf.WaitSleep) * time.Second)
	}

	fmt.Printf("Service update timeout\n")
	os.Exit(1)
}

func describeTaskDefinition(svc *ecs.ECS, family string) (*ecs.TaskDefinition, error) {
	in := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &family,
	}
	result, err := svc.DescribeTaskDefinition(in)
	if err != nil {
		return nil, err
	}
	return result.TaskDefinition, nil
}

func describeService(svc *ecs.ECS, clusterName string, serviceName string) (*ecs.Service, error) {
	in := &ecs.DescribeServicesInput{
		Cluster: &clusterName,
		Services: []*string{
			&serviceName,
		},
	}
	result, err := svc.DescribeServices(in)
	if err != nil {
		return nil, err
	}
	return result.Services[0], nil
}

func updateService(svc *ecs.ECS, clusterName string, serviceName string, taskArn string) error {
	in := &ecs.UpdateServiceInput{
		Cluster:        &clusterName,
		Service:        &serviceName,
		TaskDefinition: &taskArn,
	}
	_, err := svc.UpdateService(in)
	if err != nil {
		return err
	}
	return nil
}

func registerTaskDefinition(svc *ecs.ECS, family string, containerDefinitions []*ecs.ContainerDefinition, prev *ecs.TaskDefinition) (*ecs.TaskDefinition, error) {
	in := &ecs.RegisterTaskDefinitionInput{
		Family:               &family,
		ContainerDefinitions: containerDefinitions,
	}
	if prev.TaskRoleArn != nil {
		in.TaskRoleArn = prev.TaskRoleArn
	}
	if prev.ExecutionRoleArn != nil {
		in.ExecutionRoleArn = prev.ExecutionRoleArn
	}
	if prev.NetworkMode != nil {
		in.NetworkMode = prev.NetworkMode
	}
	if prev.Volumes != nil {
		in.Volumes = prev.Volumes
	}
	if prev.Cpu != nil {
		in.Cpu = prev.Cpu
	}
	if prev.Memory != nil {
		in.Memory = prev.Memory
	}
	if prev.PlacementConstraints != nil && len(prev.PlacementConstraints) > 0 {
		in.PlacementConstraints = prev.PlacementConstraints
	}
	if prev.RequiresCompatibilities != nil && len(prev.RequiresCompatibilities) > 0 {
		in.RequiresCompatibilities = prev.RequiresCompatibilities
	}

	result, err := svc.RegisterTaskDefinition(in)
	if err != nil {
		return nil, err
	}
	return result.TaskDefinition, nil
}
