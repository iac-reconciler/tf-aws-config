package compare

const (
	configComplianceResourceType      = "AWS::Config::ResourceCompliance"
	resourceContains                  = "Contains"
	resourceAttachedToInstance        = "Is attached to Instance"
	resourceTypeStack                 = "AWS::CloudFormation::Stack"
	resourceTypeElasticBeanstalk      = "AWS::ElasticBeanstalk::Application"
	resourceTypeRouteTable            = "AWS::EC2::RouteTable"
	resourceTypeRouteTableAssociation = "AWS::EC2::SubnetRouteTableAssociation"
	resourceTypeVPCEndpoint           = "AWS::EC2::VPCEndpoint"
	resourceTypeENI                   = "AWS::EC2::NetworkInterface"
	resourceTypeEBSVolume             = "AWS::EC2::Volume"
	resourceTypeEC2Instance           = "AWS::EC2::Instance"
	resourceTypeSecurityGroup         = "AWS::EC2::SecurityGroup"
	resourceTypeEksCluster            = "AWS::EKS::Cluster"
	resourceTypeASG                   = "AWS::AutoScaling::AutoScalingGroup"
	resourceTypeELB                   = "AWS::ElasticLoadBalancing::LoadBalancer"
	resourceTypeELBV2                 = "AWS::ElasticLoadBalancingV2::LoadBalancer"
	resourceTypeRDSInstance           = "AWS::RDS::DBInstance"
	resourceTypeLambda                = "AWS::Lambda::Function"
	resourceTypeAlarm                 = "AWS::CloudWatch::Alarm"
	eksEniOwnerTagName                = "eks:eni:owner"
	eksEniOwnerTagValue               = "eks-vpc-resource-controller"
	terraformIAMPolicyType            = "aws_iam_policy"
	nlb                               = "network_load_balancer"
	eksClusterOwnerTagNamePrefix      = "kubernetes.io/cluster/"
	owned                             = "owned"
	awsELBOwner                       = "amazon-elb"
	lambdaInterfaceType               = "lambda"
	lambdaPrefix                      = "AWS Lambda VPC ENI-"
	elbPrefix                         = "ELB "
	cloudWatchNamespaceELB            = "AWS/ELB"
	k8sInstanceTag                    = "node.k8s.amazonaws.com/instance_id"
	rdsENI                            = "RDSNetworkInterface"
	elbArnPrefix                      = "arn:aws:elasticloadbalancing"
	dimensionLoadBalancerName         = "LoadBalancerName"
	ingress                           = "ingress"
	egress                            = "egress"

	terraformAWSRegistryProvider   = `provider["registry.terraform.io/hashicorp/aws"]"`
	terraformAWSProvider           = "provider.aws"
	terraformAWSProviderSuffix     = "provider.aws"
	terraformTypeSecurityGroupRule = "aws_security_group_rule"
)
