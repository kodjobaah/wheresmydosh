import pprint
import json
import troposphere.ec2 as ec2
import jsonpickle
from troposphere import Base64, GetAtt, ImportValue
from troposphere import Join, Output, Parameter, Ref, Sub, Tags
from troposphere import cloudformation
from troposphere.cloudwatch import Alarm, MetricDimension
from troposphere.cloudformation import InitFile, InitFiles

from troposphere.route53 import RecordSetGroup

from troposphere import autoscaling
from troposphere.autoscaling import AutoScalingGroup,ScalingPolicy

from troposphere.iam import Policy
from awacs.aws import Allow, Action, Principal, Statement, PolicyDocument
from awacs.s3 import ListBucket,GetObject,GetObjectVersion,PutObject,ListBucketVersions
from awacs import iam

from cosmosTroposphere import CosmosTemplate
t = CosmosTemplate(total_az=2,region="eu-west-2",description="used for sharing credit card with family/friends", component_name="wheres-my-dosh", project_name='wheres-my-dosh')


t.parameters["DesiredCapacity"].Default = "1"

t.parameters["MaxSize"].Default = "1"

t.parameters["MinSize"].Default = "1"

t.parameters["UpdateMaxBatchSize"].Default = "1"

t.parameters["UpdateMinInService"].Default = "0"

t.parameters["UpdatePauseTime"].Default = "PT0S"

t.parameters["InstanceType"].Default = "t2.micro"

t.parameters["ImageId"].Default = "ami-09fd46debba20791d"

t.parameters["CnameEntry"].Default = "ugc-file-upload.int"

t.parameters["DomainNameBase"].Default = "c7dff5ab13c48206.xhst.bbci.co.uk."
t.parameters["ElbHealthCheckGracePeriod"].Default = 10000

componentAutoScalingGroup = t.resources['ComponentAutoScalingGroup']

#componentAutoScalingGroup.MinSize = 1
#componentAutoScalingGroup.MaxSize = 1

#componentAutoScalingGroup.DesiredCapacity = 1
#componentAutoScalingGroup.HealthCheckGracePeriod = 0
#componentAutoScalingGroup.HealthCheckType = ""
#componentAutoScalingGroup.UpdatePolicy.AutoScalingRollingUpdate.MaxBatchSize = 0
#componentAutoScalingGroup.UpdatePolicy.AutoScalingRollingUpdate.MinInstancesInService = -1

print(t.to_json())
