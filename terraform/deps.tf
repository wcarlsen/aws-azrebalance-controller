data "aws_iam_policy_document" "this" {
  statement {
    sid = "Autoscaling"
    actions = [
      "autoscaling:ResumeProcesses",
      "autoscaling:DescribeAutoScalingGroups",
      "autoscaling:SuspendProcesses",
    ]
    resources = [
      "*",
    ]
  }

  statement {
    sid = "EKS"
    actions = [
      "eks:DescribeNodegroup",
      "eks:ListNodegroups",
    ]
    resources = [
      "*",
    ]
  }
}