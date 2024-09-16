resource "aws_iam_role_policy" "this" {
  name   = var.role_name
  role   = module.this.iam_role_name
  policy = data.aws_iam_policy_document.this.json
}

module "this" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  version = "~> 5.44.0"

  role_name          = var.role_name
  policy_name_prefix = var.policy_name_prefix

  oidc_providers = {
    ex = {
      provider_arn               = var.oidc_provider_arn
      namespace_service_accounts = ["${var.namespace}:${var.serviceaccount_name}"]
    }
  }
}