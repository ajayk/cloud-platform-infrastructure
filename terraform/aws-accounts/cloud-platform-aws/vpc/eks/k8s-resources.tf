

#########################
# Pod Security Policies #
#########################

resource "null_resource" "pod_security_policy" {
  depends_on = [module.eks.cluster_id]
  provisioner "local-exec" {
    command = "kubectl apply -f ${path.module}/resources/psp/pod-security-policy.yaml"
  }

  provisioner "local-exec" {
    when    = destroy
    command = "kubectl delete --ignore-not-found -f ${path.module}/resources/psp/pod-security-policy.yaml"
  }

  triggers = {
    content = filesha1("${path.module}/resources/psp/pod-security-policy.yaml")
  }
}